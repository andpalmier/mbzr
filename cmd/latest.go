package cmd

import (
	"flag"
	"fmt"
)

// executeLatest handles the 'latest' subcommand
func executeLatest(args []string) error {
	latestCmd := flag.NewFlagSet("latest", flag.ExitOnError)
	selector := latestCmd.String("selector", "time", "Selector for latest samples: 'time' (last 60m) or '100' (last 100 samples)")

	latestCmd.Usage = func() {
		printUsageHeader("latest", "Retrieves the latest malware samples added to MalwareBazaar.")
		fmt.Println("\nFlags:")
		fmt.Println("  -selector <value>\tSelector: 'time' (default, last 60m) or '100' (last 100 samples)")
		fmt.Println("\nExamples:")
		fmt.Println("  mbzr latest")
		fmt.Println("  mbzr latest -selector 100")
	}

	if err := latestCmd.Parse(args); err != nil {
		return err
	}

	if *selector != "time" && *selector != "100" {
		printError("invalid selector: must be 'time' or '100'")
		latestCmd.Usage()
		fmt.Println()
		return fmt.Errorf("invalid selector")
	}

	client, err := getAPIClient()
	if err != nil {
		printDetailedError(err, "Failed to create API client")
		return err
	}

	ctx, cancel := getContext()
	defer cancel()

	results, err := client.QueryLatest(ctx, *selector)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		fmt.Println("No latest samples found.")
		return nil
	}

	printJSON(results)

	return nil
}
