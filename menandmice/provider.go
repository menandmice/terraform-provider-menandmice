package menandmice

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"web": &schema.Schema{
				Type: schema.TypeString,
				// Required:    true,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MENANDMICE_WEB", nil),
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
			"menandmice_dnsrecord": resourceDNSrec(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"menandmice_dnsrecord": DataSourceDNSrec(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	var diags diag.Diagnostics

	params := Cfg{
		MMWeb:      d.Get("web").(string),
		MMUsername: d.Get("username").(string),
		MMPassword: d.Get("password").(string),
		TLSVerify:  d.Get("tls_verify").(bool),
		Timeout:    d.Get("timeout").(int),
	}

	//TODO better error messages
	client, err := ClientInit(&params)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return client, diags
}
