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

	// helper function
	downloadCmd.Usage = func() {
		fmt.Println("Usage:\n  mbzr download [flags]\n")
		fmt.Println("Downloads a malware sample by its sha256 hash from MalwareBazaar.")
		fmt.Println("\nFlags:")
		fmt.Println("  -sha256 <sha256_hash>	sha256 hash of the file to download")
	}

	downloadCmd.Parse(args)

	if *sha256 == "" {
		fmt.Println("Error: please use the -sha256 flag to specify the sha256 hash of the file to download\n")
		downloadCmd.Usage()
		fmt.Println()
		return fmt.Errorf("missing required -sha256 flag")
	}

	return api.DownloadSample(*sha256, apiKey)
}
