package menandmice

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMenandmiceDHCPReservationBasic(t *testing.T) {
	name := "terraform-test-reservation"
	owner := "DHCPScopes/192.168.2.128/25"
	clientIdentifier := "44:55:66:77:88:00"
	addressess1 := `"192.168.2.138"`
	// addressess2 := `"192.168.2.138","192.168.2.139"`

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckMenandmiceDHCPReservationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMenandmiceDHCPReservationConfigBasic(name, owner, clientIdentifier, addressess1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("menandmice_dhcp_reservation.testreservation"),
				),
			},
			// TODO test update dhcp reservation
			// {
			// 	Config: testAccCheckMenandmiceDHCPReservationConfigBasic(name, owner, clientIdentifier, addressess2),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckResourceExists("menandmice_dhcp_reservation.testreservation"),
			// 	),
			// },
			{

				ResourceName:      "menandmice_dhcp_reservation.testreservation",
				ImportState:       true,
				ImportStateId:     name,
				ImportStateVerify: true,

				// owner is not stored on server, only owner-ref
				// and you can't owner is not unique
				//TODO avoid ImportStateVerifyIgnore: "owner"
				ImportStateVerifyIgnore: []string{"owner"},
			},

			{
				ResourceName:      "menandmice_dhcp_reservation.testreservation",
				ImportState:       true,
				ImportStateVerify: true,

				// owner is not stored on server, only owner-ref
				// and you can't owner is not unique
				//TODO avoid ImportStateVerifyIgnore: "owner"
				ImportStateVerifyIgnore: []string{"owner"},
			},
		},
	})
}

func testAccCheckMenandmiceDHCPReservationDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*Mmclient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "menandmice_dhcp_reservation" {
			continue
		}

		ref := rs.Primary.ID

		err := c.DeleteDHCPReservation(ref)
		if err != nil {
			return err
		}
	}

	return nil
}

func testAccCheckMenandmiceDHCPReservationConfigBasic(name, owner, clientIdentifier, addressess string) string {
	return fmt.Sprintf(`
	resource menandmice_dhcp_reservation testreservation{
		name              = "%s"
		owner             = "%s"
		client_identifier = "%s"
		addresses          = [%s]
	}
	`, name, owner, clientIdentifier, addressess)
}
