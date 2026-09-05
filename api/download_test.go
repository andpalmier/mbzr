package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const dlHash = "cd24da0069d1d15c28538a0e9cb9610d0a47d6e7b36f00ac380ddfd0362afb93"

func zipBody() []byte { return append([]byte{'P', 'K', 3, 4}, []byte("rest-of-archive")...) }

func TestDownloadSampleWritesToOutPath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(zipBody())
	}))
	defer server.Close()

	c := NewClient("test-key")
	c.baseURL = server.URL + "/"

	dest := filepath.Join(t.TempDir(), "nested-name.zip")
	got, err := c.DownloadSample(context.Background(), dlHash, dest)
	if err != nil {
		t.Fatalf("DownloadSample() error = %v", err)
	}
	if got != dest {
		t.Errorf("returned path = %q, want %q", got, dest)
	}
	b, err := os.ReadFile(dest)
	if err != nil {
		t.Fatalf("output file missing: %v", err)
	}
	if string(b) != string(zipBody()) {
		t.Errorf("written bytes = %q", b)
	}
}

// TestDownloadSampleShortRead guards the ZIP sniff against a body that
// arrives in pieces. A plain Read may hand back fewer than 4 bytes, which
// previously made a valid archive look like an error response.
func TestDownloadSampleShortRead(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			t.Skip("ResponseWriter is not a Flusher")
		}
		full := zipBody()
		_, _ = w.Write(full[:2])
		flusher.Flush()
		time.Sleep(20 * time.Millisecond)
		_, _ = w.Write(full[2:])
		flusher.Flush()
	}))
	defer server.Close()

	c := NewClient("test-key")
	c.baseURL = server.URL + "/"

	dest := filepath.Join(t.TempDir(), "chunked.zip")
	if _, err := c.DownloadSample(context.Background(), dlHash, dest); err != nil {
		t.Fatalf("a chunked ZIP was rejected: %v", err)
	}
}

func TestDownloadSampleReportsAPIStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"query_status":"file_not_found"}`))
	}))
	defer server.Close()

	c := NewClient("test-key")
	c.baseURL = server.URL + "/"

	_, err := c.DownloadSample(context.Background(), dlHash, filepath.Join(t.TempDir(), "x.zip"))
	if err == nil {
		t.Fatal("expected an error for file_not_found")
	}
	if !strings.Contains(err.Error(), "unknown to MalwareBazaar") {
		t.Errorf("error should explain the reason, got: %v", err)
	}
}

func TestDownloadSampleRefusesToOverwrite(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(zipBody())
	}))
	defer server.Close()

	c := NewClient("test-key")
	c.baseURL = server.URL + "/"

	dest := filepath.Join(t.TempDir(), "taken.zip")
	if err := os.WriteFile(dest, []byte("existing"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := c.DownloadSample(context.Background(), dlHash, dest); err == nil {
		t.Error("expected an error when the destination already exists")
	}
}

func TestDownloadSampleRejectsBadHash(t *testing.T) {
	c := NewClient("test-key")
	if _, err := c.DownloadSample(context.Background(), "../../etc/passwd", ""); err == nil {
		t.Error("expected an error for a non-SHA256 hash")
	}
}
