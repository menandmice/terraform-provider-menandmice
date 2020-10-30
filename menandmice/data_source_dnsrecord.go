package menandmice

import (
	"regexp"
	"strconv"
	"time"

	"terraform-provider-menandmice/diag"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func DataSourceDNSRec() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDNSRectRead,
		Schema: map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"view": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"zone": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`\.$`), "server should end with '.'"),
			},
			"server": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`\.$`), "server should end with '.'"),
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"ref": &schema.Schema{
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

	dnsZoneRef := tryGetString(d, "server") + ":" + tryGetString(d, "view") + ":" + tryGetString(d, "zone")

	dnsrecs, err := c.FindDNSRec(dnsZoneRef, map[string]string{
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
