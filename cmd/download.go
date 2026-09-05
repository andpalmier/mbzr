package cmd

import (
	"flag"
	"fmt"

	"github.com/andpalmier/mbzr/api"
)

// executeDownload handles the 'download' subcommand
func executeDownload(args []string) error {
	downloadCmd := flag.NewFlagSet("download", flag.ExitOnError)
	sha256 := downloadCmd.String("sha256", "", "sha256 hash of the file to download")
	out := downloadCmd.String("out", "", "Path to write the sample to (default: <sha256>.zip)")

	downloadCmd.Usage = func() {
		printUsageHeader("download", "Downloads a malware sample by its sha256 hash from MalwareBazaar.")
		fmt.Println("\nFlags:")
		fmt.Println("  -sha256 <sha256_hash>\tsha256 hash of the file to download")
		fmt.Println("  -out <path>\t\tPath to write the sample to (default: <sha256>.zip)")
		fmt.Println("\nExamples:")
		fmt.Println("  mbzr download -sha256 <hash>")
		fmt.Println("  mbzr download -sha256 <hash> -out /tmp/sample.zip")
		fmt.Printf("\nSamples are zipped and protected with the password %q.\n", api.ZipPassword)
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
		return err
	}

	ctx, cancel := getContext()
	defer cancel()

	path, err := client.DownloadSample(ctx, *sha256, *out)
	if err != nil {
		printDetailedError(err, fmt.Sprintf("Failed to download sample: %s", *sha256))
		return err
	}

	printSuccess(fmt.Sprintf("File downloaded successfully: %s", path))
	printSuccess(fmt.Sprintf("The archive is password protected: %s", api.ZipPassword))
	return nil
}
