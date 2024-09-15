package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func SendMsgsforAlerts(messages []LiveChatMessage) {
	url := "http://" + raspberryPi + ":3000/takemsgs" // Replace with your URL

	// Create the request payload
	payload := RequestPayload{Messages: messages}

	// Convert the payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Print the response
	fmt.Println("Response Code:", resp.StatusCode)
}

func Contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
