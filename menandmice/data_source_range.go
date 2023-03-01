package menandmice

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceRange() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRangeRead,
		Schema: map[string]*schema.Schema{

			// TODO add atributes: cloudAllocationPools, dhcpScopes authority ,discoveredProperties
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
			// "cidr": {
			// 	Type:         schema.TypeString,
			// 	Description:  "The CIDR of the range",
			// 	ExactlyOneOf: []string{"cidr", "from"},
			// 	Optional:     true,
			// },
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
				Description: "DDate when zone was created in Micetro.",
				Computed:    true,
			},
			"lastmodified": {
				Type:        schema.TypeString,
				Description: "Date when zone was last modified in Micetro.",
				Computed:    true,
			},
		},
	}
}

func dataSourceRangeRead(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	client := m.(*Mmclient)

	name := d.Get("name").(string)
	iprange, err := client.ReadRange(name)
	if err != nil {
		return diag.FromErr(err)
	}

	if iprange == nil {
		return diag.Errorf("range_%v does not exist", name)
	}

	writeRangeSchema(d, *iprange)
	d.SetId(iprange.Ref)

	return diags

}
