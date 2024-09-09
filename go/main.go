package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type ChatMessage struct {
	ChatID         string `json:"chatid"`
	AuthorName     string `json:"authorName"`
	AuthorPhotoURL string `json:"authorPhotoUrl"`
	MessageContent string `json:"messageContent"`
}

type MessagePayload struct {
	Messages []ChatMessage `json:"messages"`
}

var messageQueue chan ChatMessage
var speakQueue chan ChatMessage
var connections = make(map[*websocket.Conn]bool)
var mu sync.Mutex // To ensure thread-safety

func init() {
	// Initialize the message queue channel with a buffer size
	messageQueue = make(chan ChatMessage, 100)
	speakQueue = make(chan ChatMessage, 100)
}
func processQueue() {
	for {
		// Receive messages from the queue (blocking operation)
		msg := <-messageQueue

		// Process the message, only if it starts with "s"
		if len(msg.MessageContent) > 0 {
			speakQueue <- msg
			// Send to WebSocket or further processing (e.g., for HTML display)
		} else {
			fmt.Println("Ignoring message:", msg.MessageContent)
		}
	}
}
func processSpeakQueue() {
	for {
		// Receive messages from the queue (blocking operation)
		msg := <-speakQueue

		fmt.Println("Processing !speak message:", msg.MessageContent)
		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			log.Println("Error marshalling message to JSON:", err)
			continue
		}
		// Send the message to all connected WebSocket clients
		mu.Lock()
		for conn := range connections {
			if err := conn.WriteMessage(websocket.TextMessage, []byte(jsonMsg)); err != nil {
				log.Println("WebSocket write error:", err)
				conn.Close()
				delete(connections, conn)
			}
		}
		mu.Unlock()

	}
}
func main() {
	go processQueue()
	go processSpeakQueue()
	app := fiber.New()

	app.Static("/", "./static")
	app.Get("/speak", func(c *fiber.Ctx) error {
		return c.SendFile("./static/speak.html")
	})
	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
		mu.Lock()
		connections[c] = true // Track the connection
		mu.Unlock()

		defer func() {
			mu.Lock()
			delete(connections, c) // Clean up when the connection closes
			mu.Unlock()
			c.Close()
		}()

		var (
			mt  int
			msg []byte
			err error
		)
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", msg)

			// Echo the message back to the client
			if err = c.WriteMessage(mt, msg); err != nil {
				log.Println("write:", err)
				break
			}
		}
	}))

	app.Post("/takemsgs", func(c *fiber.Ctx) error {
		// Create a struct to hold the incoming data
		var chatMessages MessagePayload

		// Parse the incoming JSON request body into the struct
		if err := c.BodyParser(&chatMessages); err != nil {
			fmt.Println("Error parsing body:", err)
			return c.Status(fiber.StatusBadRequest).SendString("Failed to parse request body")
		}

		go func() {
			for _, msg := range chatMessages.Messages {
				messageQueue <- msg // Enqueue without filtering
			}
		}()
		// Return success response
		return c.Status(fiber.StatusOK).SendString("Messages received successfully")
	})

	log.Fatal(app.Listen(":3000"))
}
