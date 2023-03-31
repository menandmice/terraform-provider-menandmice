package menandmice

import (
	"context"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"endpoint": {
					Type: schema.TypeString,
					// Required:    true,  // can be set via environment variable
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("MENANDMICE_ENDPOINT", nil),
					Description: "Micetro API endpoint",
				},
				"username": {
					Type: schema.TypeString,
					// Required:    true,  // can be set via environment variable
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("MENANDMICE_USERNAME", nil),
					Description: "Micetro username",
				},
				"password": {
					Type: schema.TypeString,
					// Required:    true, // can be set via environment variable
					Optional:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("MENANDMICE_PASSWORD", nil),
					Description: "Micetro password",
				},
				"tls_verify": {
					Type:        schema.TypeBool,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("MENANDMICE_TLS_VERIFY", true),
					Description: "Micetro SSL validation",
				},
				"timeout": {
					Type:        schema.TypeInt,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("MENANDMICE_TIMEOUT", 30),
					Description: "Micetro Request timeout",
				},
				"server_timezone": {
					Type:             schema.TypeString,
					Optional:         true,
					ValidateDiagFunc: validTZ(),
					DefaultFunc:      schema.EnvDefaultFunc("MENANDMICE_SERVER_TIMEZONE", nil),
					Description: `Timezone of Mictro server.
						in IANA Time Zone format. example: America/Chicago.
						See;https://en.wikipedia.org/wiki/List_of_tz_database_time_zones .
						Default to local time zone.
						If not set correcly terraform will print wrong times for things like creation and modiviaction dates`,
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
				"menandmice_ranges":           DataSourceRanges(),
			},
		}
		p.ConfigureContextFunc = configure(version, p)
		return p
	}
}

func validTZ() schema.SchemaValidateDiagFunc {
	return func(v interface{}, k cty.Path) diag.Diagnostics {
		value, ok := v.(string)
		var diags diag.Diagnostics
		if !ok {
			return diag.Errorf("expected type of %s to be string", k)
		}
		if _, err := time.LoadLocation(value); err != nil {
			return diag.FromErr(err)
		}
		return diags
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		var diags diag.Diagnostics

		params := Cfg{
			MMEndpoint: d.Get("endpoint").(string),
			MMUsername: d.Get("username").(string),
			MMPassword: d.Get("password").(string),
			MMTimezone: d.Get("server_timezone").(string),
			TLSVerify:  d.Get("tls_verify").(bool),
			Timeout:    d.Get("timeout").(int),
			Version:    version,
			Debug:      true, // FIXME
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
		if params.MMTimezone == "" {
			tflog.Warn(ctx, "No Timezone set for Mictro server. Will assume is running in same timezone as this machine")
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
