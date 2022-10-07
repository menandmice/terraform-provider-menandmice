package menandmice

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceDHCPReservation() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDHCPResvationRead,
		Schema: map[string]*schema.Schema{

			"ref": {
				Type:        schema.TypeString,
				Description: "Internal reference to this DHCP reservation.",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the DHCP reservation you want to query.",
				Required:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "The type of this DHCP reservation. Example: DHCP , BOOTP , BOTH.",
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description for the reservation. Only applicable for MS DHCP servers.",
				Computed:    true,
			},
			"client_identifier": {
				Type:        schema.TypeString,
				Description: "The client_identifier of this reservation.",
				Computed:    true,
			},

			"reservation_method": {
				Type:        schema.TypeString,
				Description: "DHCP reservation method, Example: HardwareAddress , ClientIdentifier.",
				Computed:    true,
			},
			"addresses": {
				Type:        schema.TypeList,
				Description: "A list of IP addresses used for the reservation.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ddns_hostname": {
				Type:        schema.TypeString,
				Description: "Dynamic DNS host name for the reservation. Only applicable for ISC DHCP servers.",
				Computed:    true,
			},

			"filename": {
				Type:        schema.TypeString,
				Description: "The filename DHCP option. Only applicable for ISC DHCP servers.",
				Computed:    true,
			},
			"servername": {
				Type:        schema.TypeString,
				Description: "The server-name DHCP option. Only applicable for ISC DHCP servers.",
				Computed:    true,
			},
			"next_server": {
				Type:        schema.TypeString,
				Description: "The next-server ISC DHCP option. Only applicable for ISC DHCP servers.",
				Computed:    true,
			},
			// TODO one off dhcpserver,dhcpgroup,dhcpscope

			"owner_ref": {
				Type:        schema.TypeString,
				Description: "Internal reference to the DHCP group scope or server where this reservation is made.",
				Computed:    true,
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
	d.SetId(dhcpReservation.Ref)

	return diags

}
