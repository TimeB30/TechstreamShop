package telegram

type UpdatesResponse struct {
	OK     bool     `json:"ok"`
	Result []Update `json:"result"`
}
type Update struct {
	ID      int64            `json:"update_id"`
	Message *IncomingMessage `json:"message"`
}

type IncomingMessage struct {
	Text string `json:"text"`
	From User   `json:"from"`
	Chat Chat   `json:"chat"`
}
type User struct {
	UserName string `json:"username"`
	UserID   int64  `json:"user_id"`
}
type Chat struct {
	ChatID int64 `json:"id"`
}

type Message struct {
	ChatID               int64                `json:"chat_id"`
	Text                 string               `json:"text"`
	InlineKeyboardMarkUp InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

type InlineKeyboardMarkup [][]InlineKeyboardButton

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	CallBackData string `json:"callback_data"`
}
