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
	// Payload captured from the live get_cscb endpoint.
	jsonData := `{
		"query_status": "ok",
		"data": [
			{
				"time_stamp": "2026-03-20 15:28:46",
				"serial_number": "05a6cf9108f6941e492c3d2a1dc4dc9631b2",
				"thumbprint": "4ba01ae2c7fbedfe0e096c564665ed5d8f4b06124963a9cde140651cb29922b8",
				"thumbprint_algorithm": "SHA256",
				"subject_cn": "softportal.com",
				"issuer_cn": "R12",
				"valid_from": "2026-03-01T00:04:43Z",
				"valid_to": "2026-05-30T00:04:42Z",
				"cscb_listed": true,
				"cscb_reason": "Vidar"
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
	checks := []struct {
		name string
		got  string
		want string
	}{
		{"TimeStamp", entry.TimeStamp, "2026-03-20 15:28:46"},
		{"SerialNumber", entry.SerialNumber, "05a6cf9108f6941e492c3d2a1dc4dc9631b2"},
		{"Thumbprint", entry.Thumbprint, "4ba01ae2c7fbedfe0e096c564665ed5d8f4b06124963a9cde140651cb29922b8"},
		{"ThumbprintAlgorithm", entry.ThumbprintAlgorithm, "SHA256"},
		{"SubjectCN", entry.SubjectCN, "softportal.com"},
		{"IssuerCN", entry.IssuerCN, "R12"},
		{"ValidFrom", entry.ValidFrom, "2026-03-01T00:04:43Z"},
		{"ValidTo", entry.ValidTo, "2026-05-30T00:04:42Z"},
		{"CSCBReason", entry.CSCBReason, "Vidar"},
	}
	for _, c := range checks {
		if c.got != c.want {
			t.Errorf("%s = %q, want %q", c.name, c.got, c.want)
		}
	}
	if !entry.CSCBListed {
		t.Error("CSCBListed = false, want true")
	}
}

// TestParseFullSampleResponse guards the fields that were previously dropped
// on the floor. Payload shaped from a live get_info response.
func TestParseFullSampleResponse(t *testing.T) {
	jsonData := `{
		"query_status": "ok",
		"data": [
			{
				"sha256_hash": "cd24da0069d1d15c28538a0e9cb9610d0a47d6e7b36f00ac380ddfd0362afb93",
				"file_name": "sample.exe",
				"file_format": "PE",
				"file_arch": "AMD64",
				"magika": "pebin",
				"trid": ["33.1% (.EXE) Win64 Executable (generic)", "25.6% (.EXE) Generic Win/DOS"],
				"archive_pw": "1515",
				"comment": "top-level comment",
				"intelligence": {"clamav": ["Win.Trojan.Generic"], "downloads": "166", "uploads": "1", "mail": null},
				"file_information": [{"context": "cape", "value": "https://www.capesandbox.com/analysis/1"}],
				"ole_information": [],
				"code_sign": [
					{
						"subject_cn": "Bad Corp",
						"thumbprint": "4ba01ae2",
						"thumbprint_algorithm": "SHA256",
						"cscb_listed": true,
						"cscb_reason": "Vidar"
					}
				]
			}
		]
	}`

	resp, err := ParseAPIResponse([]byte(jsonData))
	if err != nil {
		t.Fatalf("ParseAPIResponse failed: %v", err)
	}
	if len(resp.Data) != 1 {
		t.Fatalf("Expected 1 sample, got %d", len(resp.Data))
	}
	s := resp.Data[0]

	if s.FileFormat != "PE" {
		t.Errorf("FileFormat = %q, want PE", s.FileFormat)
	}
	if s.FileArch != "AMD64" {
		t.Errorf("FileArch = %q, want AMD64", s.FileArch)
	}
	if s.Magika != "pebin" {
		t.Errorf("Magika = %q, want pebin", s.Magika)
	}
	if len(s.TrID) != 2 {
		t.Errorf("TrID has %d entries, want 2", len(s.TrID))
	}
	if s.ArchivePW != "1515" {
		t.Errorf("ArchivePW = %q, want 1515 (needed to open downloaded samples)", s.ArchivePW)
	}
	if s.Comment != "top-level comment" {
		t.Errorf("Comment = %q", s.Comment)
	}
	if s.Intelligence == nil {
		t.Fatal("Intelligence was dropped")
	}
	// The API quotes these counters; parsing them as ints would fail.
	if s.Intelligence.Downloads != "166" {
		t.Errorf("Intelligence.Downloads = %q, want \"166\"", s.Intelligence.Downloads)
	}
	if s.Intelligence.Uploads != "1" {
		t.Errorf("Intelligence.Uploads = %q, want \"1\"", s.Intelligence.Uploads)
	}
	if len(s.Intelligence.ClamAV) != 1 {
		t.Errorf("Intelligence.ClamAV = %v", s.Intelligence.ClamAV)
	}
	if len(s.FileInformation) != 1 || s.FileInformation[0].Context != "cape" {
		t.Errorf("FileInformation = %+v", s.FileInformation)
	}
	if len(s.CodeSign) != 1 {
		t.Fatalf("CodeSign has %d entries", len(s.CodeSign))
	}
	cs := s.CodeSign[0]
	if cs.Thumbprint != "4ba01ae2" || cs.ThumbprintAlgorithm != "SHA256" {
		t.Errorf("CodeSign thumbprint fields dropped: %+v", cs)
	}
	if !cs.CSCBListed || cs.CSCBReason != "Vidar" {
		t.Errorf("CodeSign CSCB fields dropped: %+v", cs)
	}
}

// TestNullableSampleFields checks that the nullable fields the API returns as
// JSON null do not break decoding.
func TestNullableSampleFields(t *testing.T) {
	jsonData := `{"query_status":"ok","data":[{
		"sha256_hash":"abc",
		"file_format":null,"file_arch":null,"magika":null,"trid":null,
		"archive_pw":null,"comment":null,"intelligence":null,
		"file_information":null,"ole_information":null
	}]}`

	resp, err := ParseAPIResponse([]byte(jsonData))
	if err != nil {
		t.Fatalf("null fields broke decoding: %v", err)
	}
	if resp.Data[0].SHA256Hash != "abc" {
		t.Errorf("unexpected sample: %+v", resp.Data[0])
	}
}
