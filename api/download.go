package api

import (
	"bytes"
	"encoding/json"
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

	respBytes := []byte(response)

	// ZIP files start with PK\x03\x04
	if len(respBytes) >= 4 && respBytes[0] == 'P' && respBytes[1] == 'K' && respBytes[2] == 3 && respBytes[3] == 4 {
		fileName := fmt.Sprintf("%s.zip", sha256)
		if err := os.WriteFile(fileName, respBytes, 0644); err != nil {
			return fmt.Errorf("saving downloaded file: %w", err)
		}
		fmt.Printf("File downloaded successfully: %s\n", fileName)
		return nil
	}

	// Try to parse JSON error response
	var js map[string]interface{}
	if err := json.Unmarshal(respBytes, &js); err == nil {
		// common key used by the API
		if v, ok := js["query_status"]; ok {
			return fmt.Errorf("%v", v)
		}
		// generic JSON error formatting
		if b, err := json.MarshalIndent(js, "", "  "); err == nil {
			return fmt.Errorf("api response: %s", string(b))
		}
	}

	// Fallback: check for known plain-text error tokens
	if bytes.Contains(respBytes, []byte("file_not_found")) ||
		bytes.Contains(respBytes, []byte("no_sha256_hash")) ||
		bytes.Contains(respBytes, []byte("illegal_sha256_hash")) ||
		bytes.Contains(respBytes, []byte("query_status")) {
		return fmt.Errorf("api error: %s", string(respBytes))
	}

	// Unknown non-zip response
	return fmt.Errorf("unexpected response when downloading sample: %s", string(respBytes))
}
