package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNangoIntegration_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNangoIntegrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNangoIntegrationConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNangoIntegrationExists("nango_integration.test"),
					resource.TestCheckResourceAttr("nango_integration.test", "name", "test-github"),
					resource.TestCheckResourceAttr("nango_integration.test", "provider_config_key", "github"),
				),
			},
		},
	})
}

func testAccNangoIntegrationConfig() string {
	return `
resource "nango_integration" "test" {
  name                = "test-github"
  provider_config_key = "github"
  oauth_client_id     = "test-client-id"
  oauth_client_secret = "test-client-secret"
  oauth_scopes        = ["repo", "user"]
}
`
}
