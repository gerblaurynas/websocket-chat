package transport

import "time"

type Input struct {
	Text string `json:"text"`
}

type Output struct {
	*Input
	Timestamp time.Time `json:"timestamp"`
	Username  string    `json:"username"`
}

type UserError struct {
	error
	userMessage string
}
