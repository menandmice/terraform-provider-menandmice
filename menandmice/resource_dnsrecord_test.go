package menandmice

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccMenandmiceDNSRecBasic(t *testing.T) {
	name := "rec1"
	date := "127.0.0.1"
	rectype := "A"
	zone := "DNSZones/217"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMenandmiceDNSRecDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMenandmiceDNSRecConfigBasic(name, date, rectype, zone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("menandmice_dnsrecord.testrec"),
				),
			},
			{
				Config: testAccCheckMenandmiceDNSRecConfigBasic(name, "::1", "AAAA", zone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("menandmice_dnsrecord.testrec"),
				),
				// TODO test minimal parameters,
				// TODO test with all parameters set to non default
				// TODO test update, and recreate with change zone
			},
		},
	})
}

func testAccCheckMenandmiceDNSRecDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*Mmclient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "menandmice_dnsrecord" {
			continue
		}

		ref := rs.Primary.ID

		err := c.DeleteDNSRec(ref)
		if err != nil {
			return err
		}
	}

	return nil
}

func testAccCheckMenandmiceDNSRecConfigBasic(name, date, rectype, zone string) string {
	return fmt.Sprintf(`
	resource menandmice_dnsrecord testrec{
		name    = "%s"
		data    = "%s"
		type    = "%s"
		dnszone = "%s"
	}
	`, name, date, rectype, zone)
}
