package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"passenger-go-cli/internal/config"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
}

func GetClient() *Client {
	configuration, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Failed to load config file")
		os.Exit(1)
	}

	baseURL := configuration.ServerURL
	if baseURL == "" {
		fmt.Println("Server URL not configured, please use `passenger server <server-url>` to configure the server URL")
		os.Exit(1)
	}

	baseURL = strings.TrimSuffix(baseURL, "/")

	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseURL:    baseURL,
	}
}

func (client *Client) buildURL(
	endpoint string,
) (string, error) {
	if !strings.HasPrefix(endpoint, "/") {
		endpoint = "/" + endpoint
	}

	fullURL, err := url.JoinPath(client.baseURL, endpoint)
	if err != nil {
		return "", fmt.Errorf("failed to build URL: %w", err)
	}

	return fullURL, nil
}

func (client *Client) Get(
	endpoint string,
) (*http.Response, error) {
	fullURL, err := client.buildURL(endpoint)
	if err != nil {
		return nil, err
	}

	resp, err := client.httpClient.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("GET request failed: %w", err)
	}

	return resp, nil
}

func (client *Client) Post(
	endpoint string,
	payload any,
) (*http.Response, error) {
	fullURL, err := client.buildURL(endpoint)
	if err != nil {
		return nil, err
	}

	var body io.Reader
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest("POST", fullURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("POST request failed: %w", err)
	}

	return resp, nil
}

func (client *Client) Put(
	endpoint string,
	payload any,
) (*http.Response, error) {
	fullURL, err := client.buildURL(endpoint)
	if err != nil {
		return nil, err
	}

	var body io.Reader
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest("PUT", fullURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("PUT request failed: %w", err)
	}

	return resp, nil
}

func (client *Client) Delete(
	endpoint string,
) (*http.Response, error) {
	fullURL, err := client.buildURL(endpoint)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("DELETE", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("DELETE request failed: %w", err)
	}

	return resp, nil
}

func (client *Client) SetAuthToken(
	token string,
) *http.Request {
	req, _ := http.NewRequest("GET", "", nil)
	req.Header.Set("Cookie", "token="+token)
	return req
}

func (client *Client) DoWithAuth(
	method string,
	endpoint string,
	token string,
	payload any,
) (*http.Response, error) {
	fullURL, err := client.buildURL(endpoint)
	if err != nil {
		return nil, err
	}

	var body io.Reader
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, fullURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if token != "" {
		req.Header.Set("Cookie", "token="+token)
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s request failed: %w", method, err)
	}

	return resp, nil
}
