package cmd

import (
	"log/slog"
	"os"
)

var logger *slog.Logger

// InitLogger initializes the structured logger
func InitLogger(verbose bool) {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	if verbose {
		opts.Level = slog.LevelDebug
	}

	// Use JSON handler for structured logging
	handler := slog.NewJSONHandler(os.Stderr, opts)
	logger = slog.New(handler)
}

// GetLogger returns the global logger instance
func GetLogger() *slog.Logger {
	if logger == nil {
		InitLogger(false)
	}
	return logger
}
