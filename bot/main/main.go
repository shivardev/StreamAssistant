package main

import (
	"bot/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/emirpasic/gods/sets/hashset"
	"github.com/playwright-community/playwright-go"
)

var (
	page_cursor playwright.Page
)
var Users = hashset.New()

func handleRequests(w http.ResponseWriter, r *http.Request) {
	// Handle your requests here
	switch r.Method {
	case "GET":
		// Respond to a GET request (you can extend this to serve dynamic data)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Server is running. Visit /start to begin Chromium automation.\n")
	case "POST":
		var msg utils.PostReq
		err := json.NewDecoder(r.Body).Decode(&msg)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Check if the message is not empty and process it
		if msg.Msg != "" {
			// Call the SendMsgToYoutube function to process the message
			SendMsgToYoutube(msg.Msg)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Message processed: %s", msg.Msg)
		} else {
			// Return an error if no message was provided
			http.Error(w, "No message provided", http.StatusBadRequest)
		}
	default:
		http.Error(w, "Unsupported HTTP method", http.StatusMethodNotAllowed)
	}
}
func main() {
	// Set up the HTTP server
	cmd := exec.Command("cmd.exe", "/c", "chromium.bat")
	err := cmd.Start()
	if err != nil {
		fmt.Println("Error starting batch script:", err)
		return
	}
	fmt.Println("Chromium started")
	time.Sleep(2 * time.Second)
	// Start Playwright instance
	utils.RunNodeScript()
	// connect to playwright which has the youtube open by now
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}
	browser, err := pw.Chromium.ConnectOverCDP("http://127.0.0.1:8989")
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}
	// defaultContext := browser.Contexts()
	// page := defaultContext[0].Pages()[0] // get the first page
	page, err := browser.Contexts()[0].NewPage()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}
	page.Goto(utils.StreamingLink)
	time.Sleep(2 * time.Second)
	iframeLocator := page.Locator("iframe#chatframe")
	if iframeExists, err := iframeLocator.Count(); err != nil {
		log.Fatalf("could not count iframes: %v", err)
	} else if iframeExists == 0 {
		fmt.Println("Iframe not found")
		return
	} else {
		fmt.Println("Iframe found")
	}

	src, err := iframeLocator.GetAttribute("src")
	if err != nil {
		log.Fatalf("could not get attribute: %v", err)
	}
	fmt.Printf("Iframe src attribute: %s\n", src)
	page.Goto("https://www.youtube.com" + src)
	page_cursor = page
	// HTTP
	http.HandleFunc("/", handleRequests)

	fmt.Println("Starting HTTP server on port 8181...")
	http.ListenAndServe(":8181", nil)
	if err != nil {
		log.Fatalf("Error starting HTTP server: %v", err)
	}
}

func SendMsgToYoutube(content string) {
	page_cursor.Locator("div#input").PressSequentially(content)
	page_cursor.Keyboard().Press("Enter")
}

type MessageProcessor struct{}
