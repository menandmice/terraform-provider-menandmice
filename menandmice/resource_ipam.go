package menandmice

import (
	"bytes"
	"context"
	"net"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceIPAMRec() *schema.Resource {
	return &schema.Resource{
		Description: "`menandmice_ipam_record` IP address managment",

		CreateContext: resourceIPAMRecCreate,
		ReadContext:   resourceIPAMRecRead,
		UpdateContext: resourceIPAMRecUpdate,
		DeleteContext: resourceIPAMRecDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceIPAMRecImport,
		},
		Schema: map[string]*schema.Schema{

			"ref": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Internal reference to ipam record",
				Computed:    true,
			},
			"free_ip": &schema.Schema{
				Type:         schema.TypeList,
				Description:  "Find a free IP address to claim",
				Optional:     true,
				ExactlyOneOf: []string{"free_ip", "address"},
				MaxItems:     1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"range": &schema.Schema{
							Type:        schema.TypeString,
							Description: "pick IP address from range with name",
							Required:    true,
						},
						"start_at": &schema.Schema{
							Type:        schema.TypeString,
							Description: "Start searching for IP from",
							Default:     "",
							Optional:    true,
							// TODO validate that its valide ip in the range of range
						},
						"ping": &schema.Schema{
							Type:        schema.TypeBool,
							Description: "Verify ip is free with Ping",
							Default:     false,
							Optional:    true,
						},
						"exclude_dhcp": &schema.Schema{
							Type:        schema.TypeBool,
							Description: "Exclude IP address that are Assigned via DHCP",
							Default:     false,
							Optional:    true,
						},

						"temporary_claim_time": &schema.Schema{
							Type:         schema.TypeInt,
							Description:  "Time in seconds to temporary claim IP address. So it won't be claimed by others, when the claim is in progess",
							Default:      60,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 300),
						},
					},
				},
			},
			"address": &schema.Schema{
				Type:         schema.TypeString,
				Description:  "The IP address to claim",
				ExactlyOneOf: []string{"free_ip", "address"},
				Optional:     true,
				ValidateFunc: validation.Any(
					validation.IsIPv4Address,
					validation.IsIPv6Address),
				DiffSuppressFunc: func(key, old, new string, d *schema.ResourceData) bool {
					if ipv6AddressDiffSuppress(key, old, new, d) {
						return true
					}
					if freeIPRead, ok := d.GetOk("free_ip"); ok {

						return inFreeIPRange(readFreeIPMap(freeIPRead), old)
					}
					return false
				},
				ForceNew: true,
			},
			// TODO might not be a good idea to make this configerable.
			// What does it mean to delete unclaimed iprecord
			"claimed": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "If address should be claimed. Default: true",
				Optional:    true,
				Default:     true,
			},
			// "dnshost": &schema.Schema{
			// },
			// "dhcp_reservations": &schema.Schema{
			// },
			// "dhcp_leases": &schema.Schema{
			// },
			"discovery_type": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Way IP address use is dicoverd. For example: None, Ping, ARP, Lease, Custom.",
				Computed:    true,
				// Optional: true,
				// Default:  "None",
				// ForceNew: true,
				// ValidateFunc: validation.StringInSlice([]string{
				// 	"None", "Ping", "ARP", "Lease", "Custom",
				// }, false),
			},
			"last_seen_date": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The date when the address was last seen during IP address discovery.",
				Computed:    true,
			},

			"last_discovery_date": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The date when the system last performed IP address discovery for this IP address.",
				Computed:    true,
			},
			"last_known_client_identifier": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The MAC address associated with the IP address discovery info.",
				Computed:    true,
			},

			"device": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The device associated with the record.",
				Computed:    true,
			},

			"interface": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The interface associated with the record.",
				Computed:    true,
			},
			"ptr_status": &schema.Schema{
				Type:        schema.TypeString,
				Description: "PTR record status. For example: Unknown, OK, Verify.",
				Computed:    true,
			},
			"extraneous_ptr": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "Contains true if there are extraneous PTR records for the record.",
				Computed:    true,
			},
			"custom_properties": &schema.Schema{
				Type:        schema.TypeMap,
				Description: "Map of custom properties associated with this IP address. You can only assign properties that are already defined via propertie devinition",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"state": &schema.Schema{
				Type:        schema.TypeString,
				Description: "state of IP addres. For exampe: Free, Assigned, Claimed, Pending, Held.",
				Computed:    true,
			},
			// Hold info is read only property. but whould be inconvient if it was part of resource definition
			// because if free_ip is used it's state whould change after creation the moment temporaryClaimTime would exipre
			//
			// "hold_info": &schema.Schema{
			// 	Type:     schema.TypeList,
			// Description: "Contains information about who holds the otherwise free IP and for how long.",
			// 	Computed: true,
			// 	MaxItems: 1,
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"expiry_time": &schema.Schema{
			// 				Type:     schema.TypeString,
			// 				Computed: true,
			// 				// ValidateFunc: validation.ValidateRFC3339TimeString,
			// 			},
			//
			// 			"username": &schema.Schema{
			// 				Type:     schema.TypeString,
			// 				Computed: true,
			// 			},
			// 		},
			// 	},
			// },

			"usage": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "IP address usage bitmask.",
				Computed:    true,
			},

			// "cloud_device_info": &schema.Schema{
			// },
		},
	}
}

// TODO write test for this
func inFreeIPRange(freeIPMap map[string]interface{}, ipStr string) bool {
	_, ipnet, err := net.ParseCIDR(freeIPMap["range"].(string))
	if err != nil {
		return false
	}
	ip := net.ParseIP(ipStr)
	if ip == nil || ipnet.Contains(ip) == false {
		return false
	}
	if minimumIPStr := freeIPMap["start_at"].(string); minimumIPStr != "" {

		minimumIP := net.ParseIP(minimumIPStr)
		if minimumIP == nil {
			return false
		}

		return bytes.Compare(ip, minimumIP) >= 0

	}

	return true

}
func writeIPAMRecSchema(d *schema.ResourceData, ipamrec IPAMRecord) {

	d.Set("ref", ipamrec.Ref)
	d.Set("address", ipamrec.Address)
	d.Set("claimed", ipamrec.Claimed)

	// d.Set("dnshost", ipamrec.DNSHost)
	// d.Set("dhcp_reservations", ipamrec.DHCPReservations)
	// d.Set("dhcp_leases", ipamrec.DHCPLeases)

	d.Set("discovery_type", ipamrec.DiscoveryType)

	d.Set("last_seen_date", ipamrec.LastSeenDate)           // TODO convert to timeformat RFC 3339
	d.Set("last_discovery_date", ipamrec.LastDiscoveryDate) // TODO convert to timeformat RFC 3339
	d.Set("last_known_client_identifier", ipamrec.LastKnownClientIdentifier)
	d.Set("device", ipamrec.Device)
	d.Set("interface", ipamrec.Interface)
	d.Set("ptr_status", ipamrec.PTRStatus)
	d.Set("extraneous_ptr", ipamrec.ExtraneousPTR)
	d.Set("custom_properties", ipamrec.CustomProperties)

	d.Set("state", ipamrec.State)

	// if ipamrec.HoldInfo != nil {
	// 	holdInfo := make(map[string]interface{})
	//
	// 	holdInfo["expiry_time"] = ipamrec.HoldInfo.ExpiryTime
	// 	holdInfo["username"] = ipamrec.HoldInfo.Username
	//
	// 	d.Set("hold_info", [](map[string]interface{}){holdInfo})
	// }
	d.Set("usage", ipamrec.Usage)
	// d.Set("cloud_device_info", ipamrec.CloudDeviceInfo)
	return
}

func readIPAMRecSchema(d *schema.ResourceData) IPAMRecord {

	var holdInfo HoldInfo
	if holdInfoList, ok := d.GetOk("hold_info"); ok {
		holdInfoMap := holdInfoList.([]interface{})[0]
		holdInfoRead := holdInfoMap.(map[string]interface{})
		holdInfo = HoldInfo{
			Username:   holdInfoRead["username"].(string),
			ExpiryTime: holdInfoRead["expiry_time"].(string),
		}
	}
	var CustomProperties = make(map[string]string)
	if customPropertiesRead, ok := d.GetOk("custom_properties"); ok {
		for key, value := range customPropertiesRead.(map[string]interface{}) {
			CustomProperties[key] = value.(string)
		}
	}

	ipamrec := IPAMRecord{
		Ref:           tryGetString(d, "ref"),
		Address:       d.Get("address").(string),
		DiscoveryType: tryGetString(d, "discovery_type"),
		// Last_seen_date			// read only
		// last_discovery_date
		// last_known_client_identifier
		PTRStatus:     tryGetString(d, "ptr_status"),
		ExtraneousPTR: d.Get("extraneous_ptr").(bool),
		Device:        tryGetString(d, "device"),
		State:         tryGetString(d, "state"),
		HoldInfo:      &holdInfo,

		// Usage read only
		IPAMProperties: IPAMProperties{
			Claimed: d.Get("claimed").(bool),
			// DNShost,
			// DHCPReservations
			// DHCPLeases
			Interface:        tryGetString(d, "interface"),
			CustomProperties: CustomProperties,
			// CloudDeviceInfo // not inpleneted
		},
	}
	return ipamrec
}

func resourceIPAMRecRead(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	client := m.(*Mmclient)

	ipamrec, err := client.ReadIPAMRec(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	writeIPAMRecSchema(d, ipamrec)

	return diags
}

func readFreeIPMap(freeIPRead interface{}) map[string]interface{} {

	freeIPReadList := freeIPRead.([]interface{})[0]
	return freeIPReadList.(map[string]interface{})
}
func resourceIPAMRecCreate(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := m.(*Mmclient)

	if _, ok := d.GetOk("address"); !ok {
		freeIPRead, ok := d.GetOk("free_ip")
		if !ok {
			return diag.Errorf("could not read address or free_ip")
		}

		freeIPMap := readFreeIPMap(freeIPRead)
		ipRange := freeIPMap["range"].(string)
		startIP := freeIPMap["start_at"].(string)
		ping := freeIPMap["ping"].(bool)
		excludeDHCP := freeIPMap["exclude_dhcp"].(bool)
		temporaryClaimTime := freeIPMap["temporary_claim_time"].(int)

		address, err := client.NextFreeAddress(ipRange, startIP, ping, excludeDHCP, temporaryClaimTime)

		if err != nil {
			return diag.Errorf("could not read Not find a free ipaddress in %s", ipRange)
		}

		d.Set("address", address)
	}
	ipamrec := readIPAMRecSchema(d)

	err := client.CreateIPAMRec(ipamrec)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(ipamrec.Address)
	diags = resourceIPAMRecRead(c, d, m)

	if diags != nil {
		return diags
	}
	d.SetId(d.Get("ref").(string))
	return diags

}

func resourceIPAMRecUpdate(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*Mmclient)
	ref := d.Id()
	ipamrec := readIPAMRecSchema(d)

	err := client.UpdateIPAMRec(ipamrec.IPAMProperties, ref)

	if err != nil {
		return diag.FromErr(err)
	}
	return resourceIPAMRecRead(c, d, m)
}

func resourceIPAMRecDelete(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*Mmclient)
	var diags diag.Diagnostics
	ref := d.Id()
	err := client.DeleteIPAMRec(ref)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
func resourceIPAMRecImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	diags := resourceIPAMRecRead(ctx, d, m)

	if err := toError(diags); err != nil {
		return nil, err
	}
	// if we had used schema.ImportStatePassthrough
	// we could not have set id to its cannical form
	d.SetId(d.Get("ref").(string))

	return []*schema.ResourceData{d}, nil
}
