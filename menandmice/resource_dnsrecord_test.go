package menandmice

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMenandmiceDNSRecBasic(t *testing.T) {
	name := "terraform-test-rec1"
	date := "192.168.2.13"
	rectype := "A"
	view := ""
	server := "mandm.example.net."
	zone := "example.net."

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMenandmiceDNSRecDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMenandmiceDNSRecConfigBasic(name, date, rectype, server, zone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("menandmice_dns_record.testrec"),
				),
			},
			{
				Config: testAccCheckMenandmiceDNSRecConfigBasic(name, "192.168.2.14", rectype, server, zone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("menandmice_dns_record.testrec"),
				),
			},
			{
				ResourceName:      "menandmice_dns_record.testrec",
				ImportState:       true,
				ImportStateVerify: true,
				//TODO avoid ImportStateVerifyIgnore: "server", "zone"
				ImportStateVerifyIgnore: []string{"server", "zone", "view"},
			},
			{
				ResourceName:      "menandmice_dns_record.testrec",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     server + ":" + view + ":" + name + "." + zone + ":" + "A",
			},
		},
	})
}

func testAccCheckMenandmiceDNSRecDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*Mmclient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "menandmice_dns_record" {
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

func testAccCheckMenandmiceDNSRecConfigBasic(name, date, rectype, server, zone string) string {
	return fmt.Sprintf(`
	resource menandmice_dns_record testrec{
		name    = "%s"
		data    = "%s"
		type    = "%s"
		server  = "%s"
		zone    = "%s"
	}
	`, name, date, rectype, server, zone)
}
