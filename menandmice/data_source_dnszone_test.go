package menandmice

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDNSRec(t *testing.T) {

	zone := "terraform-test-zone.net."
	authority := "ext-master.mmdemo.net."

	name := "terraform-test-rec1"
	data := "192.168.2.13"
	recType := "A"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMenandmiceDNSRecDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDNSZONEConfig(name, zone, authority, recType, data),
				// PlanOnly:           true,
				// ExpectNonEmptyPlan: false,
			},
		},
	})
}
func testAccDataSourceDNSZONEConfig(name, zone, authority, recType, data string) string {
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

data "menandmice_dns_record" "testrec" {
  name = menandmice_dns_record.testrec.name
  server = menandmice_dns_record.testrec.server
  zone = menandmice_dns_record.testrec.zone
  type = menandmice_dns_record.testrec.type
}
	`, zone, authority, name, data, recType, authority)
}
