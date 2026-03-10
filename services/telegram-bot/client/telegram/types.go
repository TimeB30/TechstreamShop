package telegram

type TempUpdate struct {
	OK          bool   `json:"ok"`
	Description string `json:"description,omitempty"`
}
type EditResponse struct {
	OK bool `json:"ok"`
}
type UpdatesResponse struct {
	OK          bool     `json:"ok"`
	Result      []Update `json:"result"`
	Description string   `json:"description,omitempty"`
}
type Update struct {
	ID            int64            `json:"update_id"`
	Message       *IncomingMessage `json:"message,omitempty"`
	CallBackQuery *CallBackQuery   `json:"callback_query,omitempty"`
}
type CallBackQuery struct {
	ID      string          `json:"id"`
	From    User            `json:"from"`
	Message IncomingMessage `json:"message"`
	Data    string          `json:"data"`
}
type IncomingMessage struct {
	ID   int64  `json:"message_id"`
	Text string `json:"text"`
	From User   `json:"from"`
	Chat Chat   `json:"chat"`
}
type User struct {
	UserName string `json:"username"`
	UserID   int64  `json:"id"`
}
type Chat struct {
	ChatID int64 `json:"id"`
}

type Message struct {
	ChatID      int64        `json:"chat_id"`
	Text        string       `json:"text"`
	ReplyMarkup *ReplyMarkup `json:"reply_markup,omitempty"`
}

type ReplyMarkup struct {
	InlineKeyboardMarkup interface{} `json:"inline_keyboard"`
}
type InlineKeyboardMarkup [][]InlineKeyboardButton

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	CallBackData string `json:"callback_data"`
}

type CallBackQueryAnswer struct {
	CallBackQueryID string `json:"callback_query_id"`
	Text            string `json:"text"`
}

type EditMessageReplyMarkup struct {
	ChatID      int64        `json:"chat_id"`
	MessageID   int64        `json:"message_id"`
	ReplyMarkup *ReplyMarkup `json:"reply_markup,omitempty"`
}
