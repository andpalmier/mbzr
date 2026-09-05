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

	var result struct {
		QueryStatus string `json:"query_status"`
	}
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return fmt.Errorf("error parsing response: %w", err)
	}

	if result.QueryStatus == "" {
		return fmt.Errorf("unexpected response format")
	}

	return newStatusError(result.QueryStatus, "add_comment")
}
