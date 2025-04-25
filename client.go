package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents an Aidbox API client
type Client struct {
	URL          string
	ClientID     string
	ClientSecret string
	HTTPClient   *http.Client
}

// NewClient creates a new Aidbox API client
func NewClient(config *Config) *Client {
	return &Client{
		URL:          config.URL,
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		HTTPClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

// CreateResource creates a new resource in Aidbox
func (c *Client) CreateResource(resourceType, id string, resourceJSON string) error {
	url := fmt.Sprintf("%s/%s/%s", c.URL, resourceType, id)

	req, err := http.NewRequest("PUT", url, bytes.NewBufferString(resourceJSON))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.ClientID, c.ClientSecret)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error creating resource: status %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetResource retrieves a resource from Aidbox
func (c *Client) GetResource(resourceType, id string) (string, error) {
	url := fmt.Sprintf("%s/%s/%s", c.URL, resourceType, id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.SetBasicAuth(c.ClientID, c.ClientSecret)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", nil
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("error getting resource: status %d, body: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	return string(body), nil
}

// UpdateResource updates an existing resource in Aidbox
func (c *Client) UpdateResource(resourceType, id string, resourceJSON string) error {
	return c.CreateResource(resourceType, id, resourceJSON)
}

// DeleteResource deletes a resource from Aidbox
func (c *Client) DeleteResource(resourceType, id string) error {
	url := fmt.Sprintf("%s/%s/%s", c.URL, resourceType, id)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.SetBasicAuth(c.ClientID, c.ClientSecret)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error deleting resource: status %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}
