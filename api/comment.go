package api

import (
	"context"
	"encoding/json"
	"fmt"
)

// AddComment adds a comment to a sample identified by its hash
func (c *Client) AddComment(ctx context.Context, sha256, comment string) error {
	// Validate SHA256 format
	if err := ValidateSHA256(sha256); err != nil {
		return fmt.Errorf("invalid hash: %w", err)
	}

	data := map[string]string{
		"query":       "add_comment",
		"sha256_hash": sha256,
		"comment":     comment,
	}

	response, err := c.MakeRequest(ctx, data, nil)
	if err != nil {
		return fmt.Errorf("error adding comment to sample %s: %w", sha256, err)
	}

	// Parse the response
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return fmt.Errorf("error parsing response: %w", err)
	}

	// check for query_status
	if status, ok := result["query_status"].(string); ok {
		if status == "success" {
			return nil
		}
		return fmt.Errorf("failed to add comment: %s", status)
	}

	return fmt.Errorf("unexpected response format")
}
