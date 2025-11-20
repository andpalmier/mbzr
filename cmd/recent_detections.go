package cmd

import (
	"flag"
	"fmt"
)

// executeRecentDetections handles the 'recent_detections' subcommand
func executeRecentDetections(args []string) error {
	recentCmd := flag.NewFlagSet("recent_detections", flag.ExitOnError)
	hours := recentCmd.Int("hours", 24, "Retrieve recent detections from the last n hours (default: 24)")

	recentCmd.Usage = func() {
		printUsageHeader("recent_detections", "Retrieves malware samples detected in the last specified number of hours.")
		fmt.Println("\nFlags:")
		fmt.Println("  -hours <number>\tNumber of hours to look back (default: 24)")
		fmt.Println("\nExamples:")
		fmt.Println("  mbzr recent_detections -hours 24")
	}

	if err := recentCmd.Parse(args); err != nil {
		return err
	}

	if *hours < 1 || *hours > 168 {
		printError("hours must be between 1 and 168")
		recentCmd.Usage()
		fmt.Println()
		return fmt.Errorf("invalid hours value")
	}

	client, err := getAPIClient()
	if err != nil {
		printDetailedError(err, "Failed to create API client")
		return err
	}

	ctx, cancel := getContext()
	defer cancel()

	results, err := client.GetRecentDetections(ctx, *hours)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		fmt.Println("No recent detections found.")
		return nil
	}

	printJSON(results)

	return nil
}
