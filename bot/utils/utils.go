package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"syscall"
)

func SendMsgsforAlerts(messages []LiveChatMessage) {
	url := "http://" + raspberryPi + ":3000/takemsgs"
	payload := RequestPayload{Messages: messages}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

}

func Contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func RunNodeScript() {
	workingDir := `D:\coding\streamAssistant\nodePlaywrite\src`

	cmd := exec.Command("node", "index.js")
	cmd.Dir = workingDir
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	err := cmd.Start()
	if err != nil {
		log.Printf("Error starting Node.js script: %v\n", err)
		return
	}

	log.Println("Playwright script is running in the background...")

}
