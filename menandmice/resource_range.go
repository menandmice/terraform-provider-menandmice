package menandmice

import (
	// "context"

	// "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceRange() *schema.Resource {
	return &schema.Resource{
		// CreateContext: resourceDNSZoneCreate,
		ReadContext: resourceRangeRead,
		// UpdateContext: resourceDNSZoneUpdate,
		// DeleteContext: resourceDNSZoneDelete,
		// Importer: &schema.ResourceImporter{
		// 	StateContext: resourceDNSZoneImport,
		// },
		Schema: map[string]*schema.Schema{

			"ref": {
				Type:        schema.TypeString,
				Description: "Internal references to this range.",
				Computed:    true,
			},
			"name": {
				Type:         schema.TypeString,
				Description:  "Name and CIDR of range.",
				Required:     true,
				ValidateFunc: validation.IsCIDR,
			},
			"from": {
				Type: schema.TypeBool,
				// Description: "If the DNS zone is dynamic. (Default: False)",
				Optional:     true,
				Default:      false,
				ValidateFunc: validation.IsIPAddress,
			},
			// TODO following nameing convetion it would be ad_intergrated
			"to": {
				Type: schema.TypeBool,
				// Description: "If the DNS zone is AD integrated. (Default: False)",
				Optional:     true,
				Default:      false,
				ValidateFunc: validation.IsIPAddress,
			},
		},
	}
}

func writeRangeSchema(d *schema.ResourceData, iprange Range) {

	d.Set("ref", iprange.Ref)
	d.Set("name", iprange.Name)
	d.Set("from", iprange.To)
	d.Set("to", iprange.From)

	return
}

//
// func readDNSZoneSchema(d *schema.ResourceData) DNSZone {
//
// 	if dnsViewRefsRead, ok := d.GetOk("dnsviewrefs"); ok {
// 		dnsViewRefList := dnsViewRefsRead.(*schema.Set).List()
// 		var dnsViewRefs = make([]string, len(dnsViewRefList))
// 		for i, view := range dnsViewRefList {
// 			dnsViewRefs[i] = view.(string)
// 		}
// 	}
//
// 	var CustomProperties = make(map[string]string)
// 	if customPropertiesRead, ok := d.GetOk("custom_properties"); ok {
// 		for key, value := range customPropertiesRead.(map[string]interface{}) {
// 			CustomProperties[key] = value.(string)
// 		}
// 	}
// 	dnszone := DNSZone{
// 		Ref:          tryGetString(d, "ref"),
// 		AdIntegrated: d.Get("adintegrated").(bool),
// 		Authority:    tryGetString(d, "authority"),
//
// 		// you should not set this yourself
// 		// Created:      d.Get("created").(string),
// 		// LastModified: tryGetString(d, "lastmodified"),
//
// 		DNSZoneProperties: DNSZoneProperties{
// 			Name:              d.Get("name").(string),
// 			Dynamic:           d.Get("dynamic").(bool),
// 			ZoneType:          tryGetString(d, "type"),
// 			DnssecSigned:      d.Get("dnssecsigned").(bool),
// 			KskIDs:            tryGetString(d, "kskids"),
// 			ZskIDs:            tryGetString(d, "zskids"),
// 			AdReplicationType: tryGetString(d, "adreplicationtype"),
// 			AdPartition:       tryGetString(d, "adpartition"),
// 			CustomProperties:  CustomProperties,
// 			DisplayName:       tryGetString(d, "displayname"),
// 		},
// 	}
//
// 	dnsviewref := dnszone.Authority + ":" + tryGetString(d, "view")
// 	if dnszone.AdIntegrated {
//
// 		// TODO this does not work
// 		dnszone.DNSViewRefs = []string{dnsviewref}
// 	} else {
// 		dnszone.DNSViewRef = dnsviewref
//
// 	}
//
// 	return dnszone
// }
//
// func resourceDNSZoneCreate(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
// 	client := m.(*Mmclient)
//
// 	var masters []string
// 	if mastersRead, ok := d.Get("masters").([]interface{}); ok {
// 		masters = make([]string, len(mastersRead))
// 		for i, master := range mastersRead {
// 			masters[i] = master.(string)
// 		}
// 	} else {
// 		return diag.Errorf("Could not read masters")
// 	}
//
// 	dnszone := readDNSZoneSchema(d)
//
// 	objRef, err := client.CreateDNSZone(dnszone, masters)
//
// 	if err != nil {
// 		return diag.FromErr(err)
// 	}
// 	d.SetId(objRef)
//
// 	return resourceDNSZoneRead(c, d, m)
//
// }
//
func resourceRangeRead(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := m.(*Mmclient)

	iprange, err := client.ReadRange(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if iprange == nil {
		d.SetId("")
		return diags
	}

	writeRangeSchema(d, *iprange)

	return diags
}

//
// func resourceDNSZoneUpdate(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
//
// 	//can't change read only property
// 	if d.HasChange("ref") || d.HasChange("adintegrated") ||
// 		d.HasChange("dnsviewref") || d.HasChange("dnsviewrefs") ||
// 		d.HasChange("authority") {
// 		// this can't never error can never happen because of "ForceNew: true," for these properties
// 		return diag.Errorf("can't change read-only property of DNS zone")
// 	}
// 	client := m.(*Mmclient)
// 	ref := d.Id()
// 	dnszone := readDNSZoneSchema(d)
//
// 	err := client.UpdateDNSZone(dnszone.DNSZoneProperties, ref)
//
// 	if err != nil {
// 		return diag.FromErr(err)
// 	}
// 	return resourceDNSZoneRead(c, d, m)
// }
//
// func resourceDNSZoneDelete(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
//
// 	client := m.(*Mmclient)
// 	var diags diag.Diagnostics
// 	ref := d.Id()
// 	err := client.DeleteDNSZone(ref)
// 	if err != nil {
// 		return diag.FromErr(err)
// 	}
// 	d.SetId("")
// 	return diags
// }
//
// func resourceDNSZoneImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
//
// 	diags := resourceDNSZoneRead(ctx, d, m)
// 	if err := toError(diags); err != nil {
// 		return nil, err
// 	}
//
// 	// if we had used schema.ImportStatePassthrough
// 	// we could not have set id to its canonical form
// 	d.SetId(d.Get("ref").(string))
//
// 	return []*schema.ResourceData{d}, nil
// }
