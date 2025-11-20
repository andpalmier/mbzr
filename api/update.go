package api

import (
	"context"
	"encoding/json"
	"fmt"
)

// UpdateSample updates existing entry in MalwareBazaar
func (c *Client) UpdateSample(ctx context.Context, sha256, key, value string) error {
	// Validate SHA256 format
	if err := ValidateSHA256(sha256); err != nil {
		return fmt.Errorf("invalid hash: %w", err)
	}

	data := map[string]string{
		"query":       "update",
		"sha256_hash": sha256,
		"key":         key,
		"value":       value,
	}

	response, err := c.MakeRequest(ctx, data, nil)
	if err != nil {
		return fmt.Errorf("error updating sample %s: %w", sha256, err)
	}

	// Parse the response to check status
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return fmt.Errorf("error parsing response: %w", err)
	}

	// Check for query_status
	if status, ok := result["query_status"].(string); ok {
		if status == "ok" || status == "success" {
			return nil
		}
		return fmt.Errorf("failed to update sample: %s", status)
	}

	return fmt.Errorf("unexpected response format")
}
