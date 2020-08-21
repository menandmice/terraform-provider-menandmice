package menandmice

import (
	"context"
	"strconv"

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
			},
			"savecomment": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"comment": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"aging": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"ttl": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
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

type postCreate struct {
	DNSRecords []DNSRecord `json:"dnsRecords"`
	// saveComment string //TODO
	autoAssignRangeRef string // TODO
	// dnsZoneRef string //TODO
	// forceOverrideOfNamingConflictCheck bool // TODO

}

type postCreateResponse struct {
	Result struct {
		ObjRef []string `json:"objRefs"`
		Error  []string `json:"errors"`
	} `json:"result"`
}

func resourceDNSrecCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

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

	postcreate := postCreate{
		DNSRecords: []DNSRecord{DNSRecord{
			Ref:     optionalRef, //TODO
			Name:    d.Get("name").(string),
			Rectype: d.Get("type").(string),
			Ttl:     optionalTTL,
			Data:    d.Get("data").(string),
			Comment: d.Get("comment").(string),
			Aging:   d.Get("aging").(int),
			// Enabled:    d.Get("enabled").(bool),
			DNSZoneRef: d.Get("dnszone").(string),
		}},
		// autoAssignRangeRef: ...// TODO
	}

	c := m.(*resty.Client)

	var re postCreateResponse
	diags = MmPost(c, postcreate, &re, "DNSRecords")

	if len(re.Result.ObjRef) == 1 {
		d.SetId(re.Result.ObjRef[0])
	} else {
		diags = diag.Errorf("faild to create dnsrecord")
	}
	var diags2 = resourceDNSrecRead(ctx, d, m)

	diags = append(diags, diags2...)

	return diags
}

func resourceDNSrecRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type

	var diags diag.Diagnostics

	c := m.(*resty.Client)

	var re Response
	diags = MmGet(c, &re, "dnsrecords/"+d.Id())

	d.Set("ref", re.Result.DNSRecord.Ref)
	d.Set("name", re.Result.DNSRecord.Name)
	d.Set("type", re.Result.DNSRecord.Rectype)

	if re.Result.DNSRecord.Ttl != nil {
		ttl, err := strconv.Atoi(*re.Result.DNSRecord.Ttl)
		if err == nil {

			d.Set("ttl", ttl)
		}
	}
	d.Set("enabled", re.Result.DNSRecord.Enabled)
	d.Set("dnszoneref", re.Result.DNSRecord.DNSZoneRef)

	d.Set("aging", re.Result.DNSRecord.Aging)     //TODO default is no 0
	d.Set("comment", re.Result.DNSRecord.Comment) // comment is always given, but sometimes ""

	return diags
}

func resourceDNSrecUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return resourceDNSrecRead(ctx, d, m)
}

func resourceDNSrecDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	// orderID := d.Id()

	d.SetId("")
	return diags
}
