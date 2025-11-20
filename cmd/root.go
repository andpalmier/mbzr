package cmd

import (
	"fmt"
	"os"
)

// Execute runs the root command and handles subcommands
func Execute() error {
	// Parse global flags
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-v", "--verbose":
			SetVerbose(true)
			// Remove flag from args
			args = append(args[:i], args[i+1:]...)
			i--
		case "-V", "--version":
			return executeVersion([]string{})
		}
	}

	// handle root help
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		printRootHelp()
		fmt.Println()
		return nil
	}

	// handle subcommands
	switch args[0] {
	case "version":
		return executeVersion(args[1:])
	case "update":
		return executeUpdate(args[1:])
	case "upload":
		return executeUpload(args[1:])
	case "download":
		return executeDownload(args[1:])
	case "cscb":
		return executeCSCB(args[1:])
	case "query":
		return executeQuery(args[1:])
	case "recent_detections":
		return executeRecentDetections(args[1:])
	case "comment":
		return executeComment(args[1:])
	default:
		printError(fmt.Sprintf("unknown subcommand '%s'", args[0]))
		printRootHelp()
		fmt.Println()
		os.Exit(1)
	}
	return nil
}
