package utils

type ChatMessage struct {
	AuthorName     string `json:"authorName"`
	AuthorId       string `json:"authorId"`
	AuthorPhotoURL string `json:"authorPhotoUrl"`
	CommentTime    string `json:"timestamp"`
	MessageContent string `json:"messageContent"`
}

type StatsPayload struct {
	Stats Stat `json:"stats"`
}
type Stat struct {
	Likes              int  `json:"likes"`
	PreviousLikes      int  `json:"previousLikes"`
	Viewers            int  `json:"viewers"`
	MaxLikes           int  `json:"maxLikes"`
	ShouldCongratulate bool `json:"shouldCongratulate"`
}
