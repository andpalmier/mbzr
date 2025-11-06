package api

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// UploadDirectory uploads all files from the specified directory
func UploadDirectory(dirPath, tags string, anonymous bool, apiKey string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return fmt.Errorf("Error: directory %s does not exist", dirPath)
	}

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("Error accesssing path %s: %v", path, err)
		}

		// skip subdirectories
		if !info.IsDir() {
			fmt.Printf("Uploading file %s\n", path)
			if err := UploadFile(path, tags, anonymous, apiKey); err != nil {
				fmt.Printf("Error uploading file %s: %v\n", path, err)
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error accessing directory %s: %v", dirPath, err)
	}
	return nil
}

// UploadFile uploads a single file
func UploadFile(filePath, tags string, anonymous bool, apiKey string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("Error: opening file failed: %v", err)
	}
	defer file.Close()

	// prepare JSON data
	data := map[string]interface{}{
		"tags":      strings.Split(tags, ","),
		"anonymous": anonymous,
	}

	// prepare multipart form data
	files := map[string]io.Reader{
		"file":      file,
		"json_data": bytes.NewReader(mustMarshalJSON(data)),
	}

	// send request
	response, err := MakeRequest(nil, files, apiKey)
	if err != nil {
		return fmt.Errorf("Error uploading file %s: %v", filePath, err)
	}

	fmt.Printf("File uploaded successfully: %s\n", filePath)
	dataJSON, err := PrintData(response)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println(dataJSON)
	}
	return nil
}
