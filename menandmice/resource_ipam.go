package menandmice

import (
	"terraform-provider-menandmice/diag"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceIPAMRec() *schema.Resource {
	return &schema.Resource{
		Create: resourceIPAMRecCreate,
		Read:   resourceIPAMRecRead,
		Update: resourceIPAMRecUpdate,
		Delete: resourceIPAMRecDelete,
		Schema: map[string]*schema.Schema{

			"ref": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.Any(
					validation.IsIPv4Address,
					validation.IsIPv6Address),
				DiffSuppressFunc: ipv6AddressDiffSuppress,
				ForceNew:         true,
			},
			"claimed": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			// "dnshost": &schema.Schema{
			// },
			// "dhcp_reservations": &schema.Schema{
			// },
			// "dhcp_leases": &schema.Schema{
			// },
			"discovery_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "None",
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"None", "Ping", "ARP", "Lease", "Custom",
				}, false),
			},
			"last_seen_date": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"last_discovery_date": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_known_client_identifier": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"device": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"interface": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"ptr_status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"extraneous_ptr": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"custom_properties": &schema.Schema{
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"state": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				// ValidateFunc: validation.StringInSlice([]string{
				// 	"Free", "Assigned", "Claimed", "Pending", "Held",
				// }, false),
			},
			// "hold_info": &schema.Schema{
			// 	Type:     schema.TypeList,
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
				Type:     schema.TypeBool,
				Computed: true,
			},

			// "cloud_device_info": &schema.Schema{
			// },
		},
	}
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
	d.Set("CustomProperties", ipamrec.CustomProperties)

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

func resourceIPAMRecRead(d *schema.ResourceData, m interface{}) error {

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(*Mmclient)

	ipamrec, err := c.ReadIPAMRec(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	writeIPAMRecSchema(d, ipamrec)

	return diags
}

func resourceIPAMRecCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Mmclient)

	ipamrec := readIPAMRecSchema(d)

	err := c.CreateIPAMRec(ipamrec)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(ipamrec.Address)
	err = resourceIPAMRecRead(d, m)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(d.Get("ref").(string))
	return err

}

func resourceIPAMRecUpdate(d *schema.ResourceData, m interface{}) error {

	c := m.(*Mmclient)
	ref := d.Id()
	ipamrec := readIPAMRecSchema(d)

	err := c.UpdateIPAMRec(ipamrec.IPAMProperties, ref)

	if err != nil {
		return diag.FromErr(err)
	}
	return resourceIPAMRecRead(d, m)
}

func resourceIPAMRecDelete(d *schema.ResourceData, m interface{}) error {

	c := m.(*Mmclient)
	var diags diag.Diagnostics
	ref := d.Id()
	err := c.DeleteIPAMRec(ref)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
