package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testSHA256 = "094fd325049b8a9cf6d3e5ef2a6d4cc6a567d7d49c35f8bb8dd9e3c6acf3d78d"

func TestUpdateSampleStatuses(t *testing.T) {
	tests := []struct {
		name       string
		status     string
		wantStatus string
		wantErr    bool
	}{
		// "updated" is the value the API documents for a successful update.
		{"documented success", "updated", StatusUpdated, false},
		{"legacy ok", "ok", StatusUpdated, false},
		{"legacy success", "success", StatusUpdated, false},
		{"already present", "exists", StatusExists, false},
		{"permission denied", "permission_denied", "", true},
		{"hash not found", "hash_not_found", "", true},
		{"unknown key", "unknown_key", "", true},
		{"no api key", "no_api_key", "", true},
		{"blacklisted", "user_blacklisted", "", true},
		{"unrecognised", "something_else", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = fmt.Fprintf(w, `{"query_status": %q}`, tt.status)
			}))
			defer server.Close()

			c := NewClient("test-key")
			c.baseURL = server.URL + "/"

			got, err := c.UpdateSample(context.Background(), testSHA256, "links", "https://abuse.ch")
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected an error for status %q, got none", tt.status)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for status %q: %v", tt.status, err)
			}
			if got != tt.wantStatus {
				t.Errorf("status = %q, want %q", got, tt.wantStatus)
			}
		})
	}
}

func TestUpdateSampleRejectsBadHash(t *testing.T) {
	c := NewClient("test-key")
	if _, err := c.UpdateSample(context.Background(), "not-a-hash", "links", "x"); err == nil {
		t.Error("expected an error for an invalid hash")
	}
}
