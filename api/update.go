package api

import "fmt"

// UpdateSample updates existing entry in MalwareBazaar
func UpdateSample(sha256, key, value, apiKey string) error {
	data := map[string]string{
		"query":       "update",
		"sha256_hash": sha256,
		"key":         key,
		"value":       value,
	}

	// use MakeRequest to send the update request
	response, err := MakeRequest(data, nil, apiKey)
	if err != nil {
		return fmt.Errorf("Error updating sample %s: %v", sha256, err)
	}

	fmt.Printf("Updated sample %s\n%s\n", sha256, response)
	return nil
}
