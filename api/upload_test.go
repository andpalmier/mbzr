package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestClient_Upload(t *testing.T) {
	// Create temporary file to upload
	tmpfile, err := os.CreateTemp("", "sample-*.exe")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.WriteString("test content")
	tmpfile.Close()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		// Check that it's multipart
		if r.Header.Get("Content-Type") == "" {
			t.Error("Expected Content-Type header")
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{
			"query_status": "ok",
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

	_, err = c.UploadFile(context.Background(), tmpfile.Name(), false, []string{"tag1"}, "comment", nil)
	if err != nil {
		t.Fatalf("UploadFile() error = %v", err)
	}
}
