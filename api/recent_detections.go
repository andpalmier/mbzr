package api

import (
	"context"
	"fmt"
)

// GetRecentDetections retrieves recent detections from the API
func (c *Client) GetRecentDetections(ctx context.Context, hours int) ([]MalwareSample, error) {
	data := map[string]string{
		"query": "recent_detections",
		"hours": fmt.Sprintf("%d", hours),
	}

	response, err := c.MakeRequest(ctx, data, nil)
	if err != nil {
		return nil, fmt.Errorf("error retrieving recent detections: %w", err)
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
