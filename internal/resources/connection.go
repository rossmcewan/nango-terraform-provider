package resources

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceConnection returns the resource definition for a Nango connection
func ResourceConnection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConnectionCreate,
		ReadContext:   resourceConnectionRead,
		UpdateContext: resourceConnectionUpdate,
		DeleteContext: resourceConnectionDelete,
		Schema: map[string]*schema.Schema{
			"connection_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier for the connection",
			},
			"integration_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the integration this connection belongs to",
			},
			"environment": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "production",
				Description: "The environment for this connection (e.g., production, development)",
			},
			"credentials": {
				Type:        schema.TypeMap,
				Optional:    true,
				Sensitive:   true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Credentials for the connection (for non-OAuth connections)",
			},
			// Add more fields as needed
		},
	}
}

func resourceConnectionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// client := m.(*NangoClient)

	// Extract values from schema
	connectionID := d.Get("connection_id").(string)
	integrationName := d.Get("integration_name").(string)
	environment := d.Get("environment").(string)

	// Prepare request body
	connection := map[string]interface{}{
		"connection_id":    connectionID,
		"integration_name": integrationName,
		"environment":      environment,
	}

	// Add credentials if present
	if v, ok := d.GetOk("credentials"); ok {
		credentials := make(map[string]string)
		for k, v := range v.(map[string]interface{}) {
			credentials[k] = v.(string)
		}
		connection["credentials"] = credentials
	}

	// Make API request to create connection
	// Implementation depends on Nango API

	// Set ID for the resource
	d.SetId(fmt.Sprintf("%s-%s-%d", integrationName, connectionID, time.Now().Unix()))

	return resourceConnectionRead(ctx, d, m)
}

func resourceConnectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Implementation to read connection from Nango API
	return diag.Diagnostics{}
}

func resourceConnectionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Implementation to update connection in Nango API
	return resourceConnectionRead(ctx, d, m)
}

func resourceConnectionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Implementation to delete connection from Nango API
	return diag.Diagnostics{}
}

// DataSourceConnection returns the data source definition for a Nango connection
func DataSourceConnection() *schema.Resource {
	// Similar to ResourceConnection but for data source
	return &schema.Resource{
		ReadContext: dataSourceConnectionRead,
		Schema: map[string]*schema.Schema{
			"connection_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier for the connection",
			},
			"integration_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the integration this connection belongs to",
			},
			"environment": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "production",
				Description: "The environment for this connection",
			},
			// Add more fields as needed
		},
	}
}

func dataSourceConnectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Implementation to read connection from Nango API for data source
	return diag.Diagnostics{}
}
