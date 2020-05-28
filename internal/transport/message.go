package transport

import "time"

type Message struct {
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
	Username  string    `json:"username"`
}

type UserError struct {
	error
	userMessage string
}
