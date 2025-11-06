package cmd

import (
	"flag"
	"fmt"

	"github.com/andpalmier/mbzr/api"
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

	// helper function
	updateCmd.Usage = func() {
		fmt.Println("Usage:\n  mbzr update [flags]\n")
		fmt.Println("Updates metadata for a malware sample in MalwareBazaar.")
		fmt.Println("\nFlags:")
		fmt.Println("  -sha256 <sha256_hash>	sha256 hash of the file to update")
		fmt.Println("  -key <key>				Key to update (e.g., add_tag, remove_tag, urlhaus, any_run, joe_sandbox, malpedia, twitter, links, dropped_by_md5, dropped_by_sha256, dropped_by_malware, dropping_md5, dropping_sha256, dropping_malware, comment)")
		fmt.Println("  -value <new_value>		New value for the specified key")
		fmt.Println("\nExamples:")
		fmt.Println("  mbzr update -sha256 <hash> -key add_tag -value ransomware")
		fmt.Println("  mbzr update -sha256 <hash> -key comment -value \"This sample is part of a new campaign.\"")
	}

	updateCmd.Parse(args)

	// validate the flags
	if *sha256 == "" || *key == "" || *value == "" {
		fmt.Println("You must specify -sha256, -key and -value flags\n")
		updateCmd.Usage()
		fmt.Println()
		return fmt.Errorf("missing required flags")
	}

	// validate key against allowed choices
	if !isValidKey(*key) {
		fmt.Printf("Error: invalid key '%s'. Allowed keys are: %v\n\n", *key, AllowedKeys)
		updateCmd.Usage()
		fmt.Println()
		return fmt.Errorf("invalid key specified")
	}

	return api.UpdateSample(*sha256, *key, *value, apiKey)
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
