package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func TestProvider_configure(t *testing.T) {
	provider := Provider()

	// Test with valid configuration
	config := map[string]interface{}{
		"api_key":  "test-api-key",
		"base_url": "https://api.nango.dev",
	}

	resourceData := schema.TestResourceDataRaw(t, provider.Schema, config)
	_, diags := providerConfigure(nil, resourceData)

	if diags.HasError() {
		t.Fatalf("Expected no errors, got: %v", diags)
	}
}

func TestProvider_configure_invalid_url(t *testing.T) {
	provider := Provider()

	// Test with invalid URL
	config := map[string]interface{}{
		"api_key":  "test-api-key",
		"base_url": "invalid-url",
	}

	resourceData := schema.TestResourceDataRaw(t, provider.Schema, config)
	_, diags := providerConfigure(nil, resourceData)

	if !diags.HasError() {
		t.Fatal("Expected error for invalid URL, got none")
	}
}
