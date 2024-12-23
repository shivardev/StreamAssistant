package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func SendMsgToYoutube(msg string) {
	apiURL := "http://10.0.0.128:8181" // Replace with your actual endpoint

	// Create the request body (for example, sending the message as JSON)
	data := map[string]interface{}{
		"msg": msg,
	}

	// Marshal the data into JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error marshalling data: %v", err)
		return
	}

	// Create a new POST request with the JSON data as the body
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
		return
	}

	// Set the Content-Type header for the POST request
	req.Header.Set("Content-Type", "application/json")

	// Send the POST request using the default HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
		return
	}
	defer resp.Body.Close()

	// Check the status code of the response
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Message sent successfully")
	} else {
		fmt.Printf("Failed to send message. Status: %d\n", resp.StatusCode)
	}
}
