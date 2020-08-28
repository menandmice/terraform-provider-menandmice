package menandmice

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMenandmiceDNSrecBasic(t *testing.T) {
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
					testAccCheckMenandmiceDNSRecExists("menandmice_dnsrecord.testrec"),
				),
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

func testAccCheckMenandmiceDNSRecExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Ref set")
		}

		return nil
	}
}
