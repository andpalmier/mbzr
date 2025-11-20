package cmd

import (
	"flag"
	"fmt"
)

// AllowedKeys list of valid keys for updating an entry
var AllowedKeys = []string{
	"add_tag", "remove_tag", "urlhaus", "any_run", "joe_sandbox", "malpedia", "twitter", "links", "dropped_by_md5",
	"dropped_by_sha256", "dropped_by_malware", "dropping_md5", "dropping_sha256", "dropping_malware", "comment",
}

// executeUpdate handles the 'update' subcommand
func executeUpdate(args []string) error {
	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
	sha256 := updateCmd.String("sha256", "", "sha256 hash of the file to update")
	key := updateCmd.String("key", "", "Key to update (e.g., add_tag, remove_tag, urlhaus, any_run, joe_sandbox, malpedia, twitter, links, dropped_by_md5, dropped_by_sha256, dropped_by_malware, dropping_md5, dropping_sha256, dropping_malware, comment)")
	value := updateCmd.String("value", "", "New value for the specified key")

	updateCmd.Usage = func() {
		printUsageHeader("update", "Updates metadata for a malware sample in MalwareBazaar.")
		fmt.Println("\nFlags:")
		fmt.Println("  -sha256 <sha256_hash>\t\t\tsha256 hash of the file to update")
		fmt.Println("  -key <key>\t\t\t\tKey to update (e.g., add_tag, remove_tag, urlhaus, any_run, joe_sandbox, malpedia, twitter, links, dropped_by_md5, dropped_by_sha256, dropped_by_malware, dropping_md5, dropping_sha256, dropping_malware, comment)")
		fmt.Println("  -value <new_value>\t\t\tNew value for the specified key")
		fmt.Println("\nExamples:")
		fmt.Println("  mbzr update -sha256 <hash> -key add_tag -value ransomware")
		fmt.Println("  mbzr update -sha256 <hash> -key comment -value \"This sample is part of a new campaign.\"")
	}

	if err := updateCmd.Parse(args); err != nil {
		return err
	}

	if *sha256 == "" || *key == "" || *value == "" {
		printError("you must specify -sha256, -key and -value flags")
		updateCmd.Usage()
		fmt.Println()
		return fmt.Errorf("missing required flags")
	}

	if !isValidKey(*key) {
		printError(fmt.Sprintf("invalid key '%s'. Allowed keys are: %v", *key, AllowedKeys))
		updateCmd.Usage()
		fmt.Println()
		return fmt.Errorf("invalid key specified")
	}

	client, err := getAPIClient()
	if err != nil {
		printDetailedError(err, "Failed to create API client")
		return err
	}

	ctx, cancel := getContext()
	defer cancel()

	err = client.UpdateSample(ctx, *sha256, *key, *value)
	if err != nil {
		printDetailedError(err, fmt.Sprintf("Failed to update sample %s", *sha256))
		return err
	}

	printSuccess(fmt.Sprintf("Successfully updated sample %s", *sha256))
	return nil
}

// isValidKey checks if the provided key is in the list of allowed keys
func isValidKey(key string) bool {
	for _, allowedKey := range AllowedKeys {
		if key == allowedKey {
			return true
		}
	}
	return false
}
