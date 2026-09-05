package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestClient_Upload(t *testing.T) {
	// Create temporary file to upload
	tmpfile, err := os.CreateTemp("", "sample-*.exe")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(tmpfile.Name()) }()
	_, _ = tmpfile.WriteString("test content")
	_ = tmpfile.Close()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		// Check that it's multipart
		if r.Header.Get("Content-Type") == "" {
			t.Error("Expected Content-Type header")
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintln(w, `{
			"query_status": "inserted",
			"data": [
				{
					"sha256_hash": "dummy_hash"
				}
			]
		}`)
	}))
	defer server.Close()

	c := NewClient("test-key")
	c.baseURL = server.URL + "/"

	_, err = c.UploadFile(context.Background(), tmpfile.Name(), UploadOptions{Tags: []string{"tag1"}})
	if err != nil {
		t.Fatalf("UploadFile() error = %v", err)
	}
}

// TestUploadWireFormat pins the documented multipart submission format:
// the sample as "file" (under its real name) and metadata as JSON in "json_data".
func TestUploadWireFormat(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nasty-sample.exe")
	if err := os.WriteFile(path, []byte("MZ payload"), 0o600); err != nil {
		t.Fatal(err)
	}

	var gotJSON, gotFilename, gotFileBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Errorf("request is not multipart: %v", err)
			return
		}
		gotJSON = r.FormValue("json_data")

		f, hdr, err := r.FormFile("file")
		if err != nil {
			t.Errorf("no \"file\" part: %v", err)
			return
		}
		defer func() { _ = f.Close() }()
		gotFilename = hdr.Filename
		b, _ := io.ReadAll(f)
		gotFileBody = string(b)

		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintln(w, `{"query_status": "inserted"}`)
	}))
	defer server.Close()

	c := NewClient("test-key")
	c.baseURL = server.URL + "/"

	_, err := c.UploadFile(context.Background(), path, UploadOptions{
		Anonymous:      true,
		Tags:           []string{"exe", "trojan"},
		DeliveryMethod: "email_attachment",
	})
	if err != nil {
		t.Fatalf("UploadFile() error = %v", err)
	}

	if gotFilename != "nasty-sample.exe" {
		t.Errorf("filename = %q, want the sample's real basename", gotFilename)
	}
	if gotFileBody != "MZ payload" {
		t.Errorf("file body = %q", gotFileBody)
	}
	if gotJSON == "" {
		t.Fatal("no json_data part was sent")
	}

	var env struct {
		Anonymous      int      `json:"anonymous"`
		Tags           []string `json:"tags"`
		DeliveryMethod string   `json:"delivery_method"`
	}
	if err := json.Unmarshal([]byte(gotJSON), &env); err != nil {
		t.Fatalf("json_data is not valid JSON: %v (%s)", err, gotJSON)
	}
	if env.Anonymous != 1 {
		t.Errorf("anonymous = %d, want 1", env.Anonymous)
	}
	if len(env.Tags) != 2 || env.Tags[0] != "exe" || env.Tags[1] != "trojan" {
		t.Errorf("tags = %v, want a JSON array [exe trojan]", env.Tags)
	}
	if env.DeliveryMethod != "email_attachment" {
		t.Errorf("delivery_method = %q", env.DeliveryMethod)
	}
}

func TestUploadStatusHandling(t *testing.T) {
	tests := []struct {
		status  string
		wantErr bool
	}{
		{"inserted", false},
		{"file_already_known", true},
		{"no_api_key", true},
		{"user_blacklisted", true},
		{"file_expected", true},
		{"http_post_expected", true},
		{"", true},
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "s.bin")
	if err := os.WriteFile(path, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _ = fmt.Fprintf(w, `{"query_status": %q}`, tt.status)
			}))
			defer server.Close()

			c := NewClient("test-key")
			c.baseURL = server.URL + "/"

			_, err := c.UploadFile(context.Background(), path, UploadOptions{})
			if tt.wantErr && err == nil {
				t.Errorf("status %q should be reported as a failure, got success", tt.status)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("status %q should succeed, got %v", tt.status, err)
			}
		})
	}
}

// TestUploadEnvelopeReferencesAndContext pins the nested shape the API
// documents: references as lists, context mixing lists and a plain comment.
func TestUploadEnvelopeReferencesAndContext(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "s.bin")
	if err := os.WriteFile(path, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}

	var gotJSON string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Errorf("not multipart: %v", err)
			return
		}
		gotJSON = r.FormValue("json_data")
		_, _ = fmt.Fprintln(w, `{"query_status": "inserted"}`)
	}))
	defer server.Close()

	c := NewClient("test-key")
	c.baseURL = server.URL + "/"

	_, err := c.UploadFile(context.Background(), path, UploadOptions{
		References: map[string][]string{"any_run": {"https://app.any.run/tasks/1", "https://app.any.run/tasks/2"}},
		Context:    map[string]any{"dropped_by_malware": []string{"Gozi"}, "comment": "very nasty"},
	})
	if err != nil {
		t.Fatalf("UploadFile() error = %v", err)
	}

	var env struct {
		References map[string][]string `json:"references"`
		Context    map[string]any      `json:"context"`
	}
	if err := json.Unmarshal([]byte(gotJSON), &env); err != nil {
		t.Fatalf("json_data invalid: %v (%s)", err, gotJSON)
	}
	if len(env.References["any_run"]) != 2 {
		t.Errorf("references.any_run = %v, want 2 entries", env.References["any_run"])
	}
	if _, ok := env.Context["comment"].(string); !ok {
		t.Errorf("context.comment should be a string, got %T", env.Context["comment"])
	}
	if _, ok := env.Context["dropped_by_malware"].([]any); !ok {
		t.Errorf("context.dropped_by_malware should be a list, got %T", env.Context["dropped_by_malware"])
	}
}

// An upload with no metadata must not emit empty references/context blocks.
func TestUploadEnvelopeOmitsEmptyBlocks(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "s.bin")
	if err := os.WriteFile(path, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}

	var gotJSON string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseMultipartForm(1 << 20)
		gotJSON = r.FormValue("json_data")
		_, _ = fmt.Fprintln(w, `{"query_status": "inserted"}`)
	}))
	defer server.Close()

	c := NewClient("test-key")
	c.baseURL = server.URL + "/"

	if _, err := c.UploadFile(context.Background(), path, UploadOptions{}); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"references", "context", "tags", "delivery_method"} {
		if strings.Contains(gotJSON, key) {
			t.Errorf("json_data should omit empty %q, got %s", key, gotJSON)
		}
	}
}
