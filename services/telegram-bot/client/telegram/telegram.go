package telegram

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/timeb30/techstreamshop/services/telegram-bot/lib/e"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

func New(host string, token string) Client {
	return Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) SendMessage(chatID int64, message string, inlineKeyboardMarkup InlineKeyboardMarkup) error {
	op := "client.SendMessage"
	reqBody := Message{
		ChatID:               chatID,
		Text:                 message,
		InlineKeyboardMarkUp: inlineKeyboardMarkup,
	}
	_, err := c.doRequest(sendMessageMethod, reqBody)
	if err != nil {
		return e.Wrap(op, err)
	}
	return nil
}

func (c *Client) Updates(offset int64, limit int64) (updates []Update, err error) {
	const op = "telegram.Updates"
	defer func() {
		err = e.WrapIfErr(op+":can't get updates", err)
	}()
	message := struct {
		Offset int64 `json:"offset"`
		Limit  int64 `json:"limit"`
	}{
		Offset: offset,
		Limit:  limit,
	}
	data, err := c.doRequest(getUpdatesMethod, message)
	if err != nil {
		return nil, err
	}
	var res UpdatesResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (c *Client) doRequest(method string, payload any) (data []byte, err error) {
	const op = "client.doRequest"
	defer func() {
		err = e.WrapIfErr(op+":can't do request", err)
	}()
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath + method),
	}
	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
