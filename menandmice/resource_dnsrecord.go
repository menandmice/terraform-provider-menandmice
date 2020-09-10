package menandmice

import (
	"strconv"
	"time"

	"terraform-provider-menandmice/diag"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceDNSrec() *schema.Resource {
	return &schema.Resource{
		Create: resourceDNSrecCreate,
		Read:   resourceDNSrecRead,
		Update: resourceDNSrecUpdate,
		Delete: resourceDNSrecDelete,
		Schema: map[string]*schema.Schema{

			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ref": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"data": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				// You cannot validate data here, because you dont have acces to what kind of record it is
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
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

			// TODO it is not dnszone but dnszoneref
			"dnszone": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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

	if dnsrec.Ttl != nil {
		ttl, err := strconv.Atoi(*dnsrec.Ttl)
		if err == nil {

			d.Set("ttl", ttl)
		}
	}
	d.Set("enabled", dnsrec.Enabled)
	d.Set("dnszone", dnsrec.DNSZoneRef)

	d.Set("aging", dnsrec.Aging)
	d.Set("comment", dnsrec.Comment) // comment is always given, but sometimes ""
	return

}

func readDNSRecSchema(d *schema.ResourceData) DNSRecord {

	var ref string = d.Get("ref").(string)
	var optionalRef *string
	if ref != "" {
		optionalRef = &ref
	}
	var ttl int = d.Get("ttl").(int)
	var optionalTTL *string
	if ttl != 0 {
		ttlString := strconv.Itoa(ttl)
		optionalTTL = &ttlString
	}

	dnsrec := DNSRecord{
		Ref:        optionalRef,
		DNSZoneRef: d.Get("dnszone").(string),
		DNSProperties: DNSProperties{
			Name:    d.Get("name").(string),
			Rectype: d.Get("type").(string),
			Ttl:     optionalTTL,
			Data:    d.Get("data").(string),
			Comment: d.Get("comment").(string),
			Aging:   d.Get("aging").(int), // TODO when not specified it's 0
			Enabled: d.Get("enabled").(bool),
		},
	}
	return dnsrec
}

func resourceDNSrecCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Mmclient)

	dnsrec := readDNSRecSchema(d)

	err, objRef := c.CreateDNSRec(dnsrec)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(objRef)

	return resourceDNSrecRead(d, m)

}

func resourceDNSrecRead(d *schema.ResourceData, m interface{}) error {
	// Warning or errors can be collected in a slice type

	var diags diag.Diagnostics

	c := m.(*Mmclient)

	err, dnsrec := c.ReadDNSRec(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	writeDNSRecSchema(d, dnsrec)

	return diags
}

func resourceDNSrecUpdate(d *schema.ResourceData, m interface{}) error {

	//can't change read only property zone
	if d.HasChange("dnszone") {
		// this can't never error can never happen because of "ForceNew: true," for dnszone
		// TODO this messages use dnszone but it is a dnszone ref
		return diag.Errorf("cant update dnszone of %s.%s. you could try to delete dnsrecord first", d.Get("name"), d.Get("dnszone"))
	}
	c := m.(*Mmclient)
	ref := d.Id()
	dnsrec := readDNSRecSchema(d)
	err := c.UpdateDNSRec(dnsrec.DNSProperties, ref)

	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("last_updated", time.Now().Format(time.RFC850))
	return resourceDNSrecRead(d, m)
}

func resourceDNSrecDelete(d *schema.ResourceData, m interface{}) error {
	// Warning or errors can be collected in a slice type

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
