package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// zipMagic is the local file header of a ZIP archive (PK\x03\x04).
// Downloaded samples are always zipped and protected with the password
// "infected".
var zipMagic = []byte{'P', 'K', 3, 4}

// ZipPassword is the password abuse.ch uses for every downloaded sample.
const ZipPassword = "infected"

// DownloadSample downloads a sample by its hash and writes it to outPath.
// When outPath is empty the file is written to "<sha256>.zip" in the working
// directory. It returns the path written.
func (c *Client) DownloadSample(ctx context.Context, sha256, outPath string) (string, error) {
	// Validate SHA256 format to prevent path traversal
	if err := ValidateSHA256(sha256); err != nil {
		return "", fmt.Errorf("invalid hash: %w", err)
	}

	if outPath == "" {
		outPath = fmt.Sprintf("%s.zip", sha256)
	}

	data := map[string]string{
		"query":       "get_file",
		"sha256_hash": sha256,
	}

	body, err := c.MakeRequestRaw(ctx, data, nil)
	if err != nil {
		return "", fmt.Errorf("error downloading sample: %w", err)
	}
	defer func() { _ = body.Close() }()

	// Read the ZIP magic. io.ReadFull is required here: a plain Read may
	// return fewer bytes than asked for even when more are available, which
	// would make a valid archive look like an error response.
	header := make([]byte, len(zipMagic))
	n, err := io.ReadFull(body, header)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return "", fmt.Errorf("error reading response header: %w", err)
	}

	// Anything that is not a ZIP is an error response from the API.
	if n < len(zipMagic) || string(header[:n]) != string(zipMagic) {
		return "", downloadError(header[:n], body)
	}

	// Check if file already exists
	if _, err := os.Stat(outPath); err == nil {
		return "", fmt.Errorf("file already exists: %s", outPath)
	}

	out, err := os.Create(outPath)
	if err != nil {
		return "", fmt.Errorf("error creating file: %w", err)
	}
	defer func() { _ = out.Close() }()

	// Write the header we already read
	if _, err := out.Write(header[:n]); err != nil {
		return "", fmt.Errorf("error writing file header: %w", err)
	}

	// Copy the rest of the body
	if _, err := io.Copy(out, body); err != nil {
		return "", fmt.Errorf("error saving file: %w", err)
	}

	return outPath, nil
}

// downloadError turns a non-ZIP response body into a useful error
func downloadError(read []byte, rest io.Reader) error {
	const maxErrorSize = 1024 * 1024 // 1MB
	tail, err := io.ReadAll(io.LimitReader(rest, maxErrorSize))
	if err != nil {
		return fmt.Errorf("error reading error response: %w", err)
	}
	full := append(read, tail...)

	var js struct {
		QueryStatus string `json:"query_status"`
	}
	if err := json.Unmarshal(full, &js); err == nil && js.QueryStatus != "" {
		if serr := newStatusError(js.QueryStatus, "get_file"); serr != nil {
			return serr
		}
		return fmt.Errorf("download failed: %s", js.QueryStatus)
	}

	return fmt.Errorf("download failed: %s", string(full))
}
