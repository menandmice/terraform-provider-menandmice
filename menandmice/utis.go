package menandmice

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

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

// SetFromMap set schema from map of interface
func SetFromMap(d *schema.ResourceData, m map[string]interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	for key, value := range m {

		err := d.Set(key, value)
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}
	return diags
}

func MmTimeString2rfc(timeStr string, location *time.Location) (string, error) {
	if timeStr == "" {
		return "", nil
	}
	t, err := time.ParseInLocation("Jan 2, 2006 15:04:05", timeStr, location)
	return t.In(time.Local).Format(time.RFC3339), err
}
