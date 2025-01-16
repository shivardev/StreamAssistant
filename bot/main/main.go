package main

import (
	"bot/utils"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/emirpasic/gods/sets/hashset"
	"github.com/playwright-community/playwright-go"
)

var (
	serverIP    string
	page_cursor playwright.Page
	environment *string
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
func getStreamURL() string {
	url := fmt.Sprintf("http://%s:3000/streamurl", serverIP)

	// Make the GET request
	response, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to make GET request: %v", err)
	}
	defer response.Body.Close() // Close the response body when done

	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	// Print the response body
	fmt.Println("Response Body:", string(body))
	return string(body)
}
func main() {
	env := flag.String("env", "prod", "Set the environment (dev, prod, test)")
	flag.Parse()
	switch *env {
	case "dev":
		environment = env
		serverIP = utils.DevIP // Development IP (localhost)
	default:
		serverIP = utils.RaspberryPiIP // Default to raspberry pi
	}

	fmt.Printf("Running in %s environment\n", *env)
	fmt.Printf("Server IP Address: %s\n", serverIP)

	cmd := exec.Command("cmd.exe", "/c", "chrome.bat")
	err := cmd.Start()
	if err != nil {
		fmt.Println("Error starting batch script:", err)
		return
	}
	fmt.Println("Chromium started")
	time.Sleep(2 * time.Second)
	// Start Playwright instance
	utils.RunNodeScript(env)
	time.Sleep(10 * time.Second)
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
	defaultContext := browser.Contexts()
	page := defaultContext[0].Pages()[0]
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}
	streamURL := getStreamURL()
	fmt.Print("Stream URL: ", streamURL)
	page.Goto("https://www.youtube.com/live_chat?v=" + streamURL)
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

	// if *environment == "dev" {
	// 	return
	// }
	page_cursor.Locator("div#input").PressSequentially(content)
	page_cursor.Keyboard().Press("Enter")
	fmt.Println("Sent message to youtube", content)
}

type MessageProcessor struct{}
