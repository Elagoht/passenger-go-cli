package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"passenger-go-cli/internal/config"
	"passenger-go-cli/internal/schemas"
)

// Client represents a JSON-only HTTP client
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Token      string
	Config     *config.Config
}

// ClientOption defines options for configuring the HTTP client
type ClientOption func(*Client)

// WithTimeout sets a custom timeout for the HTTP client
func WithTimeout(timeout time.Duration) ClientOption {
	return func(client *Client) {
		client.HTTPClient.Timeout = timeout
	}
}

// WithToken sets the authentication token
func WithToken(token string) ClientOption {
	return func(client *Client) {
		client.Token = token
	}
}

// WithBaseURL sets a custom base URL
func WithBaseURL(baseURL string) ClientOption {
	return func(client *Client) {
		client.BaseURL = baseURL
	}
}

// NewClient creates a new JSON HTTP client with optional configuration
func NewClient(options ...ClientOption) (*Client, error) {
	configuration, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	client := &Client{
		BaseURL: configuration.ServerURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		Config: configuration,
	}

	// Apply options
	for _, opt := range options {
		opt(client)
	}

	return client, nil
}

// buildURL constructs the full URL for an endpoint
func (client *Client) buildURL(endpoint string) (string, error) {
	if client.BaseURL == "" {
		return "", fmt.Errorf("base URL is not configured")
	}

	baseURL, err := url.Parse(client.BaseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}

	endpointURL, err := url.Parse(endpoint)
	if err != nil {
		return "", fmt.Errorf("invalid endpoint: %w", err)
	}

	return baseURL.ResolveReference(endpointURL).String(), nil
}

// prepareRequest creates and configures an HTTP request
func (client *Client) prepareRequest(
	method string,
	endpoint string,
	body any,
) (*http.Request, error) {
	fullURL, err := client.buildURL(endpoint)
	if err != nil {
		return nil, err
	}

	var requestBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		requestBody = bytes.NewBuffer(jsonData)
	}

	request, err := http.NewRequest(method, fullURL, requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	// Add authorization header if token is available
	if client.Token != "" {
		request.Header.Set("Authorization", "Bearer "+client.Token)
	}

	return request, nil
}

// executeRequest executes an HTTP request and handles the response
func (client *Client) executeRequest(
	request *http.Request,
	result any,
) error {
	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Handle error responses
	if response.StatusCode >= 400 {
		var apiError schemas.ResponseError
		if err := json.Unmarshal(body, &apiError); err != nil {
			return fmt.Errorf("HTTP %d: %s", response.StatusCode, string(body))
		}
		return fmt.Errorf("API error (%d): %s", response.StatusCode, apiError.Message)
	}

	// Parse successful response
	if result != nil && len(body) > 0 {
		if err := json.Unmarshal(body, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// Get performs a GET request
func (client *Client) Get(
	endpoint string,
	result any,
) error {
	request, err := client.prepareRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}

	return client.executeRequest(request, result)
}

// Post performs a POST request
func (client *Client) Post(
	endpoint string,
	body any,
	result any,
) error {
	request, err := client.prepareRequest(http.MethodPost, endpoint, body)
	if err != nil {
		return err
	}

	return client.executeRequest(request, result)
}

// Put performs a PUT request
func (client *Client) Put(
	endpoint string,
	body any,
	result any,
) error {
	request, err := client.prepareRequest(http.MethodPut, endpoint, body)
	if err != nil {
		return err
	}

	return client.executeRequest(request, result)
}

// Patch performs a PATCH request
func (client *Client) Patch(
	endpoint string,
	body any,
	result any,
) error {
	request, err := client.prepareRequest(http.MethodPatch, endpoint, body)
	if err != nil {
		return err
	}

	return client.executeRequest(request, result)
}

// Delete performs a DELETE request
func (client *Client) Delete(
	endpoint string,
	result any,
) error {
	request, err := client.prepareRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}

	return client.executeRequest(request, result)
}

// PostWithQuery performs a POST request with query parameters
func (client *Client) PostWithQuery(
	endpoint string,
	queryParams url.Values,
	body any,
	result any,
) error {
	if len(queryParams) > 0 {
		endpoint += "?" + queryParams.Encode()
	}

	req, err := client.prepareRequest(http.MethodPost, endpoint, body)
	if err != nil {
		return err
	}

	return client.executeRequest(req, result)
}

// GetWithQuery performs a GET request with query parameters
func (client *Client) GetWithQuery(
	endpoint string,
	queryParams url.Values,
	result any,
) error {
	if len(queryParams) > 0 {
		endpoint += "?" + queryParams.Encode()
	}

	req, err := client.prepareRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}

	return client.executeRequest(req, result)
}

// SetToken updates the authentication token
func (client *Client) SetToken(token string) {
	client.Token = token
}

// SetBaseURL updates the base URL
func (client *Client) SetBaseURL(baseURL string) {
	client.BaseURL = baseURL
}

// GetBaseURL returns the current base URL
func (client *Client) GetBaseURL() string {
	return client.BaseURL
}

// IsConfigured returns true if the client has a base URL configured
func (client *Client) IsConfigured() bool {
	return client.BaseURL != ""
}
