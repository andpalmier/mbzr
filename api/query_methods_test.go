package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

const queryMethodsSHA256 = "cd24da0069d1d15c28538a0e9cb9610d0a47d6e7b36f00ac380ddfd0362afb93"

func mbzrTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	return NewClient("test-key", WithBaseURL(server.URL+"/"))
}

func TestQueryMethodsForwardParameters(t *testing.T) {
	tests := []struct {
		name  string
		query string
		key   string
		value string
		limit int
		call  func(context.Context, *Client) ([]MalwareSample, error)
	}{
		{"hash", "get_info", "hash", queryMethodsSHA256, 0, func(ctx context.Context, c *Client) ([]MalwareSample, error) {
			return c.QueryByHash(ctx, queryMethodsSHA256, 0)
		}},
		{"tag", "get_taginfo", "tag", "Emotet", 10, func(ctx context.Context, c *Client) ([]MalwareSample, error) { return c.QueryByTag(ctx, "Emotet", 10) }},
		{"signature", "get_siginfo", "signature", "Emotet", 10, func(ctx context.Context, c *Client) ([]MalwareSample, error) {
			return c.QueryBySignature(ctx, "Emotet", 10)
		}},
		{"file type", "get_file_type", "file_type", "exe", 10, func(ctx context.Context, c *Client) ([]MalwareSample, error) {
			return c.QueryByFileType(ctx, "exe", 10)
		}},
		{"clamav", "get_clamavinfo", "clamav", "Win.Trojan", 10, func(ctx context.Context, c *Client) ([]MalwareSample, error) {
			return c.QueryByClamAV(ctx, "Win.Trojan", 10)
		}},
		{"imphash", "get_imphash", "imphash", "d41d8cd98f00b204e9800998ecf8427e", 10, func(ctx context.Context, c *Client) ([]MalwareSample, error) {
			return c.QueryByImpHash(ctx, "d41d8cd98f00b204e9800998ecf8427e", 10)
		}},
		{"tlsh", "get_tlsh", "tlsh", "T1ABC", 10, func(ctx context.Context, c *Client) ([]MalwareSample, error) { return c.QueryByTLSH(ctx, "T1ABC", 10) }},
		{"telfhash", "get_telfhash", "telfhash", "T1ABC", 10, func(ctx context.Context, c *Client) ([]MalwareSample, error) {
			return c.QueryByTelfHash(ctx, "T1ABC", 10)
		}},
		{"dhash", "get_dhash_icon", "dhash_icon", "abc", 10, func(ctx context.Context, c *Client) ([]MalwareSample, error) { return c.QueryByDHash(ctx, "abc", 10) }},
		{"gimphash", "get_gimphash", "gimphash", "abc", 10, func(ctx context.Context, c *Client) ([]MalwareSample, error) {
			return c.QueryByGimphash(ctx, "abc", 10)
		}},
		{"yara", "get_yarainfo", "yara_rule", "Rule", 10, func(ctx context.Context, c *Client) ([]MalwareSample, error) { return c.QueryByYara(ctx, "Rule", 10) }},
		{"issuer", "get_issuerinfo", "issuer_cn", "Issuer", 10, func(ctx context.Context, c *Client) ([]MalwareSample, error) {
			return c.QueryByIssuerCN(ctx, "Issuer", 10)
		}},
		{"subject", "get_subjectinfo", "subject_cn", "Subject", 10, func(ctx context.Context, c *Client) ([]MalwareSample, error) {
			return c.QueryBySubjectCN(ctx, "Subject", 10)
		}},
		{"serial", "get_certificate", "serial_number", "1234", 10, func(ctx context.Context, c *Client) ([]MalwareSample, error) {
			return c.QueryBySerialNumber(ctx, "1234", 10)
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := mbzrTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					t.Errorf("ParseForm: %v", err)
					return
				}
				if got := r.Form.Get("query"); got != tt.query {
					t.Errorf("query = %q, want %q", got, tt.query)
				}
				if got := r.Form.Get(tt.key); got != tt.value {
					t.Errorf("%s = %q, want %q", tt.key, got, tt.value)
				}
				if tt.limit > 0 {
					if got := r.Form.Get("limit"); got != fmt.Sprint(tt.limit) {
						t.Errorf("limit = %q, want %d", got, tt.limit)
					}
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = fmt.Fprintf(w, `{"query_status":"ok","data":[{"sha256_hash":%q}]}`, queryMethodsSHA256)
			})

			got, err := tt.call(context.Background(), c)
			if err != nil {
				t.Fatalf("call failed: %v", err)
			}
			if len(got) != 1 || got[0].SHA256Hash != queryMethodsSHA256 {
				t.Fatalf("unexpected result: %+v", got)
			}
		})
	}
}

func TestQueryAPIRejectsInvalidLimitBeforeNetwork(t *testing.T) {
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		_, _ = fmt.Fprint(w, `{"query_status":"ok","data":[]}`)
	}))
	defer server.Close()

	c := NewClient("test-key", WithBaseURL(server.URL+"/"))
	for _, limit := range []int{-1, 1001} {
		if _, err := c.QueryByTag(context.Background(), "Emotet", limit); err == nil {
			t.Errorf("limit=%d was accepted", limit)
		}
	}
	if got := requests.Load(); got != 0 {
		t.Errorf("API was contacted %d times for invalid limits", got)
	}
}

func TestQueryAPIStatusError(t *testing.T) {
	c := mbzrTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"query_status":"hash_not_found"}`)
	})
	if _, err := c.QueryByHash(context.Background(), queryMethodsSHA256, 0); err == nil {
		t.Fatal("expected an API status error")
	}
}

func TestAddComment(t *testing.T) {
	var gotQuery, gotHash, gotComment string
	c := mbzrTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		gotQuery, gotHash, gotComment = r.Form.Get("query"), r.Form.Get("sha256_hash"), r.Form.Get("comment")
		_, _ = fmt.Fprint(w, `{"query_status":"ok"}`)
	})
	if err := c.AddComment(context.Background(), queryMethodsSHA256, "hello"); err != nil {
		t.Fatalf("AddComment: %v", err)
	}
	if gotQuery != "add_comment" || gotHash != queryMethodsSHA256 || gotComment != "hello" {
		t.Fatalf("unexpected form: query=%q hash=%q comment=%q", gotQuery, gotHash, gotComment)
	}
}

func TestGetRecentDetectionsAndCSCB(t *testing.T) {
	t.Run("recent detections", func(t *testing.T) {
		var gotHours string
		c := mbzrTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			_ = r.ParseForm()
			gotHours = r.Form.Get("hours")
			_, _ = fmt.Fprintf(w, `{"query_status":"ok","data":[{"sha256_hash":%q}]}`, queryMethodsSHA256)
		})
		got, err := c.GetRecentDetections(context.Background(), 24)
		if err != nil || len(got) != 1 {
			t.Fatalf("GetRecentDetections = %+v, %v", got, err)
		}
		if gotHours != "24" {
			t.Errorf("hours = %q, want 24", gotHours)
		}
	})

	t.Run("invalid hours", func(t *testing.T) {
		c := NewClient("test-key")
		if _, err := c.GetRecentDetections(context.Background(), 0); err == nil {
			t.Error("hours=0 should be rejected")
		}
		if _, err := c.GetRecentDetections(context.Background(), 169); err == nil {
			t.Error("hours=169 should be rejected")
		}
	})

	t.Run("cscb", func(t *testing.T) {
		c := mbzrTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprint(w, `{"query_status":"ok","data":[{"serial_number":"abc","cscb_listed":true}]}`)
		})
		got, err := c.GetCSCB(context.Background())
		if err != nil || len(got) != 1 || got[0].SerialNumber != "abc" || !got[0].CSCBListed {
			t.Fatalf("GetCSCB = %+v, %v", got, err)
		}
	})
}
