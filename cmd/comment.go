package cmd

import (
	"flag"
	"fmt"

	"github.com/andpalmier/mbzr/api"
)

// executeComment handles the 'comment' subcommand
func executeComment(args []string) error {
	commentCmd := flag.NewFlagSet("comment", flag.ExitOnError)
	sha256 := commentCmd.String("sha256", "", "sha256 hash of the file to comment on")
	comment := commentCmd.String("comment", "", "Comment to add to the malware sample")

	// helper function
	commentCmd.Usage = func() {
		fmt.Println("Usage:\n  mbzr comment [flags]\n")
		fmt.Println("Adds a comment to a malware sample in MalwareBazaar.")
		fmt.Println("\nFlags:")
		fmt.Println("  -sha256 <sha256_hash>	sha256 hash of the file to comment on")
		fmt.Println("  -comment <comment>		Comment to add to the malware sample")
		fmt.Println("\nExamples:")
		fmt.Println("  mbzr comment -sha256 <hash> -comment \"This sample is part of a new campaign.\"")
	}

	commentCmd.Parse(args)

	// validate the flags
	if *sha256 == "" || *comment == "" {
		fmt.Println("You must specify -sha256 and -comment flag\n")
		commentCmd.Usage()
		fmt.Println()
		return fmt.Errorf("missing required flags")
	}

	return api.AddComment(*sha256, *comment, apiKey)
}
