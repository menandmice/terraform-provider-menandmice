package menandmice

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMenandmiceDNSZoneBasic(t *testing.T) {

	name1 := "terraform-test-zone1.net."
	name2 := "terraform-test-zone2.net."
	authority1 := "ext-master.mmdemo.net."
	authority2 := "dc16.mmdemo.net."
	view := ""

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckMenandmiceDNSZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMenandmiceDNSZoneConfigBasic(name1, authority1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("menandmice_dns_zone.testzone"),
				),
			},
			{
				Config: testAccCheckMenandmiceDNSZoneConfigBasic(name2, authority1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("menandmice_dns_zone.testzone"),
				),
			},
			{
				Config: testAccCheckMenandmiceDNSZoneConfigBasic(name2, authority2),
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
				ImportStateId:     authority2 + ":" + view + ":" + name2,
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
