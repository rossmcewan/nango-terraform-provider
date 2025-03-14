package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

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
	client := m.(*NangoClient)

	// Extract values from schema
	name := d.Get("name").(string)
	providerConfigKey := d.Get("provider_config_key").(string)
	oauthClientID := d.Get("oauth_client_id").(string)
	oauthClientSecret := d.Get("oauth_client_secret").(string)

	// Prepare request body
	integration := map[string]interface{}{
		"provider_config_key": providerConfigKey,
		"provider":            name, // The provider field in the API corresponds to the name in our schema
		"oauth_client_id":     oauthClientID,
		"oauth_client_secret": oauthClientSecret,
	}

	// Add scopes if present
	if v, ok := d.GetOk("oauth_scopes"); ok {
		scopes := make([]string, 0)
		for _, s := range v.([]interface{}) {
			scopes = append(scopes, s.(string))
		}
		integration["oauth_scopes"] = strings.Join(scopes, ",") // API expects comma-separated string
	}

	// Log the request for debugging
	requestJSON, _ := json.MarshalIndent(integration, "", "  ")
	fmt.Printf("Request to /config: %s\n", string(requestJSON))

	// Make API request to create integration
	resp, err := client.MakeRequest("POST", "/config", integration)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	// Read and log the response body
	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	fmt.Printf("Response from /config (%d): %s\n", resp.StatusCode, bodyString)

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error creating integration: %s - %s", resp.Status, bodyString)
	}

	// Parse the response
	var result struct {
		Config struct {
			UniqueKey string `json:"unique_key"`
			Provider  string `json:"provider"`
		} `json:"config"`
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return diag.Errorf("Error parsing response: %v - Response body: %s", err, bodyString)
	}

	// Set the ID from the response
	d.SetId(result.Config.UniqueKey)

	return resourceIntegrationRead(ctx, d, m)
}

func resourceIntegrationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*NangoClient)

	// Get the ID (unique_key in Nango API)
	id := d.Id()

	// Make API request to get integration
	resp, err := client.MakeRequest("GET", fmt.Sprintf("/integrations/%s?include=credentials", id), nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		// Integration was deleted outside of Terraform
		d.SetId("")
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return diag.Errorf("Error reading integration: %s - %s", resp.Status, string(bodyBytes))
	}

	// Parse the response
	var result struct {
		Data struct {
			UniqueKey   string `json:"unique_key"`
			DisplayName string `json:"display_name"`
			Provider    string `json:"provider"`
			Credentials struct {
				ClientID     string `json:"client_id"`
				ClientSecret string `json:"client_secret"`
				Scopes       string `json:"scopes"`
			} `json:"credentials"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	// Set the resource data from the response
	if err := d.Set("name", result.Data.Provider); err != nil {
		return diag.FromErr(err)
	}

	// The provider_config_key is not directly returned in the response
	// We can infer it from the unique_key or keep the existing value

	if err := d.Set("oauth_client_id", result.Data.Credentials.ClientID); err != nil {
		return diag.FromErr(err)
	}

	// Don't set oauth_client_secret as it might be masked in the response

	// Set scopes if present
	if result.Data.Credentials.Scopes != "" {
		scopes := strings.Split(result.Data.Credentials.Scopes, ",")
		if err := d.Set("oauth_scopes", scopes); err != nil {
			return diag.FromErr(err)
		}
	}

	return diag.Diagnostics{}
}

func resourceIntegrationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*NangoClient)

	// Extract values from schema
	name := d.Get("name").(string)
	providerConfigKey := d.Get("provider_config_key").(string)
	oauthClientID := d.Get("oauth_client_id").(string)
	oauthClientSecret := d.Get("oauth_client_secret").(string)

	// Prepare request body
	integration := map[string]interface{}{
		"provider_config_key": providerConfigKey,
		"provider":            name, // The provider field in the API corresponds to the name in our schema
		"oauth_client_id":     oauthClientID,
		"oauth_client_secret": oauthClientSecret,
	}

	// Add scopes if present
	if v, ok := d.GetOk("oauth_scopes"); ok {
		scopes := make([]string, 0)
		for _, s := range v.([]interface{}) {
			scopes = append(scopes, s.(string))
		}
		integration["oauth_scopes"] = strings.Join(scopes, ",") // API expects comma-separated string
	}

	// Make API request to update integration
	resp, err := client.MakeRequest("PUT", "/config", integration)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return diag.Errorf("Error updating integration: %s - %s", resp.Status, string(bodyBytes))
	}

	// Parse the response
	var result struct {
		Config struct {
			UniqueKey string `json:"unique_key"`
			Provider  string `json:"provider"`
		} `json:"config"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	return resourceIntegrationRead(ctx, d, m)
}

func resourceIntegrationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*NangoClient)

	// Get the provider_config_key
	providerConfigKey := d.Get("provider_config_key").(string)

	// Make API request to delete integration
	resp, err := client.MakeRequest("DELETE", fmt.Sprintf("/config/%s", providerConfigKey), nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return diag.Errorf("Error deleting integration: %s - %s", resp.Status, string(bodyBytes))
	}

	// Parse the response
	var result struct {
		Success bool `json:"success"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if !result.Success {
		return diag.Errorf("Failed to delete integration: API returned success=false")
	}

	// Set the ID to empty to mark it as deleted
	d.SetId("")

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
