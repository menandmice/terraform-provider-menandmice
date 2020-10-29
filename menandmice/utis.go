package menandmice

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func tryGetString(d *schema.ResourceData, key string) string {
	if v, ok := d.GetOk(key); ok {
		return v.(string)

	}
	return ""
}

// TODO remove. make it easy to add filters to data_source, but not practical
func ToFilter(d *schema.ResourceData) map[string]string {
	switch v := d.Get("").(type) {

	case map[string]interface{}:
		var result = make(map[string]string)
		for key := range v {
			if val, ok := d.GetOk(key); ok {
				// TODO this will only work Aggregate Types
				// TODO this wont work with if name is diffrent

				result[key] = fmt.Sprintf("%v", val)
			}
		}
		return result

	default:
		panic("you can not do this with this schema")
	}
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
