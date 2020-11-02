package menandmice

import (
	"regexp"
	"strconv"

	"terraform-provider-menandmice/diag"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceDNSRec() *schema.Resource {
	return &schema.Resource{
		Create: resourceDNSRecCreate,
		Read:   resourceDNSRecRead,
		Update: resourceDNSRecUpdate,
		Delete: resourceDNSRecDelete,
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
				Default:      0, // if set to 0 its ignored
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
			// TODO add autoAssignRangeRef
		},
	}
}

func writeDNSRecSchema(d *schema.ResourceData, dnsrec DNSRecord) {

	d.Set("ref", dnsrec.Ref)
	d.Set("name", dnsrec.Name)
	d.Set("type", dnsrec.Rectype)
	d.Set("data", dnsrec.Data)
	if dnsrec.Ttl != nil {
		ttl, err := strconv.Atoi(*dnsrec.Ttl)
		if err == nil {

			d.Set("ttl", ttl)
		}
	}
	d.Set("enabled", dnsrec.Enabled)
	d.Set("dns_zone_ref", dnsrec.DNSZoneRef)

	d.Set("aging", dnsrec.Aging)
	d.Set("comment", dnsrec.Comment) // comment is always given, but sometimes ""
	return

}

func readDNSRecSchema(d *schema.ResourceData) DNSRecord {

	var optionalTTL *string
	if ttl, ok := d.Get("ttl").(int); ok {
		ttlString := strconv.Itoa(ttl)
		optionalTTL = &ttlString
	}

	dnsrec := DNSRecord{
		Ref:        tryGetString(d, "ref"),
		DNSZoneRef: tryGetString(d, "server") + ":" + tryGetString(d, "view") + ":" + tryGetString(d, "zone"),

		Rectype: d.Get("type").(string),
		DNSProperties: DNSProperties{
			Name:    d.Get("name").(string),
			Ttl:     optionalTTL,
			Data:    d.Get("data").(string),
			Comment: d.Get("comment").(string),
			Aging:   d.Get("aging").(int), // TODO when not specified it's 0
			Enabled: d.Get("enabled").(bool),
		},
	}
	return dnsrec
}

func resourceDNSRecCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Mmclient)

	dnsrec := readDNSRecSchema(d)

	objRef, err := c.CreateDNSRec(dnsrec)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(objRef)

	return resourceDNSRecRead(d, m)

}

func resourceDNSRecRead(d *schema.ResourceData, m interface{}) error {

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(*Mmclient)

	dnsrec, err := c.ReadDNSRec(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	writeDNSRecSchema(d, dnsrec)

	return diags
}

func resourceDNSRecUpdate(d *schema.ResourceData, m interface{}) error {

	c := m.(*Mmclient)
	ref := d.Id()
	dnsrec := readDNSRecSchema(d)
	err := c.UpdateDNSRec(dnsrec.DNSProperties, ref)

	if err != nil {
		return diag.FromErr(err)
	}
	return resourceDNSRecRead(d, m)
}

func resourceDNSRecDelete(d *schema.ResourceData, m interface{}) error {

	c := m.(*Mmclient)
	var diags diag.Diagnostics
	ref := d.Id()
	err := c.DeleteDNSRec(ref)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
