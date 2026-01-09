package cmd

import (
	"fmt"
)

// Version information (set via ldflags during build)
var (
	Version   = "dev"
	Commit    = "unknown"
	BuildDate = "unknown"
)

// executeVersion handles the 'version' subcommand
func executeVersion(_ []string) error {
	fmt.Printf("mbzr version %s\n", Version)
	fmt.Printf("  commit: %s\n", Commit)
	fmt.Printf("  built: %s\n", BuildDate)
	return nil
}
