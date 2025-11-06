package api

import (
	"encoding/json"
	"fmt"
)

// AddComment adds a comment to a sample identified by its hash
func AddComment(sha256, comment, apiKey string) error {
	data := map[string]string{
		"query":       "add_comment",
		"sha256_hash": sha256,
		"comment":     comment,
	}

	response, err := MakeRequest(data, nil, apiKey)
	if err != nil {
		return fmt.Errorf("Error adding comment to sample %s: %v", sha256, err)
	}

	// Parse the response
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return fmt.Errorf("Error parsing response commenting sample %s: %v", sha256, err)
	}

	// check for query_status
	if status, ok := result["query_status"].(string); ok {
		if status == "success" {
			fmt.Printf("Successfully added comment to sample %s: '%s'\n", sha256, comment)
			return nil
		}
		return fmt.Errorf("Failed to add comment to sample %s: %s", sha256, status)
	}

	return fmt.Errorf("Unexpected response format commenting sample %s", sha256)
}
