package api

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

// API endpoint URL
const apiURL = "https://mb-api.abuse.ch/api/v1/"

// MakeRequest makes an HTTP request to the API
func MakeRequest(data map[string]string, files map[string]io.Reader, apiKey string) (string, error) {
	client := &http.Client{}
	var req *http.Request
	var err error

	if files != nil {
		// Handle file upload
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		for key, r := range files {
			var fw io.Writer
			if key == "file" {
				// file upload
				if fw, err = writer.CreateFormFile(key, "file"); err != nil {
					return "", fmt.Errorf("Error: creating form file failed: %v", err)
				}
				if _, err = io.Copy(fw, r); err != nil {
					return "", fmt.Errorf("Error: copying file data failed: %v", err)
				}
			} else {
				// other form fields
				if fw, err = writer.CreateFormField(key); err != nil {
					return "", fmt.Errorf("Error: creating form field failed: %v", err)
				}
				if _, err = io.Copy(fw, r); err != nil {
					return "", fmt.Errorf("Error: copying form field data failed: %v", err)
				}
			}
		}
		writer.Close()
		req, err = http.NewRequest("POST", apiURL, body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
	} else {
		// Handle form data
		formData := url.Values{}
		for key, value := range data {
			formData.Add(key, value)
		}
		req, err = http.NewRequest("POST", apiURL, strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	if err != nil {
		return "", fmt.Errorf("Error: creating HTTP request failed: %v", err)
	}

	req.Header.Set("Auth-Key", apiKey)
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error: making HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error: reading response body failed: %v", err)
	}
	return string(body), nil
}
