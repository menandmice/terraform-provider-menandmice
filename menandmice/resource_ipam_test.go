package menandmice

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMenandmiceIPAMRcBasic(t *testing.T) {
	address := "192.168.2.15"
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

			// TODO add test for find free ip
			{
				ResourceName:      "menandmice_ipam_record.testipam",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "::5",
			},

			{
				ResourceName:      "menandmice_ipam_record.testipam",
				ImportState:       true,
				ImportStateVerify: true,
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
	resource menandmice_ipam_record testipam {
		address= "%s"
		custom_properties = {"location":"%s"}
		claimed = %t
	}
	`, address, location, claimed)
}
