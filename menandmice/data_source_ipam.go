package menandmice

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIPAMRec() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIPAMRecRead,
		Schema: map[string]*schema.Schema{

			"ref": {
				Type:        schema.TypeString,
				Description: "Internal reference for the IP address.",
				Computed:    true,
			},
			"address": {
				Type:        schema.TypeString,
				Description: "The IP address.",
				Required:    true,
				ValidateFunc: validation.Any(
					validation.IsIPv4Address,
					validation.IsIPv6Address),
				ForceNew: true,
			},
			"claimed": {
				Type:        schema.TypeBool,
				Description: "If the IP address is claimed.",
				Computed:    true,
			},
			// "dnshost": &schema.Schema{
			// },
			// "dhcp_reservations": &schema.Schema{
			// },
			// "dhcp_leases": &schema.Schema{
			// },
			"discovery_type": {
				Type:        schema.TypeString,
				Description: "The discovery method of the IP address. Example: None, Ping, ARP, Lease, Custom.",
				Computed:    true,
			},
			"last_seen_date": {
				Type:        schema.TypeString,
				Description: "The date when the address was last seen during IP address discovery in rfc3339 time format.",
				Computed:    true,
			},

			"last_discovery_date": {
				Type:        schema.TypeString,
				Description: "The date when the system last performed IP address discovery for this IP address rfc3339 time format.",
				Computed:    true,
			},
			"last_known_client_identifier": {
				Type:        schema.TypeString,
				Description: "The last known MAC address associated with the IP address discovery information.",
				Computed:    true,
			},

			"device": {
				Type:        schema.TypeString,
				Description: "The device associated with the object.",
				Computed:    true,
			},

			"interface": {
				Type:        schema.TypeString,
				Description: "The interface associated with the object.",
				Computed:    true,
			},
			"ptr_status": {
				Type:        schema.TypeString,
				Description: "PTR record status. Example: Unknown, OK, Verify.",
				Computed:    true,
			},
			"extraneous_ptr": {
				Type:        schema.TypeBool,
				Description: "'True' if there are extraneous PTR records for the object.",
				Computed:    true,
			},
			"custom_properties": {
				Type:        schema.TypeMap,
				Description: "Map of custom properties associated with this IP address.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"state": {
				Type:        schema.TypeString,
				Description: "The state of the IP address. Example: Free, Assigned, Claimed, Pending, Held.",
				Computed:    true,
			},
			"hold_info": {
				Type:        schema.TypeList,
				Description: "Contains information about who holds the otherwise free IP, and for how long.",
				Computed:    true,
				// MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"expiry_time": {
							Type:     schema.TypeString,
							Computed: true,
							// ValidateFunc: validation.ValidateRFC3339TimeString,
						},

						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"usage": {
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
	diags = writeIPAMRecSchema(d, ipam, client.serverLocation)
	d.SetId(ipam.Ref)

	return diags

}
