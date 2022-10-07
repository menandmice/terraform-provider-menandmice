package menandmice

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDNSRec(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMenandmiceDNSRecDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: testAccCreateNomadAclPolicies(t, numPolicies),
			},
		},
	})
}

func testAccDataSourceDNSZONEConfig(name, zone, recType, data string) string {
	return fmt.Sprintf(`
data "menandmice_dns_zone" "testzone" {
  name = "%s"
  zone =  "%s"
  type = "%s"
  data = "%s"
}
	`, name, zone, recType, data)
}
