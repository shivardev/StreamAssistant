package utils

const (
	// API Endpoint for sending messages
	streamingPC = "10.0.0.213"
	raspberryPi = "10.0.0.236"
	GamingPC    = "10.0.0.128"
	APIEndpoint = "http://10.0.0.128:3000/takemsgs"
)

func getAPIEndpoint() string {
	return APIEndpoint
}
