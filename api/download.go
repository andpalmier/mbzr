package api

import (
	"fmt"
	"os"
)

// DownloadSample downloads a sample by its hash
func DownloadSample(sha256 string, apiKey string) error {
	data := map[string]string{
		"query":       "get_file",
		"sha256_hash": sha256,
	}

	response, err := MakeRequest(data, nil, apiKey)
	if err != nil {
		return fmt.Errorf("Error: downloading sample: %v", err)
	}

	// save file as zip
	fileName := fmt.Sprintf("%s.zip", sha256)
	err = os.WriteFile(fileName, []byte(response), 0644)
	if err != nil {
		return fmt.Errorf("Error: saving downloaded file: %v", err)
	}

	fmt.Printf("File downloaded successfully: %s\n", fileName)
	return nil
}
