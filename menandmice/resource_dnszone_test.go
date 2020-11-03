package menandmice

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMenandmiceDNSZoneBasic(t *testing.T) {

	name := "zone1."
	authority := "mandm.example.net."

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMenandmiceDNSZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMenandmiceDNSZoneConfigBasic(name, authority),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("menandmice_dns_zone.testzone"),
				),
			},
			// {
			// 	Config: testAccCheckMenandmiceDNSZoneConfigBasic(name, "mandm.example.com."),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckResourceExists("menandmice_dns_zone.testzone"),
			// 	),
			// },

			// TODO test minimal parameters,
			// TODO test with all parameters set to non default
			// TODO test update, and recreate
		},
	})
}

func testAccCheckMenandmiceDNSZoneDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*Mmclient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "menandmice_dns_zone" {
			continue
		}

		ref := rs.Primary.ID

		err := c.DeleteDNSZone(ref)
		if err != nil {
			return err
		}
	}

	return nil
}

func testAccCheckMenandmiceDNSZoneConfigBasic(name, authority string) string {
	return fmt.Sprintf(`
	resource menandmice_dns_zone testzone{
		name    = "%s"
		authority   = "%s"
	}
	`, name, authority)
}
