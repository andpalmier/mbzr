package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDownloadExportWritesFile(t *testing.T) {
	original := exportsURL
	defer func() { exportsURL = original }()

	var gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_, _ = fmt.Fprint(w, "recent-export")
	}))
	defer server.Close()
	exportsURL = server.URL + "/"

	c := NewClient("test-key")
	dest := filepath.Join(t.TempDir(), "recent.csv")
	got, err := c.DownloadExport(context.Background(), "recent.csv", dest)
	if err != nil {
		t.Fatal(err)
	}
	if got != dest {
		t.Errorf("path = %q, want %q", got, dest)
	}
	if gotPath != "/test-key/recent.csv" {
		t.Errorf("request path = %q", gotPath)
	}
	b, err := os.ReadFile(dest)
	if err != nil || string(b) != "recent-export" {
		t.Errorf("file = %q, %v", b, err)
	}
}

func TestDownloadExportRejectsExistingAndInvalidName(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "recent.csv")
	if err := os.WriteFile(dest, []byte("existing"), 0o600); err != nil {
		t.Fatal(err)
	}
	c := NewClient("test-key")
	if _, err := c.DownloadExport(context.Background(), "recent.csv", dest); err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("existing file error = %v", err)
	}
	if _, err := c.DownloadExport(context.Background(), "../recent.csv", filepath.Join(dir, "x")); err == nil {
		t.Error("path traversal export name accepted")
	}
}

func TestDownloadExportRemovesPartialFile(t *testing.T) {
	original := exportsURL
	defer func() { exportsURL = original }()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "100")
		_, _ = fmt.Fprint(w, "partial")
	}))
	defer server.Close()
	exportsURL = server.URL + "/"

	dest := filepath.Join(t.TempDir(), "recent.csv")
	c := NewClient("test-key")
	if _, err := c.DownloadExport(context.Background(), "recent.csv", dest); err == nil {
		t.Fatal("expected a truncated download error")
	}
	if _, err := os.Stat(dest); !os.IsNotExist(err) {
		t.Fatalf("partial export was not removed: %v", err)
	}
}

func TestExportErrorExplainsAPIMessage(t *testing.T) {
	original := exportsURL
	defer func() { exportsURL = original }()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprint(w, `{"message":"Invalid auth token"}`)
	}))
	defer server.Close()
	exportsURL = server.URL + "/"

	dest := filepath.Join(t.TempDir(), "recent.csv")
	c := NewClient("test-key")
	_, err := c.DownloadExport(context.Background(), "recent.csv", dest)
	if err == nil || !strings.Contains(err.Error(), "Invalid auth token") {
		t.Fatalf("error = %v, want API message", err)
	}
	if _, statErr := os.Stat(dest); !os.IsNotExist(statErr) {
		t.Fatalf("failed export left a file behind: %v", statErr)
	}
}
