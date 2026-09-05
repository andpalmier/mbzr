package api

import (
	"context"
	"encoding/json"
	"fmt"
)

// StatusUpdated is returned when an entry was changed.
// StatusExists is returned when the key/value was already present, which the
// API reports as a distinct, non-fatal outcome.
const (
	StatusUpdated = "updated"
	StatusExists  = "exists"
)

// UpdateSample updates an existing entry in MalwareBazaar.
// It returns the query_status reported by the API on success, which is either
// "updated" or "exists".
func (c *Client) UpdateSample(ctx context.Context, sha256, key, value string) (string, error) {
	// Validate SHA256 format
	if err := ValidateSHA256(sha256); err != nil {
		return "", fmt.Errorf("invalid hash: %w", err)
	}

	data := map[string]string{
		"query":       "update",
		"sha256_hash": sha256,
		"key":         key,
		"value":       value,
	}

	response, err := c.MakeRequest(ctx, data, nil)
	if err != nil {
		return "", fmt.Errorf("error updating sample %s: %w", sha256, err)
	}

	var result struct {
		QueryStatus string `json:"query_status"`
	}
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return "", fmt.Errorf("error parsing response: %w", err)
	}

	switch result.QueryStatus {
	// "updated" is the documented success value; "ok" and "success" are
	// accepted defensively since the API is not consistent across endpoints.
	case StatusUpdated, "ok", "success":
		return StatusUpdated, nil
	case StatusExists:
		return StatusExists, nil
	case "":
		return "", fmt.Errorf("unexpected response format")
	default:
		return "", newStatusError(result.QueryStatus, "update")
	}
}
