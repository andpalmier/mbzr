package api

import (
	"context"
	"encoding/json"
	"fmt"
)

// GetCSCB retrieves the MalwareBazaar Code Signing Certificate Blocklist
func (c *Client) GetCSCB(ctx context.Context) ([]CSCBEntry, error) {
	data := map[string]string{
		"query": "get_cscb",
	}

	response, err := c.MakeRequest(ctx, data, nil)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Code Signing Certificate Blocklist: %w", err)
	}

	var resp CSCBResponse
	if err := json.Unmarshal([]byte(response), &resp); err != nil {
		return nil, fmt.Errorf("error parsing CSCB response: %w", err)
	}

	if resp.QueryStatus != "ok" {
		return nil, fmt.Errorf("API returned status: %s", resp.QueryStatus)
	}

	return resp.Data, nil
}
