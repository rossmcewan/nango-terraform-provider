package resources

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// NangoClient handles API communication with Nango
type NangoClient struct {
	APIKey  string
	BaseURL string
	Client  *http.Client
}

// NewClient creates a new Nango API client
func NewClient(apiKey, baseURL string) *NangoClient {
	return &NangoClient{
		APIKey:  apiKey,
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}

// MakeRequest makes an HTTP request to the Nango API
func (c *NangoClient) MakeRequest(method, path string, body interface{}) (*http.Response, error) {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, err
		}
	}

	url := fmt.Sprintf("%s%s", c.BaseURL, path)
	req, err := http.NewRequest(method, url, &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))

	return c.Client.Do(req)
}
