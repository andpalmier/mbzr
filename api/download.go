package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// DownloadSample downloads a sample by its hash
func (c *Client) DownloadSample(ctx context.Context, sha256 string) error {
	// Validate SHA256 format to prevent path traversal
	if err := ValidateSHA256(sha256); err != nil {
		return fmt.Errorf("invalid hash: %w", err)
	}

	data := map[string]string{
		"query":       "get_file",
		"sha256_hash": sha256,
	}

	body, err := c.MakeRequestRaw(ctx, data, nil)
	if err != nil {
		return fmt.Errorf("error downloading sample: %w", err)
	}
	defer body.Close()

	// Read first 4 bytes to check for ZIP header
	header := make([]byte, 4)
	n, err := body.Read(header)
	if err != nil && err != io.EOF {
		return fmt.Errorf("error reading response header: %w", err)
	}

	// Check if it's a ZIP file (PK\x03\x04)
	if n >= 4 && header[0] == 'P' && header[1] == 'K' && header[2] == 3 && header[3] == 4 {
		fileName := fmt.Sprintf("%s.zip", sha256)

		// Check if file already exists
		if _, err := os.Stat(fileName); err == nil {
			return fmt.Errorf("file already exists: %s", fileName)
		}

		out, err := os.Create(fileName)
		if err != nil {
			return fmt.Errorf("error creating file: %w", err)
		}
		defer out.Close()

		// Write the header we already read
		if _, err := out.Write(header[:n]); err != nil {
			return fmt.Errorf("error writing file header: %w", err)
		}

		// Copy the rest of the body
		if _, err := io.Copy(out, body); err != nil {
			return fmt.Errorf("error saving file: %w", err)
		}

		return nil
	}

	// If not a ZIP, read the rest to parse error (limit to 1MB)
	const maxErrorSize = 1024 * 1024 // 1MB
	limitedReader := io.LimitReader(body, maxErrorSize)
	rest, err := io.ReadAll(limitedReader)
	if err != nil {
		return fmt.Errorf("error reading error response: %w", err)
	}

	fullResponse := append(header[:n], rest...)

	// Try to parse JSON error response
	var js map[string]interface{}
	if err := json.Unmarshal(fullResponse, &js); err == nil {
		if status, ok := js["query_status"].(string); ok {
			return fmt.Errorf("download failed: %s", status)
		}
	}

	// Return the raw response as error
	return fmt.Errorf("download failed: %s", string(fullResponse))
}
