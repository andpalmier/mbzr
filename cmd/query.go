package cmd

import (
	"errors"
	"flag"
	"fmt"

	"github.com/andpalmier/mbzr/api"
)

// QueryOption query option structure with flag, description and example
type QueryOption struct {
	Flag        string
	Description string
	Example     string
}

// executeQuery handles the 'query' subcommand
func executeQuery(args []string) error {
	queryCmd := flag.NewFlagSet("query", flag.ExitOnError)

	// define query options
	queryOptions := []QueryOption{
		{"hash", "Retrieve info about a malware sample by its hash (sha1, sha256 or md5).", "mbzr query -hash <file_hash>"},
		{"tag", "Query malware sample associated with a tag.", "mbzr query -tag TrickBot -limit 10"},
		{"signature", "Query malware samples associated with a signature.", "mbzr query -signature 'Emotet' -limit 10"},
		{"filetype", "Query malware samples by filetype.", "mbzr query -filetype 'exe' -limit 10"},
		{"clamav", "Query malware samples by ClamAV signature.", "mbzr query -clamav 'Win.Trojan.Emotet-1234567' -limit 10"},
		{"imphash", "Query malware samples by imphash.", "mbzr query -imphash <imphash_value> -limit 10"},
		{"tlsh", "Query malware samples by tlsh hash.", "mbzr query -tlsh <tlsh_value> -limit 10"},
		{"telfhash", "Query malware samples by telfhash.", "mbzr query -telfhash <telfhash_value> -limit 10"},
		{"dhash", "Query malware samples by dhash icon.", "mbzr query -dhash <dhash_value> -limit 10"},
		{"gimphash", "Query malware samples by gimphash.", "mbzr query -gimphash <gimphash_value> -limit 10"},
		{"yara", "Query malware samples by YARA rule name.", "mbzr query -yara 'NETexecutableMicrosoft' -limit 10"},
		{"cert_issuer", "Query code signing certificates by Issuer Common Name.", "mbzr query -cert_issuer 'Sectigo RSA Code Signing CA' -limit 10"},
		{"cert_subject", "Query code signing certificates by Subject Common Name.", "mbzr query -cert_subject 'Microsoft Corporation' -limit 10"},
		{"cert_serial", "Query code signing certificates by Serial Number.", "mbzr query -cert_serial <cert_serial> -limit 10"},
	}

	// define flags
	queryParams := make(map[string]*string)
	for _, option := range queryOptions {
		queryParams[option.Flag] = queryCmd.String(option.Flag, "", option.Description)
	}
	limit := queryCmd.Int("limit", 100, "Number of results to return (default: 100, max: 1000)")

	// override help message
	queryCmd.Usage = func() {
		fmt.Printf("Usage:\n  mbzr query [flags]\n\n")
		fmt.Println("Query MalwareBazaar for information about malware samples.")
		fmt.Println("\nQuery options:")
		for _, option := range queryOptions {
			fmt.Printf("  -%-15s %s\n", option.Flag, option.Description)
		}
		fmt.Println("\nOptional flags:")
		fmt.Println("  -limit        Number of results to return (default: 100, max: 1000)")

		fmt.Println("\nExamples:")
		for _, option := range queryOptions {
			fmt.Printf("  %s\n", option.Example)
		}
	}

	queryCmd.Parse(args)

	// only one query allowed
	var selectedQuery string
	for key, value := range queryParams {
		if *value != "" {
			if selectedQuery != "" {
				fmt.Println("Error: please specify only one query parameter at a time")
				queryCmd.Usage()
				return fmt.Errorf("multiple query parameters specified")
			}
			selectedQuery = key
		}
	}

	if selectedQuery == "" {
		fmt.Println("Error: please specify one query flag")
		queryCmd.Usage()
		return fmt.Errorf("no query parameter specified")
	}

	if *limit < 1 || *limit > 1000 {
		fmt.Println("Error: limit must be between 1 and 1000\n")
		queryCmd.Usage()
		fmt.Println()
		return errors.New("limit must be between 1 and 1000")
	}

	// Map query types with the corresponding API function
	queryHandlers := map[string]func(string, int, string) error{
		"hash":         api.QueryByHash,
		"tag":          api.QueryByTag,
		"signature":    api.QueryBySignature,
		"filetype":     api.QueryByFileType,
		"clamav":       api.QueryByClamAV,
		"imphash":      api.QueryByImpHash,
		"tlsh":         api.QueryByTLSH,
		"telfhash":     api.QueryByTelfHash,
		"dhash":        api.QueryByDHash,
		"gimphash":     api.QueryByGimphash,
		"yara":         api.QueryByYara,
		"cert_issuer":  api.QueryByIssuerCN,
		"cert_subject": api.QueryBySubjectCN,
		"cert_serial":  api.QueryBySerialNumber,
	}

	// Execute the appropriate query handler
	handler, exists := queryHandlers[selectedQuery]
	if !exists {
		return fmt.Errorf("no handler for query type %s", selectedQuery)
	}

	// for "hash" limit is not applicable
	if selectedQuery == "hash" {
		return handler(*queryParams[selectedQuery], 0, apiKey)
	}

	return handler(*queryParams[selectedQuery], *limit, apiKey)
}
