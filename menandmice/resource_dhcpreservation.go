package menandmice

import (
	"terraform-provider-menandmice/diag"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceDHCPReservation() *schema.Resource {
	return &schema.Resource{
		Create: resourceDHCPResvationCreate,
		Read:   resourceDHCPResvationRead,
		Update: resourceDHCPResvationUpdate,
		Delete: resourceDHCPResvationDelete,
		Schema: map[string]*schema.Schema{

			"ref": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				ValidateFunc: validation.StringInSlice([]string{
					"DHCP", "BOOTP", "BOTH",
				}, false),
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_identifier": &schema.Schema{

				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsMACAddress, // TODO only when reservation_method is mac
			},

			"reservation_method": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "HardwareAddress", // TODO maybe ClientIdentifier is better for terraform
				ValidateFunc: validation.StringInSlice([]string{
					"HardwareAddress", "ClientIdentifier",
				}, false),
			},
			"addresses": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.Any(

						validation.IsIPv4Address,
						validation.IsIPv6Address),
					DiffSuppressFunc: ipv6AddressDiffSuppress,
				},
			},
			"ddns_hostname": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"filename": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"servername": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"next_server": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			// TODO one off dhcpserver,dhcpgroup,dhcpscope

			"owner": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"owner_ref": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func writeDHCPReservationSchema(d *schema.ResourceData, dhcpReservation DHCPReservation) {

	d.Set("ref", dhcpReservation.Ref)
	d.Set("name", dhcpReservation.Name)
	d.Set("client_identifier", dhcpReservation.ClientIdentifier)
	d.Set("reservation_method", dhcpReservation.ReservationMethod)

	d.Set("addresses", dhcpReservation.Addresses)
	d.Set("ddns_hostname", dhcpReservation.DDNSHostName)
	d.Set("filename", dhcpReservation.Filename)
	d.Set("servername", dhcpReservation.ServerName)
	d.Set("next_server", dhcpReservation.NextServer)
	d.Set("owner_ref", dhcpReservation.OwnerRef)
	return
}

func readDHCPReservationSchema(d *schema.ResourceData) DHCPReservation {

	addressesRead := d.Get("addresses").([]interface{}) //TODO check succes
	var addresses = make([]string, len(addressesRead))
	for i, address := range addressesRead {
		addresses[i] = address.(string)
	}
	dhcpReservation := DHCPReservation{
		Ref:      tryGetString(d, "ref"),
		OwnerRef: tryGetString(d, "owner_ref"),
		DHCPReservationPropertie: DHCPReservationPropertie{
			Name:              tryGetString(d, "name"),
			Type:              tryGetString(d, "type"),
			Description:       tryGetString(d, "description"),
			ClientIdentifier:  tryGetString(d, "client_identifier"),
			ReservationMethod: tryGetString(d, "reservation_method"),
			Addresses:         addresses,
			DDNSHostName:      tryGetString(d, "ddns_hostname"),
			Filename:          tryGetString(d, "filename"),
			ServerName:        tryGetString(d, "servername"),
			NextServer:        tryGetString(d, "next_server"),
		},
	}
	return dhcpReservation
}

func resourceDHCPResvationRead(d *schema.ResourceData, m interface{}) error {

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(*Mmclient)

	dhcpReservation, err := c.ReadDHCPReservation(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	writeDHCPReservationSchema(d, dhcpReservation)

	return diags
}

func resourceDHCPResvationCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Mmclient)

	dhcpReservation := readDHCPReservationSchema(d)

	ref, err := c.CreateDHCPReservation(dhcpReservation, tryGetString(d, "owner"))

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(ref)
	return resourceDHCPResvationRead(d, m)

}

func resourceDHCPResvationUpdate(d *schema.ResourceData, m interface{}) error {

	c := m.(*Mmclient)
	ref := d.Id()
	dhcpReservation := readDHCPReservationSchema(d)
	err := c.UpdateDHCPReservation(dhcpReservation.DHCPReservationPropertie, ref)

	if err != nil {
		return diag.FromErr(err)
	}
	return resourceDHCPResvationRead(d, m)
}

func resourceDHCPResvationDelete(d *schema.ResourceData, m interface{}) error {

	c := m.(*Mmclient)

	var diags diag.Diagnostics

	ref := d.Id()
	err := c.DeleteDHCPReservation(ref)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
