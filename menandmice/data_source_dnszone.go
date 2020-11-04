package menandmice

import (
	"context"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// schema for DNSZone resource
func DataSourceDNSZone() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDNSZoneRead,
		Schema: map[string]*schema.Schema{

			"ref": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"server": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`\.$`), "server should end with '.'"),
			},
			"view": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`\.$`), "name must end with '.'"),
				Required:     true,
			},
			"dynamic": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"adintegrated": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},

			"dnsviewref": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
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

			"customp_properties": &schema.Schema{
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},

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

func dataSourceDNSZoneRead(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	client := m.(*Mmclient)

	server := tryGetString(d, "server")
	view := tryGetString(d, "view")
	name := tryGetString(d, "name")
	dnsZoneRef := server + ":" + view + ":" + name
	dnszone, err := client.ReadDNSZone(dnsZoneRef)

	if err != nil {
		return diag.FromErr(err)
	}

	if dnszone == nil {
		if view == "" {
			return diag.Errorf("dnszone %v does not exist on server %v", name, server)
		} else {
			return diag.Errorf("dnszone %v does not exist in view %v on %v", name, view, server)
		}
	}
	writeDNSZoneSchema(d, *dnszone)
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags

}
