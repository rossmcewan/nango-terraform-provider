package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rossmcewan/nango-terraform-provider/internal/resources"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("NANGO_API_KEY", nil),
				Description: "The API key for Nango API operations. Can also be set via NANGO_API_KEY environment variable.",
			},
			"base_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NANGO_BASE_URL", "https://api.nango.dev"),
				Description: "The base URL of the Nango API. Useful for self-hosted instances or custom endpoints. Can also be set via NANGO_BASE_URL environment variable.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"nango_integration": resources.ResourceIntegration(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"nango_integration": resources.DataSourceIntegration(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiKey := d.Get("api_key").(string)
	baseURL := d.Get("base_url").(string)

	// Validate and normalize the base URL
	baseURL = strings.TrimRight(baseURL, "/")
	if !strings.HasPrefix(baseURL, "http://") && !strings.HasPrefix(baseURL, "https://") {
		return nil, diag.Errorf("base_url must be a valid HTTP or HTTPS URL")
	}

	// Initialize the client
	client := resources.NewNangoClient(apiKey, baseURL)

	return client, diag.Diagnostics{}
}
