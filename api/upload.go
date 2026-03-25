package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// UploadFile uploads a file to MalwareBazaar
func (c *Client) UploadFile(ctx context.Context, filePath string, anonymous bool, tags []string, deliveryMethod string, contextInfo map[string]string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("error opening file: %w", err)
	}
	defer func() { _ = file.Close() }()

	files := map[string]io.Reader{
		"file": file,
	}

	// Prepare other form fields
	data := map[string]string{
		"anonymous": "0",
	}
	if anonymous {
		data["anonymous"] = "1"
	}

	if len(tags) > 0 {
		data["tags"] = strings.Join(tags, ",")
	}

	if deliveryMethod != "" {
		data["delivery_method"] = deliveryMethod
	}

	// Add context info if provided
	if contextInfo != nil {
		for k, v := range contextInfo {
			data[k] = v
		}
	}

	response, err := c.MakeRequest(ctx, data, files)
	if err != nil {
		return "", fmt.Errorf("error uploading file: %w", err)
	}

	// Parse the response to check status
	resp, err := ParseAPIResponse([]byte(response))
	if err != nil {
		return "", err
	}

	// If there's data, try to pretty print it
	if len(resp.Data) > 0 {
		dataJSON, err := json.MarshalIndent(resp.Data, "", "    ")
		if err == nil {
			return fmt.Sprintf("Query Status: %s\nData:\n%s", resp.QueryStatus, string(dataJSON)), nil
		}
	}

	return fmt.Sprintf("Query Status: %s", resp.QueryStatus), nil
}
