package cmd

import (
	"errors"
	"flag"
	"fmt"

	"github.com/andpalmier/mbzr/api"
)

// executeUpload handles the 'upload' subcommand
func executeUpload(args []string) error {
	uploadCmd := flag.NewFlagSet("upload", flag.ExitOnError)
	file := uploadCmd.String("file", "", "File to upload")
	dir := uploadCmd.String("dir", "", "Directory containing files to upload")
	tags := uploadCmd.String("tags", "", "Comma separated list of tags associated with the files to upload")
	anonymous := uploadCmd.Bool("anonymous", false, "Upload files anonymously (no user association)")

	// helper function
	uploadCmd.Usage = func() {
		fmt.Println("Usage:\n  mbzr upload [flags]\n")
		fmt.Println("Uploads a file or all files in a directory to MalwareBazaar.")
		fmt.Println("\nFlags:")
		fmt.Println("  -file <file_path>		File to upload")
		fmt.Println("  -dir <directory_path>	Directory containing files to upload")
		fmt.Println("  -tags <tag1,tag2,...>	Comma separated list of tags associated with the files to upload")
		fmt.Println("  -anonymous				Upload files anonymously (no user association)")
		fmt.Println("\nExamples:")
		fmt.Println("  mbzr upload -file sample.exe -tags trojan,banker")
		fmt.Println("  mbzr upload -dir /path/to/malware_samples -anonymous")
	}

	uploadCmd.Parse(args)

	// validate input
	if *file == "" && *dir == "" {
		fmt.Println("Error: please specify a file to upload using -file or a directory using -dir\n")
		uploadCmd.Usage()
		fmt.Println()
		return errors.New("no file or directory specified for upload")
	}

	if *file != "" && *dir != "" {
		fmt.Println("Error: please specify either -file or -dir, not both\n")
		uploadCmd.Usage()
		fmt.Println()
		return errors.New("both file and directory specified for upload")
	}

	// file upload
	if *file != "" {
		return api.UploadFile(*file, *tags, *anonymous, apiKey)
	}

	if *dir != "" {
		return api.UploadDirectory(*dir, *tags, *anonymous, apiKey)
	}

	return nil
}
