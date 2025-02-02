package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

// reads the json data in responce.json and returns a RelangiData struct
func GetRelangiJSON() RelangiData {

	data, err := ioutil.ReadFile("responce.json")
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	var response RelangiData
	err = json.Unmarshal(data, &response)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	return response
}

func SendMsgToYoutube(msg string) {
	go func() {
		apiURL := "http://10.0.0.128:8181" // Replace with your actual endpoint

		data := map[string]interface{}{
			"msg": msg,
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Fatalf("Error marshalling data: %v", err)
			return
		}

		req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Fatalf("Error creating request: %v", err)
			return
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Error sending request: %v", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			fmt.Println("Message sent successfully")
		} else {
			fmt.Printf("Failed to send message. Status: %d\n", resp.StatusCode)
		}
	}()
}
