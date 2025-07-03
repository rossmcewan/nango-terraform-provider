# Basic Nango Integration Example

This example demonstrates how to use the Nango Terraform provider to create a GitHub integration using the official Nango API.

## Usage

1. Copy `terraform.tfvars.example` to `terraform.tfvars` and fill in your actual values.

2. Initialize Terraform with `terraform init`.

3. Apply the configuration with `terraform apply`.

4. To destroy the created resources when you're done, run `terraform destroy`.

## Configuration

The example creates a GitHub integration with the following configuration:

- **unique_key**: A unique identifier for the integration (`github-nango-community`)
- **provider_name**: The provider name (`github`)
- **display_name**: The display name shown in the Nango dashboard (`GitHub`)
- **credentials**: OAuth2 credentials including client ID, client secret, and scopes

## API URL Configuration

You can override the default Nango API URL for self-hosted instances or custom endpoints:

```terraform
provider "nango" {
  api_key  = var.nango_api_key
  base_url = "https://your-custom-nango-instance.com"  # Optional
}
```

Or set via environment variable:
```bash
export NANGO_BASE_URL="https://your-custom-nango-instance.com"
```

## Requirements

- Terraform >= 0.13.x
- A Nango account with API access
- GitHub OAuth application credentials

## API Endpoints Used

This provider uses the following Nango API endpoints:
- `POST /integrations` - Create integration
- `GET /integrations/{unique_key}` - Read integration
- `PATCH /integrations/{unique_key}` - Update integration
- `DELETE /integrations/{unique_key}` - Delete integration
