package messageservice

type MessageItem struct {
	Id   int    `json:"-" db:"id"`
	Text string `json:"text" binding:"required"`
	Time string `json:"time" binding:"required"`
}
