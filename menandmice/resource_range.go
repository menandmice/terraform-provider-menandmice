package menandmice

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceRange() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRangeCreate,
		ReadContext:   resourceRangeRead,
		UpdateContext: resourceRangeUpdate,
		DeleteContext: resourceRangeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRangeImport,
		},
		Schema: map[string]*schema.Schema{

			"ref": {
				Type:        schema.TypeString,
				Description: "Internal references to this range.",
				Computed:    true,
			},

			"name": {
				Type:        schema.TypeString,
				Description: "The CIDR of the range, or from-to address range.",
				Computed:    true,
			},
			"cidr": {
				Type:         schema.TypeString,
				Description:  "The CIDR of the range",
				ExactlyOneOf: []string{"cidr", "from", "free_range"},
				ForceNew:     true,
				Optional:     true,
				DiffSuppressFunc: func(key, old, new string, d *schema.ResourceData) bool {
					return new == ""
				},
			},
			"free_range": &schema.Schema{
				Type:         schema.TypeList,
				Description:  "Find a free IP address to claim.",
				Optional:     true,
				ExactlyOneOf: []string{"cidr", "from", "free_range"},
				MaxItems:     1,
				// TODO add ForceNew, do we ignore changes?
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// TODO user range_ref here
						"range": &schema.Schema{
							Type:        schema.TypeString,
							Description: "Pick IP address from range with name",
							Required:    true,
						},
						"start_at": &schema.Schema{
							Type:        schema.TypeString,
							Description: "Start searching for IP address from",
							Default:     "",
							Optional:    true,
							// TODO validate that its valide ip in the range of range
						},
						"size": &schema.Schema{
							Type:        schema.TypeInt,
							Description: "The minimum size of the address blocks, specified as the number of addresses",
							Default:     255,
							Optional:    true,
						},

						"ignore_subnet_flag": &schema.Schema{
							Type:        schema.TypeBool,
							Description: "Exclude IP addresses that are assigned via DHCP",
							Default:     false,
							Optional:    true,
						},
						"temporary_claim_time": &schema.Schema{
							Type:         schema.TypeInt,
							Description:  "Time in seconds to temporarily claim IP address, so it isn't claimed by others while the claim is in progess.",
							Default:      60,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 300),
						},
					},
				},
			},
			"from": {
				Type:         schema.TypeString,
				Description:  "The starting IP address of the range.",
				RequiredWith: []string{"to"},
				ExactlyOneOf: []string{"cidr", "from", "free_range"},
				ForceNew:     true,
				Optional:     true,
				ValidateFunc: validation.IsIPAddress,
				DiffSuppressFunc: func(key, old, new string, d *schema.ResourceData) bool {
					// ignore if it is not set in config but know
					return new == ""
				},
			},

			"to": {
				Type:         schema.TypeString,
				Description:  "The ending IP address of the range.",
				RequiredWith: []string{"from"},
				Optional:     true,
				ForceNew:     true,
				// TODO validate ip is higher then from
				ValidateFunc: validation.IsIPAddress,
				DiffSuppressFunc: func(key, old, new string, d *schema.ResourceData) bool {

					// ignore if It is not set in config but know
					return new == ""
				},
			},
			"parent_ref": {
				Type:        schema.TypeString,
				Description: "A reference to the range that contains the subranges",
				Computed:    true,
			},

			"ad_site_ref": {
				Type:        schema.TypeString,
				Description: "Internal reference of the AD site to which the the range belongs.",
				Computed:    true,
			},

			"ad_site_display_name": {
				Type:        schema.TypeString,
				Description: "The display name of the AD site to which the range belongs.",
				Computed:    true,
			},
			// TODO
			// "childRanges": {
			// 	Type:        schema.TypeList,
			// 	Description: "An list of child ranges of the range.",
			// 	Computed:    true, //TODO

			// redundant
			// IsLeaf            bool       `json:"isLeaf"`
			// NumChildren int        `json:"numchildren"`

			// TODO
			// "dhcpScopes": {
			// 	Type:        schema.TypeList,
			// 	Description:
			// 	Computed:    true, //TODO
			// 	// Default:      false,
			// },
			// "authority": {
			// 	Type:        schema.TypeList,
			// 	Description:
			// 	Computed:    true, //TODO
			// },

			"subnet": {
				Type:        schema.TypeBool,
				Description: "Determines if the range is defined as a subnet.",
				Computed:    true,
			},

			"locked": {
				Type:        schema.TypeBool,
				Description: "Determines if the range is defined as a subnet.",
				Default:     false,
				Optional:    true,
			},
			"auto_assign": {
				Type:        schema.TypeBool,
				Description: "Determines if it should be possible to automatically assign IP addresses from the range.",
				// Computed:    true, // TODO
				Default:  true,
				Optional: true,
			},
			"has_schedule": {
				Type:        schema.TypeBool,
				Description: "Determines if a discovery schedule has been set for the range.",
				Computed:    true,
			},
			"has_monitor": {
				Type:        schema.TypeBool,
				Description: "Determines if a discovery schedule has been set for the range.",
				Computed:    true,
			},

			"title": {
				Type:        schema.TypeString,
				Description: "The title of the Range",
				Required:    true,
				// Default:      false,
			},

			"description": {
				Type:        schema.TypeString,
				Description: "Description of the range",
				Optional:    true,
				// Default:      false,
			},

			"custom_properties": &schema.Schema{
				Type:        schema.TypeMap,
				Description: "Map of custom properties associated with this range. You can only assign properties that are already defined in Micetro.",

				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"inherit_access": {
				Type:        schema.TypeBool,
				Description: "If this range should inherit its access bits from its parent range.",
				Computed:    true,
			},
			"is_container": {
				Type:        schema.TypeBool,
				Description: "Set to true to create a container instead of a range.",
				Computed:    true, // TODO
			},

			"utilization_percentage": {
				Type:        schema.TypeInt,
				Description: "Utilization percentage for range.",
				Computed:    true,
			},

			"has_rogue_addresses": {
				Type:        schema.TypeBool,
				Description: "Set to true to create a container instead of a range.",
				Computed:    true,
			},

			"cloud_network_ref": {
				Type:        schema.TypeString,
				Description: "A internal reference to its cloud network",
				Computed:    true,
			},

			// TODO
			// "cloudAllocationPools": {
			// Type:        schema.TypeList,
			// Optional:    true,
			// Elem: &schema.Resource{
			// 	Schema: map[string]*schema.Schema{
			// },

			// TODO
			// "discoveredProperties": {
			// Type:        schema.TypeList,
			// Optional:    true,
			// Elem: &schema.Resource{
			// 	Schema: map[string]*schema.Schema{
			// },

			"created": &schema.Schema{
				Type:        schema.TypeString,
				Description: "DDate when zone was created in Micetro.",
				Computed:    true,
			},
			"lastmodified": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Date when zone was last modified in Micetro.",
				Computed:    true,
			},
			"discovery": &schema.Schema{
				Type:        schema.TypeList,
				Description: "Used for discovery of ranges or scopes.",
				Optional:    true,
				MaxItems:    1,
				ForceNew:    true, //TODO can we make this update
				// default does not work for list
				// Default:     [1]map[string]interface{}{{"enabled": false}},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"interval": &schema.Schema{
							Type:        schema.TypeInt,
							Description: "The interval between runs for the schedule.",
							Optional:    true,
							// TODO Default
						},
						"unit": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Unit of time for interval. One of: Minutes, Hours, Days, Weeks, Months",
							ValidateFunc: validation.StringInSlice([]string{
								"Minutes", "Hours", "Days", "Weeks", "Months",
							}, false),
						},
						"enabled": &schema.Schema{
							Type:        schema.TypeBool,
							Description: "Pick IP address from range with name",
							Optional:    true,
							Default:     false,
						},
						// TODO "start_time" : &schema.Schema{
						// },
					},
				},
			},
		},
	}
}

func writeRangeSchema(d *schema.ResourceData, iprange Range) {

	d.Set("ref", iprange.Ref)
	d.Set("name", iprange.Name)

	if _, _, err := net.ParseCIDR(iprange.Name); err == nil {
		d.Set("cidr", iprange.Name)
	} else {
		d.Set("cidr", "")
	}

	d.Set("from", iprange.From)
	d.Set("to", iprange.To)
	d.Set("parent_ref", iprange.ParentRef)
	d.Set("ad_site_ref", iprange.AdSiteRef)
	d.Set("ad_site_display_name", iprange.AdSiteDisplayName)

	// TODO childRanges, dhcpScopes,authority
	d.Set("subnet", iprange.Subnet)
	d.Set("locked", iprange.Locked)
	d.Set("auto_assign", iprange.AutoAssign)
	d.Set("has_schedule", iprange.HasSchedule)
	d.Set("has_monitor", iprange.HasMonitor)

	// api exposes this as custom properties
	d.Set("title", iprange.RangeProperties.CustomProperties["Title"])
	d.Set("description", iprange.RangeProperties.CustomProperties["Description"])

	delete(iprange.CustomProperties, "Title")
	delete(iprange.CustomProperties, "Description")
	d.Set("custom_properties", iprange.CustomProperties)

	d.Set("is_container", iprange.IsContainer)
	d.Set("utilization_percentage", iprange.UtilizationPercentage)
	d.Set("has_rogue_addresses", iprange.HasRogueAddresses)
	d.Set("cloud_network_ref", iprange.CloudNetworkRef)

	d.Set("created", iprange.Created)           // TODO convert to timeformat RFC 3339
	d.Set("lastmodified", iprange.LastModified) // TODO convert to timeformat RFC 3339

	// TODO	 discovery, discoveredProperties,cloudAllocationPools
	return
}

func readAvailableAddressBlocksRequest(freeRange interface{}) AvailableAddressBlocksRequest {

	freeRangeInterface := freeRange.([]interface{})[0].(map[string]interface{})
	availableAddressBlocksRequest := AvailableAddressBlocksRequest{
		RangeRef:           freeRangeInterface["range"].(string),
		StartAddress:       freeRangeInterface["start_at"].(string),
		Size:               freeRangeInterface["size"].(int),
		Limit:              1,
		IgnoreSubnetFlag:   freeRangeInterface["ignore_subnet_flag"].(bool),
		TemporaryClaimTime: freeRangeInterface["temporary_claim_time"].(int),
	}

	return availableAddressBlocksRequest

}
func readDiscoverySchema(discovery_schemas interface{}) Discovery {

	schemas := discovery_schemas.([]interface{})
	discoveryMap := schemas[0].(map[string]interface{})
	discovery := Discovery{

		Interval: discoveryMap["interval"].(int),
		Unit:     discoveryMap["unit"].(string),
		Enabled:  discoveryMap["enabled"].(bool),
		// TODO "StartTime" :
	}
	return discovery
}

func readRangeSchema(d *schema.ResourceData) Range {

	var name, from, to string
	if cidr, ok := d.GetOk("cidr"); ok {
		name = cidr.(string)
		to = ""
		from = ""
	} else {
		name = tryGetString(d, "from") + "-" + tryGetString(d, "to")
		from = tryGetString(d, "from")
		to = tryGetString(d, "to")
	}
	var customProperties = make(map[string]string)
	if customPropertiesRead, ok := d.GetOk("custom_properties"); ok {
		for key, value := range customPropertiesRead.(map[string]interface{}) {
			customProperties[key] = value.(string)
		}
	}

	if description, ok := d.GetOk("description"); ok {
		customProperties["description"] = description.(string)
	}

	if title, ok := d.GetOk("title"); ok {
		customProperties["title"] = title.(string)
	}

	iprange := Range{
		Ref:  tryGetString(d, "ref"),
		Name: name,
		From: from,
		To:   to,

		// you should not set this yourself
		// Created:      d.Get("created").(string),
		// LastModified: tryGetString(d, "lastmodified"),

		InheritAccess: d.Get("inherit_access").(bool),
		RangeProperties: RangeProperties{

			Locked:           d.Get("locked").(bool),
			AutoAssign:       d.Get("auto_assign").(bool),
			CustomProperties: customProperties,
		},
	}
	return iprange
}

func resourceRangeCreate(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Mmclient)

	if freeRangeMap, ok := d.GetOk("free_range"); ok {
		availableAddressBlocksRequest := readAvailableAddressBlocksRequest(freeRangeMap)

		tflog.Debug(c, fmt.Sprintf("Request a available AddressBlock"))
		AvailableAddressBlocks, err := client.AvailableAddressBlocks(availableAddressBlocksRequest)

		if err != nil {
			return diag.FromErr(err)
		}
		if len(AvailableAddressBlocks) <= 0 {
			// TODO better messages
			return diag.Errorf("No available address blocks found")
		}
		d.Set("from", AvailableAddressBlocks[0].From)
		d.Set("to", AvailableAddressBlocks[0].To)

	}

	discovery_schemas, ok := d.GetOk("discovery")
	var discovery Discovery
	if ok {
		discovery = readDiscoverySchema(discovery_schemas)
	} else {
		// default is defined here. because can't be done in schema
		discovery = Discovery{Enabled: false}
	}

	iprange := readRangeSchema(d)
	objRef, err := client.CreateRange(iprange, discovery)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(objRef)
	return resourceRangeRead(c, d, m)

}

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

func resourceRangeUpdate(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	//TODO maybe check more
	if d.HasChange("discovery") {
		return diag.Errorf("can't change read-only property of DNS zone")
	}
	client := m.(*Mmclient)
	ref := d.Id()

	iprange := readRangeSchema(d)

	err := client.UpdateRange(iprange.RangeProperties, ref)

	if err != nil {
		return diag.FromErr(err)
	}
	return resourceRangeRead(c, d, m)
}

func resourceRangeDelete(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*Mmclient)
	var diags diag.Diagnostics
	ref := d.Id()
	err := client.DeleteRange(ref)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}

func resourceRangeImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	diags := resourceRangeRead(ctx, d, m)
	if err := toError(diags); err != nil {
		return nil, err
	}

	// if we had used schema.ImportStatePassthrough
	// we could not have set id to its canonical form

	ref := d.Get("ref").(string)
	if ref == "" {
		return nil, errors.New("Import failed")
	}
	d.SetId(ref)

	return []*schema.ResourceData{d}, nil
}
