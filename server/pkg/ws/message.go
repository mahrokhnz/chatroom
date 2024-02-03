package ws

import "time"

type Message struct {
	Username  string    `json:"username"`
	Type      string    `json:"type"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
}
