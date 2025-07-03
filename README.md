# Nango Terraform Provider

This Terraform provider allows you to manage your Nango resources as Infrastructure as Code (IaC). With this provider, you can create, update, and delete Nango integrations programmatically through Terraform using the official Nango API.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
- [Go](https://golang.org/doc/install) >= 1.18 (to build the provider plugin)
- A Nango account with API access

## Installation

### Using Terraform Registry (Once Published)
```terraform
terraform {
  required_providers {
    nango = {
      source  = "[username]/nango"
      version = "~> 1.0.0"
    }
  }
}

provider "nango" {
  api_key = var.nango_api_key
}
```
### Local Development Installation

# Build the provider
```bash
go build -o terraform-provider-nango
```

# Create a local development directory
```bash
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/[username]/nango/1.0.0/[OS_ARCH]
```

# Copy the provider to the development directory
```bash
cp terraform-provider-nango ~/.terraform.d/plugins/registry.terraform.io/[username]/nango/1.0.0/[OS_ARCH]/terraform-provider-nango_v1.0.0
```
Replace `[OS_ARCH]` with your system's OS and architecture (e.g., `linux_amd64`, `darwin_amd64`).
Replace `[username]` with your active username.

## Authentication

The provider requires a Nango API key for authentication. You can provide this in several ways:

1. Set the `api_key` parameter in the provider configuration:
```terraform
provider "nango" {
  api_key = "your-api-key"
}
```

2. Set the `NANGO_API_KEY` environment variable:
```bash
export NANGO_API_KEY="your-api-key"
```

## Provider Configuration
```terraform  
provider "nango" {
  api_key  = var.nango_api_key                # Required: Your Nango API key
  base_url = "https://api.nango.dev"          # Optional: Defaults to https://api.nango.dev
}
```

## Resources

### `nango_integration`

Manages a Nango integration using the official Nango API.

#### Example Usage
```terraform
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
```
#### Argument Reference

* `unique_key` - (Required) A unique identifier for the integration.
* `provider_name` - (Required) The provider name (e.g., 'slack', 'github').
* `display_name` - (Required) The display name for the integration.
* `credentials` - (Required) OAuth credentials configuration block:
  * `type` - (Required) The type of credentials (e.g., 'OAUTH1', 'OAUTH2').
  * `client_id` - (Required) The OAuth client ID.
  * `client_secret` - (Required) The OAuth client secret.
  * `scopes` - (Optional) The OAuth scopes as a comma-separated string.

#### Attribute Reference

* `id` - The unique key of the integration.

## Data Sources

### `nango_integration`

Retrieves information about a Nango integration.

#### Example Usage
```terraform
data "nango_integration" "github" {
  unique_key = "github-nango-community"
}
```

```terraform    
output "github_integration_provider" {
  value = data.nango_integration.github.provider_name
}
```

## Complete Example
```terraform
terraform {
  required_providers {
    nango = {
      source  = "[username]/nango"
      version = "~> 1.0.0"
    }
  }
}

provider "nango" {
  api_key = var.nango_api_key
}
```
# Create a GitHub integration
```terraform
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
```

# Output the integration ID
```terraform
output "github_integration_id" {
  value = nango_integration.github.id
}
```

## API Endpoints

This provider uses the following Nango API endpoints:
- `POST /integrations` - Create integration
- `GET /integrations/{unique_key}` - Read integration
- `PATCH /integrations/{unique_key}` - Update integration
- `DELETE /integrations/{unique_key}` - Delete integration

## Development

### Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
- [Go](https://golang.org/doc/install) >= 1.18

### Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the Go `build` command:
```bash
go build -o terraform-provider-nango
```

### Testing The Provider

Set the required environment variables:
```bash
export NANGO_API_KEY="your-api-key"
```

Run the tests:
```bash 
# Run unit tests
go test ./...

# Run acceptance tests (creates real resources)
TF_ACC=1 go test ./... -v
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This provider is distributed under the [MIT License](LICENSE).
