package resources

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// NangoClient is the client for interacting with the Nango API
type NangoClient struct {
	APIKey  string
	BaseURL string
	Client  *http.Client
}

// NewNangoClient creates a new Nango API client
func NewNangoClient(apiKey, baseURL string) *NangoClient {
	// Create an HTTP client with reasonable timeouts
	client := &http.Client{
		Timeout: time.Second * 30,
	}

	return &NangoClient{
		APIKey:  apiKey,
		BaseURL: baseURL,
		Client:  client, // Initialize the HTTP client
	}
}

// MakeRequest makes a request to the Nango API
func (c *NangoClient) MakeRequest(method, path string, body interface{}) (*http.Response, error) {
	// Ensure the client is initialized
	if c.Client == nil {
		c.Client = &http.Client{
			Timeout: time.Second * 30,
		}
	}

	// Construct the full URL
	url := fmt.Sprintf("%s%s", c.BaseURL, path)

	// Marshal the body to JSON if it's not nil
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	// Create the request
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))

	// Make the request
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	return resp, nil
}
