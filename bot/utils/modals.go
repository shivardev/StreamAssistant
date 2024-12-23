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

// Action struct for streamer bot

// Define the struct for relangi ma a response
type RelangiData struct {
	Cmds  []string `json:"cmds"`
	Mama  []string `json:" mama "`
	Op    []string `json:" op "`
	NT    []string `json:" nt "`
	Hi    []string `json:" hi "`
	Sleep []string `json:" sleep"`
	Bye   []string `json:" bye"`
	Bf    []string `json:" bf"`
	Food  []string `json:" food "`
	Games []string `json:" games "`
	Blog  []string `json:" blog"`
}
type User struct {
	ID          int    `json:"id"`
	Username    string `json:"username"`
	Points      int    `json:"points"`
	JoinedDate  string `json:"joinedDate"`
	LastComment string `json:"lastComment"`
}

// Dummy implementation for MessageProcessor
type MessageProcessor struct{}

// Relangi server Modals

// hTTP post requst body struct

type PostReq struct {
	Msg string `json:"msg"`
}
