package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"nango": testAccProvider,
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("NANGO_API_KEY"); v == "" {
		t.Fatal("NANGO_API_KEY must be set for acceptance tests")
	}
}

func testAccCheckNangoIntegrationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Implementation to check if the integration exists
		return nil
	}
}

func testAccCheckNangoIntegrationDestroy(s *terraform.State) error {
	// Implementation to check if the integration was destroyed
	return nil
}
