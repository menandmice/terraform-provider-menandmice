package menandmice

import (
	"terraform-provider-menandmice/diag"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"endpoint": &schema.Schema{
				Type: schema.TypeString,
				// Required:    true,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MENANDMICE_ENDPOINT", nil),
				Description: "Men&Mice Web API endpoint",
			},
			"username": &schema.Schema{
				Type: schema.TypeString,
				// Required:    true,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MENANDMICE_USERNAME", nil),
				Description: "Men&Mice username",
			},
			"password": &schema.Schema{
				Type: schema.TypeString,
				// Required:    true,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("MENANDMICE_PASSWORD", nil),
				Description: "Men&Mice password",
			},
			"tls_verify": &schema.Schema{
				Type: schema.TypeBool,
				// Required:    true,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MENANDMICE_TLS_VERIFY", true),
				Description: "Men&Mice SSL validation",
			},
			"timeout": &schema.Schema{
				Type: schema.TypeInt,
				// Required:    true,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MENANDMICE_TIMEOUT", 30),
				Description: "Men&Mice Request timeout",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"menandmice_dns_record":  resourceDNSRec(),
			"menandmice_dns_zone":    resourceDNSZone(),
			"menandmice_ipam_record": resourceIPAMRec(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"menandmice_dns_record":  DataSourceDNSRec(),
			"menandmice_dns_zone":    DataSourceDNSZone(),
			"menandmice_ipam_record": DataSourceIPAMRec(),
		},
		ConfigureFunc: providerConfigure,
	}
}
func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	var diags diag.Diagnostics

	params := Cfg{
		MMEndpoint: d.Get("endpoint").(string),
		MMUsername: d.Get("username").(string),
		MMPassword: d.Get("password").(string),
		TLSVerify:  d.Get("tls_verify").(bool),
		Timeout:    d.Get("timeout").(int),
	}

	if params.MMEndpoint == "" {
		diags = diag.Append(diags, diag.Errorf("REST API endpoint set for provider menandmice."))
	}
	if params.MMUsername == "" {
		diags = diag.Append(diags, diag.Errorf("No username set for provider menandmice."))
	}
	if params.MMPassword == "" {
		diags = diag.Append(diags, diag.Errorf("No password set for provider menandmice."))
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
