package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
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
		// Provide meaningful error messages for common network issues
		return client.wrapNetworkError(err, request.URL.Host)
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
			// If we can't parse as JSON, provide a clean error message based on status code
			return client.formatHTTPError(response.StatusCode, request.URL.Path)
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

// wrapNetworkError provides user-friendly error messages for network issues
func (client *Client) wrapNetworkError(err error, host string) error {
	if err == nil {
		return nil
	}

	errMsg := err.Error()

	// DNS resolution failure
	if strings.Contains(errMsg, "no such host") {
		return fmt.Errorf("❌ Cannot reach server at '%s'. Please check:\n"+
			"  • The server URL is correct\n"+
			"  • You have internet connectivity\n"+
			"  • The domain exists and is accessible\n\n"+
			"Use 'passenger server <correct-url>' to update the server URL", host)
	}

	// Connection refused
	if strings.Contains(errMsg, "connection refused") {
		return fmt.Errorf("❌ Server at '%s' refused the connection. Please check:\n"+
			"  • The server is running\n"+
			"  • The port is correct\n"+
			"  • No firewall is blocking the connection", host)
	}

	// Timeout errors
	if strings.Contains(errMsg, "timeout") || strings.Contains(errMsg, "context deadline exceeded") {
		return fmt.Errorf("❌ Request to '%s' timed out. Please check:\n"+
			"  • Your internet connection\n"+
			"  • The server is responding\n"+
			"  • Try again in a moment", host)
	}

	// SSL/TLS errors
	if strings.Contains(errMsg, "certificate") || strings.Contains(errMsg, "tls") {
		return fmt.Errorf("❌ SSL/TLS error connecting to '%s'. Please check:\n"+
			"  • The server has a valid SSL certificate\n"+
			"  • The URL uses the correct protocol (https/http)", host)
	}

	// Network unreachable
	if strings.Contains(errMsg, "network is unreachable") {
		return fmt.Errorf("❌ Network unreachable. Please check your internet connection")
	}

	// Parse URL to see if it's a URL format issue
	if _, urlErr := url.Parse(client.BaseURL); urlErr != nil {
		return fmt.Errorf("❌ Invalid server URL format: '%s'\n"+
			"Please use a valid URL like: https://example.com or http://localhost:8080\n"+
			"Use 'passenger server <correct-url>' to update the server URL", client.BaseURL)
	}

	// Check if the error is a network error
	if netErr, ok := err.(net.Error); ok {
		if netErr.Timeout() {
			return fmt.Errorf("❌ Network timeout connecting to '%s'. Please try again", host)
		}
		return fmt.Errorf("❌ Network error connecting to '%s': %s", host, netErr.Error())
	}

	// Default network error message
	return fmt.Errorf("❌ Failed to connect to server at '%s': %s\n"+
		"Use 'passenger server <correct-url>' to update the server URL if needed", host, err.Error())
}

// formatHTTPError provides user-friendly error messages for HTTP status codes
func (client *Client) formatHTTPError(statusCode int, path string) error {
	switch statusCode {
	case 400:
		return fmt.Errorf("❌ Bad request to '%s'. Please check your input parameters", path)
	case 401:
		return fmt.Errorf("❌ Authentication required. Please login first with 'passenger login'")
	case 403:
		return fmt.Errorf("❌ Access denied. You don't have permission to access '%s'", path)
	case 404:
		return fmt.Errorf("❌ Endpoint '%s' not found. Please check:\n"+
			"  • The server URL is correct\n"+
			"  • You're using a compatible Passenger Go server\n"+
			"  • The API endpoint exists on this server", path)
	case 405:
		return fmt.Errorf("❌ Method not allowed for '%s'", path)
	case 409:
		return fmt.Errorf("❌ Conflict: The request could not be completed due to a conflict")
	case 422:
		return fmt.Errorf("❌ Invalid data provided. Please check your input")
	case 429:
		return fmt.Errorf("❌ Too many requests. Please wait before trying again")
	case 500:
		return fmt.Errorf("❌ Server error. The Passenger Go server encountered an internal error")
	case 502:
		return fmt.Errorf("❌ Bad gateway. The server is temporarily unavailable")
	case 503:
		return fmt.Errorf("❌ Service unavailable. The server is temporarily down for maintenance")
	case 504:
		return fmt.Errorf("❌ Gateway timeout. The server took too long to respond")
	default:
		if statusCode >= 400 && statusCode < 500 {
			return fmt.Errorf("❌ Client error (HTTP %d) accessing '%s'", statusCode, path)
		} else if statusCode >= 500 {
			return fmt.Errorf("❌ Server error (HTTP %d). Please try again later", statusCode)
		}
		return fmt.Errorf("❌ Unexpected HTTP status %d for '%s'", statusCode, path)
	}
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
