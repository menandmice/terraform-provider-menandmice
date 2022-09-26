package menandmice

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMenandmiceDNSZoneBasic(t *testing.T) {

	name := "terraform-test-zone.net."
	authority1 := "ext-master.mmdemo.net."
	authority2 := "dc16.mmdemo.net."
	view := ""

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMenandmiceDNSZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMenandmiceDNSZoneConfigBasic(name, authority1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("menandmice_dns_zone.testzone"),
				),
			},
			{
				Config: testAccCheckMenandmiceDNSZoneConfigBasic(name, authority2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("menandmice_dns_zone.testzone"),
				),
			},

			{
				ResourceName:      "menandmice_dns_zone.testzone",
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				ResourceName:      "menandmice_dns_zone.testzone",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     authority1 + ":" + view + ":" + name,
			},
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
