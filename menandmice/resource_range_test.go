package menandmice

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMenandmiceRangeCIDR(t *testing.T) {

	cidr1 := "192.168.2.0/24"
	cidr2 := "192.168.2.0/25"
	title1 := "Terraform acceptance testrange #1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckMenandmiceRangeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMenandmiceRangeConfigCIDR(cidr1, title1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("menandmice_range.testrange"),
				),
			},
			{
				Config: testAccCheckMenandmiceRangeConfigCIDR(cidr2, title1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("menandmice_range.testrange"),
				),
			},
			{
				ResourceName:      "menandmice_range.testrange",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "menandmice_range.testrange",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     cidr2,
			},
		},
	})
}

func testAccCheckMenandmiceRangeConfigCIDR(cidr, title string) string {
	return fmt.Sprintf(`
	resource menandmice_range testrange{
		cidr = "%s"
		title = "%s"
	}
	`, cidr, title)
}
func TestAccMenandmiceRangeToFrom(t *testing.T) {

	from1 := "192.168.2.0"
	from2 := "192.168.2.20"
	to1 := "192.168.2.255"
	to2 := "192.168.2.100"
	locked1 := false
	locked2 := true
	autoAssign1 := false
	autoAssign2 := true
	title1 := "Terraform acceptance testrange #1"
	title2 := "Terraform acceptance testrange #2"
	description1 := ""
	description2 := title2

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckMenandmiceRangeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMenandmiceRangeConfigToFrom(from1, to1, title1, description1, locked1, autoAssign1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("menandmice_range.testrange"),
				),
			},
			{
				Config: testAccCheckMenandmiceRangeConfigToFrom(from1, to2, title1, description1, locked1, autoAssign1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("menandmice_range.testrange"),
				),
			},
			{
				Config: testAccCheckMenandmiceRangeConfigToFrom(from2, to2, title1, description1, locked1, autoAssign1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("menandmice_range.testrange"),
				),
			},
			{
				Config: testAccCheckMenandmiceRangeConfigToFrom(from2, to2, title1, description1, locked2, autoAssign1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("menandmice_range.testrange"),
				),
			},

			{
				Config: testAccCheckMenandmiceRangeConfigToFrom(from2, to2, title1, description1, locked2, autoAssign2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("menandmice_range.testrange"),
				),
			},

			{
				Config: testAccCheckMenandmiceRangeConfigToFrom(from2, to2, title2, description1, locked2, autoAssign2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("menandmice_range.testrange"),
				),
			},
			{
				Config: testAccCheckMenandmiceRangeConfigToFrom(from2, to2, title2, description2, locked2, autoAssign2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("menandmice_range.testrange"),
				),
			},
			{
				ResourceName:      "menandmice_range.testrange",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "menandmice_range.testrange",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     from2 + "-" + to2,
			},
		},
	})
}

func testAccCheckMenandmiceRangeConfigToFrom(from, to, title, description string, locked, autoAssign bool) string {
	return fmt.Sprintf(`
	resource menandmice_range testrange{
		to = "%s"
		from= "%s"
		title = "%s"
		description = "%s"
		locked = "%t"
		auto_assign = "%t"
	}
	`, to, from, title, description, locked, autoAssign)
}

func TestAccMenandmiceRangeFreeRange(t *testing.T) {

	parentRange1 := "192.168.2.0/24"
	startAt1 := "192.168.2.0"
	startAt2 := "192.168.2.30"
	size1 := 40
	size2 := 20
	subnet1 := true
	subnet2 := true
	mask1 := 28
	mask2 := 27
	title := "Terraform acceptionat testrange #1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckMenandmiceRangeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMenandmiceRangeConfigFreeRangeSize(parentRange1, startAt1, size1, title),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("menandmice_range.testrange"),
				),
			},
			{
				Config: testAccCheckMenandmiceRangeConfigFreeRangeSize(parentRange1, startAt2, size1, title),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("menandmice_range.testrange"),
				),
			},
			{
				Config: testAccCheckMenandmiceRangeConfigFreeRangeSize(parentRange1, startAt2, size2, title),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("menandmice_range.testrange"),
				),
			},
			{
				Config: testAccCheckMenandmiceRangeConfigFreeRangeMask(parentRange1, mask1, subnet1, title),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("menandmice_range.testrange"),
				),
			},
			{
				Config: testAccCheckMenandmiceRangeConfigFreeRangeMask(parentRange1, mask1, subnet1, title),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("menandmice_range.testrange"),
				),
			},
			{
				Config: testAccCheckMenandmiceRangeConfigFreeRangeMask(parentRange1, mask2, subnet2, title),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("menandmice_range.testrange"),
				),
			},
		},
	})
}

func testAccCheckMenandmiceRangeConfigFreeRangeSize(parentRange, startAt string, size int, title string) string {
	return fmt.Sprintf(`

resource "menandmice_range" "super_range" {
  cidr = "%s"
  title = "terraform acceptance test parentRange"
}

resource "menandmice_range" "testrange" {
  free_range {
    range = menandmice_range.super_range.name
    start_at = "%s"
    size = %v
    temporary_claim_time = 1
}

  title       = "%s"
}
	`, parentRange, startAt, size, title)
}

func testAccCheckMenandmiceRangeConfigFreeRangeMask(parentRange string, mask int, subnet bool, title string) string {
	return fmt.Sprintf(`

resource "menandmice_range" "super_range" {
  cidr = "%s"
  title = "terraform acceptance test parentRange"
}

resource "menandmice_range" "testrange" {
  free_range {
    range = menandmice_range.super_range.name
    mask = %v
    temporary_claim_time = 1
  }

  subnet = %t
  title       = "%s"
}
	`, parentRange, mask, subnet, title)
}

func testAccCheckMenandmiceRangeDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*Mmclient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "menandmice_range" {
			continue
		}

		ref := rs.Primary.ID

		err := c.DeleteRange(ref)
		if err != nil {
			return err
		}
	}

	return nil
}
