package api

import (
	"encoding/json"
	"testing"
)

func TestParseAPIResponse(t *testing.T) {
	jsonData := `{
		"query_status": "ok",
		"data": [
			{
				"sha256_hash": "094fd325049b8a9cf6d3e5ef2a6d4cc6a567d7d49c35f8bb8dd9e3c6acf3d78d",
				"file_name": "sample.exe",
				"file_size": 1024,
				"signature": "TrickBot",
				"tags": ["exe", "trojan"],
				"anonymous": 0,
				"first_seen": "2021-01-01 12:00:00"
			}
		]
	}`

	resp, err := ParseAPIResponse([]byte(jsonData))
	if err != nil {
		t.Fatalf("ParseAPIResponse failed: %v", err)
	}

	if resp.QueryStatus != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", resp.QueryStatus)
	}

	if len(resp.Data) != 1 {
		t.Fatalf("Expected 1 sample, got %d", len(resp.Data))
	}

	sample := resp.Data[0]
	if sample.SHA256Hash != "094fd325049b8a9cf6d3e5ef2a6d4cc6a567d7d49c35f8bb8dd9e3c6acf3d78d" {
		t.Errorf("Incorrect SHA256 hash: %s", sample.SHA256Hash)
	}
	if sample.FileName != "sample.exe" {
		t.Errorf("Incorrect file name: %s", sample.FileName)
	}
	if sample.FileSize != 1024 {
		t.Errorf("Incorrect file size: %d", sample.FileSize)
	}
	if sample.Signature != "TrickBot" {
		t.Errorf("Incorrect signature: %s", sample.Signature)
	}
	if len(sample.Tags) != 2 {
		t.Errorf("Incorrect tags length: %d", len(sample.Tags))
	}
}

func TestParseCSCBResponse(t *testing.T) {
	jsonData := `{
		"query_status": "ok",
		"data": [
			{
				"subject_cn": "Bad Corp",
				"serial_number": "123456",
				"reason": "Malware distribution"
			}
		]
	}`

	var resp CSCBResponse
	if err := json.Unmarshal([]byte(jsonData), &resp); err != nil {
		t.Fatalf("Unmarshal CSCBResponse failed: %v", err)
	}

	if resp.QueryStatus != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", resp.QueryStatus)
	}

	if len(resp.Data) != 1 {
		t.Fatalf("Expected 1 entry, got %d", len(resp.Data))
	}

	entry := resp.Data[0]
	if entry.SubjectCN != "Bad Corp" {
		t.Errorf("Incorrect SubjectCN: %s", entry.SubjectCN)
	}
	if entry.Reason != "Malware distribution" {
		t.Errorf("Incorrect Reason: %s", entry.Reason)
	}
}
