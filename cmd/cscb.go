package cmd

import (
	"flag"
	"fmt"
)

// executeCSCB handles the 'cscb' subcommand
func executeCSCB(args []string) error {
	cscbCmd := flag.NewFlagSet("cscb", flag.ExitOnError)

	cscbCmd.Usage = func() {
		printUsageHeader("cscb", "Queries the Code Signing Certificate Blocklist (CSCB) from MalwareBazaar.")
		fmt.Println("\nExample:")
		fmt.Println("  mbzr cscb")
	}

	if err := cscbCmd.Parse(args); err != nil {
		return err
	}

	client, err := getAPIClient()
	if err != nil {
		printDetailedError(err, "Failed to create API client")
		return err
	}

	ctx, cancel := getContext()
	defer cancel()

	results, err := client.GetCSCB(ctx)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		fmt.Println("No CSCB entries found.")
		return nil
	}

	printJSON(results)

	return nil
}
