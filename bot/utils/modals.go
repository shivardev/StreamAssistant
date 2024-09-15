package utils

type LiveChatMessage struct {
	ChatID         string `json:"chatid"`
	AuthorName     string `json:"authorName"`
	AuthorPhotoURL string `json:"authorPhotoUrl"`
	Timestamp      string `json:"timestamp,omitempty"`
	MessageContent string `json:"messageContent"`
}

// Define the struct for the request payload
type RequestPayload struct {
	Messages []LiveChatMessage `json:"messages"`
}

// Dummy implementation for MessageProcessor
type MessageProcessor struct{}
