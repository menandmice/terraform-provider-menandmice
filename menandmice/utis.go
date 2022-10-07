package menandmice

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func tryGetString(d *schema.ResourceData, key string) string {
	if v, ok := d.GetOk(key); ok {
		return v.(string)

	}
	return ""
}

func testAccCheckResourceExists(n string) resource.TestCheckFunc {
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

// convert structure to map ignore `json:"omitempty"`
func toMap(item interface{}) (map[string]interface{}, error) {

	var properties map[string]interface{}
	serialized, err := json.Marshal(item)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(serialized, &properties)

	if err != nil {
		return nil, err
	}
	return properties, nil
}

func ipv6AddressDiffSuppress(_, old, new string, _ *schema.ResourceData) bool {
	oldIP := net.ParseIP(old)
	newIP := net.ParseIP(new)

	return oldIP.Equal(newIP)
}
func toError(diags diag.Diagnostics) error {

	if diags.HasError() {

		var message string
		for _, diag := range diags {
			message += diag.Summary + "\n"

		}
		return fmt.Errorf(message)

	}
	return nil
}
