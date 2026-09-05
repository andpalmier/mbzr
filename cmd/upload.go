package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/andpalmier/mbzr/api"
)

// executeUpload handles the 'upload' subcommand
func executeUpload(args []string) error {
	uploadCmd := flag.NewFlagSet("upload", flag.ExitOnError)
	file := uploadCmd.String("file", "", "File to upload")
	dir := uploadCmd.String("dir", "", "Directory containing files to upload")
	tags := uploadCmd.String("tags", "", "Comma separated list of tags associated with the files to upload")
	anonymous := uploadCmd.Bool("anonymous", false, "Upload files anonymously (no user association)")
	deliveryMethod := uploadCmd.String("delivery_method", "", "How the sample was distributed: "+strings.Join(deliveryMethods, ", "))
	references := newKVFlag("reference", referenceKeys)
	uploadCmd.Var(references, "reference", "Reference as key=value, repeatable (e.g. -reference any_run=https://...)")
	contextInfo := newKVFlag("context", contextKeys)
	uploadCmd.Var(contextInfo, "context", "Context as key=value, repeatable (e.g. -context dropped_by_malware=Gozi)")

	uploadCmd.Usage = func() {
		printUsageHeader("upload", "Uploads a file or all files in a directory to MalwareBazaar.")
		fmt.Println("\nFlags:")
		fmt.Println("  -file <file_path>\t\tFile to upload")
		fmt.Println("  -dir <directory_path>\tDirectory containing files to upload")
		fmt.Println("  -tags <tag1,tag2,...>\tComma separated list of tags associated with the files to upload")
		fmt.Println("  -anonymous\t\t\tUpload files anonymously (no user association)")
		fmt.Println("  -delivery_method <method>\tOne of: " + strings.Join(deliveryMethods, ", "))
		fmt.Println("  -reference <key=value>\tRepeatable. Keys: " + strings.Join(referenceKeys, ", "))
		fmt.Println("  -context <key=value>\t\tRepeatable. Keys: " + strings.Join(contextKeys, ", "))
		fmt.Println("\n📖 Examples:")
		fmt.Println("  # Upload a single file with tags")
		fmt.Println("  mbzr upload -file malware.exe -tags trojan,banker")
		fmt.Println()
		fmt.Println("  # Upload all files in a directory anonymously")
		fmt.Println("  mbzr upload -dir /path/to/samples -anonymous")
		fmt.Println()
		fmt.Println("  # Upload with provenance")
		fmt.Println("  mbzr upload -file loader.exe -delivery_method email_attachment \\")
		fmt.Println("    -reference any_run=https://app.any.run/tasks/1 \\")
		fmt.Println("    -context dropped_by_malware=Gozi")
	}

	if err := uploadCmd.Parse(args); err != nil {
		return err
	}

	// Parse tags
	var tagList []string
	if *tags != "" {
		tagList = strings.Split(*tags, ",")
	}

	if *deliveryMethod != "" && !isDeliveryMethod(*deliveryMethod) {
		printError(fmt.Sprintf("invalid delivery method %q. Allowed: %s", *deliveryMethod, strings.Join(deliveryMethods, ", ")))
		uploadCmd.Usage()
		return fmt.Errorf("invalid delivery method")
	}

	// validate input
	if *file == "" && *dir == "" {
		printError("you must specify either a file (-file) or a directory (-dir) to upload")
		uploadCmd.Usage()
		return fmt.Errorf("missing file or directory argument")
	}

	if *file != "" && *dir != "" {
		printError("you cannot specify both a file and a directory")
		uploadCmd.Usage()
		return fmt.Errorf("cannot specify both file and directory")
	}

	client, err := getAPIClient()
	if err != nil {
		printDetailedError(err, "Failed to create API client")
		return err
	}

	opts := api.UploadOptions{
		Anonymous:      *anonymous,
		Tags:           tagList,
		DeliveryMethod: *deliveryMethod,
	}
	if !references.empty() {
		opts.References = references.pairs
	}
	if !contextInfo.empty() {
		opts.Context = contextValues(contextInfo)
	}

	// file upload
	if *file != "" {
		if verbose {
			printVerbose(fmt.Sprintf("Uploading file: %s", *file))
		}

		ctx, cancel := getContext()
		defer cancel()

		result, err := client.UploadFile(ctx, *file, opts)
		if err != nil {
			printDetailedError(err, fmt.Sprintf("Failed to upload file: %s", *file))
			return err
		}
		printSuccess(fmt.Sprintf("Uploaded: %s - Status: %s", *file, result.QueryStatus))
		if len(result.Data) > 0 {
			printJSON(result.Data)
		}
		return nil
	}

	// directory upload
	if verbose {
		printVerbose(fmt.Sprintf("Uploading directory: %s", *dir))
	}

	// Walk through directory and upload files
	failed := 0
	err = filepath.Walk(*dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip subdirectories and hidden files
		if info.IsDir() || strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		if verbose {
			printVerbose(fmt.Sprintf("Uploading file: %s", path))
		}

		// Each file gets its own timeout, so a large directory does not run
		// out of time partway through the walk.
		ctx, cancel := getContext()
		defer cancel()

		result, err := client.UploadFile(ctx, path, opts)
		if err != nil {
			failed++
			printError(fmt.Sprintf("Failed to upload %s: %v", path, err))
			return nil
		}
		printSuccess(fmt.Sprintf("Uploaded: %s - Status: %s", path, result.QueryStatus))
		return nil
	})

	if err != nil {
		printDetailedError(err, fmt.Sprintf("Failed to walk directory: %s", *dir))
		return err
	}

	if failed > 0 {
		return fmt.Errorf("%d file(s) failed to upload", failed)
	}
	return nil
}

// deliveryMethods, referenceKeys and contextKeys mirror the values documented
// at https://bazaar.abuse.ch/api/
var (
	deliveryMethods = []string{"email_attachment", "email_link", "web_download", "web_drive-by", "multiple", "other"}
	referenceKeys   = []string{"urlhaus", "any_run", "joe_sandbox", "malpedia", "twitter", "links"}
	contextKeys     = []string{"dropped_by_md5", "dropped_by_sha256", "dropped_by_malware", "dropping_md5", "dropping_sha256", "dropping_malware", "comment"}
)

func isDeliveryMethod(m string) bool {
	for _, d := range deliveryMethods {
		if d == m {
			return true
		}
	}
	return false
}

// contextValues shapes the context block the way the API documents it: every
// key carries a list, except "comment", which is a single string.
func contextValues(k *kvFlag) map[string]any {
	out := make(map[string]any, len(k.pairs))
	for key, vals := range k.pairs {
		if key == "comment" {
			out[key] = strings.Join(vals, " ")
			continue
		}
		out[key] = vals
	}
	return out
}
