package main

import (
	"bytes"
	"encoding/json"
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
	accessToken  string
}

// NewClient creates a new Aidbox API client
func NewClient(config *Config) *Client {
	client := &Client{
		URL:          config.URL,
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		HTTPClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}
	if err := client.acquireToken(); err != nil {
		panic(fmt.Sprintf("failed to acquire Aidbox token: %v", err))
	}
	return client
}

// acquireToken fetches an OAuth2 token using client credentials grant
func (c *Client) acquireToken() error {
	form := []byte("grant_type=client_credentials")
	req, err := http.NewRequest("POST", c.URL+"/auth/token", bytes.NewBuffer(form))
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.ClientID, c.ClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to get token: status %d, body: %s", resp.StatusCode, string(body))
	}

	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}
	if result.AccessToken == "" {
		return fmt.Errorf("no access_token in response")
	}
	c.accessToken = result.AccessToken
	return nil
}

// CreateResource creates a new resource in Aidbox
func (c *Client) CreateResource(resourceType, id string, resourceJSON string) error {
	url := fmt.Sprintf("%s/%s/%s", c.URL, resourceType, id)

	req, err := http.NewRequest("PUT", url, bytes.NewBufferString(resourceJSON))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.accessToken)

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

	req.Header.Set("Authorization", "Bearer "+c.accessToken)

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

	req.Header.Set("Authorization", "Bearer "+c.accessToken)

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
