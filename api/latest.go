package api

import (
	"context"
	"fmt"
)

// QueryLatest retrieves the latest malware samples added to MalwareBazaar
// selector can be "time" (last 60 minutes) or "100" (last 100 samples)
func (c *Client) QueryLatest(ctx context.Context, selector string) ([]MalwareSample, error) {
	if selector != "time" && selector != "100" {
		return nil, fmt.Errorf("invalid selector: must be 'time' or '100'")
	}

	data := map[string]string{
		"query":    "get_recent",
		"selector": selector,
	}

	response, err := c.MakeRequest(ctx, data, nil)
	if err != nil {
		return nil, fmt.Errorf("error retrieving latest samples: %w", err)
	}

	resp, err := ParseAPIResponse([]byte(response))
	if err != nil {
		return nil, err
	}

	if resp.QueryStatus != "ok" {
		return nil, fmt.Errorf("API returned status: %s", resp.QueryStatus)
	}

	return resp.Data, nil
}
