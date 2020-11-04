package menandmice

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceDHCPReservation() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDHCPResvationRead,
		Schema: map[string]*schema.Schema{

			"ref": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_identifier": &schema.Schema{

				Type:     schema.TypeString,
				Computed: true,
			},

			"reservation_method": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"addresses": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ddns_hostname": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"filename": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"servername": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"next_server": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			// TODO one off dhcpserver,dhcpgroup,dhcpscope

			"owner_ref": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDHCPResvationRead(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	client := m.(*Mmclient)

	name := d.Get("name").(string)
	dhcpReservation, err := client.ReadDHCPReservation(name)
	if err != nil {
		return diag.FromErr(err)
	}

	if dhcpReservation == nil {
		return diag.Errorf("dhcp_reservation %v does not exist", name)
	}

	writeDHCPReservationSchema(d, *dhcpReservation)
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags

}
