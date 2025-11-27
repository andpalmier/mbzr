package api

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// defaultAPIURL is the default MalwareBazaar API endpoint
const defaultAPIURL = "https://mb-api.abuse.ch/api/v1/"

// RateLimiter implements a simple token bucket rate limiter
type RateLimiter struct {
	tokens chan struct{}
	ticker *time.Ticker
	stop   chan struct{}
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(requestsPerSecond int) *RateLimiter {
	rl := &RateLimiter{
		tokens: make(chan struct{}, requestsPerSecond),
		ticker: time.NewTicker(time.Second / time.Duration(requestsPerSecond)),
		stop:   make(chan struct{}),
	}

	// Fill initial tokens
	for i := 0; i < requestsPerSecond; i++ {
		rl.tokens <- struct{}{}
	}

	// Start refill loop
	go func() {
		for {
			select {
			case <-rl.ticker.C:
				select {
				case rl.tokens <- struct{}{}:
				default:
				}
			case <-rl.stop:
				rl.ticker.Stop()
				return
			}
		}
	}()

	return rl
}

// Wait blocks until a token is available
func (rl *RateLimiter) Wait(ctx context.Context) error {
	select {
	case <-rl.tokens:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Stop stops the rate limiter
func (rl *RateLimiter) Stop() {
	close(rl.stop)
}

// Client interacts with the MalwareBazaar API
type Client struct {
	apiKey      string
	baseURL     string
	httpClient  *http.Client
	rateLimiter *RateLimiter
}

// Option configures the Client
type Option func(*Client)

// WithTimeout sets the HTTP client timeout
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// WithBaseURL sets the API base URL
func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
	}
}

// NewClient creates a new MalwareBazaar API client
func NewClient(apiKey string, options ...Option) *Client {
	c := &Client{
		apiKey:      apiKey,
		baseURL:     defaultAPIURL,
		rateLimiter: NewRateLimiter(10), // 10 requests per second
		httpClient: &http.Client{
			Timeout: 30 * time.Second, // Default timeout for security
		},
	}
	for _, opt := range options {
		opt(c)
	}
	return c
}

// buildRequest creates an HTTP request with the given data and files
func (c *Client) buildRequest(ctx context.Context, data map[string]string, files map[string]io.Reader) (*http.Request, error) {
	var req *http.Request
	var err error

	if files != nil {
		// Handle file upload
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		for key, r := range files {
			var fw io.Writer
			if key == "file" {
				// file upload
				if fw, err = writer.CreateFormFile(key, "file"); err != nil {
					return nil, fmt.Errorf("error creating form file: %w", err)
				}
			} else {
				// other form fields
				if fw, err = writer.CreateFormField(key); err != nil {
					return nil, fmt.Errorf("error creating form field: %w", err)
				}
			}
			if _, err = io.Copy(fw, r); err != nil {
				return nil, fmt.Errorf("error copying data: %w", err)
			}
		}
		writer.Close()
		req, err = http.NewRequestWithContext(ctx, "POST", c.baseURL, body)
		if err != nil {
			return nil, fmt.Errorf("error creating request: %w", err)
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())
	} else {
		// Handle form data
		formData := url.Values{}
		for key, value := range data {
			formData.Add(key, value)
		}
		req, err = http.NewRequestWithContext(ctx, "POST", c.baseURL, strings.NewReader(formData.Encode()))
		if err != nil {
			return nil, fmt.Errorf("error creating request: %w", err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	req.Header.Set("User-Agent", "mbzr-client/1.0")

	if c.apiKey != "" {
		req.Header.Set("Auth-Key", c.apiKey)
	}

	return req, nil
}

// MakeRequest makes an HTTP request to the API and returns the response as a string
func (c *Client) MakeRequest(ctx context.Context, data map[string]string, files map[string]io.Reader) (string, error) {
	// Wait for rate limiter
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return "", fmt.Errorf("rate limiter error: %w", err)
	}

	req, err := c.buildRequest(ctx, data, files)
	if err != nil {
		return "", err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned non-OK status: %s", resp.Status)
	}

	// Limit response size to 10MB to prevent OOM attacks
	const maxResponseSize = 10 * 1024 * 1024 // 10MB
	limitedReader := io.LimitReader(resp.Body, maxResponseSize)

	body, err := io.ReadAll(limitedReader)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	// Check if we hit the limit
	if len(body) == maxResponseSize {
		return "", fmt.Errorf("response too large: exceeded %d bytes", maxResponseSize)
	}

	return string(body), nil
}

// MakeRequestRaw makes an HTTP request and returns the raw response body.
// The caller is responsible for closing the response body.
func (c *Client) MakeRequestRaw(ctx context.Context, data map[string]string, files map[string]io.Reader) (io.ReadCloser, error) {
	// Wait for rate limiter
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limiter error: %w", err)
	}

	req, err := c.buildRequest(ctx, data, files)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("API returned non-OK status: %s", resp.Status)
	}

	return resp.Body, nil
}
