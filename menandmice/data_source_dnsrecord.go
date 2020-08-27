package menandmice

import (
	"context"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceDNSrec() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDNSrectRead,
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
			"dnszoneref": &schema.Schema{ //FIXME nameing confention
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDNSrectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	c := m.(*resty.Client)

	err, re := ReadDNSRec(c, d.Get("domain").(string))

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("ref", re.Ref)
	d.Set("name", re.Name)
	d.Set("type", re.Rectype)
	if re.Ttl != nil {
		ttl, err := strconv.Atoi(*re.Ttl)
		if err == nil {
			d.Set("ttl", ttl)
		}
	}
	d.Set("enabled", re.Enabled)
	d.Set("dnszoneref", re.DNSZoneRef)

	d.Set("aging", re.Aging)     //TODO default is no 0
	d.Set("comment", re.Comment) // comment is always given, but sometimes ""
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags

}
