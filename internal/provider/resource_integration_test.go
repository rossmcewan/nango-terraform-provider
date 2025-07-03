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
					resource.TestCheckResourceAttr("nango_integration.test", "unique_key", "test-github-integration"),
					resource.TestCheckResourceAttr("nango_integration.test", "provider_name", "github"),
					resource.TestCheckResourceAttr("nango_integration.test", "display_name", "Test GitHub"),
				),
			},
		},
	})
}

func testAccNangoIntegrationConfig() string {
	return `
resource "nango_integration" "test" {
  unique_key   = "test-github-integration"
  provider_name = "github"
  display_name = "Test GitHub"
  
  credentials {
    type          = "OAUTH2"
    client_id     = "test-client-id"
    client_secret = "test-client-secret"
    scopes        = "repo,user"
  }
}
`
}
