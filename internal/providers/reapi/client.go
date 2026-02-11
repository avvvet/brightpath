package reapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	baseURL        = "https://api.realestateapi.com"
	defaultTimeout = 30 * time.Second
	rateLimitDelay = 500 * time.Millisecond // 500ms between requests
)

type Client struct {
	apiKey     string
	httpClient *http.Client
	lastCall   time.Time
}

// NewClient creates a new REAPI client
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		lastCall: time.Time{},
	}
}

// rateLimit ensures we don't exceed API rate limits
func (c *Client) rateLimit() {
	if !c.lastCall.IsZero() {
		elapsed := time.Since(c.lastCall)
		if elapsed < rateLimitDelay {
			time.Sleep(rateLimitDelay - elapsed)
		}
	}
	c.lastCall = time.Now()
}

// doRequest performs an HTTP request with authentication and rate limiting
func (c *Client) doRequest(method, path string, body interface{}) ([]byte, error) {
	c.rateLimit()

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, baseURL+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)

	// Execute request with retry logic
	var resp *http.Response
	maxRetries := 3
	for attempt := 1; attempt <= maxRetries; attempt++ {
		resp, err = c.httpClient.Do(req)
		if err != nil {
			if attempt == maxRetries {
				return nil, fmt.Errorf("request failed after %d attempts: %w", maxRetries, err)
			}
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}
		break
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for API errors
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}
