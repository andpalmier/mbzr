package api

import "fmt"

// GetRecentDetections retrieves recent detections from the API
func GetRecentDetections(hours int, apiKey string) error {
	data := map[string]string{
		"query": "recent_detections",
		"hours": fmt.Sprintf("%d", hours),
	}

	response, err := MakeRequest(data, nil, apiKey)
	if err != nil {
		return fmt.Errorf("Error retrieving recent detections: %v", err)
	}

	fmt.Println("Recent detections:")
	dataJSON, err := PrintData(response)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println(dataJSON)
	}
	return nil
}
