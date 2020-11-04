package menandmice

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDHCPReservation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDHCPResvationCreate,
		ReadContext:   resourceDHCPResvationRead,
		UpdateContext: resourceDHCPResvationUpdate,
		DeleteContext: resourceDHCPResvationDelete,
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
				ForceNew: true,
				Default:  "HardwareAddress", // TODO maybe ClientIdentifier is better for terraform
				ValidateFunc: validation.StringInSlice([]string{
					"HardwareAddress", "ClientIdentifier",
				}, false),
			},
			"addresses": &schema.Schema{
				Type:     schema.TypeList,
				ForceNew: true,
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
		Ref:               tryGetString(d, "ref"),
		OwnerRef:          tryGetString(d, "owner_ref"),
		ReservationMethod: tryGetString(d, "reservation_method"),
		Addresses:         addresses,
		DHCPReservationPropertie: DHCPReservationPropertie{
			Name:             tryGetString(d, "name"),
			Type:             tryGetString(d, "type"),
			Description:      tryGetString(d, "description"),
			ClientIdentifier: tryGetString(d, "client_identifier"),
			DDNSHostName:     tryGetString(d, "ddns_hostname"),
			Filename:         tryGetString(d, "filename"),
			ServerName:       tryGetString(d, "servername"),
			NextServer:       tryGetString(d, "next_server"),
		},
	}
	return dhcpReservation
}

func resourceDHCPResvationRead(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	client := m.(*Mmclient)

	dhcpReservation, err := client.ReadDHCPReservation(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if dhcpReservation == nil {
		d.SetId("")
		return diags
	}
	writeDHCPReservationSchema(d, *dhcpReservation)

	return diags
}

func resourceDHCPResvationCreate(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Mmclient)

	dhcpReservation := readDHCPReservationSchema(d)

	ref, err := client.CreateDHCPReservation(dhcpReservation, tryGetString(d, "owner"))

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(ref)
	return resourceDHCPResvationRead(c, d, m)

}

func resourceDHCPResvationUpdate(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*Mmclient)
	ref := d.Id()
	dhcpReservation := readDHCPReservationSchema(d)
	err := client.UpdateDHCPReservation(dhcpReservation.DHCPReservationPropertie, ref)

	if err != nil {
		return diag.FromErr(err)
	}
	return resourceDHCPResvationRead(c, d, m)
}

func resourceDHCPResvationDelete(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*Mmclient)

	var diags diag.Diagnostics

	ref := d.Id()
	err := client.DeleteDHCPReservation(ref)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
