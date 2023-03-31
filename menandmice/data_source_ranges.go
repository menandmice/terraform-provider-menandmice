package menandmice

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceRanges() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRangesRead,
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
				Description: "Folder from which to get ranges.",
			},

			"custom_properties": {
				Type:        schema.TypeMap,
				Description: "Search for zones with these custom_properties",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},

			"is_container": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Filter on if range is a container",
			},

			"subnet": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Filter on if range is a subnet",
			},

			// TODO filter on parent_range, adSiteRef, cloudNetworkRef
			"ranges": {
				Type:        schema.TypeList,
				Description: "Ranges found with described properties.",
				Computed:    true,
				Elem: &schema.Resource{

					// TODO add atributes: cloudAllocationPools, dhcpScopes authority ,discoveredProperties
					Schema: map[string]*schema.Schema{

						"ref": {
							Type:        schema.TypeString,
							Description: "Internal references to this range.",
							Computed:    true,
						},

						"name": {
							Type:        schema.TypeString,
							Description: "The CIDR of the range, or from-to address range.",
							Required:    true,
						},
						"cidr": {
							Type:        schema.TypeString,
							Description: "The CIDR of the range",
							Optional:    true,
						},
						"from": {
							Type:        schema.TypeString,
							Description: "The starting IP address of the range.",
							Computed:    true,
						},
						"to": {
							Type:        schema.TypeString,
							Description: "The ending IP address of the range.",
							Computed:    true,
						},
						"parent_ref": {
							Type:        schema.TypeString,
							Description: "A reference to the range that contains the subranges",
							Computed:    true,
						},

						"ad_site_ref": {
							Type:        schema.TypeString,
							Description: "Internal reference of the AD site to which the the range belongs.",
							Computed:    true,
						},

						"ad_site_display_name": {
							Type:        schema.TypeString,
							Description: "The display name of the AD site to which the range belongs.",
							Computed:    true,
						},
						"child_ranges": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "An list of child ranges of the range.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ref": {
										Type:        schema.TypeString,
										Description: "Internal references to this child range.",
										Computed:    true,
									},

									"name": {
										Type:        schema.TypeString,
										Description: "Name to this child range.",
										Computed:    true,
									},
								},
							},
						},
						// "dhcpScopes": {
						// 	Type:        schema.TypeList,
						// 	Description:
						// 	Computed:    true,
						// 	// Default:      false,
						// },
						// "authority": {
						// 	Type:        schema.TypeList,
						// 	Description:
						// 	Computed:    true,
						// },

						"subnet": {
							Type:        schema.TypeBool,
							Description: "Determines if the range is defined as a subnet.",
							Computed:    true,
						},

						"locked": {
							Type:        schema.TypeBool,
							Description: "Determines if the range is locked.",
							Computed:    true,
						},
						"auto_assign": {
							Type:        schema.TypeBool,
							Description: "Determines if it should be possible to automatically assign IP addresses from the range.",
							Computed:    true,
						},
						"has_schedule": {
							Type:        schema.TypeBool,
							Description: "Determines if a discovery schedule has been set for the range.",
							Computed:    true,
						},
						"has_monitor": {
							Type:        schema.TypeBool,
							Description: "Determines if a discovery schedule has been set for the range.",
							Computed:    true,
						},

						"title": {
							Type:        schema.TypeString,
							Description: "The title of the Range",
							Computed:    true,
						},

						"description": {
							Type:        schema.TypeString,
							Description: "Description of the range",
							Computed:    true,
						},

						"custom_properties": {
							Type:        schema.TypeMap,
							Description: "Map of custom properties associated with this range. You can only assign properties that are already defined in Micetro.",

							Elem: &schema.Schema{
								Type: schema.TypeString,
							},

							Computed: true,
						},
						"inherit_access": {
							Type:        schema.TypeBool,
							Description: "If this range should inherit its access bits from its parent range.",
							Computed:    true,
						},
						"is_container": {
							Type:        schema.TypeBool,
							Description: "Set to true to create a container instead of a range.",
							Computed:    true,
						},

						"utilization_percentage": {
							Type:        schema.TypeInt,
							Description: "Utilization percentage for range.",
							Computed:    true,
						},

						"has_rogue_addresses": {
							Type:        schema.TypeBool,
							Description: "Set to true to create a container instead of a range.",
							Computed:    true,
						},

						"cloud_network_ref": {
							Type:        schema.TypeString,
							Description: "A internal reference to its cloud network",
							Computed:    true,
						},

						// "cloudAllocationPools": {
						// Type:        schema.TypeList,
						// Optional:    true,
						// Elem: &schema.Resource{
						// 	Schema: map[string]*schema.Schema{
						// },

						// "discoveredProperties": {
						// Type:        schema.TypeList,
						// Optional:    true,
						// Elem: &schema.Resource{
						// 	Schema: map[string]*schema.Schema{
						// },

						"created": {
							Type:        schema.TypeString,
							Description: "Date when range was created in Micetro in rfc3339 time format",
							Computed:    true,
						},
						"lastmodified": {
							Type:        schema.TypeString,
							Description: "Date when range was last modified in Micetro rfc3339 time format",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func flattenRanges(ipranges []Range, tz *time.Location) ([]interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	if ipranges == nil {
		return make([]interface{}, 0), diags
	}
	flattend := make([]interface{}, len(ipranges))

	for i, iprange := range ipranges {

		flattend[i], diags = flattenRange(iprange, tz)
		if diags != nil {
			return nil, diags
		}
	}
	return flattend, diags
}

func dataSourceRangesRead(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	client := m.(*Mmclient)

	limit := d.Get("limit").(int)

	filter := map[string]interface{}{}

	if folder, ok := d.GetOk("folder"); ok {
		filter["folderRef"] = folder
	}
	if customProperties, ok := d.GetOk("custom_properties"); ok {
		for key, val := range customProperties.(map[string]interface{}) {
			filter[key] = val
		}
	}

	if isContainer, ok := d.GetOk("is_container"); ok {
		filter["isContainer"] = isContainer
	}

	if subnet, ok := d.GetOk("subnet"); ok {
		filter["subnet"] = subnet
	}

	if rawFilter, ok := d.GetOk("filter"); ok {
		filter["filter"] = rawFilter
	}

	ipranges, err := client.FindRanges(limit, filter)
	if err != nil {
		return diag.FromErr(err)
	}
	ranges, diags := flattenRanges(ipranges, client.serverLocation)
	if diags != nil {
		return diags
	}
	if err := d.Set("ranges", ranges); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(len(ipranges)))

	return diags

}
