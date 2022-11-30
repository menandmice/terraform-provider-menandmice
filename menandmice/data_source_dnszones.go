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
			// TODO disabled for no. not realy needed. And need to be tested first
			// "filter": {
			// 	Type:        schema.TypeString,
			// 	Description: "Raw quickfilter String. Can be used to create more complex filter with >= etz",
			// 	Optional:    true,
			// },
			"folder": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Folder from which to get zones.",
			},
			"server": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Fully qualified name of the DNS server where the record is stored, ending with the trailing dot '.'.",
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`\.$`), "server should end with '.'"),
			},
			"view": {
				Type:        schema.TypeString,
				Description: "Name of the view this DNS zone is in.",
				Optional:    true,
				// Default:     "",
			},
			"dynamic": {
				Type:        schema.TypeBool,
				Description: "If the DNS zone is dynamic.",
				Optional:    true,
			},
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
			"dnssec_signed": {
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

						"ad_integrated": {
							Type:        schema.TypeBool,
							Description: "If the DNS zone is AD integrated.",
							Computed:    true,
						},
						"dns_view_ref": {
							Type:        schema.TypeString,
							Description: "Interal references to views.",
							Computed:    true,
						},
						"dns_view_refs": {
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
						"dnssec_signed": {
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
						"display_name": {
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
		flat["ad_integrated"] = zone.AdIntegrated
		flat["dns_view_ref"] = zone.DNSViewRef
		flat["dns_view_refs"] = zone.DNSViewRefs
		flat["authority"] = zone.Authority
		flat["created"] = zone.Created
		flat["lastmodified"] = zone.LastModified

		flat["dynamic"] = zone.Dynamic
		flat["type"] = zone.ZoneType
		flat["dnssec_signed"] = zone.DnssecSigned
		flat["kskids"] = zone.KskIDs
		flat["zskids"] = zone.ZskIDs
		flat["customp_properties"] = zone.CustomProperties
		// are not api respones
		// flat["adReplicationType"] = zone.AdReplicationType
		// flat["adPartition"] = zone.AdPartition
		flat["display_name"] = zone.DisplayName
		flattend[i] = flat
	}
	return flattend
}
func dataSourceDNSZonesRead(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	client := m.(*Mmclient)
	limit := d.Get("limit").(int)

	filter := map[string]interface{}{}
	if folder, ok := d.GetOk("folder"); ok {
		filter["folderRef"] = folder
	}

	if server, ok := d.GetOk("server"); ok {
		filter["dnsServerRef"] = server
	}

	if customProperties, ok := d.GetOk("custom_properties"); ok {
		for key, val := range customProperties.(map[string]interface{}) {
			filter[key] = val
		}
	}

	if view, ok := d.GetOk("view"); ok {
		filter["dnsViewRef"] = view
	}

	if zoneType, ok := d.GetOk("type"); ok {
		filter["type"] = zoneType
	}

	if rawFilter, ok := d.GetOk("filter"); ok {
		filter["filter"] = rawFilter
	}

	if dnssecSigned, ok := d.GetOk("dnssec_signed"); ok {
		filter["dnssecSigned"] = dnssecSigned
	}

	if dynamic, ok := d.GetOk("dynamic"); ok {
		filter["dynamic"] = dynamic
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
}
