package menandmice

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMenandmiceDHCPReservationBasic(t *testing.T) {
	name := "testresrvation"
	owner := "mandm.example.net."
	clientIdentifier := "44:55:66:77:88:00"
	addressess := `"172.16.17.9"`

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMenandmiceDHCPReservationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMenandmiceDHCPReservationConfigBasic(name, owner, clientIdentifier, addressess),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("menandmice_dhcp_reservation.testresrvation"),
				),
			},
			{

				Config: testAccCheckMenandmiceDHCPReservationConfigBasic(name, owner, clientIdentifier, `"172.16.17.5","172.16.17.6"`),
				Check: resource.ComposeTestCheckFunc(

					testAccCheckResourceExists("menandmice_dhcp_reservation.testresrvation"),
				),
				// TODO test minimal parameters,
				// TODO test with all parameters set to non default
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

func testAccCheckMenandmiceDHCPReservationConfigBasic(name, owner, client_identifier, addressess string) string {
	return fmt.Sprintf(`
	resource menandmice_dhcp_reservation testresrvation{
		name              = "%s"
		owner             = "%s"
		client_identifier = "%s"
		addresses          = [%s]
	}
	`, name, owner, client_identifier, addressess)
}
