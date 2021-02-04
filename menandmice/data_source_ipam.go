package menandmice

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIPAMRec() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIPAMRecRead,
		Schema: map[string]*schema.Schema{

			"ref": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Internal reference to ipam record",
				Computed:    true,
			},
			"address": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The IP address",
				Required:    true,
				ValidateFunc: validation.Any(
					validation.IsIPv4Address,
					validation.IsIPv6Address),
				ForceNew: true,
			},
			"claimed": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "If this address is claimed",
				Computed:    true,
			},
			// "dnshost": &schema.Schema{
			// },
			// "dhcp_reservations": &schema.Schema{
			// },
			// "dhcp_leases": &schema.Schema{
			// },
			"discovery_type": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Way IP address use is dicoverd. For example: None, Ping, ARP, Lease, Custom.",
				Computed:    true,
			},
			"last_seen_date": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The date when the address was last seen during IP address discovery.",
				Computed:    true,
			},

			"last_discovery_date": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The date when the system last performed IP address discovery for this IP address.",
				Computed:    true,
			},
			"last_known_client_identifier": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The MAC address associated with the IP address discovery info.",
				Computed:    true,
			},

			"device": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The device associated with the record.",
				Computed:    true,
			},

			"interface": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The interface associated with the record.",
				Computed:    true,
			},
			"ptr_status": &schema.Schema{
				Type:        schema.TypeString,
				Description: "PTR record status. For example: Unknown, OK, Verify.",
				Computed:    true,
			},
			"extraneous_ptr": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "Contains true if there are extraneous PTR records for the record.",
				Computed:    true,
			},
			"custom_properties": &schema.Schema{
				Type:        schema.TypeMap,
				Description: "Map of custom properties associated with this IP address.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"state": &schema.Schema{
				Type:        schema.TypeString,
				Description: "state of IP addres. For exampe: Free, Assigned, Claimed, Pending, Held.",
				Computed:    true,
			},
			"hold_info": &schema.Schema{
				Type:        schema.TypeList,
				Description: "Contains information about who holds the otherwise free IP and for how long.",
				Computed:    true,
				// MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"expiry_time": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
							// ValidateFunc: validation.ValidateRFC3339TimeString,
						},

						"username": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"usage": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "IP address usage bitmask.",
				Computed:    true,
			},

			// "cloud_device_info": &schema.Schema{
			// },
		},
	}
}

func dataSourceIPAMRecRead(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	client := m.(*Mmclient)

	ipam, err := client.ReadIPAMRec(d.Get("address").(string))

	if err != nil {
		return diag.FromErr(err)
	}
	writeIPAMRecSchema(d, ipam)
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags

}
