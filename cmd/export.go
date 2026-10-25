package cmd

import (
	"flag"
	"fmt"

	"github.com/andpalmier/mbzr/api"
)

// executeExport handles the 'export' subcommand
func executeExport(args []string) error {
	exportCmd := flag.NewFlagSet("export", flag.ExitOnError)
	file := exportCmd.String("file", "", "Export file to download, e.g. recent.csv")
	out := exportCmd.String("out", "", "Path to write the export to (default: the export file name)")

	exportCmd.Usage = func() {
		printUsageHeader("export", "Downloads a bulk data export from MalwareBazaar.")
		fmt.Println("\nFlags:")
		fmt.Println("  -file <name>             Export file to download (required)")
		fmt.Println("  -out <path>              Path to write the export to (default: the file name)")
		fmt.Println("\nFiles:")
		for _, f := range api.ExportFiles {
			fmt.Printf("  %-24s %s\n", f.Name, f.Description)
		}
		fmt.Println("\nNote:")
		fmt.Println("  Recent exports are regenerated every 5 minutes and full dumps hourly; do not fetch them more often.")
		fmt.Println("  Full dumps are large: raise the timeout, e.g. mbzr -t 10m export -file full.csv.zip")
	}

	if err := exportCmd.Parse(args); err != nil {
		return err
	}

	if *file == "" {
		printError("you must specify an export file using -file")
		exportCmd.Usage()
		fmt.Println()
		return fmt.Errorf("missing export file")
	}

	client, err := getAPIClient()
	if err != nil {
		printDetailedError(err, "Failed to create API client")
		return err
	}

	ctx, cancel := getContext()
	defer cancel()

	path, err := client.DownloadExport(ctx, *file, *out)
	if err != nil {
		printDetailedError(err, fmt.Sprintf("Failed to download export: %s", *file))
		return err
	}

	printSuccess(fmt.Sprintf("Export downloaded successfully: %s", path))
	return nil
}
