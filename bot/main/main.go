package main

import (
	"bot/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
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

func main() {
	utils.GetActionList()
	initMessageHandlers()
	utils.DataBaseConnection()
	utils.InsertUser("alice", 100, "2023-01-15", "This is Alice's comment.")
	cmd := exec.Command("cmd.exe", "/c", "chromium.bat")
	err := cmd.Start()
	if err != nil {
		fmt.Println("Error starting batch script:", err)
		return
	}
	fmt.Println("Chromium started")
	time.Sleep(5 * time.Second)
	// Get Relangi JSON
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}
	browser, err := pw.Chromium.ConnectOverCDP("http://127.0.0.1:8989")
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}
	defaultContext := browser.Contexts()
	page := defaultContext[0].Pages()[0]
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
	for {
		fetchNewMessages(page)
		time.Sleep(2 * time.Second)
	}
}
func fetchNewMessages(page playwright.Page) {
	elements, err := page.QuerySelectorAll("yt-live-chat-text-message-renderer")
	if err != nil {
		log.Printf("error querying chat messages: %v", err)
		return
	}

	newMessages := []utils.LiveChatMessage{}
	foundLastChatID := false

	for _, element := range elements {
		messageID, err := element.GetAttribute("id")
		if err != nil {
			log.Printf("error getting attribute 'id': %v", err)
			continue
		}

		if lastChatID == "" {
			chatContent, err := NewLiveChatMessage(element)
			if err != nil {
				log.Printf("error creating LiveChatMessage: %v", err)
				continue
			}
			newMessages = append(newMessages, *chatContent)
			lastChatID = messageID
		} else {
			if messageID == lastChatID {
				foundLastChatID = true
			} else if foundLastChatID {
				chatContent, err := NewLiveChatMessage(element)
				if err != nil {
					log.Printf("error creating LiveChatMessage: %v", err)
					continue
				}
				newMessages = append(newMessages, *chatContent)
			}
		}
	}

	if len(elements) > 0 {
		lastChatID, _ = elements[len(elements)-1].GetAttribute("id")
	} else {
		fmt.Println("No new messages found")
	}

	if len(newMessages) > 0 {
		msgProcessor := MessageProcessor{}
		msgProcessor.processEachMessage(newMessages)
	}
}

func SendMsgToYoutube(content string) {
	page_cursor.Locator("div#input").PressSequentially(content)
	page_cursor.Keyboard().Press("Enter")
}

type MessageProcessor struct{}

func (mp *MessageProcessor) processEachMessage(messages []utils.LiveChatMessage) {

	// send messages to alerts go app for speak commands and other stuff
	if len(messages) > 0 {
		fmt.Println("Sending messages to alerts go app")
		SendMsgsforAlerts(messages)
	}
	// process bot messages
	for _, message := range messages {
		// Add user to the Database
		if !ignoredUsers.Contains(message.AuthorName) {
			ProcessUser(message)
		}
		userPoints, err := utils.GetUserPoints(message.AuthorName)
		if err != nil {
			if userPoints%10 == 0 {
				message := message.AuthorName + "Congralations for getting " + strconv.Itoa(userPoints) + " points"
				SendMsgToYoutube(message)
			}
		}
		println(message.AuthorName + " ->-> " + message.MessageContent)
		if !ignoredUsers.Contains(message.AuthorName) && strings.HasPrefix(message.MessageContent, "!point") {
			userPoints, err := utils.GetUserPoints(message.AuthorName)
			if err != nil {
				message := "You have " + strconv.Itoa(userPoints) + " points."
				SendMsgToYoutube(message)
				return
			}
		}
		if !Users.Contains(message.AuthorName) && !ignoredUsers.Contains(message.AuthorName) {
			Users.Add(message.AuthorName)
			SendMsgToYoutube(relangiData.Hi[rand.Intn(len(relangiData.Hi))] + " " + message.AuthorName)
			return
		}
		if !ignoredUsers.Contains(message.AuthorName) {
			lowerMessage := strings.ToLower(message.MessageContent)
			if !ignoredUsers.Contains(message.AuthorName) {
				for keyword, handler := range messageHandlers {
					if strings.Contains(lowerMessage, keyword) {
						handler()
						break // Assuming only one handler is needed per message
					}
				}
			}
			return
		} else {
			fmt.Println("Text from Ignored Author", message.AuthorName)
		}
	}
}
func ProcessUser(message utils.LiveChatMessage) {
	err := utils.InsertOrUpdateUser(message.AuthorName, message.MessageContent)
	if err != nil {
		log.Printf("Error handling user %s: %v", message.AuthorName, err)
	}
}
func SendMsgsforAlerts(messages []utils.LiveChatMessage) {
	url := "http://10.0.0.236:3000/takemsgs" // Replace with your URL

	// Create the request payload
	payload := utils.RequestPayload{Messages: messages}

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
}

// var messageHandlers = map[string]func(){
// 	"blog": func() {
// 		SendMsgToYoutube(relangiData.Cmds[rand.Intn(len(relangiData.Blog))])
// 	},
// 	"cmd": func() {
// 		SendMsgToYoutube(relangiData.Cmds[rand.Intn(len(relangiData.Cmds))])
// 	},
// }

func initMessageHandlers() {
	val := reflect.ValueOf(relangiData)
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fieldName := field.Tag.Get("json") // Get the JSON tag to use as the key
		messageHandlers[fieldName] = func(f reflect.Value) func() {
			return func() {
				if f.Len() > 0 { // Check if there are any messages available
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
