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

type DNSRecord struct { // TODO do we neet point if omit empty
	Ref        *string `json:"ref,omitempty"`
	Name       string  `json:"name"`
	Rectype    string  `json:"type"`
	Ttl        *string `json:"ttl,omitempty"`
	Data       string  `json:"data"`
	Comment    string  `json:"comment,omitempty"`
	Aging      int     `json:"aging,omitempty"`
	Enabled    bool    `json:"enabled,omitempty"`
	DNSZoneRef string  `json:"dnsZoneRef"`
}

type Response struct {
	Result struct {
		DNSRecord `json:"dnsRecord"`
	} `json:"result"`
}

func dataSourceDNSrectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	c := m.(*resty.Client)

	var re Response
	diags = MmGet(c, &re, "dnsrecords/"+(d.Get("domain").(string)))

	d.Set("ref", re.Result.DNSRecord.Ref)
	d.Set("name", re.Result.DNSRecord.Name)
	d.Set("type", re.Result.DNSRecord.Rectype)
	if re.Result.DNSRecord.Ttl != nil {
		ttl, err := strconv.Atoi(*re.Result.DNSRecord.Ttl)
		if err == nil {
			d.Set("ttl", ttl)
		}
	}
	d.Set("enabled", re.Result.DNSRecord.Enabled)
	d.Set("dnszoneref", re.Result.DNSRecord.DNSZoneRef)

	d.Set("aging", re.Result.DNSRecord.Aging)     //TODO default is no 0
	d.Set("comment", re.Result.DNSRecord.Comment) // comment is always given, but sometimes ""
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags

}
