package menandmice

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"endpoint": &schema.Schema{
					Type: schema.TypeString,
					// Required:    true,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("MENANDMICE_ENDPOINT", nil),
					Description: "Micetro API endpoint",
				},
				"username": &schema.Schema{
					Type: schema.TypeString,
					// Required:    true,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("MENANDMICE_USERNAME", nil),
					Description: "Micetro username",
				},
				"password": &schema.Schema{
					Type: schema.TypeString,
					// Required:    true,
					Optional:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("MENANDMICE_PASSWORD", nil),
					Description: "Micetro password",
				},
				"tls_verify": &schema.Schema{
					Type: schema.TypeBool,
					// Required:    true,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("MENANDMICE_TLS_VERIFY", true),
					Description: "Micetro SSL validation",
				},
				"timeout": &schema.Schema{
					Type: schema.TypeInt,
					// Required:    true,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("MENANDMICE_TIMEOUT", 30),
					Description: "Micetro Request timeout",
				},
			},
			ResourcesMap: map[string]*schema.Resource{
				"menandmice_dns_record":       resourceDNSRec(),
				"menandmice_dns_zone":         resourceDNSZone(),
				"menandmice_ipam_record":      resourceIPAMRec(),
				"menandmice_dhcp_reservation": resourceDHCPReservation(),
				"menandmice_range":            resourceRange(),
			},
			DataSourcesMap: map[string]*schema.Resource{
				"menandmice_dns_record":       DataSourceDNSRec(),
				"menandmice_dns_zone":         DataSourceDNSZone(),
				"menandmice_dns_zones":        DataSourceDNSZones(),
				"menandmice_ipam_record":      DataSourceIPAMRec(),
				"menandmice_dhcp_reservation": DataSourceDHCPReservation(),
				"menandmice_dhcp_scope":       DataSourceDHCPScope(),
				"menandmice_range":            DataSourceRange(),
			},
		}
		p.ConfigureContextFunc = configure(version, p)
		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		var diags diag.Diagnostics

		params := Cfg{
			MMEndpoint: d.Get("endpoint").(string),
			MMUsername: d.Get("username").(string),
			MMPassword: d.Get("password").(string),
			TLSVerify:  d.Get("tls_verify").(bool),
			Timeout:    d.Get("timeout").(int),
			Version:    version,
			// Debug:      true,
		}

		if params.MMEndpoint == "" {
			diags = append(diags, diag.Errorf("No REST API endpoint set for provider menandmice.")...)
		}
		if params.MMUsername == "" {
			diags = append(diags, diag.Errorf("No username set for provider menandmice.")...)
		}
		if params.MMPassword == "" {
			diags = append(diags, diag.Errorf("No password set for provider menandmice.")...)
		}
		if diags != nil {
			return nil, diags
		}
		client, err := ClientInit(&params)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		return client, diags
	}
}
