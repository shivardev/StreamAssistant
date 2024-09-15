package main

import (
	"bot/utils"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/emirpasic/gods/sets/hashset"
	"github.com/playwright-community/playwright-go"
)

var (
	lastChatID  string
	page_cursor playwright.Page
)
var ignoredAuthors = map[string]struct{}{
	"Nightbot":    {},
	"YouTube":     {},
	"BlazingBane": {},
}
var Users = hashset.New()

func main() {
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
	page.Goto("https://www.youtube.com/watch?v=kiQHT_9y1FE")
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

				if _, found := ignoredAuthors[chatContent.AuthorName]; !found {
					newMessages = append(newMessages, *chatContent)
				}
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
		utils.SendMsgsforAlerts(messages)
		fmt.Print("Sent messages to bot go app")
	}
	// process bot messages
	for _, message := range messages {
		print(message.AuthorName)
		if !Users.Contains(message.AuthorName) {
			Users.Add(message.AuthorName)
			SendMsgToYoutube("Welcome home " + message.AuthorName)
		}
		if _, found := ignoredAuthors[message.AuthorName]; !found {
			// Check if the message content matches any predefined types
			lowerMessage := strings.ToLower(message.MessageContent)
			for keyword, handler := range messageHandlers {
				if strings.Contains(lowerMessage, keyword) {
					handler()
					break // Assuming only one handler is needed per message
				}
			}
		} else {
			fmt.Println("Text from Ignored Author")
		}
	}

}

var messageHandlers = map[string]func(){
	"obs": func() {
		SendMsgToYoutube(utils.OBS)
	},
	"blog": func() {
		SendMsgToYoutube(utils.Blog)
	},
	"cmd": func() {
		SendMsgToYoutube(utils.Cmds)
	},
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
