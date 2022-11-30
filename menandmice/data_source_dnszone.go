package menandmice

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceDNSZone() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDNSZoneRead,
		Schema: map[string]*schema.Schema{

			"ref": {
				Type:        schema.TypeString,
				Description: "Internal references to this DNS zone.",
				Computed:    true,
			},

			"server": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Fully qualified name of the DNS server where the record is stored, ending with the trailing dot '.'.",
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`\.$`), "server should end with '.'"),
			},
			"view": {
				Type:        schema.TypeString,
				Description: "Name of the view this DNS zone is in.",
				Optional:    true,
				Default:     "",
			},
			"name": {
				Type:         schema.TypeString,
				Description:  "Fully qualified name of DNS zone, ending with the trailing dot '.'.",
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`\.$`), "name must end with '.'"),
				Required:     true,
			},
			"dynamic": {
				Type:        schema.TypeBool,
				Description: "If the DNS zone is dynamic.",
				Computed:    true,
			},
			"adintegrated": {
				Type:        schema.TypeBool,
				Description: "If the DNS zone is AD integrated.",
				Deprecated:  "use ad_integrated instead",
				Computed:    true,
			},
			"ad_integrated": {
				Type:        schema.TypeBool,
				Description: "If the DNS zone is AD integrated.",
				Computed:    true,
			},

			// TODO unify dnsviewref dnsviewrefs
			"dnsviewref": {
				Type:        schema.TypeString,
				Description: "Interal references to views.",
				Deprecated:  "use dns_view_ref instead",
				Computed:    true,
			},

			"dns_view_ref": {
				Type:        schema.TypeString,
				Description: "Interal references to views.",
				Deprecated:  "use dns_view_ref instead",
				Computed:    true,
			},

			"dnsviewrefs": {
				Type:        schema.TypeSet,
				Description: "Interal references to views. Only used with Active Directory.",
				Deprecated:  "use dns_view_refs instead",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"dns_view_refs": {
				Type:        schema.TypeSet,
				Description: "Interal references to views. Only used with Active Directory.",
				Deprecated:  "use dns_view_refs instead",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "The type of the DNS zone. Example: Master, Slave, Hint, Stub, Forward.",
				Computed:    true,
			},
			"authority": {
				Type:        schema.TypeString,
				Description: "The authoritative DNS server for this zone.",
				Computed:    true,
			},

			"dnssec_signed": {
				Type:        schema.TypeBool,
				Description: "If DNS signing is enabled.",
				Optional:    true,
				Default:     false,
			},
			"dnssecsigned": {
				Type:        schema.TypeBool,
				Description: "If DNS signing is enabled.",
				Computed:    true,
				Deprecated:  "use dnssec_signed instead",
			},
			"kskids": {
				Type:        schema.TypeString,
				Description: "A comma-separated string of IDs of KSKs. Starting with active keys, then inactive keys in parenthesis.",
				Computed:    true,
			},
			"zskids": {
				Type:        schema.TypeString,
				Description: "A comma-separated string of IDs of ZSKs. Starting with active keys, then inactive keys in parenthesis.",
				Computed:    true,
			},

			"customp_properties": {
				Type:        schema.TypeMap,
				Description: "Map of custom properties associated with this DNS zone.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},

			"created": {
				Type:        schema.TypeString,
				Description: "Date when zone was created in Micetro.",
				Computed:    true,
			},
			"lastmodified": {
				Type:        schema.TypeString,
				Description: "Date when zone was last modified in Micetro.",
				Computed:    true,
			},
			"displayname": {
				Type:        schema.TypeString,
				Description: "A display name to distinguish the zone from other, identically named zone instances.",
				Computed:    true,
				Deprecated:  "use display_name instead",
			},
			"display_name": {
				Type:        schema.TypeString,
				Description: "A display name to distinguish the zone from other, identically named zone instances.",
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
			return diag.Errorf("The DNS zone %v does not exist on server %v", name, server)
		}
		return diag.Errorf("The DNS zone %v does not exist in view %v on %v", name, view, server)
	}
	writeDNSZoneSchema(d, *dnszone)
	d.SetId(dnszone.Ref)

	return diags

}
