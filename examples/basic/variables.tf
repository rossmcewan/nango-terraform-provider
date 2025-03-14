variable "nango_api_key" {
  description = "Nango API key"
  type        = string
  sensitive   = true
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
