package menandmice

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccMenandmiceIPAMRcBasic(t *testing.T) {
	address := "10.0.0.1"
	location := "here"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMenandmiceIPAMRecDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMenandmiceIPAMRecConfigBasic(address, location, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("menandmice_ipam_record.testipam"),
				),
			},
			{
				Config: testAccCheckMenandmiceIPAMRecConfigBasic(address, location, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("menandmice_ipam_record.testipam"),
				),
			},
			{
				Config: testAccCheckMenandmiceIPAMRecConfigBasic("::5", location, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("menandmice_ipam_record.testipam"),
				),
			},
		},
	})
}

func testAccCheckMenandmiceIPAMRecDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*Mmclient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "menandmice_ipam_record" {
			continue
		}

		ref := rs.Primary.ID

		err := c.DeleteIPAMRec(ref)
		if err != nil {
			return err
		}
	}

	return nil
}

func testAccCheckMenandmiceIPAMRecConfigBasic(address, location string, claimed bool) string {
	return fmt.Sprintf(`
	resource menandmice_ipam_record testipam{
		address= "%s"
		custom_properties = {"location":"%s"}
		claimed = %t
	}
	`, address, location, claimed)
}
