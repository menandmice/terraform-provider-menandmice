package menandmice

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDNSZone() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSZoneCreate,
		ReadContext:   resourceDNSZoneRead,
		UpdateContext: resourceDNSZoneUpdate,
		DeleteContext: resourceDNSZoneDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDNSZoneImport,
		},
		Schema: map[string]*schema.Schema{

			"ref": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Internal references to this DNS zone.",
				Computed:    true,
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Description:  "Fully qualified name of DNS zone, ending with the trailing dot '.'.",
				Required:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`\.$`), "Name must end with '.'"),
			},
			"dynamic": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "If the DNS zone is dynamic. (Default: False)",
				Optional:    true,
				Default:     false,
			},
			// TODO following nameing convetion it would be ad_intergrated
			"adintegrated": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "If the DNS zone is AD integrated. (Default: False)",
				Optional:    true,
				Default:     false,
				ForceNew:    true,
			},
			"view": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Name of the view this DNS zone is in.",
				Optional:    true,
				Default:     "",
			},
			"dnsviewref": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Interal references to views.",
				Computed:    true,
			},
			"dnsviewrefs": &schema.Schema{
				Type:        schema.TypeSet,
				Description: "Interal references to views. Only used with Active Directory.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The type of the DNS zone. Example: Master, Slave, Hint, Stub, Forward. (Default: Master)",
				Optional:    true,
				Default:     "Master",
				ValidateFunc: validation.StringInSlice([]string{
					"Master", "Slave", "Hint", "Stub", "Forward",
				}, false),
			},
			"masters": &schema.Schema{
				Type:        schema.TypeList,
				Description: "List of IP addresses of all master zones, for slave zones.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.Any(
						validation.IsIPv4Address,
						validation.IsIPv6Address),
				},
				ForceNew: true,
				Optional: true,
			},

			"authority": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The authoritative DNS server for this zone. Requires FQDN with the trailing dot '.'.",
				ForceNew:    true,
				Required:    true,
				// TODO can also be a AD authority
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`\.$`), "authority should end with '.'"),
			},

			// TODO  following naming convention whould be dnssec_signed
			"dnssecsigned": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"kskids": &schema.Schema{
				Type:        schema.TypeString,
				Description: "A comma-separated string of IDs of KSKs. Starting with active keys, then inactive keys in parenthesis.",
				Optional:    true,
			},

			"zskids": &schema.Schema{
				Type:        schema.TypeString,
				Description: "A comma-separated string of IDs of ZSKs. Starting with active keys, then inactive keys in parenthesis.",
				Optional:    true,
			},

			"custom_properties": &schema.Schema{
				Type:        schema.TypeMap,
				Description: "Map of custom properties associated with this DNS zone.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"adreplicationtype": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Replication type if the zone is AD integrated.",
				ValidateFunc: validation.StringInSlice([]string{
					"None", "To_All_DNS_Servers_In_AD_Forrest",
					"To_All_DNS_Servers_In_AD_Domain", "To_All_Domain_Controllers_In_AD_Domain",
					"To_All_Domain_Controllers_In_Specified_Partition", "Unavailable",
				}, false),
			},
			"adpartition": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The AD partition if the zone is AD integrated.",
				Optional:    true,
			},
			"created": &schema.Schema{
				Type:        schema.TypeString,
				Description: "DDate when zone was created in Micetro.",
				Computed:    true,
			},
			"lastmodified": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Date when zone was last modified in Micetro.",
				Computed:    true,
			},
			"displayname": &schema.Schema{
				Type:        schema.TypeString,
				Description: "A display name to distinguish the zone from other, identically named zone instances.",
				Optional:    true,
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

	if dnsViewRefsRead, ok := d.GetOk("dnsviewrefs"); ok {
		dnsViewRefList := dnsViewRefsRead.(*schema.Set).List()
		var dnsViewRefs = make([]string, len(dnsViewRefList))
		for i, view := range dnsViewRefList {
			dnsViewRefs[i] = view.(string)
		}
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

		// you should not set this yourself
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

		// TODO this does not work
		dnszone.DNSViewRefs = []string{dnsviewref}
	} else {
		dnszone.DNSViewRef = dnsviewref

	}

	return dnszone
}

func resourceDNSZoneCreate(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Mmclient)

	var masters []string
	if mastersRead, ok := d.Get("masters").([]interface{}); ok {
		masters = make([]string, len(mastersRead))
		for i, master := range mastersRead {
			masters[i] = master.(string)
		}
	} else {
		return diag.Errorf("Could not read masters")
	}

	dnszone := readDNSZoneSchema(d)

	objRef, err := client.CreateDNSZone(dnszone, masters)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(objRef)

	return resourceDNSZoneRead(c, d, m)

}

func resourceDNSZoneRead(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := m.(*Mmclient)

	dnszone, err := client.ReadDNSZone(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if dnszone == nil {
		d.SetId("")
		return diags
	}

	writeDNSZoneSchema(d, *dnszone)

	return diags
}

func resourceDNSZoneUpdate(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	//can't change read only property
	if d.HasChange("ref") || d.HasChange("adintegrated") ||
		d.HasChange("dnsviewref") || d.HasChange("dnsviewrefs") ||
		d.HasChange("authority") {
		// this can't never error can never happen because of "ForceNew: true," for these properties
		return diag.Errorf("can't change read-only property of DNS zone")
	}
	client := m.(*Mmclient)
	ref := d.Id()
	dnszone := readDNSZoneSchema(d)

	err := client.UpdateDNSZone(dnszone.DNSZoneProperties, ref)

	if err != nil {
		return diag.FromErr(err)
	}
	return resourceDNSZoneRead(c, d, m)
}

func resourceDNSZoneDelete(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*Mmclient)
	var diags diag.Diagnostics
	ref := d.Id()
	err := client.DeleteDNSZone(ref)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}

func resourceDNSZoneImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	diags := resourceDNSZoneRead(ctx, d, m)
	if err := toError(diags); err != nil {
		return nil, err
	}

	// if we had used schema.ImportStatePassthrough
	// we could not have set id to its canonical form
	d.SetId(d.Get("ref").(string))

	return []*schema.ResourceData{d}, nil
}
