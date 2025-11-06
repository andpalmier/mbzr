package cmd

import (
	"errors"
	"flag"
	"fmt"

	"github.com/andpalmier/mbzr/api"
)

// executeRecentDetections handles the recent_detections command
func executeRecentDetections(args []string) error {
	recentCmd := flag.NewFlagSet("recent_detections", flag.ExitOnError)
	hours := recentCmd.Int("hours", 48, "How many hours to look back (default: 48, max: 168)")

	// helper function
	recentCmd.Usage = func() {
		fmt.Println("Usage:\n  mbzr recent_detections [flags]\n")
		fmt.Println("Retrieves malware samples detected in the last specified number of hours.")
		fmt.Println("\nFlags:")
		fmt.Println("  -hours <number>	Number of hours to look back (default: 48, max: 168)")
		fmt.Println("\nExamples:")
		fmt.Println("  mbzr recent_detections -hours 24")
	}

	recentCmd.Parse(args)

	if *hours < 1 || *hours > 168 {
		fmt.Println("Error: hours must be between 1 and 168\n")
		recentCmd.Usage()
		fmt.Println()
		return errors.New("hours must be between 1 and 168")
	}

	return api.GetRecentDetections(*hours, apiKey)
}
