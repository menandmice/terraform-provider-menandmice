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
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if old == new+"." {
						return true
					}
					return false
				},
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
				ForceNew: true,
			},

			// TODO maybe choose between ref or refs automatic
			"dnsviewref": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				//TODO Default:  "DNSView/1", ?
			},
			"dnsviewrefs": &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ForceNew: true,
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
			"masters": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.Any(
						validation.IsIPv4Address,
						validation.IsIPv4Address),
				},
				ForceNew: true,
				Optional: true,
			},

			"authority": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
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
			"adreplicationtype": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"None", "To_All_DNS_Servers_In_AD_Forrest",
					"To_All_DNS_Servers_In_AD_Domain", "To_All_Domain_Controllers_In_AD_Domain",
					"To_All_Domain_Controllers_In_Specified_Partition", "Unavailable",
				}, false),
			},
			"adpartition": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
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

	d.Set("dnsviewref", dnszone.DnsViewRef)
	d.Set("dnsviewrefs", dnszone.DnsViewRefs)
	d.Set("authority", dnszone.Authority)
	d.Set("type", dnszone.ZoneType)
	d.Set("dnssecsigned", dnszone.DnssecSigned)
	d.Set("kskids", dnszone.KskIDs)
	d.Set("zskids", dnszone.ZskIDs)
	// TODO set customProperties

	d.Set("adreplicationtype", dnszone.AdReplicationType)
	d.Set("adpartition", dnszone.AdPartition)
	d.Set("created", dnszone.Created)
	d.Set("lastmodified", dnszone.LastModified)
	d.Set("displayname", dnszone.DisplayName)
	return

}

func readDNSzoneSchema(d *schema.ResourceData) DNSzone {
	// TODO  check dnsViewRef and dnsViewRefs are not both set

	dnsViewRefsRead := d.Get("dnsviewrefs").(*schema.Set).List() //TODO check succes
	var dnsViewRefs = make([]string, len(dnsViewRefsRead))
	for i, view := range dnsViewRefsRead {
		dnsViewRefs[i] = view.(string)
	}
	dnszone := DNSzone{
		Ref:          tryGetString(d, "ref"),
		AdIntegrated: d.Get("adintegrated").(bool),
		DnsViewRef:   tryGetString(d, "dnsviewref"),
		DnsViewRefs:  dnsViewRefs,
		Authority:    tryGetString(d, "authority"),

		DNSZoneProperties: DNSZoneProperties{
			Name:              d.Get("name").(string),
			Dynamic:           d.Get("dynamic").(bool),
			ZoneType:          tryGetString(d, "type"),
			DnssecSigned:      d.Get("dnssecsigned").(bool),
			KskIDs:            tryGetString(d, "kskids"),
			ZskIDs:            tryGetString(d, "zskids"),
			AdReplicationType: tryGetString(d, "adreplicationtype"),
			AdPartition:       tryGetString(d, "adpartition"),
			Created:           d.Get("created").(string),       // TODO convert to timeformat RFC 3339
			LastModified:      tryGetString(d, "lastmodified"), // TODO convert to timeformat RFC 3339
			DisplayName:       tryGetString(d, "displayname"),
		},
	}
	return dnszone
}

func resourceDNSzoneCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Mmclient)

	var masters []string
	if mastersRead, ok := d.Get("masters").([]interface{}); ok {
		masters = make([]string, len(mastersRead))
		for i, master := range mastersRead {
			masters[i] = master.(string)
		}
	}

	dnszone := readDNSzoneSchema(d)

	err, objRef := c.CreateDNSzone(dnszone, masters)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(objRef)

	return resourceDNSzoneRead(ctx, d, m)

}

func resourceDNSzoneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	c := m.(*Mmclient)

	err, dnszone := c.ReadDNSzone(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	writeDNSzoneSchema(d, dnszone)

	return diags
}

func resourceDNSzoneUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	//can't change read only property
	if d.HasChange("ref") || d.HasChange("adintegrated") ||
		d.HasChange("dnsviewref") || d.HasChange("dnsviewrefs") ||
		d.HasChange("authority") {
		// this can't never error can never happen because of "ForceNew: true," for these properties
		return diag.Errorf("can't change readonly property, of DNSZone")
	}
	c := m.(*Mmclient)
	ref := d.Id()
	dnszone := readDNSzoneSchema(d)

	err := c.UpdateDNSZone(dnszone.DNSZoneProperties, ref)

	if err != nil {
		return diag.FromErr(err)
	}
	return resourceDNSzoneRead(ctx, d, m)
}

func resourceDNSzoneDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*Mmclient)
	var diags diag.Diagnostics
	ref := d.Id()
	err := c.DeleteDNSZone(ref)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
