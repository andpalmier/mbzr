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

// API constants
const (
	defaultAPIURL   = "https://mb-api.abuse.ch/api/v1/"
	defaultTimeout  = 30 * time.Second
	maxResponseSize = 10 * 1024 * 1024 // prevents OOM from large responses (10MB)
)

// Client interacts with the MalwareBazaar API
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	interval   time.Duration
	lastReq    time.Time
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
// Note: API key is required
func NewClient(apiKey string, options ...Option) *Client {
	c := &Client{
		apiKey:   apiKey,
		baseURL:  defaultAPIURL,
		interval: 100 * time.Millisecond, // 10 requests per second
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}

	for _, opt := range options {
		opt(c)
	}

	return c
}

// wait handles simple rate limiting
func (c *Client) wait() {
	elapsed := time.Since(c.lastReq)
	if elapsed < c.interval {
		time.Sleep(c.interval - elapsed)
	}
	c.lastReq = time.Now()
}

// buildRequest creates an HTTP request with the given data and files
func (c *Client) buildRequest(ctx context.Context, data map[string]string, files map[string]io.Reader) (*http.Request, error) {
	var body io.Reader
	var contentType string

	if files != nil {
		buf := &bytes.Buffer{}
		writer := multipart.NewWriter(buf)
		for key, r := range files {
			var fw io.Writer
			var err error
			if key == "file" {
				fw, err = writer.CreateFormFile(key, "file")
			} else {
				fw, err = writer.CreateFormField(key)
			}
			if err != nil {
				return nil, fmt.Errorf("creating form field %q: %w", key, err)
			}
			if _, err = io.Copy(fw, r); err != nil {
				return nil, fmt.Errorf("copying data for %q: %w", key, err)
			}
		}
		if err := writer.Close(); err != nil {
			return nil, fmt.Errorf("closing multipart writer: %w", err)
		}
		body = buf
		contentType = writer.FormDataContentType()
	} else {
		formData := url.Values{}
		for key, value := range data {
			formData.Add(key, value)
		}
		body = strings.NewReader(formData.Encode())
		contentType = "application/x-www-form-urlencoded"
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("User-Agent", "mbzr-client/1.0")
	if c.apiKey != "" {
		req.Header.Set("Auth-Key", c.apiKey)
	}

	return req, nil
}

// MakeRequest makes an HTTP request to the API and returns the response as a string
func (c *Client) MakeRequest(ctx context.Context, data map[string]string, files map[string]io.Reader) (string, error) {
	c.wait()

	req, err := c.buildRequest(ctx, data, files)
	if err != nil {
		return "", err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status %s", resp.Status)
	}

	limitedReader := io.LimitReader(resp.Body, maxResponseSize)
	body, err := io.ReadAll(limitedReader)
	if err != nil {
		return "", fmt.Errorf("reading response: %w", err)
	}

	if len(body) == maxResponseSize {
		return "", fmt.Errorf("response too large: exceeded %d bytes", maxResponseSize)
	}

	return string(body), nil
}

// MakeRequestRaw makes an HTTP request and returns the raw response body
// The caller is responsible for closing the response body
func (c *Client) MakeRequestRaw(ctx context.Context, data map[string]string, files map[string]io.Reader) (io.ReadCloser, error) {
	c.wait()

	req, err := c.buildRequest(ctx, data, files)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("API returned status %s", resp.Status)
	}

	return resp.Body, nil
}
