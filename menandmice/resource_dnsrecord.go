package menandmice

import (
	"context"
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

			"ref": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"data": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				// You cannot validate data here, because you dont have acces to what kind of record it is
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Default:  "A",
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
			"comment": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"aging": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"ttl": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"server": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`\.$`), "server should end with '.'"),
			},
			"view": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				ForceNew: true,
			},

			"zone": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`\.$`), "server should end with '.'"),
			},

			"dns_zone_ref": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			// TODO add force oferwrite
		},
	}
}

func writeDNSRecSchema(d *schema.ResourceData, dnsrec DNSRecord) {

	d.Set("ref", dnsrec.Ref)
	d.Set("name", dnsrec.Name)
	d.Set("type", dnsrec.Rectype)
	d.Set("data", dnsrec.Data)
	if dnsrec.Ttl != "" {
		ttl, err := strconv.Atoi(dnsrec.Ttl)
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
	return

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
			Ttl:     ttlString,
			Data:    d.Get("data").(string),
			Comment: d.Get("comment").(string),
			Aging:   d.Get("aging").(int), // when not specified it's 0 witch will be ignored
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

	if parts := strings.Split(arg, ":"); len(parts) == 4 {

		// format is authority:view:fqdn
		d.Set("server", parts[0])
		d.Set("view", parts[1])
		fqdn := strings.SplitN(parts[2], ".", 2)

		if len(fqdn) != 2 {
			return nil, fmt.Errorf("could not parse fqdn %s", parts[2])
		}
		d.Set("name", fqdn[0])
		d.Set("zone", fqdn[1])
		d.Set("type", parts[3])

		diags := dataSourceDNSRectRead(ctx, d, m)
		if err := toError(diags); err != nil {
			return nil, err
		}
		d.SetId(d.Get("ref").(string))

	} else {
		// otherwise format is dnsrecords/<id> which would work for read to

		diags := resourceDNSRecRead(ctx, d, m)
		if err := toError(diags); err != nil {
			return nil, err
		}
		d.SetId(d.Get("ref").(string))
	}

	return []*schema.ResourceData{d}, nil
}
