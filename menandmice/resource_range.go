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

			// TODO add atributes: cloudAllocationPools, dhcpScopes authority ,discoveredProperties
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
				DiffSuppressFunc: func(key, old, new string, d *schema.ResourceData) bool {
					return new == ""
				},
				ForceNew: true,
				Optional: true,
			},
			"free_range": {
				Type:         schema.TypeList,
				Description:  "Find a free IP address to claim.",
				Optional:     true,
				ExactlyOneOf: []string{"cidr", "from", "free_range"},
				MaxItems:     1,
				ForceNew:     true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"range": {
							Type:         schema.TypeString,
							Description:  "Pick IP address from range with name",
							ExactlyOneOf: []string{"free_range.0.range", "free_range.0.ranges"},
							Optional:     true,
						},
						"ranges": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description:  "Pick IP address from one of these range with name",
							ExactlyOneOf: []string{"free_range.0.range", "free_range.0.ranges"},
							Optional:     true,
						},
						"start_at": {
							Type:          schema.TypeString,
							Description:   "Start searching for IP address from",
							ConflictsWith: []string{"free_range.0.mask"},
							Default:       "",
							Optional:      true,
							ForceNew:      true,
							// TODO validate that its valide ip in the range of range
						},
						"size": {
							Type:          schema.TypeInt,
							ExactlyOneOf:  []string{"free_range.0.mask"},
							Description:   "The minimum size of the address blocks, specified as the number of addresses",
							ConflictsWith: []string{"subnet"},
							Optional:      true,
							ForceNew:      true,
						},

						"mask": {
							Type:         schema.TypeInt,
							ExactlyOneOf: []string{"free_range.0.mask", "free_range.0.size"},
							Description:  "The minimum size of the address blocks, specified as a subnet mask.",
							// Default:     24, // setting default here gives problem if user also set size
							Optional: true,
							ForceNew: true,
						},

						"ignore_subnet_flag": {
							Type:        schema.TypeBool,
							Description: "Exclude IP addresses that are assigned via DHCP",
							Default:     false,
							Optional:    true,
						},
						"temporary_claim_time": {
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

			"child_ranges": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An list of child ranges of the range.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ref": {
							Type:        schema.TypeString,
							Description: "Internal references to this child range.",
							Computed:    true,
						},

						"name": {
							Type:        schema.TypeString,
							Description: "Name to this child range.",
							Computed:    true,
						},
					},
				},
			},
			// "dhcpScopes": {
			// 	Type:        schema.TypeList,
			// 	Description:
			// 	// Default:      false,
			// },
			// "authority": {
			// 	Type:        schema.TypeList,
			// 	Description:
			// },

			"subnet": {
				Type:          schema.TypeBool,
				Description:   "Determines if the range is defined as a subnet.",
				Default:       false,
				Optional:      true,
				ConflictsWith: []string{"from", "to"},
			},

			"locked": {
				Type:        schema.TypeBool,
				Description: "Determines if the range is locked.",
				Default:     false,
				Optional:    true,
			},
			"auto_assign": {
				Type:        schema.TypeBool,
				Description: "Determines if it should be possible to automatically assign IP addresses from the range.",
				Default:     false,
				Optional:    true,
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
			},

			"description": {
				Type:        schema.TypeString,
				Description: "Description of the range",
				Optional:    true,
			},

			// TODO make custom_properties case insensitive
			// TODO defaults custom_properties are set during creation by server. but terraform does not know about those
			//		so we could fix this behavior:
			//	    * we have to change that behavior, extra api call at create
			//		* fetch defaults maybe at DefaultFunc or when doing update
			//		* ignore changes from default
			"custom_properties": {
				Type:        schema.TypeMap,
				Description: "Map of custom properties associated with this range. You can only assign properties that are already defined in Micetro.",

				// DiffSuppressFunc: func(_key, _old, _new string, d *schema.ResourceData) bool {
				// 	// this is need because DiffSuppressFunc does not support maps
				// 	// https://github.com/hashicorp/terraform-plugin-sdk/issues/477
				// 	// this solution is based on @diabloneo  from that issue
				//	//  This does not work for configerdInterface/new
				// 	oldInterface, configerdInterface := d.GetChange("custom_properties")
				// 	old := oldInterface.(map[string]interface{})
				// 	configerd := configerdInterface.(map[string]interface{})
				// 	suppressDiff := false
				//
				// 	fmt.Printf("####### pre %v:%v\n\n", oldInterface, configerdInterface)
				// 	for key, valOld := range old {
				//
				// 		if valNew, ok := configerd[key]; ok {
				//
				// 			fmt.Printf("####### %v:%v\n\n", valNew, valOld)
				// 			if valNew != valOld {
				// 				// panic(fmt.Sprintf("%v:%v", valNew, valOld))
				//
				// 				return false
				// 			}
				// 		} else {
				// 			suppressDiff = true
				// 			_ = suppressDiff
				// 			// var client Mmclient
				// 			// d.GetProviderMeta(client)
				// 		}
				// 	}
				// 	return true //suppressDiff
				// },

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
				Optional:    true,
				Default:     false,
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

			// "cloudAllocationPools": {
			// Type:        schema.TypeList,
			// Optional:    true,
			// Elem: &schema.Resource{
			// 	Schema: map[string]*schema.Schema{
			// },

			// "discoveredProperties": {
			// Type:        schema.TypeList,
			// Optional:    true,
			// Elem: &schema.Resource{
			// 	Schema: map[string]*schema.Schema{
			// },

			"created": {
				Type:        schema.TypeString,
				Description: "DDate when zone was created in Micetro.",
				Computed:    true,
			},
			"lastmodified": {
				Type:        schema.TypeString,
				Description: "Date when zone was last modified in Micetro.",
				Computed:    true,
			},
			// "discovery": &schema.Schema{
			// 	Type:        schema.TypeList,
			// 	Description: "Used for discovery of ranges or scopes.",
			// 	Computed:    true,
			// 	// Optional:    true,
			// 	// ForceNew:    true,
			// 	// MaxItems: 1,
			// 	// default does not work for list
			// 	// Default:     [1]map[string]interface{}{{"enabled": false}},
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"interval": &schema.Schema{
			// 				Type:        schema.TypeInt,
			// 				Description: "The interval between runs for the schedule.",
			// 				Optional:    true,
			// 			},
			// 			"unit": &schema.Schema{
			// 				Type:        schema.TypeString,
			// 				Optional:    true,
			// 				Description: "Unit of time for interval. One of: Minutes, Hours, Days, Weeks, Months",
			// 				ValidateFunc: validation.StringInSlice([]string{
			// 					"Minutes", "Hours", "Days", "Weeks", "Months",
			// 				}, false),
			// 			},
			// 			"enabled": &schema.Schema{
			// 				Type:        schema.TypeBool,
			// 				Description: "Pick IP address from range with name",
			// 				Optional:    true,
			// 				Default:     false,
			// 			},
			// 			// "start_time" : &schema.Schema{
			// 			// },
			// 		},
			// 	},
			// },
		},
	}
}

func writeRangeSchema(d *schema.ResourceData, iprange Range) {
	// We schema indirecly via creating a Map here.
	// So we can share code with data_source_ranges.
	// data_source_ranges has to set same fields but for every range it finds
	SetFromMap(d, flattenRange(iprange))
}

func flattenRange(iprange Range) map[string]interface{} {
	var m = map[string]interface{}{}
	m["ref"] = iprange.Ref
	m["name"] = iprange.Name

	if _, _, err := net.ParseCIDR(iprange.Name); err == nil {
		m["cidr"] = iprange.Name
	} else {
		m["cidr"] = ""
	}

	m["from"] = iprange.From
	m["to"] = iprange.To
	m["parent_ref"] = iprange.ParentRef
	m["ad_site_ref"] = iprange.AdSiteRef
	m["ad_site_display_name"] = iprange.AdSiteDisplayName

	m["subnet"] = iprange.Subnet
	m["locked"] = iprange.Locked
	m["auto_assign"] = iprange.AutoAssign
	m["has_schedule"] = iprange.HasSchedule
	m["has_monitor"] = iprange.HasMonitor

	// api exposes this as custom properties
	m["title"] = iprange.RangeProperties.CustomProperties["Title"]
	m["description"] = iprange.RangeProperties.CustomProperties["Description"]

	delete(iprange.CustomProperties, "Title")
	delete(iprange.CustomProperties, "Description")
	m["custom_properties"] = iprange.CustomProperties

	m["is_container"] = iprange.IsContainer
	m["utilization_percentage"] = iprange.UtilizationPercentage
	m["has_rogue_addresses"] = iprange.HasRogueAddresses
	m["cloud_network_ref"] = iprange.CloudNetworkRef

	m["created"] = iprange.Created           // TODO convert to timeformat RFC 3339
	m["lastmodified"] = iprange.LastModified // TODO convert to timeformat RFC 3339

	var namedRefs = make([]map[string]interface{}, len(iprange.ChildRanges))
	for i, namedRef := range iprange.ChildRanges {
		namedRefs[i] = map[string]interface{}{
			"ref":  namedRef.Ref,
			"name": namedRef.Name,
		}
	}
	m["child_ranges"] = namedRefs
	// TODO	 discovery, discoveredProperties,cloudAllocationPools
	return m
}

func readAvailableAddressBlocksRequest(freeRange map[string]interface{}) AvailableAddressBlocksRequest {

	availableAddressBlocksRequest := AvailableAddressBlocksRequest{
		RangeRef:           freeRange["range"].(string),
		StartAddress:       freeRange["start_at"].(string),
		Size:               freeRange["size"].(int),
		Mask:               freeRange["mask"].(int),
		Limit:              1,
		IgnoreSubnetFlag:   freeRange["ignore_subnet_flag"].(bool),
		TemporaryClaimTime: freeRange["temporary_claim_time"].(int),
	}

	return availableAddressBlocksRequest

}
func readDiscoverySchema(discoverySchemas interface{}) Discovery {

	schemas := discoverySchemas.([]interface{})
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

		InheritAccess: d.Get("inherit_access").(bool),
		RangeProperties: RangeProperties{

			Locked:           d.Get("locked").(bool),
			AutoAssign:       d.Get("auto_assign").(bool),
			Subnet:           d.Get("subnet").(bool),
			IsContainer:      d.Get("is_container").(bool),
			CustomProperties: customProperties,
		},
	}
	return iprange
}

func resourceRangeCreate(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Mmclient)
	var err error
	var AvailableAddressBlocks []AddressBlock

	if freeRange, ok := d.GetOk("free_range"); ok {
		freeRangeMap := freeRange.([]interface{})[0].(map[string]interface{})
		var ranges []interface{}
		// convert to list if range, otherwise pick from ranges
		if rangeName, ok := freeRangeMap["range"]; ok {
			ranges = []interface{}{rangeName}
		} else {
			ranges = []interface{}{freeRangeMap["ranges"]}
		}

		// For readAvailableAddressBlocksRequest "range" has to be set.
		// Which is not the case if ranges was used
		freeRangeMap["range"] = ranges[0]
		availableAddressBlocksRequest := readAvailableAddressBlocksRequest(freeRangeMap)
		for _, iprange := range ranges {
			rangeName := iprange.(string)

			tflog.Debug(c, "Request a available AddressBlock")
			availableAddressBlocksRequest.RangeRef = rangeName
			AvailableAddressBlocks, err = client.AvailableAddressBlocks(availableAddressBlocksRequest)
			if err != nil {
				return diag.FromErr(err)
				// TODO do we want to make errors allowed as long there is one succesfull AvailableAddressBlocksRequest
			}

			if len(AvailableAddressBlocks) >= 1 {
				break
			}
			tflog.Info(c, fmt.Sprintf("No suitable unclaimed range found in %v", rangeName))
		}

		if len(AvailableAddressBlocks) <= 0 {
			// TODO better messages
			return diag.Errorf("No available address blocks found")
		}

		d.Set("from", AvailableAddressBlocks[0].From)
		d.Set("to", AvailableAddressBlocks[0].To)
	}

	// TODO discovery
	// discovery_schemas, ok := d.GetOk("discovery")
	// var discovery Discovery
	// if ok {
	// 	discovery = readDiscoverySchema(discovery_schemas)
	// } else {
	// 	// default is defined here. because can't be done in schema
	discovery := Discovery{Enabled: false}
	// }

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
	// TODO discovery?

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
