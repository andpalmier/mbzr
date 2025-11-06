package cmd

import (
	"fmt"
	"os"
)

var apiKey string

func Execute() error {
	// check for API key in env variables
	apiKey = os.Getenv("MBZR_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("Error: please set MBZR_API_KEY environment variable")
	}

	// handle root help
	if len(os.Args) < 2 || os.Args[1] == "-h" || os.Args[1] == "--help" || os.Args[1] == "help" {
		printRootHelp()
		fmt.Println()
		return nil
	}

	// if no subcommand is provided
	if len(os.Args) < 2 {
		printRootHelp()
		fmt.Println()
		return fmt.Errorf("subcommand not specified")
	}

	// handle subcommands
	switch os.Args[1] {
	case "update":
		return executeUpdate(os.Args[2:])
	case "upload":
		return executeUpload(os.Args[2:])
	case "download":
		return executeDownload(os.Args[2:])
	case "cscb":
		return executeCSCB(os.Args[2:])
	case "query":
		return executeQuery(os.Args[2:])
	case "recent_detections":
		return executeRecentDetections(os.Args[2:])
	case "comment":
		return executeComment(os.Args[2:])
	default:
		fmt.Printf("Error: unknown subcommand '%s'\n\n", os.Args[1])
		printRootHelp()
		fmt.Println()
		os.Exit(1)
	}
	return nil
}
