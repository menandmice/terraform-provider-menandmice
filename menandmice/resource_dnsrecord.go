package menandmice

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDNSRec() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSRecCreate,
		ReadContext:   resourceDNSRecRead,
		UpdateContext: resourceDNSRecUpdate,
		DeleteContext: resourceDNSRecDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDNSRecImport,
		},
		Schema: map[string]*schema.Schema{

			"ref": {
				Type:        schema.TypeString,
				Description: "Internal reference to this DNS record.",
				Computed:    true,
			},
			"name": {
				Type:         schema.TypeString,
				Description:  "The DNS record name.",
				Required:     true,
				ValidateFunc: validation.StringDoesNotMatch(regexp.MustCompile(`\.$`), "Hostname should not end with '.'"),
			},
			"data": {
				Type:         schema.TypeString,
				Description:  "The data stored in the DNS record.",
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				// You cannot validate data here, because you dont have access to the record type
			},
			"type": {
				Type:        schema.TypeString,
				Description: "The DNS recod type. Accepted types: A, AAAA, CNAME, DNAME, DLV, DNSKEY, DS, HINFO, LOC, MX, NAPTR, NS, NSEC3PARAM, PTR, RP, SOA, SPF, SRV, SSHFP, TLSA, TXT. (Default: A)",
				ForceNew:    true,
				Optional:    true,
				Default:     "A",
				ValidateFunc: validation.StringInSlice([]string{
					"A", "AAAA", "CNAME",
					"DNAME", "DLV", "DNSKEY",
					"DS", "HINFO", "LOC",
					"MX", "NAPTR", "NS", "NSEC3PARAM",
					"PTR", "RP", "SOA",
					"SPF", "SRV", "SSHFP",
					"TLSA", "TXT",
				}, false),
			},
			"comment": {
				Type:        schema.TypeString,
				Description: "Contains the comment string for the record. Only records in static DNS zones can have a comment string. Some cloud DNS provides do not support comments.",
				Optional:    true,
			},
			"aging": {
				Type:         schema.TypeInt,
				Description:  "The aging timestamp of dynamic records in AD integrated zones. Hours since January 1, 1601, UTC. Providing a non-zero value creates a dynamic record.",
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"ttl": {
				Type:         schema.TypeInt,
				Description:  "The DNS record's Time To Live value in seconds, setting how long the record is allowed to be cached.",
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "If the DNS record is enabled. (Default: True)",
				Optional:    true,
				Default:     true,
			},

			"server": {
				Type:         schema.TypeString,
				Description:  "The DNS server where the DNS record is stored. Requires FQDN with the trialing dot '.'.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`\.$`), "Server name should end with '.'"),
			},
			"view": {
				Type:        schema.TypeString,
				Description: "The view of the DNS record. Example: internal.",
				Optional:    true,
				Default:     "",
				ForceNew:    true,
			},

			"zone": {
				Type:         schema.TypeString,
				Description:  "The DNS zone where the record is stored. Requires a trailing dot '.'.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`\.$`), "server should end with '.'"),
			},

			"dns_zone_ref": {
				Type:        schema.TypeString,
				Description: "Internal reference to the zone where this DNS record is stored.",
				Computed:    true,
			},
			"fqdn": {
				Type:        schema.TypeString,
				Description: "Fully qualified domain name of this DNS record.",
				Computed:    true,
			},

			// TODO add force overwrite
		},
	}
}

func writeDNSRecSchema(d *schema.ResourceData, dnsrec DNSRecord) {

	d.Set("ref", dnsrec.Ref)
	d.Set("name", dnsrec.Name)
	d.Set("type", dnsrec.Rectype)
	d.Set("data", dnsrec.Data)
	if dnsrec.TTL != "" {
		ttl, err := strconv.Atoi(dnsrec.TTL)
		if err == nil {
			d.Set("ttl", ttl)
		}
	}

	d.Set("dns_zone_ref", dnsrec.DNSZoneRef)

	// TODO does not set server and view

	if dnsrec.Aging != 0 {
		d.Set("aging", dnsrec.Aging)
	}
	d.Set("enabled", dnsrec.Enabled)
	d.Set("comment", dnsrec.Comment) // comment is always given, but sometimes ""

	// set fqdn for user convience. this not information from the api
	if zone, ok := d.Get("zone").(string); ok && zone != "" {
		d.Set("fqdn", dnsrec.Name+"."+zone) // FIXME
	}
}

func readDNSRecSchema(d *schema.ResourceData) DNSRecord {

	var ttlString string
	if ttl, ok := d.Get("ttl").(int); ok && ttl != 0 {
		ttlString = strconv.Itoa(ttl)
	}

	dnsrec := DNSRecord{
		Ref:        tryGetString(d, "ref"),
		DNSZoneRef: tryGetString(d, "server") + ":" + tryGetString(d, "view") + ":" + tryGetString(d, "zone"),

		Rectype: d.Get("type").(string),
		DNSProperties: DNSProperties{
			Name:    d.Get("name").(string),
			TTL:     ttlString,
			Data:    d.Get("data").(string),
			Comment: d.Get("comment").(string),
			Aging:   d.Get("aging").(int), // when not specified it's 0 which will be ignored
			Enabled: d.Get("enabled").(bool),
		},
	}
	return dnsrec
}

func resourceDNSRecCreate(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Mmclient)

	dnsrec := readDNSRecSchema(d)

	objRef, err := client.CreateDNSRec(dnsrec)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(objRef)

	return resourceDNSRecRead(c, d, m)

}

func resourceDNSRecRead(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	client := m.(*Mmclient)

	dnsrec, err := client.ReadDNSRec(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if dnsrec == nil {
		d.SetId("")
		return diags
	}
	writeDNSRecSchema(d, *dnsrec)

	return diags
}

func resourceDNSRecUpdate(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*Mmclient)
	ref := d.Id()
	dnsrec := readDNSRecSchema(d)
	err := client.UpdateDNSRec(dnsrec.DNSProperties, ref)

	if err != nil {
		return diag.FromErr(err)
	}
	return resourceDNSRecRead(c, d, m)
}

func resourceDNSRecDelete(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*Mmclient)
	var diags diag.Diagnostics
	ref := d.Id()
	err := client.DeleteDNSRec(ref)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}

func resourceDNSRecImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	arg := d.Id()

	var diags diag.Diagnostics
	if parts := strings.Split(arg, ":"); len(parts) == 4 {

		// format is authority:view:fqdn:type
		d.Set("server", parts[0])
		d.Set("view", parts[1])
		fqdn := strings.SplitN(parts[2], ".", 2)

		if len(fqdn) != 2 {
			return nil, fmt.Errorf("Could not parse FQDN %s", parts[2])
		}
		d.Set("name", fqdn[0]) // TODO this not always true
		d.Set("zone", fqdn[1])
		d.Set("type", parts[3])

		// TODO this fragile. call client directly
		diags = dataSourceDNSRectRead(ctx, d, m)
	} else {

		diags = resourceDNSRecRead(ctx, d, m)
	}
	if err := toError(diags); err != nil {
		return nil, err
	}

	ref := d.Get("ref").(string)
	if ref == "" {
		return nil, errors.New("Import failed")
	}
	d.SetId(ref)

	return []*schema.ResourceData{d}, nil
}
