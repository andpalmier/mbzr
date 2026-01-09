package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_QueryByHash(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{
			"query_status": "ok",
			"data": [
				{
					"sha256_hash": "88d862aeb067278155c67a6d4e5be927b36f08149c950d75a3a419eb20560aa1",
					"file_name": "malware.exe",
					"file_size": 12345,
					"file_type_mime": "application/x-dosexec",
					"signature": "Emotet"
				}
			]
		}`)
	}))
	defer server.Close()

	c := NewClient("test-key")
	c.baseURL = server.URL + "/"

	data, err := c.QueryByHash(context.Background(), "88d862aeb067278155c67a6d4e5be927b36f08149c950d75a3a419eb20560aa1", 0)
	if err != nil {
		t.Fatalf("QueryByHash() error = %v", err)
	}

	if len(data) != 1 {
		t.Errorf("Expected 1 result, got %d", len(data))
	}
	if data[0].SHA256Hash != "88d862aeb067278155c67a6d4e5be927b36f08149c950d75a3a419eb20560aa1" {
		t.Errorf("Got wrong hash: %s", data[0].SHA256Hash)
	}
}

func TestClient_QueryByTag(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			t.Fatal(err)
		}
		if r.Form.Get("query") != "get_taginfo" {
			t.Errorf("Expected query=get_taginfo, got %s", r.Form.Get("query"))
		}
		if r.Form.Get("tag") != "Emotet" {
			t.Errorf("Expected tag=Emotet, got %s", r.Form.Get("tag"))
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{
			"query_status": "ok",
			"data": [
				{
					"sha256_hash": "88d862aeb067278155c67a6d4e5be927b36f08149c950d75a3a419eb20560aa1"
				}
			]
		}`)
	}))
	defer server.Close()

	c := NewClient("test-key")
	c.baseURL = server.URL + "/"

	data, err := c.QueryByTag(context.Background(), "Emotet", 10)
	if err != nil {
		t.Fatalf("QueryByTag() error = %v", err)
	}
	if len(data) != 1 {
		t.Errorf("Expected 1 result, got %d", len(data))
	}
}

func TestClient_QueryLatest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		if r.Form.Get("query") != "get_recent" {
			t.Errorf("Expected query=get_recent, got %s", r.Form.Get("query"))
		}
		if r.Form.Get("selector") != "100" {
			t.Errorf("Expected selector=100, got %s", r.Form.Get("selector"))
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{
			"query_status": "ok",
			"data": []
		}`)
	}))
	defer server.Close()

	c := NewClient("test-key")
	c.baseURL = server.URL + "/"

	_, err := c.QueryLatest(context.Background(), "100")
	if err != nil {
		t.Fatalf("QueryLatest() error = %v", err)
	}
}
