package menandmice

import (
	"context"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDNSrec() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSrecCreate,
		ReadContext:   resourceDNSrecRead,
		UpdateContext: resourceDNSrecUpdate,
		DeleteContext: resourceDNSrecDelete,
		//TODO for now 1 records at a time
		Schema: map[string]*schema.Schema{

			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// TODO add force oferwrite
			// TODO add autoAssignRangeRef
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"data": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				//TODO validate
			},
			"comment": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"aging": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				//TODO valiate
			},
			"ttl": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				//TODO validate
			},
			"ref": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			// TODO should this be dnszoneref
			"dnszone": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}
func schema2dnsrecord(d *schema.ResourceData) DNSRecord {

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
			Aging:   d.Get("aging").(int),
			// Enabled:    d.Get("enabled").(bool), //TODO
		},
	}
	return dnsrec
}

func resourceDNSrecCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*resty.Client)

	dnsrec := schema2dnsrecord(d)

	err, objRef := CreateDNSRec(c, dnsrec)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(objRef)

	return resourceDNSrecRead(ctx, d, m)

}

func resourceDNSrecRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type

	var diags diag.Diagnostics

	c := m.(*resty.Client)

	err, dnsrec := ReadDNSRec(c, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	// TODO  remove duplcation dataSourceDNSrectRead

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
	d.Set("dnszoneref", dnsrec.DNSZoneRef)

	d.Set("aging", dnsrec.Aging)     //TODO default
	d.Set("comment", dnsrec.Comment) // comment is always given, but sometimes ""

	return diags
}

func resourceDNSrecUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	//can't change read only property zone
	if d.HasChange("dnszone") {
		return diag.Errorf("cant update dnszone of %s.%s. you could try to delete dnsrecord first", d.Get("name"), d.Get("dnszone")) // TODO find work around
	}
	c := m.(*resty.Client)
	ref := d.Id()
	dnsrec := schema2dnsrecord(d)
	err := UpdateDNSRec(c, dnsrec.DNSProperties, ref)

	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("last_updated", time.Now().Format(time.RFC850))
	return resourceDNSrecRead(ctx, d, m)
}

func resourceDNSrecDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type

	c := m.(*resty.Client)
	var diags diag.Diagnostics
	ref := d.Id()
	err := DeleteDNSRec(c, ref)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
