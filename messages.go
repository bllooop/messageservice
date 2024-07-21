package messageservice

import "time"

type MessageItem struct {
	Id   int       `json:"-" db:"id"`
	Text string    `json:"text" binding:"required"`
	Time time.Time `json:"time" binding:"required"`
}
