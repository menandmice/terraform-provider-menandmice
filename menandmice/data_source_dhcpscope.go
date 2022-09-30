package menandmice

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceDHCPScope() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDHCPScopeRead,
		Schema: map[string]*schema.Schema{

			"ref": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Internal reference to this DHCP reservation.",
				Computed:    true,
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The name of the DHCP scope you want to query.",
				Computed:    true,
			},

			"cidr": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The cidr of DHCPScope.",
				// TODO validate
				Required: true,
			},

			"dhcp_server": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The DHCP server of this scope.",
				// TODO validate
				Optional: true,
			},

			"superscope": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The name of the superscope for the DHCP scope. Only applicable for MS DHCP servers.",
				Computed:    true,
			},

			"description": &schema.Schema{
				Type:        schema.TypeString,
				Description: "A description for the DHCP scope.",
				Computed:    true,
			},

			"available": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "Number of available addresses in the address pool(s) of the scope.",
				Computed:    true,
			},

			"enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "If this scope is enabled",
				Computed:    true,
			},
		},
	}
}

func writeDHCPScopeSchema(d *schema.ResourceData, dhcpScope DHCPScope) {

	d.Set("ref", dhcpScope.Ref)
	d.Set("name", dhcpScope.Name)
	d.Set("cidr", dhcpScope.RangeRef)
	d.Set("dhcp_server", dhcpScope.DHCPServerRef)
	d.Set("superscope", dhcpScope.Superscope)
	d.Set("description", dhcpScope.Description)
	d.Set("available", dhcpScope.Available)
	d.Set("enabled", dhcpScope.Enabled)
	return
}
func dataSourceDHCPScopeRead(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	client := m.(*Mmclient)

	filter := map[string]string{"RangeRef": d.Get("cidr").(string)}
	if dhcpServerRef, ok := d.GetOk("dhcp_server"); ok {
		filter["dhcpServerRef"] = dhcpServerRef.(string)
	}
	dhcpScopes, err := client.FindDHCPScope(filter)

	if err != nil {
		return diag.FromErr(err)
	}

	switch {
	case len(dhcpScopes) <= 0:
		return diag.Errorf("No matching DHCP scopes were found.")
		// TODO comment why not needed
		// case len(dnsrecs) > 1:
		// 	return diag.Errorf("%v DNSRecords found matching you criteria, but should be only 1", len(dnsrecs))
	}

	writeDHCPScopeSchema(d, dhcpScopes[0])
	d.SetId(dhcpScopes[0].Ref)

	return diags

}
