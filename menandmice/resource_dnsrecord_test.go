package menandmice

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMenandmiceDNSRecBasic(t *testing.T) {

	zone := "terraform-test-zone.net."
	authority := "ext-master.mmdemo.net."

	name := "terraform-test-rec1"
	data1 := "192.168.2.13"
	data2 := "192.168.2.14"
	rectype := "A"
	view := ""

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMenandmiceDNSRecDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMenandmiceDNSRecConfigBasic(name, data1, rectype, authority, zone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("menandmice_dns_record.testrec"),
				),
			},
			{
				Config: testAccCheckMenandmiceDNSRecConfigBasic(name, data2, rectype, authority, zone),
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
				ResourceName:            "menandmice_dns_record.testrec",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateId:           authority + ":" + view + ":" + name + "." + zone + ":" + "A",
				ImportStateVerifyIgnore: []string{"view"},
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

	resource menandmice_dns_zone testzone{
		name    = "%s"
		authority   = "%s"
	}

	resource menandmice_dns_record testrec{
		name    = "%s"
		data    = "%s"
		type    = "%s"
		server  = "%s"
		zone    = menandmice_dns_zone.testzone.name
	}
	`, zone, server, name, date, rectype, server)
}
