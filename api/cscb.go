package api

import "fmt"

// GetCSCB retrieves the MalwareBazaar Code Signing Certificate Blocklist
func GetCSCB(apiKey string) error {
	data := map[string]string{
		"query": "get_cscb",
	}

	response, err := MakeRequest(data, nil, apiKey)
	if err != nil {
		return fmt.Errorf("Error retrieving Code Signing Certificate Blocklist: %v", err)
	}

	fmt.Println("Code Signing Certificate Blocklist:")
	dataJSON, err := PrintData(response)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println(dataJSON)
	}
	return nil
}
