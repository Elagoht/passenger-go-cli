package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"passenger-go-cli/internal/auth"
	"passenger-go-cli/internal/config"
	"passenger-go-cli/internal/schemas"
)

// ApiCaller handles HTTP requests to the Passenger Go API
type ApiCaller struct {
	baseURL string
	client  *http.Client
}

// NewApiCaller creates a new API caller instance
func NewApiCaller() (*ApiCaller, error) {
	configuration, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	if configuration.ServerURL == "" {
		return nil, fmt.Errorf("server URL not configured, use 'passenger-go server <url>' to set it")
	}

	return &ApiCaller{
		baseURL: strings.TrimSuffix(
			strings.TrimSuffix(configuration.ServerURL, "/"),
			"/api", // Remove /api given by the user
		) + "/api", // Add /api to the base URL
		client: &http.Client{},
	}, nil
}

// RequestConfig holds configuration for HTTP requests
type RequestConfig struct {
	Method      string
	Endpoint    string
	Body        any
	FilePath    string
	FileField   string
	ContentType string
}

// DoRequest performs HTTP request with generic response handling
func DoRequest[T any](endpoint string, config RequestConfig) (*T, []byte, error) {
	caller, err := NewApiCaller()
	if err != nil {
		return nil, nil, err
	}

	// Build URL
	requestURL := fmt.Sprintf("%s%s", caller.baseURL, endpoint)

	// Create request
	var request *http.Request
	if config.FilePath != "" {
		// Handle multipart/form-data for file uploads
		request, err = caller.createMultipartRequest(requestURL, config)
	} else {
		// Handle application/json
		request, err = caller.createJSONRequest(requestURL, config)
	}

	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication cookie if token is available
	err = caller.addAuthCookie(request)
	if err != nil {
		// If we can't get the token, proceed without auth
		// The API will return a meaningful error message
	}

	// Perform request
	resp, err := caller.client.Do(request)
	if err != nil {
		return nil, nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode >= 400 {
		fmt.Printf("Response: %s\n", string(body))
		var errorResponse schemas.ResponseError
		if err := json.Unmarshal(body, &errorResponse); err != nil {
			// If we can't parse the error response as JSON, use the raw response
			return nil, body, fmt.Errorf("server error (%d): %s", resp.StatusCode, string(body))
		}
		return nil, body, fmt.Errorf("%s", errorResponse.Message)
	}

	// Handle 204 No Content responses
	if resp.StatusCode == 204 {
		return nil, body, nil
	}

	// Handle response based on content type
	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		// Handle empty JSON response bodies
		if len(body) == 0 || string(body) == "" {
			return nil, body, nil
		}

		// Parse JSON response
		var result T
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, body, fmt.Errorf("failed to parse JSON response: %w", err)
		}
		return &result, body, nil
	}

	// Return raw bytes for non-JSON responses
	return nil, body, nil
}

func (client *ApiCaller) createJSONRequest(
	url string,
	config RequestConfig,
) (*http.Request, error) {
	var body io.Reader

	if config.Body != nil {
		jsonData, err := json.Marshal(config.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal JSON: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	request, err := http.NewRequest(config.Method, url, body)
	if err != nil {
		return nil, err
	}

	if config.Body != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	return request, nil
}

func (client *ApiCaller) createMultipartRequest(
	url string,
	config RequestConfig,
) (*http.Request, error) {
	file, err := os.Open(config.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	fileField := config.FileField
	if fileField == "" {
		fileField = "file"
	}

	part, err := writer.CreateFormFile(fileField, filepath.Base(config.FilePath))
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	if config.Body != nil {
		if formData, ok := config.Body.(map[string]string); ok {
			for key, value := range formData {
				writer.WriteField(key, value)
			}
		}
	}

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	req, err := http.NewRequest(config.Method, url, &body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

func (client *ApiCaller) addAuthCookie(request *http.Request) error {
	token, err := auth.GetToken()
	if err != nil {
		return err // Token not available
	}

	if token != "" {
		cookie := &http.Cookie{
			Name:  "token",
			Value: token, // Send just the token, not "Bearer <token>"
		}
		request.AddCookie(cookie)
	}

	return nil
}

func Get[T any](endpoint string) (*T, []byte, error) {
	return DoRequest[T](endpoint, RequestConfig{
		Method:   "GET",
		Endpoint: endpoint,
	})
}

func Post[T any](endpoint string, body any) (*T, []byte, error) {
	return DoRequest[T](endpoint, RequestConfig{
		Method:   "POST",
		Endpoint: endpoint,
		Body:     body,
	})
}

func PostFile[T any](
	endpoint string,
	filePath string,
	fileField string,
	formData map[string]string,
) (*T, []byte, error) {
	return DoRequest[T](endpoint, RequestConfig{
		Method:    "POST",
		Endpoint:  endpoint,
		FilePath:  filePath,
		FileField: fileField,
		Body:      formData,
	})
}

func Patch[T any](endpoint string, body any) (*T, []byte, error) {
	return DoRequest[T](endpoint, RequestConfig{
		Method:   "PATCH",
		Endpoint: endpoint,
		Body:     body,
	})
}

func Put[T any](endpoint string, body interface{}) (*T, []byte, error) {
	return DoRequest[T](endpoint, RequestConfig{
		Method:   "PUT",
		Endpoint: endpoint,
		Body:     body,
	})
}

func Delete[T any](endpoint string) (*T, []byte, error) {
	return DoRequest[T](endpoint, RequestConfig{
		Method:   "DELETE",
		Endpoint: endpoint,
	})
}
