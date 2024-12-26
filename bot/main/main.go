package main

import (
	"bot/utils"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os/exec"
	"reflect"
	"time"

	"github.com/emirpasic/gods/sets/hashset"
	"github.com/playwright-community/playwright-go"
)

var (
	lastChatID  string
	page_cursor playwright.Page
	relangiData = utils.GetRelangiJSON()
)
var Users = hashset.New()
var ignoredUsers = hashset.New("Nightbot", "YouTube", "Blazing Bane", "Relangi mama")
var messageHandlers = map[string]func(){}
var urlsToMonitor = [2]string{"https://www.youtube.com/youtubei/v1/live_chat/get_live_chat?prettyPrint=false", "https://www.youtube.com/youtubei/v1/updated_metadata?prettyPrint=false"}

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
	initMessageHandlers()
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
	page, err := browser.Contexts()[0].NewPage() // get the first page
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

// func fetchNewMessages(page playwright.Page) {
// 	elements, err := page.QuerySelectorAll("yt-live-chat-text-message-renderer")
// 	if err != nil {
// 		log.Printf("error querying chat messages: %v", err)
// 		return
// 	}

// 	newMessages := []utils.LiveChatMessage{}
// 	foundLastChatID := false

// 	for _, element := range elements {
// 		messageID, err := element.GetAttribute("id")
// 		if err != nil {
// 			log.Printf("error getting attribute 'id': %v", err)
// 			continue
// 		}

// 		if lastChatID == "" {
// 			chatContent, err := NewLiveChatMessage(element)
// 			if err != nil {
// 				log.Printf("error creating LiveChatMessage: %v", err)
// 				continue
// 			}
// 			newMessages = append(newMessages, *chatContent)
// 			lastChatID = messageID
// 		} else {
// 			if messageID == lastChatID {
// 				foundLastChatID = true
// 			} else if foundLastChatID {
// 				chatContent, err := NewLiveChatMessage(element)
// 				if err != nil {
// 					log.Printf("error creating LiveChatMessage: %v", err)
// 					continue
// 				}
// 				newMessages = append(newMessages, *chatContent)
// 			}
// 		}
// 	}

// 	if len(elements) > 0 {
// 		lastChatID, _ = elements[len(elements)-1].GetAttribute("id")
// 	} else {
// 		fmt.Println("No new messages found")
// 	}

// 	if len(newMessages) > 0 {
// 		msgProcessor := MessageProcessor{}
// 		msgProcessor.processEachMessage(newMessages)
// 	}
// }

func SendMsgToYoutube(content string) {
	page_cursor.Locator("div#input").PressSequentially(content)
	page_cursor.Keyboard().Press("Enter")
}

type MessageProcessor struct{}

func initMessageHandlers() {
	val := reflect.ValueOf(relangiData)
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fieldName := field.Tag.Get("json")
		messageHandlers[fieldName] = func(f reflect.Value) func() {
			return func() {
				if f.Len() > 0 {
					SendMsgToYoutube(f.Index(rand.Intn(f.Len())).Interface().(string))
				}
			}
		}(val.Field(i))
	}
}
func NewLiveChatMessage(element playwright.ElementHandle) (*utils.LiveChatMessage, error) {
	chatID, err := element.GetAttribute("id")
	if err != nil {
		return nil, err
	}

	authorPhotoURL, err := element.QuerySelector("yt-img-shadow img")
	if err != nil {
		return nil, err
	}
	photoURL, err := authorPhotoURL.GetAttribute("src")
	if err != nil {
		return nil, err
	}

	authorNameElement, err := element.QuerySelector("yt-live-chat-author-chip #author-name")
	if err != nil {
		return nil, err
	}
	authorName, err := authorNameElement.InnerText()
	if err != nil {
		return nil, err
	}

	messageContentElement, err := element.QuerySelector("#message")
	if err != nil {
		return nil, err
	}
	messageContent, err := messageContentElement.InnerText()
	if err != nil {
		return nil, err
	}

	return &utils.LiveChatMessage{
		ChatID:         chatID,
		AuthorName:     authorName,
		AuthorPhotoURL: photoURL,
		MessageContent: messageContent,
	}, nil
}
