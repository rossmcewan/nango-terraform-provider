variable "nango_api_key" {
  description = "Nango API key"
  type        = string
  sensitive   = true
}

variable "nango_base_url" {
  description = "Nango API base URL (optional, defaults to https://api.nango.dev)"
  type        = string
  default     = "https://api.nango.dev"
}

variable "github_client_id" {
  description = "GitHub OAuth client ID"
  type        = string
}

variable "github_client_secret" {
  description = "GitHub OAuth client secret"
  type        = string
  sensitive   = true
}
