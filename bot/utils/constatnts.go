package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

const (
	// API Endpoint for sending messages
	streamingPC   = "10.0.0.213"
	raspberryPi   = "10.0.0.236"
	GamingPC      = "10.0.0.128"
	StreamingLink = "https://www.youtube.com/@blazingbane5565/live"
	// StreamingLink = "https://www.youtube.com/watch?v=pfiCNAc2AgU&ab_channel=LofiGirl"
)

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
