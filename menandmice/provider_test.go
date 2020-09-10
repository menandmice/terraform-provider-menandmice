package menandmice

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]terraform.ResourceProvider{
		"menandmice": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if err := os.Getenv("MENANDMICE_USERNAME"); err == "" {
		t.Fatal("MENANDMICE_USERNAME must be set for acceptance tests")
	}
	if err := os.Getenv("MENANDMICE_PASSWORD"); err == "" {
		t.Fatal("MENANDMICE_PASSWORD must be set for acceptance tests")
	}

	if err := os.Getenv("MENANDMICE_ENDPOINT"); err == "" {
		t.Fatal("MENANDMICE_ENDPOINT must be set for acceptance tests")
	}
}
