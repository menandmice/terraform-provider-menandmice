package menandmice

import (
	"context"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceDNSZones() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDNSZonesRead,
		Schema: map[string]*schema.Schema{
			"limit": {
				Type:        schema.TypeInt,
				Description: "The number of zones to return.",
				Optional:    true,
			},
			"filter": {
				Type:        schema.TypeString,
				Description: "Raw filter String. Can be used to create more complex filter with >= etz",
				Optional:    true,
			},

			"folder": {
				Type: schema.TypeString,
				// TODO
				Optional:    true,
				Description: "Reference to a folder from which to get zones.",
			},
			"server": {
				Type: schema.TypeString,
				// TODO
				Optional:     true,
				Description:  "Fully qualified name of the DNS server where the record is stored, ending with the trailing dot '.'.",
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`\.$`), "server should end with '.'"),
			},
			"view": {
				Type:        schema.TypeString,
				Description: "Name of the view this DNS zone is in.",
				// TODO
				Optional: true,
				// Default:     "",
			},
			// "name": {
			// 	Type:         schema.TypeString,
			// 	Description:  "Fully qualified name of DNS zone, ending with the trailing dot '.'.",
			// 	ValidateFunc: validation.StringMatch(regexp.MustCompile(`\.$`), "name must end with '.'"),
			// 	Required:     true,
			// },
			"dynamic": {
				Type:        schema.TypeBool,
				Description: "If the DNS zone is dynamic.",
				Optional:    true,
			},
			// TODO following nameing convetion it would be ad_intergrated
			// "adintegrated": {
			// 	Type:        schema.TypeBool,
			// 	Description: "If the DNS zone is AD integrated.",
			// 	Computed:    true,
			// },
			//
			// // // TODO unify dnsviewref dnsviewrefs
			// // TODO following nameing convetion it would be dns_view_ref
			// "dnsviewref": {
			// 	Type:        schema.TypeString,
			// 	Description: "Interal references to views.",
			// 	Computed:    true,
			// },
			//
			// // TODO following nameing convetion it would be dns_view_refs
			// "dnsviewrefs": {
			// 	Type:        schema.TypeSet,
			// 	Description: "Interal references to views. Only used with Active Directory.",
			// 	Elem: &schema.Schema{
			// 		Type: schema.TypeString,
			// 	},
			// 	Computed: true,
			// },
			"type": {
				Type:        schema.TypeString,
				Description: "The type of the DNS zone. Example: Master, Slave, Hint, Stub, Forward.",
				Optional:    true,
			},
			"authority": {
				Type:        schema.TypeString,
				Description: "The authoritative DNS server for this zone.",
				Optional:    true,
			},
			"dnssecsigned": {
				Type:        schema.TypeBool,
				Description: "If DNS signing is enabled.",
				Optional:    true,
			},
			"custom_properties": {
				Type:        schema.TypeMap,
				Description: "Map of custom properties associated with this DNS zone.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"zones": {
				Type:        schema.TypeList,
				Description: "Zones found with described properties",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"ref": {
							Type:        schema.TypeString,
							Description: "Internal references to this DNS zone.",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Fully qualified name of DNS zone, ending with the trailing dot '.'.",
							Computed:    true,
						},

						"dynamic": {
							Type:        schema.TypeBool,
							Description: "If the DNS zone is dynamic.",
							Computed:    true,
						},

						// TODO following nameing convetion it would be ad_intergrated
						"adintegrated": {
							Type:        schema.TypeBool,
							Description: "If the DNS zone is AD integrated.",
							Computed:    true,
						},

						// TODO unify dnsviewref dnsviewrefs
						// TODO following nameing convetion it would be dns_view_ref
						"dnsviewref": {
							Type:        schema.TypeString,
							Description: "Interal references to views.",
							Computed:    true,
						},

						// TODO following nameing convetion it would be dns_view_refs
						"dnsviewrefs": {
							Type:        schema.TypeSet,
							Description: "Interal references to views. Only used with Active Directory.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed: true,
						},
						"type": {
							Type:        schema.TypeString,
							Description: "The type of the DNS zone. Example: Master, Slave, Hint, Stub, Forward.",
							Computed:    true,
						},
						"authority": {
							Type:        schema.TypeString,
							Description: "The authoritative DNS server for this zone.",
							Computed:    true,
						},
						"dnssecsigned": {
							Type:        schema.TypeBool,
							Description: "If DNS signing is enabled.",
							Computed:    true,
						},
						"kskids": {
							Type:        schema.TypeString,
							Description: "A comma-separated string of IDs of KSKs. Starting with active keys, then inactive keys in parenthesis.",
							Computed:    true,
						},
						"zskids": {
							Type:        schema.TypeString,
							Description: "A comma-separated string of IDs of ZSKs. Starting with active keys, then inactive keys in parenthesis.",
							Computed:    true,
						},

						"customp_properties": {
							Type:        schema.TypeMap,
							Description: "Map of custom properties associated with this DNS zone.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed: true,
						},

						"created": {
							Type:        schema.TypeString,
							Description: "Date when zone was created in Micetro.",
							Computed:    true,
						},
						"lastmodified": {
							Type:        schema.TypeString,
							Description: "Date when zone was last modified in Micetro.",
							Computed:    true,
						},
						"displayname": {
							Type:        schema.TypeString,
							Description: "A display name to distinguish the zone from other, identically named zone instances.",
							Computed:    true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}
func flattenZones(zones []DNSZone) []interface{} {
	if zones == nil {
		return make([]interface{}, 0)
	}
	flattend := make([]interface{}, len(zones))

	for i, zone := range zones {

		flat := make(map[string]interface{})
		flat["ref"] = zone.Ref
		flat["name"] = zone.Name
		// flat["ad_intergrated"] = zone.AdIntegrated //TODO
		// panic(zone.DNSViewRef)
		flat["dnsviewref"] = zone.DNSViewRef
		flat["dnsviewrefs"] = zone.DNSViewRefs
		flat["authority"] = zone.Authority
		flat["created"] = zone.Created
		flat["lastmodified"] = zone.LastModified

		flat["dynamic"] = zone.Dynamic
		flat["type"] = zone.ZoneType
		flat["dnssecsigned"] = zone.DnssecSigned
		flat["kskids"] = zone.KskIDs
		flat["zskids"] = zone.ZskIDs
		flat["customp_properties"] = zone.CustomProperties
		// flat["adReplicationType"] = zone.AdReplicationType
		// flat["adPartition"] = zone.AdPartition
		flat["displayname"] = zone.DisplayName
		flattend[i] = flat
	}
	return flattend
}
func dataSourceDNSZonesRead(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	client := m.(*Mmclient)
	limit := d.Get("limit").(int)

	// TODO implement raw filter
	filter := map[string]string{}
	// if filter,ok :+
	if folder, ok := d.GetOk("folder"); ok {
		filter["folderRef"] = folder.(string)
	}

	if server, ok := d.GetOk("server"); ok {
		filter["dnsServerRef"] = server.(string)
	}

	if customProperties, ok := d.GetOk("custom_properties"); ok {
		for key, val := range customProperties.(map[string]interface{}) {
			filter[key] = val.(string)
		}
	}

	if view, ok := d.GetOk("view"); ok {
		filter["dnsViewRef"] = view.(string)
	}

	if zoneType, ok := d.GetOk("type"); ok {
		filter["type"] = zoneType.(string)
	}

	dnszones, err := client.FindDNSZones(limit, filter)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("type", "test")
	if err := d.Set("zones", flattenZones(dnszones)); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(len(dnszones)))
	return diags
	//
	// server := tryGetString(d, "server")
	// view := tryGetString(d, "view")
	// name := tryGetString(d, "name")
	// dnsZoneRef := server + ":" + view + ":" + name
	// dnszone, err := client.ReadDNSZone(dnsZoneRef)
	//
	// if err != nil {
	// 	return diag.FromErr(err)
	// }
	//
	// if dnszone == nil {
	// 	if view == "" {
	// 		return diag.Errorf("The DNS zone %v does not exist on server %v", name, server)
	// 	}
	// 	return diag.Errorf("The DNS zone %v does not exist in view %v on %v", name, view, server)
	// }
	// writeDNSZoneSchema(d, *dnszone)
	// d.SetId(dnszone.Ref)
	//
	// return diags
}
