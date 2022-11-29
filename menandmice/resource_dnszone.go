package menandmice

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-log/tflog"
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

			"ref": {
				Type:        schema.TypeString,
				Description: "Internal references to this DNS zone.",
				Computed:    true,
			},
			"name": {
				Type:         schema.TypeString,
				Description:  "Fully qualified name of DNS zone, ending with the trailing dot '.'.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`\.$`), "Name must end with '.'"),
			},
			"dynamic": {
				Type:        schema.TypeBool,
				Description: "If the DNS zone is dynamic. (Default: False)",
				Optional:    true,
				Default:     false,
			},
			"ad_integrated": {
				Type:          schema.TypeBool,
				Description:   "If the DNS zone is AD integrated. (Default: False)",
				ConflictsWith: []string{"adintegrated"},
				Optional:      true,
				Default:       false,
				ForceNew:      true,
			},
			"adintegrated": {
				Type:          schema.TypeBool,
				Description:   "If the DNS zone is AD integrated. (Default: False)",
				ConflictsWith: []string{"ad_integrated"},
				Deprecated:    "use ad_integrated instead",
				Optional:      true,
				Default:       false,
				ForceNew:      true,
			},
			"view": {
				Type:        schema.TypeString,
				Description: "Name of the view this DNS zone is in.",
				Optional:    true,
				Default:     "",
			},
			"dns_view_ref": {
				Type:        schema.TypeString,
				Description: "Interal references to views.",
				Computed:    true,
			},
			"dnsviewref": {
				Type:        schema.TypeString,
				Deprecated:  "use dns_view_ref instead",
				Description: "Interal references to views.",
				Computed:    true,
			},

			"dns_view_refs": {
				Type:        schema.TypeSet,
				Description: "Interal references to views. Only used with Active Directory.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"dnsviewrefs": {
				Type:        schema.TypeSet,
				Description: "Interal references to views. Only used with Active Directory.",
				Deprecated:  "use dns_view_refs instead",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"type": {
				Type:        schema.TypeString,
				Description: "The type of the DNS zone. Example: Master, Slave, Hint, Stub, Forward. (Default: Master)",
				Optional:    true,
				Default:     "Master",
				ValidateFunc: validation.StringInSlice([]string{
					"Master", "Slave", "Hint", "Stub", "Forward",
				}, false),
			},
			"masters": {
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

			"authority": {
				Type:        schema.TypeString,
				Description: "The authoritative DNS server for this zone. Requires FQDN with the trailing dot '.'.",
				ForceNew:    true,
				Required:    true,
				// TODO can also be a AD authority
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`\.$`), "authority should end with '.'"),
			},

			// TODO  following naming convention whould be dnssec_signed
			"dnssecsigned": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"kskids": {
				Type:        schema.TypeString,
				Description: "A comma-separated string of IDs of KSKs. Starting with active keys, then inactive keys in parenthesis.",
				Optional:    true,
			},

			"zskids": {
				Type:        schema.TypeString,
				Description: "A comma-separated string of IDs of ZSKs. Starting with active keys, then inactive keys in parenthesis.",
				Optional:    true,
			},
			// TODO make custom_properties case insensitive
			"custom_properties": {
				Type:        schema.TypeMap,
				Description: "Map of custom properties associated with this DNS zone.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},

			"ad_replication_type": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Replication type if the zone is AD integrated.",
				ConflictsWith: []string{"adreplicationtype"},
				ValidateFunc: validation.StringInSlice([]string{
					"None", "To_All_DNS_Servers_In_AD_Forrest",
					"To_All_DNS_Servers_In_AD_Domain", "To_All_Domain_Controllers_In_AD_Domain",
					"To_All_Domain_Controllers_In_Specified_Partition", "Unavailable",
				}, false),
			},
			"adreplicationtype": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Replication type if the zone is AD integrated.",
				Deprecated:    "use ad_replication_type instead",
				ConflictsWith: []string{"ad_replication_type"},
				ValidateFunc: validation.StringInSlice([]string{
					"None", "To_All_DNS_Servers_In_AD_Forrest",
					"To_All_DNS_Servers_In_AD_Domain", "To_All_Domain_Controllers_In_AD_Domain",
					"To_All_Domain_Controllers_In_Specified_Partition", "Unavailable",
				}, false),
			},
			"ad_partition": {
				Type:          schema.TypeString,
				Description:   "The AD partition if the zone is AD integrated.",
				ConflictsWith: []string{"adpartition"},
				Optional:      true,
			},
			"adpartition": {
				Type:          schema.TypeString,
				Description:   "The AD partition if the zone is AD integrated.",
				Deprecated:    "use ad_partition instead",
				ConflictsWith: []string{"ad_partition"},
				Optional:      true,
			},
			"created": {
				Type:        schema.TypeString,
				Description: "DDate when zone was created in Micetro.",
				Computed:    true,
			},
			"lastmodified": {
				Type:        schema.TypeString,
				Description: "Date when zone was last modified in Micetro.",
				Computed:    true,
			},
			// TODO display might need to be computed not optional
			"display_name": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"displayname"},
				Description:   "A display name to distinguish the zone from other, identically named zone instances.",
				Optional:      true,
			},

			// TODO display might need to be computed not optional
			"displayname": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"display_name"},
				Deprecated:    "use displayname instead",
				Description:   "A display name to distinguish the zone from other, identically named zone instances.",
				Optional:      true,
			},
		},
	}
}

func writeDNSZoneSchema(d *schema.ResourceData, dnszone DNSZone) {

	d.Set("ref", dnszone.Ref)
	d.Set("name", dnszone.Name)
	d.Set("dynamic", dnszone.Dynamic)
	d.Set("adintegrated", dnszone.AdIntegrated)
	d.Set("ad_integrated", dnszone.AdIntegrated)

	d.Set("dnsviewref", dnszone.DNSViewRef)
	d.Set("dns_view_ref", dnszone.DNSViewRef)
	d.Set("dnsviewrefs", dnszone.DNSViewRefs)
	d.Set("dns_view_refs", dnszone.DNSViewRefs)
	d.Set("authority", dnszone.Authority)
	d.Set("type", dnszone.ZoneType)
	d.Set("dnssecsigned", dnszone.DnssecSigned)
	d.Set("kskids", dnszone.KskIDs)
	d.Set("zskids", dnszone.ZskIDs)
	d.Set("custom_properties", dnszone.CustomProperties)

	d.Set("adreplicationtype", dnszone.AdReplicationType)
	d.Set("ad_replication_type", dnszone.AdReplicationType)
	d.Set("adpartition", dnszone.AdPartition)
	d.Set("ad_partition", dnszone.AdPartition)
	d.Set("created", dnszone.Created)           // TODO convert to timeformat RFC 3339
	d.Set("lastmodified", dnszone.LastModified) // TODO convert to timeformat RFC 3339
	d.Set("displayname", dnszone.DisplayName)
}

func readDNSZoneSchema(d *schema.ResourceData) DNSZone {

	var dnsViewRefsRead interface{}
	var ok bool
	if dnsViewRefsRead, ok = d.GetOk("dnsviewrefs"); !ok {
		dnsViewRefsRead, ok = d.GetOk("dns_view_refs")
	}

	if ok {
		dnsViewRefList := dnsViewRefsRead.(*schema.Set).List()
		var dnsViewRefs = make([]string, len(dnsViewRefList))
		for i, view := range dnsViewRefList {
			dnsViewRefs[i] = view.(string)
		}
	}

	var customProperties = make(map[string]string)
	if customPropertiesRead, ok := d.GetOk("custom_properties"); ok {
		for key, value := range customPropertiesRead.(map[string]interface{}) {
			customProperties[key] = value.(string)
		}
	}

	adPartition, ok := d.GetOk("adpartition")
	if !ok {
		adPartition = d.Get("ad_partition")
	}

	adReplicationType, ok := d.GetOk("adreplicationtype")
	if !ok {
		adReplicationType = d.Get("ad_replication_type")
	}

	displayName, ok := d.GetOk("displayname")
	if !ok {
		displayName = d.Get("display_name")
	}

	dnszone := DNSZone{
		Ref:          tryGetString(d, "ref"),
		Name:         d.Get("name").(string),
		AdIntegrated: d.Get("adintegrated").(bool),
		Authority:    tryGetString(d, "authority"),

		// you should not set this yourself
		// Created:      d.Get("created").(string),
		// LastModified: tryGetString(d, "lastmodified"),

		DNSZoneProperties: DNSZoneProperties{
			Dynamic:           d.Get("dynamic").(bool),
			ZoneType:          tryGetString(d, "type"),
			DnssecSigned:      d.Get("dnssecsigned").(bool),
			KskIDs:            tryGetString(d, "kskids"),
			ZskIDs:            tryGetString(d, "zskids"),
			AdReplicationType: adReplicationType.(string),
			AdPartition:       adPartition.(string),
			CustomProperties:  customProperties,
			DisplayName:       displayName.(string),
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

func resourceDNSZoneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

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
	if d.HasChange("ref") ||
		d.HasChange("adintegrated") || d.HasChange("ad_integrated") ||
		d.HasChange("dnsviewref") || d.HasChange("dns_view_ref") ||
		d.HasChange("dnsviewrefs") || d.HasChange("dns_view_refs") ||
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
	ref := d.Get("ref").(string)
	if ref == "" {
		tflog.Debug(ctx, fmt.Sprintf("%v", d))
		return nil, errors.New("Import failed")
	}
	d.SetId(ref)

	return []*schema.ResourceData{d}, nil
}
