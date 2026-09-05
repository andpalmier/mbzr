package api

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// UploadOptions carries the metadata submitted alongside a sample.
// DeliveryMethod, References and Context are part of the documented API but
// are not yet exposed as CLI flags.
type UploadOptions struct {
	Anonymous      bool
	Tags           []string
	DeliveryMethod string
	References     map[string][]string
	Context        map[string]any
}

// UploadResult is the parsed outcome of a successful submission
type UploadResult struct {
	QueryStatus string          `json:"query_status"`
	Data        []MalwareSample `json:"data,omitempty"`
}

// uploadEnvelope is the json_data part of the multipart submission
type uploadEnvelope struct {
	Anonymous      int                 `json:"anonymous"`
	Tags           []string            `json:"tags,omitempty"`
	DeliveryMethod string              `json:"delivery_method,omitempty"`
	References     map[string][]string `json:"references,omitempty"`
	Context        map[string]any      `json:"context,omitempty"`
}

// UploadFile uploads a file to MalwareBazaar.
// The API expects a multipart POST carrying the sample as "file" and its
// metadata as a JSON object in "json_data".
func (c *Client) UploadFile(ctx context.Context, filePath string, opts UploadOptions) (*UploadResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer func() { _ = file.Close() }()

	envelope := uploadEnvelope{
		Tags:           opts.Tags,
		DeliveryMethod: opts.DeliveryMethod,
		References:     opts.References,
		Context:        opts.Context,
	}
	if opts.Anonymous {
		envelope.Anonymous = 1
	}

	jsonData, err := json.Marshal(envelope)
	if err != nil {
		return nil, fmt.Errorf("encoding upload metadata: %w", err)
	}

	files := map[string]FormFile{
		"file": {Name: filepath.Base(filePath), Reader: file},
	}
	data := map[string]string{
		"json_data": string(jsonData),
	}

	response, err := c.MakeRequest(ctx, data, files)
	if err != nil {
		return nil, fmt.Errorf("error uploading file: %w", err)
	}

	resp, err := ParseAPIResponse([]byte(response))
	if err != nil {
		return nil, fmt.Errorf("error parsing upload response: %w", err)
	}

	if resp.QueryStatus == "" {
		return nil, fmt.Errorf("upload failed: response contained no query_status")
	}
	if err := newStatusError(resp.QueryStatus, "upload"); err != nil {
		return nil, err
	}

	return &UploadResult{QueryStatus: resp.QueryStatus, Data: resp.Data}, nil
}
