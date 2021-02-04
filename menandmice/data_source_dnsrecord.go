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

func DataSourceDNSRec() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDNSRectRead,
		Schema: map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The DNS record name.",
				Required:    true,
			},
			"view": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The view of DNS record. For example internal.",
				Optional:    true,
				Default:     "",
			},
			"zone": &schema.Schema{
				Type:         schema.TypeString,
				Description:  "The DNS zone were record is in.",
				Required:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`\.$`), "server should end with '.'"),
			},
			"server": &schema.Schema{
				Type:         schema.TypeString,
				Description:  "The DNS server where DNS record is stored",
				Required:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`\.$`), "server should end with '.'"),
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The DNS recod type. for example: CNAME",
				Required:    true,
			},

			"ref": &schema.Schema{
				Type:        schema.TypeString,
				Description: "internal reference to this DNS record",
				Computed:    true,
			},
			"ttl": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "The DNS recod Time To Live. How long in seconds the record is allowed to be cached",
				Computed:    true,
			},
			"aging": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "The aging timestamp of dynamic records in AD integrated zones. Hours since January 1, 1601, UTC. Providing a non-zero value creates a dynamic record.",
				Computed:    true,
			},
			"data": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The data stored in the record",
				Computed:    true,
			},
			"comment": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Contains the comment string for the record. Note that only records in static DNS zones can have a comment string.",
				Computed:    true,
			},
			"enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "If DNS record is enabled",
				Computed:    true,
			},

			"dns_zone_ref": &schema.Schema{
				Type:        schema.TypeString,
				Description: "internal reference to zone where this DNS record is store",
				Computed:    true,
			},
		},
	}
}

func dataSourceDNSRectRead(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	client := m.(*Mmclient)

	dnsZoneRef := tryGetString(d, "server") + ":" + tryGetString(d, "view") + ":" + tryGetString(d, "zone")

	dnsrecs, err := client.FindDNSRec(dnsZoneRef, map[string]string{
		"name": tryGetString(d, "name"),
		"type": tryGetString(d, "type"),
	})

	if err != nil {
		return diag.FromErr(err)
	}

	switch {
	case len(dnsrecs) <= 0:
		return diag.Errorf("no DNSRecord found matching you criteria")
	case len(dnsrecs) > 1:
		return diag.Errorf("%v DNSRecords found matching you criteria, but should be only 1", len(dnsrecs))
	}
	writeDNSRecSchema(d, dnsrecs[0])
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags

}
