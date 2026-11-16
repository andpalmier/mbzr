package cmd

import (
	"strings"
	"testing"
)

func TestCommandValidationRejectsBeforeNetwork(t *testing.T) {
	tests := []struct {
		name string
		fn   func([]string) error
		args []string
	}{
		{"download missing hash", executeDownload, nil},
		{"export missing file", executeExport, nil},
		{"comment missing flags", executeComment, nil},
		{"update missing flags", executeUpdate, nil},
		{"update invalid key", executeUpdate, []string{"-sha256", strings.Repeat("a", 64), "-key", "bad", "-value", "x"}},
		{"upload missing file", executeUpload, nil},
		{"upload both file and dir", executeUpload, []string{"-file", "a", "-dir", "b"}},
		{"upload invalid delivery", executeUpload, []string{"-file", "a", "-delivery_method", "carrier_pigeon"}},
		{"latest invalid selector", executeLatest, []string{"-selector", "bad"}},
		{"recent invalid hours", executeRecentDetections, []string{"-hours", "0"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fn(tt.args); err == nil {
				t.Error("expected validation error")
			}
		})
	}
}

func TestIsValidKey(t *testing.T) {
	for _, key := range AllowedKeys {
		if !isValidKey(key) {
			t.Errorf("allowed key %q rejected", key)
		}
	}
	if isValidKey("not_a_key") {
		t.Error("unknown key accepted")
	}
}

func TestKVFlagStringAndEmpty(t *testing.T) {
	k := newKVFlag("reference", referenceKeys)
	if !k.empty() || k.String() != "" {
		t.Fatalf("new flag should be empty")
	}
	for _, v := range []string{"twitter=https://b", "any_run=https://a"} {
		if err := k.Set(v); err != nil {
			t.Fatal(err)
		}
	}
	if k.empty() {
		t.Error("flag should not be empty")
	}
	if got := k.String(); got != "any_run=https://a,twitter=https://b" {
		t.Errorf("String() = %q", got)
	}
}
