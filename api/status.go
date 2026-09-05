package api

import "fmt"

// StatusError is returned when the API answers with a query_status that
// indicates failure. It carries the raw status so callers can match on it,
// and renders a plain-English explanation.
type StatusError struct {
	Status string
	Query  string
}

func (e *StatusError) Error() string {
	if msg, ok := statusMessages[e.Status]; ok {
		return msg
	}
	if e.Query != "" {
		return fmt.Sprintf("the API rejected the %s query with status %q", e.Query, e.Status)
	}
	return fmt.Sprintf("the API returned status %q", e.Status)
}

// newStatusError returns nil when the status reports success, and a
// *StatusError describing the failure otherwise.
func newStatusError(status, query string) error {
	switch status {
	case "ok", "success", "inserted", "updated", "exists":
		return nil
	}
	return &StatusError{Status: status, Query: query}
}

// statusMessages explains every query_status documented at
// https://bazaar.abuse.ch/api/ in terms a user can act on.
var statusMessages = map[string]string{
	// Authentication and transport
	"no_api_key":         "no API key was accepted: set ABUSECH_API_KEY (get one at https://auth.abuse.ch/)",
	"user_blacklisted":   "this API key is blacklisted: contact https://www.spamhaus.com/#contact-form",
	"http_post_expected": "the API expected an HTTP POST request",

	// Nothing matched
	"no_results":          "the query returned no results",
	"hash_not_found":      "that hash is unknown to MalwareBazaar",
	"file_not_found":      "that sample is unknown to MalwareBazaar, or is no longer available",
	"tag_not_found":       "that tag is unknown to MalwareBazaar",
	"signature_not_found": "that signature is unknown to MalwareBazaar",
	"clamav_not_found":    "that ClamAV signature is unknown to MalwareBazaar",
	"yara_not_found":      "that YARA rule is unknown to MalwareBazaar",

	// Malformed input
	"illegal_hash":        "the API rejected that hash as malformed",
	"illegal_sha256_hash": "the API rejected that SHA256 hash as malformed",
	"illegal_tag":         "the API rejected that tag as malformed",
	"illegal_signature":   "the API rejected that signature as malformed",
	"illegal_file_type":   "the API rejected that file type as malformed",
	"illegal_clamav":      "the API rejected that ClamAV signature as malformed",
	"illegal_imphash":     "the API rejected that imphash as malformed",
	"illegal_tlsh":        "the API rejected that TLSH as malformed",
	"illegal_telfhash":    "the API rejected that telfhash as malformed",
	"illegal_gimphash":    "the API rejected that gimphash as malformed",
	"illegal_dhash_icon":  "the API rejected that icon dhash as malformed",
	"illegal_yara_rule":   "the API rejected that YARA rule name as malformed",
	"illegal_issuer_cn":   "the API rejected that issuer CN as malformed",
	"illegal_hours":       "hours must be a number between 1 and 168",

	// Missing input
	"no_hash_provided":      "no hash was sent to the API",
	"no_sha256_hash":        "no SHA256 hash was sent to the API",
	"no_tag_provided":       "no tag was sent to the API",
	"no_signature_provided": "no signature was sent to the API",
	"no_file_type":          "no file type was sent to the API",
	"no_clamav_provided":    "no ClamAV signature was sent to the API",
	"no_imphash":            "no imphash was sent to the API",
	"no_tlsh":               "no TLSH was sent to the API",
	"no_telfhash":           "no telfhash was sent to the API",
	"no_gimphash":           "no gimphash was sent to the API",
	"no_yara_rule_provided": "no YARA rule was sent to the API",
	"no_issuer_cn":          "no issuer CN was sent to the API",
	"no_selector":           "no selector was sent: use 'time' or '100'",
	"unknown_selector":      "unknown selector: use 'time' or '100'",

	// Upload
	"file_already_known": "this sample is already known to MalwareBazaar",
	"file_expected":      "no file reached the API",
	// Not in the published docs; returned for an empty or unreadable upload.
	"illegal_file": "the API rejected the file (an empty or unreadable sample)",

	// Update
	"permission_denied": "you can only update entries created by your own account",
	"unknown_key":       "the API does not recognise that key",
}
