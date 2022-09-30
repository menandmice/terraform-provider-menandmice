package menandmice

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMenandmiceIPAMRcBasic(t *testing.T) {

	address1 := "192.168.2.15"
	// address2 := "::192.168.2.15" //TODO test ipv6
	location := "here"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMenandmiceIPAMRecDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMenandmiceIPAMRecConfigBasic(address1, location, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("menandmice_ipam_record.testipam"),
				),
			},
			{
				Config: testAccCheckMenandmiceIPAMRecConfigBasic(address1, location, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("menandmice_ipam_record.testipam"),
				),
			},
			// TODO
			// {
			// 	Config: testAccCheckMenandmiceIPAMRecConfigBasic(address2, location, true),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckResourceExists("menandmice_ipam_record.testipam"),
			// 	),
			// },

			// TODO add test for find free ip
			{
				ResourceName:      "menandmice_ipam_record.testipam",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     address1,
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
		custom_properties = {"Location":"%s"}
		claimed = %t
	}
	`, address, location, claimed)
}
