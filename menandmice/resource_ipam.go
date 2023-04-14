package menandmice

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
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

			"ref": {
				Type:        schema.TypeString,
				Description: "Internal reference for the IP address.",
				Computed:    true,
			},
			"free_ip": {

				// TODO add ForceNew, do we ignore changes?
				Type:         schema.TypeList,
				Description:  "Find a free IP address to claim.",
				Optional:     true,
				ExactlyOneOf: []string{"free_ip", "address"},
				MaxItems:     1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// TODO user range_ref here
						"range": {
							Type:        schema.TypeString,
							Description: "Pick IP address from range with name",
							Required:    true,
						},
						"start_at": {
							Type:        schema.TypeString,
							Description: "Start searching for IP address from",
							Default:     "",
							Optional:    true,
							// TODO validate that its valide ip in the range of range
						},
						"ping": {
							Type:        schema.TypeBool,
							Description: "Verify IP address is free with ping",
							Default:     false,
							Optional:    true,
						},
						"exclude_dhcp": {
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

			"address": {
				Type:         schema.TypeString,
				Description:  "The IP address to claim.",
				ExactlyOneOf: []string{"free_ip", "address"},
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.IsIPAddress,
				DiffSuppressFunc: func(key, old, new string, d *schema.ResourceData) bool {
					if ipv6AddressDiffSuppress(key, old, new, d) {
						return true
					}
					if _, ok := d.GetOk("free_ip"); ok {
						return true
					}
					return false
				},
				ForceNew: true,
			},
			"current_address": {
				Type:        schema.TypeString,
				Description: "Address currently used.",
				Deprecated:  "user address instead",
				Computed:    true,
			},
			// TODO might not be a good idea to make this configerable.
			// What does it mean to delete unclaimed iprecord
			"claimed": {
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
			"discovery_type": {
				Type:        schema.TypeString,
				Description: "The discovery method of the IP address. Example: None, Ping, ARP, Lease, Custom.",
				Computed:    true,
				// Optional: true,
				// Default:  "None",
				// ForceNew: true,
				// ValidateFunc: validation.StringInSlice([]string{
				// 	"None", "Ping", "ARP", "Lease", "Custom",
				// }, false),
			},
			"last_seen_date": {
				Type:        schema.TypeString,
				Description: "The date when the address was last seen during IP address discovery in rfc3339 time format.",
				Computed:    true,
			},

			"last_discovery_date": {
				Type:        schema.TypeString,
				Description: "The date when the system last performed IP address discovery for this IP address rfc3339 time format.",
				Computed:    true,
			},
			"last_known_client_identifier": {
				Type:        schema.TypeString,
				Description: "The last known MAC address associated with the IP address discovery information.",
				Computed:    true,
			},

			"device": {
				Type:        schema.TypeString,
				Description: "The device associated with the object.",
				Computed:    true,
			},

			"interface": {
				Type:        schema.TypeString,
				Description: "The interface associated with the object.",
				Computed:    true,
			},
			"ptr_status": {
				Type:        schema.TypeString,
				Description: "PTR record status. Example: Unknown, OK, Verify.",
				Computed:    true,
			},
			"extraneous_ptr": {
				Type:        schema.TypeBool,
				Description: "'True' if there are extraneous PTR records for the object.",
				Computed:    true,
			},
			// TODO make custom_properties case insensitive
			"custom_properties": {
				Type:        schema.TypeMap,
				Description: "Map of custom properties associated with this IP address. You can only assign properties that are already defined in Micetro.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"state": {
				Type:        schema.TypeString,
				Description: "The state of the IP address. Example: Free, Assigned, Claimed, Pending, Held.",
				Computed:    true,
			},
			// Hold info is read only property. But would be inconvenient if it was part of resource definition
			// because if free_ip is used its state whould change after creation the moment temporaryClaimTime would expire
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

			"usage": {
				Type:        schema.TypeInt,
				Description: "IP address usage bitmask.",
				Computed:    true,
			},

			// "cloud_device_info": &schema.Schema{
			// },
		},
	}
}

func writeIPAMRecSchema(d *schema.ResourceData, ipamrec IPAMRecord, tz *time.Location) diag.Diagnostics {
	var diags diag.Diagnostics
	d.Set("ref", ipamrec.Ref)
	d.Set("address", ipamrec.Address)
	d.Set("current_address", ipamrec.Address)
	d.Set("claimed", ipamrec.Claimed)

	// d.Set("dnshost", ipamrec.DNSHost)
	// d.Set("dhcp_reservations", ipamrec.DHCPReservations)
	// d.Set("dhcp_leases", ipamrec.DHCPLeases)

	d.Set("discovery_type", ipamrec.DiscoveryType)

	lastSeenDate, err := MmTimeString2rfc(ipamrec.LastSeenDate, tz)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("last_seen_date", lastSeenDate)

	lastDiscoveryDate, err := MmTimeString2rfc(ipamrec.LastDiscoveryDate, tz)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("last_discovery_date", lastDiscoveryDate)

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
	return diags
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
	diags = writeIPAMRecSchema(d, ipamrec, client.serverLocation)

	return diags
}

func readNextFreeIPRequest(freeIPRead interface{}) NextFreeAddressRequest {

	freeIPReadIntface := freeIPRead.([]interface{})[0].(map[string]interface{})
	nextFreeIPRequest := NextFreeAddressRequest{

		RangeRef:           freeIPReadIntface["range"].(string),
		StartAddress:       freeIPReadIntface["start_at"].(string),
		Ping:               freeIPReadIntface["ping"].(bool),
		ExcludeDHCP:        freeIPReadIntface["exclude_dhcp"].(bool),
		TemporaryClaimTime: freeIPReadIntface["temporary_claim_time"].(int),
	}
	return nextFreeIPRequest
}
func resourceIPAMRecCreate(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := m.(*Mmclient)

	if freeIPMap, ok := d.GetOk("free_ip"); ok {

		nextFreeIPRequest := readNextFreeIPRequest(freeIPMap)

		tflog.Debug(c, "Request next fee address")
		address, err := client.NextFreeAddress(nextFreeIPRequest)

		if err != nil {
			return diag.Errorf("No free IPs available in %s", nextFreeIPRequest.RangeRef)
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
	// we could not have set id to its canonical form
	d.SetId(d.Get("ref").(string))

	return []*schema.ResourceData{d}, nil
}
