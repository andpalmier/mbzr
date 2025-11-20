package cmd

import (
	"context"
	"testing"
	"time"
)

func TestGetContext(t *testing.T) {
	ctx, cancel := getContext()
	defer cancel()

	if ctx == nil {
		t.Fatal("getContext() returned nil")
	}

	// Check that context has a deadline
	deadline, ok := ctx.Deadline()
	if !ok {
		t.Error("Context should have a deadline")
	}

	// Check that deadline is in the future
	if time.Until(deadline) <= 0 {
		t.Error("Context deadline should be in the future")
	}

	// Check that deadline is approximately 30 seconds
	expectedDeadline := 30 * time.Second
	actualDeadline := time.Until(deadline)

	// Allow 1 second tolerance
	if actualDeadline < expectedDeadline-time.Second || actualDeadline > expectedDeadline+time.Second {
		t.Errorf("Expected deadline around %v, got %v", expectedDeadline, actualDeadline)
	}
}

func TestSetVerbose(t *testing.T) {
	// Save original state
	original := IsVerbose()
	defer SetVerbose(original)

	SetVerbose(true)
	if !IsVerbose() {
		t.Error("SetVerbose(true) did not set verbose flag")
	}

	SetVerbose(false)
	if IsVerbose() {
		t.Error("SetVerbose(false) did not unset verbose flag")
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		s      string
		substr string
		want   bool
	}{
		{"hello world", "world", true},
		{"hello world", "hello", true},
		{"hello world", "lo wo", true},
		{"hello world", "xyz", false},
		{"", "", true},
		{"hello", "", true},
		{"", "hello", false},
	}

	for _, tt := range tests {
		got := contains(tt.s, tt.substr)
		if got != tt.want {
			t.Errorf("contains(%q, %q) = %v, want %v", tt.s, tt.substr, got, tt.want)
		}
	}
}

func TestGetContextCancellation(t *testing.T) {
	ctx, cancel := getContext()
	defer cancel()

	// Wait for context to be done (should timeout after 30 seconds)
	select {
	case <-ctx.Done():
		// Context should eventually be done due to timeout
		if ctx.Err() != context.DeadlineExceeded {
			t.Errorf("Expected DeadlineExceeded error, got %v", ctx.Err())
		}
	case <-time.After(31 * time.Second):
		t.Error("Context did not timeout within expected time")
	}
}
