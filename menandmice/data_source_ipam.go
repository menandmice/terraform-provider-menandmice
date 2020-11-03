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
				Type:     schema.TypeString,
				Computed: true,
			},
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.Any(
					validation.IsIPv4Address,
					validation.IsIPv6Address),
				ForceNew: true,
			},
			"claimed": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			// "dnshost": &schema.Schema{
			// },
			// "dhcp_reservations": &schema.Schema{
			// },
			// "dhcp_leases": &schema.Schema{
			// },
			"discovery_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_seen_date": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"last_discovery_date": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_known_client_identifier": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"device": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"interface": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"ptr_status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"extraneous_ptr": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"custom_properties": &schema.Schema{
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"state": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				// ValidateFunc: validation.StringInSlice([]string{
				// 	"Free", "Assigned", "Claimed", "Pending", "Held",
				// }, false),
			},
			"hold_info": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
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
				Type:     schema.TypeBool,
				Computed: true,
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
