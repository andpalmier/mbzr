package api

import (
	"context"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	apiKey := "test-api-key"
	client := NewClient(apiKey)

	if client == nil {
		t.Fatal("NewClient returned nil")
	}

	if client.apiKey != apiKey {
		t.Errorf("Expected API key %q, got %q", apiKey, client.apiKey)
	}

	if client.baseURL != defaultAPIURL {
		t.Errorf("Expected base URL %q, got %q", defaultAPIURL, client.baseURL)
	}

	if client.httpClient == nil {
		t.Error("HTTP client should not be nil")
	}

	if client.httpClient.Timeout != 30*time.Second {
		t.Errorf("Expected timeout 30s, got %v", client.httpClient.Timeout)
	}
}

func TestNewClientWithOptions(t *testing.T) {
	apiKey := "test-api-key"
	customTimeout := 60 * time.Second
	customURL := "https://custom.api.url/"

	client := NewClient(apiKey,
		WithTimeout(customTimeout),
		WithBaseURL(customURL),
	)

	if client.httpClient.Timeout != customTimeout {
		t.Errorf("Expected timeout %v, got %v", customTimeout, client.httpClient.Timeout)
	}

	if client.baseURL != customURL {
		t.Errorf("Expected base URL %q, got %q", customURL, client.baseURL)
	}
}

func TestWithTimeout(t *testing.T) {
	client := NewClient("test-key")
	customTimeout := 45 * time.Second

	opt := WithTimeout(customTimeout)
	opt(client)

	if client.httpClient.Timeout != customTimeout {
		t.Errorf("Expected timeout %v, got %v", customTimeout, client.httpClient.Timeout)
	}
}

func TestWithBaseURL(t *testing.T) {
	client := NewClient("test-key")
	customURL := "https://test.example.com/"

	opt := WithBaseURL(customURL)
	opt(client)

	if client.baseURL != customURL {
		t.Errorf("Expected base URL %q, got %q", customURL, client.baseURL)
	}
}

func TestBuildRequest(t *testing.T) {
	client := NewClient("test-key")
	ctx := context.Background()

	data := map[string]string{
		"query": "get_info",
		"hash":  "abc123",
	}

	req, err := client.buildRequest(ctx, data, nil)
	if err != nil {
		t.Fatalf("buildRequest failed: %v", err)
	}

	if req == nil {
		t.Fatal("buildRequest returned nil request")
	}

	if req.Method != "POST" {
		t.Errorf("Expected POST method, got %s", req.Method)
	}

	if req.URL.String() != client.baseURL {
		t.Errorf("Expected URL %q, got %q", client.baseURL, req.URL.String())
	}

	// Check User-Agent header
	userAgent := req.Header.Get("User-Agent")
	if userAgent != "mbzr-client/1.0" {
		t.Errorf("Expected User-Agent 'mbzr-client/1.0', got %q", userAgent)
	}

	// Check Auth-Key header
	authKey := req.Header.Get("Auth-Key")
	if authKey != "test-key" {
		t.Errorf("Expected Auth-Key 'test-key', got %q", authKey)
	}
}

func TestBuildRequestWithoutAPIKey(t *testing.T) {
	client := NewClient("")
	ctx := context.Background()

	data := map[string]string{
		"query": "get_info",
	}

	req, err := client.buildRequest(ctx, data, nil)
	if err != nil {
		t.Fatalf("buildRequest failed: %v", err)
	}

	// Auth-Key header should not be set when API key is empty
	authKey := req.Header.Get("Auth-Key")
	if authKey != "" {
		t.Errorf("Expected no Auth-Key header, got %q", authKey)
	}
}

func TestBuildRequestWithContext(t *testing.T) {
	client := NewClient("test-key")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data := map[string]string{
		"query": "get_info",
	}

	req, err := client.buildRequest(ctx, data, nil)
	if err != nil {
		t.Fatalf("buildRequest failed: %v", err)
	}

	// Check that request has the context
	if req.Context() != ctx {
		t.Error("Request context does not match provided context")
	}
}
