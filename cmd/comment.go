package cmd

import (
	"flag"
	"fmt"
)

// executeComment handles the 'comment' subcommand
func executeComment(args []string) error {
	commentCmd := flag.NewFlagSet("comment", flag.ExitOnError)
	sha256 := commentCmd.String("sha256", "", "sha256 hash of the file to comment on")
	comment := commentCmd.String("comment", "", "Comment to add to the malware sample")

	commentCmd.Usage = func() {
		printUsageHeader("comment", "Adds a comment to a malware sample in MalwareBazaar.")
		fmt.Println("\nFlags:")
		fmt.Println("  -sha256 <sha256_hash>\tsha256 hash of the file to comment on")
		fmt.Println("  -comment <comment>\t\tComment to add to the malware sample")
		fmt.Println("\nExamples:")
		fmt.Println("  mbzr comment -sha256 <hash> -comment \"This sample is part of a new campaign.\"")
	}

	if err := commentCmd.Parse(args); err != nil {
		return err
	}

	if *sha256 == "" || *comment == "" {
		printError("you must specify -sha256 and -comment flags")
		commentCmd.Usage()
		fmt.Println()
		return fmt.Errorf("missing required flags")
	}

	client, err := getAPIClient()
	if err != nil {
		printDetailedError(err, "Failed to create API client")
		return err
	}

	ctx, cancel := getContext()
	defer cancel()

	err = client.AddComment(ctx, *sha256, *comment)
	if err != nil {
		printDetailedError(err, fmt.Sprintf("Failed to add comment to %s", *sha256))
		return err
	}

	printSuccess(fmt.Sprintf("Successfully added comment to sample %s", *sha256))
	return nil
}
