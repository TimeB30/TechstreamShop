package kafkaclient

import "time"

type Message struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
}
