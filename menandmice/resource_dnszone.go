package menandmice

import (
	"terraform-provider-menandmice/diag"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceDNSZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceDNSZoneCreate,
		Read:   resourceDNSZoneRead,
		Update: resourceDNSZoneUpdate,
		Delete: resourceDNSZoneDelete,
		Schema: map[string]*schema.Schema{

			"ref": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				// TODO beter to force name ending with .
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
			"view": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"dnsviewref": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"dnsviewrefs": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
				// TODO Requires . at end
				ValidateFunc: validation.StringIsNotEmpty,
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

			"custom_properties": &schema.Schema{
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
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

func writeDNSZoneSchema(d *schema.ResourceData, dnszone DNSZone) {

	d.Set("ref", dnszone.Ref)
	d.Set("name", dnszone.Name)
	d.Set("dynamic", dnszone.Dynamic)
	d.Set("adintegrated", dnszone.AdIntegrated)

	d.Set("dnsviewref", dnszone.DNSViewRef)
	d.Set("dnsviewrefs", dnszone.DNSViewRefs)
	d.Set("authority", dnszone.Authority)
	d.Set("type", dnszone.ZoneType)
	d.Set("dnssecsigned", dnszone.DnssecSigned)
	d.Set("kskids", dnszone.KskIDs)
	d.Set("zskids", dnszone.ZskIDs)
	d.Set("custom_properties", dnszone.CustomProperties)

	d.Set("adreplicationtype", dnszone.AdReplicationType)
	d.Set("adpartition", dnszone.AdPartition)
	d.Set("created", dnszone.Created)           // TODO convert to timeformat RFC 3339
	d.Set("lastmodified", dnszone.LastModified) // TODO convert to timeformat RFC 3339
	d.Set("displayname", dnszone.DisplayName)
	return

}

func readDNSZoneSchema(d *schema.ResourceData) DNSZone {

	dnsViewRefsRead := d.Get("dnsviewrefs").(*schema.Set).List() //TODO check succes
	var dnsViewRefs = make([]string, len(dnsViewRefsRead))
	for i, view := range dnsViewRefsRead {
		dnsViewRefs[i] = view.(string)
	}

	var CustomProperties = make(map[string]string)
	if customPropertiesRead, ok := d.GetOk("custom_properties"); ok {
		for key, value := range customPropertiesRead.(map[string]interface{}) {
			CustomProperties[key] = value.(string)
		}
	}
	dnszone := DNSZone{
		Ref:          tryGetString(d, "ref"),
		AdIntegrated: d.Get("adintegrated").(bool),
		Authority:    tryGetString(d, "authority"),

		// you should not set this youself
		// Created:      d.Get("created").(string),
		// LastModified: tryGetString(d, "lastmodified"),

		DNSZoneProperties: DNSZoneProperties{
			Name:              d.Get("name").(string),
			Dynamic:           d.Get("dynamic").(bool),
			ZoneType:          tryGetString(d, "type"),
			DnssecSigned:      d.Get("dnssecsigned").(bool),
			KskIDs:            tryGetString(d, "kskids"),
			ZskIDs:            tryGetString(d, "zskids"),
			AdReplicationType: tryGetString(d, "adreplicationtype"),
			AdPartition:       tryGetString(d, "adpartition"),
			CustomProperties:  CustomProperties,
			DisplayName:       tryGetString(d, "displayname"),
		},
	}

	dnsviewref := dnszone.Authority + ":" + tryGetString(d, "view")
	if dnszone.AdIntegrated {

		dnszone.DNSViewRefs = []string{dnsviewref}
	} else {
		dnszone.DNSViewRef = dnsviewref

	}

	return dnszone
}

func resourceDNSZoneCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Mmclient)

	var masters []string
	if mastersRead, ok := d.Get("masters").([]interface{}); ok {
		masters = make([]string, len(mastersRead))
		for i, master := range mastersRead {
			masters[i] = master.(string)
		}
	}

	dnszone := readDNSZoneSchema(d)

	objRef, err := c.CreateDNSZone(dnszone, masters)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(objRef)

	return resourceDNSZoneRead(d, m)

}

func resourceDNSZoneRead(d *schema.ResourceData, m interface{}) error {

	var diags diag.Diagnostics

	c := m.(*Mmclient)

	dnszone, err := c.ReadDNSZone(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	writeDNSZoneSchema(d, dnszone)

	return diags
}

func resourceDNSZoneUpdate(d *schema.ResourceData, m interface{}) error {

	//can't change read only property
	if d.HasChange("ref") || d.HasChange("adintegrated") ||
		d.HasChange("dnsviewref") || d.HasChange("dnsviewrefs") ||
		d.HasChange("authority") {
		// this can't never error can never happen because of "ForceNew: true," for these properties
		return diag.Errorf("can't change readonly property, of DNSZone")
	}
	c := m.(*Mmclient)
	ref := d.Id()
	dnszone := readDNSZoneSchema(d)

	err := c.UpdateDNSZone(dnszone.DNSZoneProperties, ref)

	if err != nil {
		return diag.FromErr(err)
	}
	return resourceDNSZoneRead(d, m)
}

func resourceDNSZoneDelete(d *schema.ResourceData, m interface{}) error {

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
