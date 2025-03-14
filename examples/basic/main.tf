terraform {
  required_providers {
    nango = {
      source  = "rossmcewan/nango"
      version = "~> 1.0.0"
    }
  }
}

provider "nango" {
  api_key = var.nango_api_key
}

resource "nango_integration" "github" {
  name                = "github"
  provider_config_key = "github"
  oauth_client_id     = var.github_client_id
  oauth_client_secret = var.github_client_secret
  oauth_scopes        = ["repo", "user"]
}

output "github_integration_id" {
  value = nango_integration.github.id
}
