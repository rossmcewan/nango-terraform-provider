package resources

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceIntegration returns the resource definition for a Nango integration
func ResourceIntegration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIntegrationCreate,
		ReadContext:   resourceIntegrationRead,
		UpdateContext: resourceIntegrationUpdate,
		DeleteContext: resourceIntegrationDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the integration",
			},
			"provider_config_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The provider configuration key",
			},
			"oauth_client_id": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The OAuth client ID",
			},
			"oauth_client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The OAuth client secret",
			},
			"oauth_scopes": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The OAuth scopes",
			},
			// Add more fields as needed
		},
	}
}

func resourceIntegrationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// client := m.(*NangoClient)

	// Extract values from schema
	name := d.Get("name").(string)
	providerConfigKey := d.Get("provider_config_key").(string)
	oauthClientID := d.Get("oauth_client_id").(string)
	oauthClientSecret := d.Get("oauth_client_secret").(string)

	// Prepare request body
	integration := map[string]interface{}{
		"name":                name,
		"provider_config_key": providerConfigKey,
		"oauth_client_id":     oauthClientID,
		"oauth_client_secret": oauthClientSecret,
	}

	// Add scopes if present
	if v, ok := d.GetOk("oauth_scopes"); ok {
		scopes := make([]string, 0)
		for _, s := range v.([]interface{}) {
			scopes = append(scopes, s.(string))
		}
		integration["oauth_scopes"] = scopes
	}

	// Make API request to create integration
	// Implementation depends on Nango API

	// Set ID for the resource
	d.SetId(fmt.Sprintf("%s-%d", name, time.Now().Unix()))

	return resourceIntegrationRead(ctx, d, m)
}

func resourceIntegrationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Implementation to read integration from Nango API
	return diag.Diagnostics{}
}

func resourceIntegrationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Implementation to update integration in Nango API
	return resourceIntegrationRead(ctx, d, m)
}

func resourceIntegrationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Implementation to delete integration from Nango API
	return diag.Diagnostics{}
}

// DataSourceIntegration returns the data source definition for a Nango integration
func DataSourceIntegration() *schema.Resource {
	// Similar to ResourceIntegration but for data source
	return &schema.Resource{
		ReadContext: dataSourceIntegrationRead,
		Schema:      map[string]*schema.Schema{
			// Schema definition
		},
	}
}

func dataSourceIntegrationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Implementation to read integration from Nango API for data source
	return diag.Diagnostics{}
}
