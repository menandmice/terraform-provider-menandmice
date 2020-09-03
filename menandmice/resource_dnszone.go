package menandmice

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDNSzone() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSzoneCreate,
		ReadContext:   resourceDNSzoneRead,
		UpdateContext: resourceDNSzoneUpdate,
		DeleteContext: resourceDNSzoneDelete,
		Schema: map[string]*schema.Schema{

			"ref": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"dynamic": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"adintegrated": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			// TODO maybe choose between ref or refs automatic
			"dnsviewref": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				// Default:  "DNSView/1", //TODO
			},
			"dnsviewrefs": &schema.Schema{
				Type: schema.TypeList, // TODO TypeSet
				Elem: &schema.Schema{
					Type: schema.TypeString,
					// TODO Default:
				},
				Optional: true,
			},

			"type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Master",
				ValidateFunc: validation.StringInSlice([]string{
					"Master", "Slave", "Hint", "Stub", "Forward",
				}, false),
			},

			"authority": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"dnssecsigned": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"kskids": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"zskids": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			// TODO  "customProperties": &schema.Schema{
			// 	Type:     ?
			// 	Computed: true,
			// }
			// TODO adReplicationType
			// TODO adPartition
			"created": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"lastmodified": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"displayname": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func writeDNSzoneSchema(d *schema.ResourceData, dnszone DNSzone) {

	d.Set("ref", dnszone.Ref)
	d.Set("name", dnszone.Name)
	d.Set("dynamic", dnszone.Dynamic)
	d.Set("adintegrated", dnszone.AdIntegrated)

	d.Set("dnsViewRef", dnszone.DnsViewRef)
	d.Set("dnsViewRefs", dnszone.DnsViewRefs)
	// if len(dnszone.DnsViewRefs) <= 0 {
	// 	d.Set("dnsViewRefs", []string{dnszone.DnsViewRef})
	// } else {
	// 	d.Set("dnsViewRefs", dnszone.DnsViewRefs)
	// }
	d.Set("dnsViewRefs", dnszone.DnsViewRefs)
	d.Set("authority", dnszone.Authority)
	d.Set("type", dnszone.ZoneType)
	d.Set("dnssecsigned", dnszone.DnssecSigned)
	d.Set("kskids", dnszone.KskIDs)
	d.Set("zskids", dnszone.ZskIDs)
	// TODO set customProperties
	d.Set("created", dnszone.Created)
	d.Set("lastmodified", dnszone.LastModified)
	d.Set("displayname", dnszone.DisplayName)
	return

}

func readDNSzoneSchema(d *schema.ResourceData) DNSzone {
	// TODO  check dnsViewRef and dnsViewRefs are not both set

	var ref string = d.Get("ref").(string)
	var optionalRef *string
	if ref != "" {
		optionalRef = &ref
	}

	dnsViewRefsRead := d.Get("dnsviewrefs").([]interface{})
	var dnsViewRefs = make([]string, len(dnsViewRefsRead))
	for i, view := range dnsViewRefsRead {
		dnsViewRefs[i] = view.(string)
	}
	dnszone := DNSzone{
		Ref:          optionalRef,
		Name:         d.Get("name").(string),
		Dynamic:      d.Get("dynamic").(bool),
		AdIntegrated: d.Get("adintegrated").(bool),

		DnsViewRef:   tryGetString(d, "dnsviewref"),
		DnsViewRefs:  dnsViewRefs,
		Authority:    tryGetString(d, "Authority"),
		ZoneType:     tryGetString(d, "type"),
		DnssecSigned: d.Get("dnssecsigned").(bool),
		KskIDs:       tryGetString(d, "kskids"),
		ZskIDs:       tryGetString(d, "zskids"),
		Created:      d.Get("created").(string),
		LastModified: tryGetString(d, "lastmodified"),
		DisplayName:  tryGetString(d, "displayname"),
	}
	return dnszone
}

func resourceDNSzoneCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Mmclient)

	dnszone := readDNSzoneSchema(d)

	err, objRef := c.CreateDNSzone(dnszone)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(objRef)

	return resourceDNSzoneRead(ctx, d, m)

}

func resourceDNSzoneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	// c := m.(*Mmclient)
	//
	// err, dnszone := ReadDNSzone(c, d.Id())
	// if err != nil {
	// 	return diag.FromErr(err)
	// }
	// writeDNSzoneSchema(d, dnszone)

	return diags
}

func resourceDNSzoneUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// c := m.(*Mmclient)
	// ref := d.Id()
	// dnszone := readDNSzoneSchema(d)
	// err := UpdateDNSzone(c, dnszone.DNSProperties, ref)
	//
	// if err != nil {
	// 	return diag.FromErr(err)
	// }
	return resourceDNSzoneRead(ctx, d, m)
}

func resourceDNSzoneDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*Mmclient)
	var diags diag.Diagnostics
	ref := d.Id()
	err := c.DeleteDNSzone(ref)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
