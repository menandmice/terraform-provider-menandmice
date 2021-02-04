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
				Type:        schema.TypeString,
				Description: "Internal references to this DNS zone",
				Computed:    true,
			},

			"server": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "DNS server where record is stored. DNS server name should with '.' ",
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`\.$`), "server should end with '.'"),
			},
			"view": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Name of the view this DNS zone is in",
				Optional:    true,
				Default:     "",
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Description:  "Name of DNS zone. Name must and with '.' ",
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`\.$`), "name must end with '.'"),
				Required:     true,
			},
			"dynamic": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "If DNS zone Dynamic",
				Computed:    true,
			},
			// TODO following nameing convetion it would be ad_intergrated
			"adintegrated": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "If DNS zone is intergrated with Active Directory.",
				Computed:    true,
			},

			"dnsviewref": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Interal references to views.",
				Computed:    true,
			},
			"dnsviewrefs": &schema.Schema{
				Type:        schema.TypeSet,
				Description: "Interal references to views. Only used with Active Directory.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Description: "the DNS zone type.For example: Master, Slave, Hint, Stub, Forward.",
				Computed:    true,
			},
			"authority": &schema.Schema{
				Type:        schema.TypeString,
				Description: "the DNS authoritive server for this zone",
				Computed:    true,
			},
			"dnssecsigned": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "If DNS signing is enabled",
				Computed:    true,
			},
			"kskids": &schema.Schema{
				Type:        schema.TypeString,
				Description: "A comma separated string of IDs of KSKs, starting with active keys, then inactive keys in parenthesis.",
				Computed:    true,
			},
			"zskids": &schema.Schema{
				Type:        schema.TypeString,
				Description: "A comma separated string of IDs of ZSKs, starting with active keys, then inactive keys in parenthesis.",
				Computed:    true,
			},

			"customp_properties": &schema.Schema{
				Type:        schema.TypeMap,
				Description: "Map of custom properties associated with this DNS zone.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},

			"created": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Date when zone was created in the suite.",
				Computed:    true,
			},
			"lastmodified": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Date when zone was last modified in the suite.",
				Computed:    true,
			},
			"displayname": &schema.Schema{
				Type:        schema.TypeString,
				Description: "A name that can distinguish the zone from other zone instances with the same name.",
				Computed:    true,
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
