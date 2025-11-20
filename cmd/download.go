package cmd

import (
	"flag"
	"fmt"
)

// executeDownload handles the 'download' subcommand
func executeDownload(args []string) error {
	downloadCmd := flag.NewFlagSet("download", flag.ExitOnError)
	sha256 := downloadCmd.String("sha256", "", "sha256 hash of the file to download")

	downloadCmd.Usage = func() {
		printUsageHeader("download", "Downloads a malware sample by its sha256 hash from MalwareBazaar.")
		fmt.Println("\nFlags:")
		fmt.Println("  -sha256 <sha256_hash>\tsha256 hash of the file to download")
	}

	if err := downloadCmd.Parse(args); err != nil {
		return err
	}

	if *sha256 == "" {
		printError("you must specify a sha256 hash using -sha256")
		downloadCmd.Usage()
		fmt.Println()
		return fmt.Errorf("missing sha256 hash")
	}

	client, err := getAPIClient()
	if err != nil {
		printDetailedError(err, "Failed to create API client")
		return fmt.Errorf("MALWAREBAZAAR_API_KEY environment variable is required for downloading samples")
	}

	ctx, cancel := getContext()
	defer cancel()

	err = client.DownloadSample(ctx, *sha256)
	if err != nil {
		printDetailedError(err, fmt.Sprintf("Failed to download sample: %s", *sha256))
		return err
	}

	printSuccess(fmt.Sprintf("File downloaded successfully: %s.zip", *sha256))
	return nil
}
