package menandmice

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceDNSRec() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDNSRectRead,
		Schema: map[string]*schema.Schema{

			"name": {
				Type:         schema.TypeString,
				Description:  "The DNS record name.",
				Required:     true,
				ValidateFunc: validation.StringDoesNotMatch(regexp.MustCompile(`\.$`), "Hostname should not end with '.'"),
			},
			"view": {
				Type:        schema.TypeString,
				Description: "The view of the DNS record. Example: internal.",
				Optional:    true,
				Default:     "",
			},
			"zone": {
				Type:         schema.TypeString,
				Description:  "The DNS zone where the record is stored.",
				Required:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`\.$`), "server should end with '.'"),
			},
			"server": {
				Type:         schema.TypeString,
				Description:  "The DNS server where the DNS record is stored.",
				Required:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`\.$`), "server should end with '.'"),
			},
			"type": {
				Type:        schema.TypeString,
				Description: "The DNS record type. Example: CNAME.",
				Required:    true,
			},

			"ref": {
				Type:        schema.TypeString,
				Description: "Internal reference to this DNS record.",
				Computed:    true,
			},
			"ttl": {
				Type:        schema.TypeInt,
				Description: "The DNS record's Time To Live value in seconds, setting how long the record is allowed to be cached.",
				Computed:    true,
			},
			"aging": {
				Type:        schema.TypeInt,
				Description: "The aging timestamp of dynamic records in AD integrated zones. Hours since January 1, 1601, UTC. Providing a non-zero value creates a dynamic record.",
				Computed:    true,
			},
			"data": {
				Type:        schema.TypeString,
				Description: "The data stored in the record",
				Computed:    true,
			},
			"comment": {
				Type:        schema.TypeString,
				Description: "Contains the comment string for the record. Only records in static DNS zones can have a comment string. Some cloud DNS provides do not support comments.",
				Computed:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "If the DNS record is enabled.",
				Computed:    true,
			},
			"dns_zone_ref": {
				Type:        schema.TypeString,
				Description: "Internal reference to the zone where this DNS record is stored.",
				Computed:    true,
			},
			"fqdn": {
				Type:        schema.TypeString,
				Description: "Fully qualified domain name of this DNS record.",
				Computed:    true,
			},
		},
	}
}

func dataSourceDNSRectRead(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	client := m.(*Mmclient)

	dnsZoneRef := tryGetString(d, "server") + ":" + tryGetString(d, "view") + ":" + tryGetString(d, "zone")

	dnsrecs, err := client.FindDNSRec(dnsZoneRef, map[string]interface{}{
		"name": tryGetString(d, "name"),
		"type": tryGetString(d, "type"),
	})

	if err != nil {
		return diag.FromErr(err)
	}

	switch {
	case len(dnsrecs) <= 0:
		return diag.Errorf("No matching DNS record was found")
	case len(dnsrecs) > 1:
		return diag.Errorf("Found %v DNS records matching you criteria, but should be only 1", len(dnsrecs))
	}
	writeDNSRecSchema(d, dnsrecs[0])
	d.SetId(dnsrecs[0].Ref)

	return diags

}
