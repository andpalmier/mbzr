package api

import (
	"errors"
	"strings"
	"testing"
)

func TestNewStatusError(t *testing.T) {
	successes := []string{"ok", "success", "inserted", "updated", "exists"}
	for _, s := range successes {
		if err := newStatusError(s, "q"); err != nil {
			t.Errorf("status %q should be success, got %v", s, err)
		}
	}

	failures := []string{"hash_not_found", "no_api_key", "user_blacklisted", "illegal_tag", "no_results"}
	for _, s := range failures {
		err := newStatusError(s, "q")
		if err == nil {
			t.Errorf("status %q should be an error", s)
			continue
		}
		var se *StatusError
		if !errors.As(err, &se) {
			t.Errorf("status %q should yield a *StatusError, got %T", s, err)
			continue
		}
		if se.Status != s {
			t.Errorf("StatusError.Status = %q, want %q", se.Status, s)
		}
		// The whole point is that the user learns the reason, so the message
		// must not just echo the raw status back.
		if strings.TrimSpace(err.Error()) == s {
			t.Errorf("status %q has no explanation", s)
		}
	}
}

func TestStatusErrorExplainsUnknownStatus(t *testing.T) {
	err := newStatusError("brand_new_status", "get_taginfo")
	if err == nil {
		t.Fatal("an unrecognised status should still be an error")
	}
	msg := err.Error()
	if !strings.Contains(msg, "brand_new_status") || !strings.Contains(msg, "get_taginfo") {
		t.Errorf("unknown status should name both the status and the query, got: %s", msg)
	}
}

// Every documented status the client can encounter should have a message.
func TestDocumentedStatusesHaveMessages(t *testing.T) {
	documented := []string{
		"no_api_key", "user_blacklisted", "http_post_expected",
		"hash_not_found", "file_not_found", "tag_not_found", "signature_not_found",
		"clamav_not_found", "yara_not_found", "no_results",
		"illegal_hash", "illegal_sha256_hash", "illegal_tag", "illegal_signature",
		"illegal_file_type", "illegal_clamav", "illegal_imphash", "illegal_tlsh",
		"illegal_telfhash", "illegal_gimphash", "illegal_dhash_icon",
		"illegal_yara_rule", "illegal_issuer_cn", "illegal_hours",
		"no_hash_provided", "no_sha256_hash", "no_tag_provided", "no_signature_provided",
		"no_file_type", "no_clamav_provided", "no_imphash", "no_tlsh", "no_telfhash",
		"no_gimphash", "no_yara_rule_provided", "no_issuer_cn",
		"no_selector", "unknown_selector",
		"file_already_known", "file_expected", "permission_denied", "unknown_key",
	}
	for _, s := range documented {
		if _, ok := statusMessages[s]; !ok {
			t.Errorf("documented status %q has no explanation", s)
		}
	}
}
