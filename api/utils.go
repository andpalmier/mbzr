package api

import (
	"encoding/json"
	"fmt"
)

// mustMarshalJSON marshals data or panic
func mustMarshalJSON(data interface{}) []byte {
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(fmt.Sprintf("Error: marshaling JSON failed: %v", err))
	}
	return jsonData
}

// PrintData returns the pretty-printed JSON contained in the "data" section as a string
func PrintData(response string) (string, error) {
	// Create a generic map to hold the unmarshalled JSON
	var result map[string]interface{}

	// Unmarshal the JSON response
	err := json.Unmarshal([]byte(response), &result)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	// Extracting data and marshalling it back to a pretty-printed JSON string
	if data, ok := result["data"]; ok {
		dataJSON, err := json.MarshalIndent(data, "", "    ")
		if err != nil {
			return "", fmt.Errorf("error marshalling data: %v", err)
		}
		return string(dataJSON), nil
	}

	return "", fmt.Errorf("no data found in the response")
}
