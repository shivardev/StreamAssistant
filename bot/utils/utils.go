package utils

import (
	"log"
	"os/exec"
	"syscall"
)

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
