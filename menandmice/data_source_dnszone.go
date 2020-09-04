package menandmice

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceDNSzone() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDNSzoneRead,
		Schema: map[string]*schema.Schema{

			"domain": &schema.Schema{
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
			"dynamic": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"adintegrated": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},

			// TODO maybe choose between ref or refs automatic
			"dnsviewref": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"dnsviewrefs": &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"authority": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"dnssecsigned": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"kskids": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"zskids": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			// "customProperties": &schema.Schema{
			// 	Type:     ?
			// 	Computed: true,
			// }
			"created": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"lastmodified": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"displayname": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDNSzoneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	c := m.(*Mmclient)

	err, dnszone := c.ReadDNSzone(d.Get("domain").(string))

	if err != nil {
		return diag.FromErr(err)
	}
	writeDNSzoneSchema(d, dnszone)
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags

}
