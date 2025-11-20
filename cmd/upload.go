package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// executeUpload handles the 'upload' subcommand
func executeUpload(args []string) error {
	uploadCmd := flag.NewFlagSet("upload", flag.ExitOnError)
	file := uploadCmd.String("file", "", "File to upload")
	dir := uploadCmd.String("dir", "", "Directory containing files to upload")
	tags := uploadCmd.String("tags", "", "Comma separated list of tags associated with the files to upload")
	anonymous := uploadCmd.Bool("anonymous", false, "Upload files anonymously (no user association)")

	uploadCmd.Usage = func() {
		printUsageHeader("upload", "Uploads a file or all files in a directory to MalwareBazaar.")
		fmt.Println("\nFlags:")
		fmt.Println("  -file <file_path>\t\tFile to upload")
		fmt.Println("  -dir <directory_path>\tDirectory containing files to upload")
		fmt.Println("  -tags <tag1,tag2,...>\tComma separated list of tags associated with the files to upload")
		fmt.Println("  -anonymous\t\t\tUpload files anonymously (no user association)")
		fmt.Println("\n📖 Examples:")
		fmt.Println("  # Upload a single file with tags")
		fmt.Println("  mbzr upload -file malware.exe -tags trojan,banker")
		fmt.Println()
		fmt.Println("  # Upload all files in a directory anonymously")
		fmt.Println("  mbzr upload -dir /path/to/samples -anonymous")
	}

	if err := uploadCmd.Parse(args); err != nil {
		return err
	}

	// Parse tags
	var tagList []string
	if *tags != "" {
		tagList = strings.Split(*tags, ",")
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

	// file upload
	if *file != "" {
		if verbose {
			printVerbose(fmt.Sprintf("Uploading file: %s", *file))
		}

		ctx, cancel := getContext()
		defer cancel()

		result, err := client.UploadFile(ctx, *file, *anonymous, tagList, "", nil)
		if err != nil {
			printDetailedError(err, fmt.Sprintf("Failed to upload file: %s", *file))
			return err
		}
		printSuccess(result)
		return nil
	}

	// directory upload
	if *dir != "" {
		if verbose {
			printVerbose(fmt.Sprintf("Uploading directory: %s", *dir))
		}

		ctx, cancel := getContext()
		defer cancel()

		// Walk through directory and upload files
		err := filepath.Walk(*dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Skip subdirectories and hidden files
			if info.IsDir() || info.Name()[0] == '.' {
				return nil
			}

			if verbose {
				printVerbose(fmt.Sprintf("Uploading file: %s", path))
			}

			result, err := client.UploadFile(ctx, path, *anonymous, tagList, "", nil)
			if err != nil {
				printError(fmt.Sprintf("Failed to upload %s: %v", path, err))
			} else {
				printSuccess(fmt.Sprintf("Uploaded: %s - Status: %s", path, result))
			}
			return nil
		})

		if err != nil {
			printDetailedError(err, fmt.Sprintf("Failed to walk directory: %s", *dir))
			return err
		}
		return nil
	}

	return nil
}
