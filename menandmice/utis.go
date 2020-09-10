package menandmice

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func tryGetString(d *schema.ResourceData, key string) string {
	if v, ok := d.GetOk(key); ok {
		return v.(string)
	}
	return ""
}
