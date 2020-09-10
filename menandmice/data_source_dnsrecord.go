package menandmice

import (
	"strconv"
	"time"

	"terraform-provider-menandmice/diag"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func DataSourceDNSRec() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDNSRectRead,
		Schema: map[string]*schema.Schema{

			"fqdn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"ref": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"ttl": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"aging": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"data": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"comment": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"dnszoneref": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDNSRectRead(d *schema.ResourceData, m interface{}) error {

	var diags diag.Diagnostics
	c := m.(*Mmclient)

	err, dnsrec := c.ReadDNSRec(d.Get("fqdn").(string))

	if err != nil {
		return diag.FromErr(err)
	}
	writeDNSRecSchema(d, dnsrec)
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags

}
