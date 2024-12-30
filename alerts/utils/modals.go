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
