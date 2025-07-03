package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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
			"unique_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A unique identifier for the integration",
			},
			"provider_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The provider name (e.g., 'slack', 'github')",
			},
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The display name for the integration",
			},
			"credentials": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "OAuth credentials configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The type of credentials (e.g., 'OAUTH1', 'OAUTH2')",
						},
						"client_id": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "The OAuth client ID",
						},
						"client_secret": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "The OAuth client secret",
						},
						"scopes": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The OAuth scopes as a comma-separated string",
						},
					},
				},
			},
		},
	}
}

func resourceIntegrationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*NangoClient)

	// Extract values from schema
	uniqueKey := d.Get("unique_key").(string)
	providerName := d.Get("provider_name").(string)
	displayName := d.Get("display_name").(string)
	credentialsList := d.Get("credentials").([]interface{})

	if len(credentialsList) == 0 {
		return diag.Errorf("credentials block is required")
	}

	credentials := credentialsList[0].(map[string]interface{})
	credType := credentials["type"].(string)
	clientID := credentials["client_id"].(string)
	clientSecret := credentials["client_secret"].(string)

	// Prepare request body according to API documentation
	integration := map[string]interface{}{
		"unique_key":   uniqueKey,
		"provider":     providerName,
		"display_name": displayName,
		"credentials": map[string]interface{}{
			"type":          credType,
			"client_id":     clientID,
			"client_secret": clientSecret,
		},
	}

	// Add scopes if present
	if scopes, ok := credentials["scopes"].(string); ok && scopes != "" {
		integration["credentials"].(map[string]interface{})["scopes"] = scopes
	}

	// Log the request for debugging
	requestJSON, _ := json.MarshalIndent(integration, "", "  ")
	fmt.Printf("Request to /integrations: %s\n", string(requestJSON))

	// Make API request to create integration
	resp, err := client.MakeRequest("POST", "/integrations", integration)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	// Read and log the response body
	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	fmt.Printf("Response from /integrations (%d): %s\n", resp.StatusCode, bodyString)

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error creating integration: %s - %s", resp.Status, bodyString)
	}

	// Parse the response
	var result struct {
		Data struct {
			UniqueKey   string `json:"unique_key"`
			Provider    string `json:"provider"`
			DisplayName string `json:"display_name"`
		} `json:"data"`
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return diag.Errorf("Error parsing response: %v - Response body: %s", err, bodyString)
	}

	// Set the ID from the response
	d.SetId(result.Data.UniqueKey)

	return resourceIntegrationRead(ctx, d, m)
}

func resourceIntegrationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*NangoClient)

	// Get the unique_key
	uniqueKey := d.Id()

	// Make API request to get integration
	resp, err := client.MakeRequest("GET", fmt.Sprintf("/integrations/%s", uniqueKey), nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	// Read and log the response body
	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	fmt.Printf("Response from /integrations/%s (%d): %s\n", uniqueKey, resp.StatusCode, bodyString)

	if resp.StatusCode == http.StatusNotFound {
		// Integration was deleted outside of Terraform
		d.SetId("")
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error reading integration: %s - %s", resp.Status, bodyString)
	}

	// Parse the response
	var result struct {
		Data struct {
			UniqueKey   string `json:"unique_key"`
			Provider    string `json:"provider"`
			DisplayName string `json:"display_name"`
		} `json:"data"`
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return diag.Errorf("Error parsing response: %v - Response body: %s", err, bodyString)
	}

	// Set the resource data from the response
	if err := d.Set("unique_key", result.Data.UniqueKey); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("provider_name", result.Data.Provider); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("display_name", result.Data.DisplayName); err != nil {
		return diag.FromErr(err)
	}

	// Note: We don't set credentials as they're likely not returned by the API for security reasons

	return diag.Diagnostics{}
}

func resourceIntegrationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*NangoClient)

	// Extract values from schema
	uniqueKey := d.Get("unique_key").(string)
	providerName := d.Get("provider_name").(string)
	displayName := d.Get("display_name").(string)
	credentialsList := d.Get("credentials").([]interface{})

	if len(credentialsList) == 0 {
		return diag.Errorf("credentials block is required")
	}

	credentials := credentialsList[0].(map[string]interface{})
	credType := credentials["type"].(string)
	clientID := credentials["client_id"].(string)
	clientSecret := credentials["client_secret"].(string)

	// Prepare request body according to API documentation
	integration := map[string]interface{}{
		"unique_key":   uniqueKey,
		"provider":     providerName,
		"display_name": displayName,
		"credentials": map[string]interface{}{
			"type":          credType,
			"client_id":     clientID,
			"client_secret": clientSecret,
		},
	}

	// Add scopes if present
	if scopes, ok := credentials["scopes"].(string); ok && scopes != "" {
		integration["credentials"].(map[string]interface{})["scopes"] = scopes
	}

	// Make API request to update integration
	resp, err := client.MakeRequest("PATCH", fmt.Sprintf("/integrations/%s", uniqueKey), integration)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return diag.Errorf("Error updating integration: %s - %s", resp.Status, string(bodyBytes))
	}

	return resourceIntegrationRead(ctx, d, m)
}

func resourceIntegrationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*NangoClient)

	// Get the unique_key
	uniqueKey := d.Id()

	// Make API request to delete integration
	resp, err := client.MakeRequest("DELETE", fmt.Sprintf("/integrations/%s", uniqueKey), nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return diag.Errorf("Error deleting integration: %s - %s", resp.Status, string(bodyBytes))
	}

	// Set the ID to empty to mark it as deleted
	d.SetId("")

	return diag.Diagnostics{}
}

// DataSourceIntegration returns the data source definition for a Nango integration
func DataSourceIntegration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIntegrationRead,
		Schema: map[string]*schema.Schema{
			"unique_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique key of the integration to retrieve",
			},
			"provider_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The provider name",
			},
			"display_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The display name of the integration",
			},
		},
	}
}

func dataSourceIntegrationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*NangoClient)

	// Get the unique_key
	uniqueKey := d.Get("unique_key").(string)

	// Make API request to get integration
	resp, err := client.MakeRequest("GET", fmt.Sprintf("/integrations/%s", uniqueKey), nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return diag.Errorf("Integration with unique_key '%s' not found", uniqueKey)
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return diag.Errorf("Error reading integration: %s - %s", resp.Status, string(bodyBytes))
	}

	// Parse the response
	var result struct {
		Data struct {
			UniqueKey   string `json:"unique_key"`
			Provider    string `json:"provider"`
			DisplayName string `json:"display_name"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	// Set the data source attributes
	d.SetId(result.Data.UniqueKey)
	d.Set("provider_name", result.Data.Provider)
	d.Set("display_name", result.Data.DisplayName)

	return diag.Diagnostics{}
}
