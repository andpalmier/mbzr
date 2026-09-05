package cmd

import (
	"context"
	"errors"
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

func TestGetContextCancellation(t *testing.T) {
	ctx, cancel := getContext()

	// Cancelling must propagate immediately rather than waiting out the
	// deadline, so this asserts the behaviour without sleeping for 30s.
	cancel()

	select {
	case <-ctx.Done():
		if !errors.Is(ctx.Err(), context.Canceled) {
			t.Errorf("Expected Canceled error, got %v", ctx.Err())
		}
	case <-time.After(time.Second):
		t.Error("Context was not done after cancel()")
	}
}
