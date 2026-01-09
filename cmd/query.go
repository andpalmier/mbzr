package cmd

import (
	"context"
	"flag"
	"fmt"

	"github.com/andpalmier/mbzr/api"
)

// executeQuery handles the 'query' subcommand
func executeQuery(args []string) error {
	queryCmd := flag.NewFlagSet("query", flag.ExitOnError)
	hash := queryCmd.String("hash", "", "Query by hash (SHA256, SHA1, MD5)")
	tag := queryCmd.String("tag", "", "Query by tag")
	signature := queryCmd.String("signature", "", "Query by signature")
	filetype := queryCmd.String("file_type", "", "Query by file type")
	clamav := queryCmd.String("clamav", "", "Query by ClamAV signature")
	imphash := queryCmd.String("imphash", "", "Query by Imphash")
	tlsh := queryCmd.String("tlsh", "", "Query by TLSH")
	telfhash := queryCmd.String("telfhash", "", "Query by Telfhash")
	dhash := queryCmd.String("dhash", "", "Query by Dhash")
	gimphash := queryCmd.String("gimphash", "", "Query by Gimphash")
	yara := queryCmd.String("yara", "", "Query by YARA rule")
	issuerCN := queryCmd.String("issuer_cn", "", "Query by Issuer Common Name")
	subjectCN := queryCmd.String("subject_cn", "", "Query by Subject Common Name")
	serialNumber := queryCmd.String("serial_number", "", "Query by Serial Number")
	limit := queryCmd.Int("limit", 100, "Limit the number of results")

	if len(args) < 1 {
		printError("expected query arguments")
		queryCmd.Usage()
		return fmt.Errorf("expected query arguments")
	}

	if err := queryCmd.Parse(args); err != nil {
		return err
	}

	// Map flags to their values
	queryParams := map[string]*string{
		"hash":          hash,
		"tag":           tag,
		"signature":     signature,
		"file_type":     filetype,
		"clamav":        clamav,
		"imphash":       imphash,
		"tlsh":          tlsh,
		"telfhash":      telfhash,
		"dhash":         dhash,
		"gimphash":      gimphash,
		"yara":          yara,
		"issuer_cn":     issuerCN,
		"subject_cn":    subjectCN,
		"serial_number": serialNumber,
	}

	var selectedQuery string
	for key, val := range queryParams {
		if *val != "" {
			selectedQuery = key
			break
		}
	}

	if selectedQuery == "" {
		return fmt.Errorf("please provide a query parameter (e.g., -hash, -tag)")
	}

	client, err := getAPIClient()
	if err != nil {
		// API key is optional for some queries, so we create a client without it
		// But if we fail to get it from environment, we still need a client
		// Some queries MIGHT require API key, but client will handle 401
		client = api.NewClient("")
	}

	// Define handlers using the client methods
	queryHandlers := map[string]func(context.Context, *api.Client, string, int) ([]api.MalwareSample, error){
		"hash": func(ctx context.Context, c *api.Client, val string, lim int) ([]api.MalwareSample, error) {
			return c.QueryByHash(ctx, val, lim)
		},
		"tag": func(ctx context.Context, c *api.Client, val string, lim int) ([]api.MalwareSample, error) {
			return c.QueryByTag(ctx, val, lim)
		},
		"signature": func(ctx context.Context, c *api.Client, val string, lim int) ([]api.MalwareSample, error) {
			return c.QueryBySignature(ctx, val, lim)
		},
		"file_type": func(ctx context.Context, c *api.Client, val string, lim int) ([]api.MalwareSample, error) {
			return c.QueryByFileType(ctx, val, lim)
		},
		"clamav": func(ctx context.Context, c *api.Client, val string, lim int) ([]api.MalwareSample, error) {
			return c.QueryByClamAV(ctx, val, lim)
		},
		"imphash": func(ctx context.Context, c *api.Client, val string, lim int) ([]api.MalwareSample, error) {
			return c.QueryByImpHash(ctx, val, lim)
		},
		"tlsh": func(ctx context.Context, c *api.Client, val string, lim int) ([]api.MalwareSample, error) {
			return c.QueryByTLSH(ctx, val, lim)
		},
		"telfhash": func(ctx context.Context, c *api.Client, val string, lim int) ([]api.MalwareSample, error) {
			return c.QueryByTelfHash(ctx, val, lim)
		},
		"dhash": func(ctx context.Context, c *api.Client, val string, lim int) ([]api.MalwareSample, error) {
			return c.QueryByDHash(ctx, val, lim)
		},
		"gimphash": func(ctx context.Context, c *api.Client, val string, lim int) ([]api.MalwareSample, error) {
			return c.QueryByGimphash(ctx, val, lim)
		},
		"yara": func(ctx context.Context, c *api.Client, val string, lim int) ([]api.MalwareSample, error) {
			return c.QueryByYara(ctx, val, lim)
		},
		"issuer_cn": func(ctx context.Context, c *api.Client, val string, lim int) ([]api.MalwareSample, error) {
			return c.QueryByIssuerCN(ctx, val, lim)
		},
		"subject_cn": func(ctx context.Context, c *api.Client, val string, lim int) ([]api.MalwareSample, error) {
			return c.QueryBySubjectCN(ctx, val, lim)
		},
		"serial_number": func(ctx context.Context, c *api.Client, val string, lim int) ([]api.MalwareSample, error) {
			return c.QueryBySerialNumber(ctx, val, lim)
		},
	}

	handler, exists := queryHandlers[selectedQuery]
	if !exists {
		printError(fmt.Sprintf("no handler for query type %s", selectedQuery))
		return fmt.Errorf("no handler for query type %s", selectedQuery)
	}

	ctx, cancel := getContext()
	defer cancel()

	var results []api.MalwareSample
	if selectedQuery == "hash" {
		// Hash query usually returns one result, limit might not apply or is 1
		results, err = handler(ctx, client, *queryParams[selectedQuery], 0)
	} else {
		results, err = handler(ctx, client, *queryParams[selectedQuery], *limit)
	}

	if err != nil {
		printDetailedError(err, fmt.Sprintf("Failed to query by %s", selectedQuery))
		return err
	}

	printJSON(results)
	return nil
}

// executeLatest handles the 'latest' subcommand
func executeLatest(args []string) error {
	latestCmd := flag.NewFlagSet("latest", flag.ExitOnError)
	selector := latestCmd.String("selector", "time", "Selector for latest samples: 'time' (last 60m) or '100' (last 100 samples)")

	latestCmd.Usage = func() {
		printUsageHeader("latest", "Retrieves the latest malware samples added to MalwareBazaar.")
		fmt.Println("\nFlags:")
		fmt.Println("  -selector <value>\tSelector: 'time' (default, last 60m) or '100' (last 100 samples)")
		fmt.Println("\nExamples:")
		fmt.Println("  mbzr latest")
		fmt.Println("  mbzr latest -selector 100")
	}

	if err := latestCmd.Parse(args); err != nil {
		return err
	}

	if *selector != "time" && *selector != "100" {
		printError("invalid selector: must be 'time' or '100'")
		latestCmd.Usage()
		fmt.Println()
		return fmt.Errorf("invalid selector")
	}

	client, err := getAPIClient()
	if err != nil {
		printDetailedError(err, "Failed to create API client")
		return err
	}

	ctx, cancel := getContext()
	defer cancel()

	results, err := client.QueryLatest(ctx, *selector)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		fmt.Println("No latest samples found.")
		return nil
	}

	printJSON(results)

	return nil
}

// executeRecentDetections handles the 'recent_detections' subcommand
func executeRecentDetections(args []string) error {
	recentCmd := flag.NewFlagSet("recent_detections", flag.ExitOnError)
	hours := recentCmd.Int("hours", 24, "Retrieve recent detections from the last n hours (default: 24)")

	recentCmd.Usage = func() {
		printUsageHeader("recent_detections", "Retrieves malware samples detected in the last specified number of hours.")
		fmt.Println("\nFlags:")
		fmt.Println("  -hours <number>\tNumber of hours to look back (default: 24)")
		fmt.Println("\nExamples:")
		fmt.Println("  mbzr recent_detections -hours 24")
	}

	if err := recentCmd.Parse(args); err != nil {
		return err
	}

	if *hours < 1 || *hours > 168 {
		printError("hours must be between 1 and 168")
		recentCmd.Usage()
		fmt.Println()
		return fmt.Errorf("invalid hours value")
	}

	client, err := getAPIClient()
	if err != nil {
		printDetailedError(err, "Failed to create API client")
		return err
	}

	ctx, cancel := getContext()
	defer cancel()

	results, err := client.GetRecentDetections(ctx, *hours)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		fmt.Println("No recent detections found.")
		return nil
	}

	printJSON(results)

	return nil
}

// executeCSCB handles the 'cscb' subcommand
func executeCSCB(args []string) error {
	cscbCmd := flag.NewFlagSet("cscb", flag.ExitOnError)

	cscbCmd.Usage = func() {
		printUsageHeader("cscb", "Queries the Code Signing Certificate Blocklist (CSCB) from MalwareBazaar.")
		fmt.Println("\nExample:")
		fmt.Println("  mbzr cscb")
	}

	if err := cscbCmd.Parse(args); err != nil {
		return err
	}

	client, err := getAPIClient()
	if err != nil {
		printDetailedError(err, "Failed to create API client")
		return err
	}

	ctx, cancel := getContext()
	defer cancel()

	results, err := client.GetCSCB(ctx)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		fmt.Println("No CSCB entries found.")
		return nil
	}

	printJSON(results)

	return nil
}
