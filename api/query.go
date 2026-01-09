package api

import (
	"context"
	"encoding/json"
	"fmt"
)

// queryAPI helper to make API requests and return results
func (c *Client) queryAPI(ctx context.Context, queryType, queryKey, queryValue string, limit int) ([]MalwareSample, error) {
	data := map[string]string{
		"query": queryType,
	}

	if queryKey != "" && queryValue != "" {
		data[queryKey] = queryValue
	}

	if limit > 0 {
		data["limit"] = fmt.Sprintf("%d", limit)
	}

	response, err := c.MakeRequest(ctx, data, nil)
	if err != nil {
		return nil, fmt.Errorf("error retrieving samples for %s %s: %w", queryType, queryValue, err)
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

// QueryByHash retrieves malware samples by hash
func (c *Client) QueryByHash(ctx context.Context, hash string, limit int) ([]MalwareSample, error) {
	// Validate hash format (SHA256 or MD5)
	if err := ValidateSHA256(hash); err != nil {
		if err := ValidateMD5(hash); err != nil {
			return nil, fmt.Errorf("invalid hash format: must be SHA256 (64 hex) or MD5 (32 hex)")
		}
	}
	return c.queryAPI(ctx, "get_info", "hash", hash, limit)
}

// QueryByTag retrieves malware samples with a specific tag
func (c *Client) QueryByTag(ctx context.Context, tag string, limit int) ([]MalwareSample, error) {
	if err := ValidateTag(tag); err != nil {
		return nil, fmt.Errorf("invalid tag: %w", err)
	}
	return c.queryAPI(ctx, "get_taginfo", "tag", tag, limit)
}

// QueryBySignature retrieves malware samples with a specific signature
func (c *Client) QueryBySignature(ctx context.Context, signature string, limit int) ([]MalwareSample, error) {
	return c.queryAPI(ctx, "get_siginfo", "signature", signature, limit)
}

// QueryByFileType retrieves malware samples of a specific file type
func (c *Client) QueryByFileType(ctx context.Context, filetype string, limit int) ([]MalwareSample, error) {
	return c.queryAPI(ctx, "get_file_type", "file_type", filetype, limit)
}

// QueryByClamAV retrieves samples by ClamAV signature
func (c *Client) QueryByClamAV(ctx context.Context, clamav string, limit int) ([]MalwareSample, error) {
	return c.queryAPI(ctx, "get_clamavinfo", "clamav", clamav, limit)
}

// QueryByImpHash retrieves samples by Imphash signature
func (c *Client) QueryByImpHash(ctx context.Context, imphash string, limit int) ([]MalwareSample, error) {
	return c.queryAPI(ctx, "get_imphash", "imphash", imphash, limit)
}

// QueryByTLSH retrieves samples by TLSH signature
func (c *Client) QueryByTLSH(ctx context.Context, tlshash string, limit int) ([]MalwareSample, error) {
	return c.queryAPI(ctx, "get_tlsh", "tlsh", tlshash, limit)
}

// QueryByTelfHash retrieves samples by Telfhash signature
func (c *Client) QueryByTelfHash(ctx context.Context, telfhash string, limit int) ([]MalwareSample, error) {
	return c.queryAPI(ctx, "get_telfhash", "telfhash", telfhash, limit)
}

// QueryByDHash retrieves samples by DHash signature
func (c *Client) QueryByDHash(ctx context.Context, dhash string, limit int) ([]MalwareSample, error) {
	return c.queryAPI(ctx, "get_dhash_icon", "dhash_icon", dhash, limit)
}

// QueryByGimphash retrieves samples by Gimphash signature
func (c *Client) QueryByGimphash(ctx context.Context, dhash string, limit int) ([]MalwareSample, error) {
	return c.queryAPI(ctx, "get_gimphash", "gimphash", dhash, limit)
}

// QueryByYara retrieves samples by Yara Rule
func (c *Client) QueryByYara(ctx context.Context, yara string, limit int) ([]MalwareSample, error) {
	return c.queryAPI(ctx, "get_yarainfo", "yara_rule", yara, limit)
}

// QueryByIssuerCN retrieves samples by Issuer Common Name
func (c *Client) QueryByIssuerCN(ctx context.Context, issuerCN string, limit int) ([]MalwareSample, error) {
	return c.queryAPI(ctx, "get_issuerinfo", "issuer_cn", issuerCN, limit)
}

// QueryBySubjectCN retrieves samples by Subject Common Name
func (c *Client) QueryBySubjectCN(ctx context.Context, subjectCN string, limit int) ([]MalwareSample, error) {
	return c.queryAPI(ctx, "get_subjectinfo", "subject_cn", subjectCN, limit)
}

// QueryBySerialNumber retrieves samples by Serial Number
func (c *Client) QueryBySerialNumber(ctx context.Context, serialNumber string, limit int) ([]MalwareSample, error) {
	return c.queryAPI(ctx, "get_certificate", "serial_number", serialNumber, limit)
}

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
