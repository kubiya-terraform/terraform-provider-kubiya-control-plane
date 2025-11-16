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
	logger := kubiyasentry.GetLogger()
	var bodyReader io.Reader
	var jsonBody []byte

	if body != nil {
		var err error
		jsonBody, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	fullURL := c.BaseURL + path
	req, err := http.NewRequest(method, fullURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	startTime := time.Now()
	resp, err := c.HTTPClient.Do(req)
	duration := time.Since(startTime)

	if err != nil {
		logger.Error("HTTP request failed",
			"method", method,
			"url", fullURL,
			"error", err.Error(),
		)
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// Log request details if there's an error response
	if resp.StatusCode >= 400 {
		logFile := os.Getenv("KUBIYA_API_LOG_FILE")
		if logFile == "" {
			logFile = "/tmp/kubiya_api_errors.log"
		}

		f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			defer f.Close()
			fmt.Fprintf(f, "\n========== API ERROR ==========\n")
			fmt.Fprintf(f, "Time: %s\n", time.Now().Format(time.RFC3339))
			fmt.Fprintf(f, "Method: %s\n", method)
			fmt.Fprintf(f, "URL: %s\n", fullURL)
			fmt.Fprintf(f, "Status Code: %d\n", resp.StatusCode)
			fmt.Fprintf(f, "Duration: %dms\n", duration.Milliseconds())
			fmt.Fprintf(f, "\n--- Request Headers ---\n")
			for k, v := range req.Header {
				if k != "Authorization" { // Don't log auth token
					fmt.Fprintf(f, "%s: %v\n", k, v)
				}
			}
			if len(jsonBody) > 0 {
				fmt.Fprintf(f, "\n--- Request Body ---\n%s\n", string(jsonBody))
			}
			fmt.Fprintf(f, "===============================\n\n")
		}

		logger.Error("API Error - Full Request Details",
			"method", method,
			"url", fullURL,
			"status_code", resp.StatusCode,
			"duration_ms", duration.Milliseconds(),
			"request_body", string(jsonBody),
			"log_file", logFile,
		)
	}

	return resp, nil
}

// ParseResponse parses the HTTP response into the provided interface
func ParseResponse(resp *http.Response, target interface{}) error {
	logger := kubiyasentry.GetLogger()
	defer func() { _ = resp.Body.Close() }()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		logFile := os.Getenv("KUBIYA_API_LOG_FILE")
		if logFile == "" {
			logFile = "/tmp/kubiya_api_errors.log"
		}

		f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			defer f.Close()
			fmt.Fprintf(f, "--- Response Status ---\n")
			fmt.Fprintf(f, "Status Code: %d\n", resp.StatusCode)
			fmt.Fprintf(f, "\n--- Response Headers ---\n")
			for k, v := range resp.Header {
				fmt.Fprintf(f, "%s: %v\n", k, v)
			}
			fmt.Fprintf(f, "\n--- Response Body ---\n%s\n", string(bodyBytes))
			fmt.Fprintf(f, "===============================\n\n")
		}

		logger.Error("API Error - Full Response Details",
			"status_code", resp.StatusCode,
			"response_body", string(bodyBytes),
			"response_headers", fmt.Sprintf("%v", resp.Header),
			"content_type", resp.Header.Get("Content-Type"),
			"log_file", logFile,
		)
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	if target != nil && len(bodyBytes) > 0 {
		if err := json.Unmarshal(bodyBytes, target); err != nil {
			logger.Error("Failed to parse response body",
				"error", err.Error(),
				"response_body", string(bodyBytes),
			)
			return fmt.Errorf("failed to parse response: %w", err)
		}
	}

	return nil
}
