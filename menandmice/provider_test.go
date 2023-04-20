package menandmice

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAccProviders = map[string]func() (*schema.Provider, error){
	"menandmice": func() (*schema.Provider, error) {
		return Provider("test")(), nil
	},
}

// This provider can be used in testing code for API calls without requiring
// the use of saving and referencing specific ProviderFactories instances.
//
// PreCheck(t) must be called before using this provider instance
var testAccProvider *schema.Provider = Provider("test")()

func TestProvider(t *testing.T) {
	if provider := Provider("test"); provider != nil {
		t.Fatalf("could not initialise provier")
	}
}
func testAccPreCheck(t *testing.T) {

	// Might exist better solution: ref https://github.com/hashicorp/terraform-provider-scaffolding/issues/79
	testAccProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(nil))

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
