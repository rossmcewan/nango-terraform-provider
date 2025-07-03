terraform {
  required_providers {
    nango = {
      source  = "rossmcewan/nango"
      version = "~> 1.0.0"
    }
  }
}

provider "nango" {
  api_key  = var.nango_api_key
  base_url = var.nango_base_url  # Optional: defaults to https://api.nango.dev
  # For self-hosted instances, use: base_url = "https://your-nango-instance.com"
}

resource "nango_integration" "github" {
  unique_key   = "github-nango-community"
  provider_name = "github"
  display_name = "GitHub"
  
  credentials {
    type          = "OAUTH2"
    client_id     = var.github_client_id
    client_secret = var.github_client_secret
    scopes        = "repo,user"
  }
}

output "github_integration_id" {
  value = nango_integration.github.id
}
