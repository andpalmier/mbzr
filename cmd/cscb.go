package cmd

import (
	"flag"
	"fmt"

	"github.com/andpalmier/mbzr/api"
)

// executeCSCB handles the cscb command
func executeCSCB(args []string) error {
	cscbCmd := flag.NewFlagSet("cscb", flag.ExitOnError)
	cscbCmd.Usage = func() {
		fmt.Println("Usage:\n  mbzr cscb [flags]\n")
		fmt.Println("Query the Code Signing Certificate Blocklist from MalwareBazaar.")
		fmt.Println("\nExample:")
		fmt.Println("  mbzr cscb")
	}
	cscbCmd.Parse(args)

	return api.GetCSCB(apiKey)
}
