package cmd

import (
	"strings"
	"testing"
)

func TestExecuteQueryRejectsMultipleSelectors(t *testing.T) {
	err := executeQuery([]string{"-hash", strings.Repeat("a", 64), "-tag", "Emotet"})
	if err == nil {
		t.Fatal("expected multiple query selectors to be rejected")
	}
}
