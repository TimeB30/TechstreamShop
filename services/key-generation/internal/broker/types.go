package broker

type KeyMessage struct {
	OrderMessage *OrderMessage `json:"order_message"`
	Key          string        `json:"key"`
}

type OrderMessage struct {
	UserID     int64  `json:"user_id"`
	ChatID     int64  `json:"chat_id"`
	SoftwareID string `json:"software_id"`
	Version    int64  `json:"version"`  // index
	Duration   int64  `json:"duration"` //days

}
