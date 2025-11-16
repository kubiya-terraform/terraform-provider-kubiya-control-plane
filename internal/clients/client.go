package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	kubiyasentry "kubiya-control-plane/internal/sentry"

	"github.com/getsentry/sentry-go"
)

type Client struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
}

// New creates a new Control Plane API client
func New(apiKey string) (*Client, error) {
	// Get logger
	logger := kubiyasentry.GetLogger()

	if apiKey == "" {
		logger.Error("Failed to create client", "error", "API key is required")
		return nil, fmt.Errorf("API key is required")
	}

	baseURL := getBaseURL()

	// Create HTTP client with Sentry tracing transport
	httpClient := &http.Client{
		Timeout:   60 * time.Second,
		Transport: kubiyasentry.NewHTTPTransport(http.DefaultTransport),
	}

	client := &Client{
		APIKey:     apiKey,
		BaseURL:    baseURL,
		HTTPClient: httpClient,
	}

	logger.Info("Created Kubiya Control Plane client",
		"base_url", baseURL,
	)

	kubiyasentry.AddBreadcrumb("client", "Kubiya Control Plane client created", sentry.LevelInfo, map[string]interface{}{
		"base_url": baseURL,
	})

	return client, nil
}

// getBaseURL returns the base URL for the Control Plane API
// Default: https://control-plane.kubiya.ai
// Override with KUBIYA_CONTROL_PLANE_BASE_URL environment variable
func getBaseURL() string {
	// Check for custom base URL from environment variable
	if customURL := os.Getenv("KUBIYA_CONTROL_PLANE_BASE_URL"); customURL != "" {
		return customURL
	}

	// Default production URL
	return "https://control-plane.kubiya.ai"
}

// DoRequest performs an HTTP request with proper headers
func (c *Client) DoRequest(method, path string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, c.BaseURL+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

// ParseResponse parses the HTTP response into the provided interface
func ParseResponse(resp *http.Response, target interface{}) error {
	defer func() { _ = resp.Body.Close() }()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	if target != nil && len(bodyBytes) > 0 {
		if err := json.Unmarshal(bodyBytes, target); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}
	}

	return nil
}
