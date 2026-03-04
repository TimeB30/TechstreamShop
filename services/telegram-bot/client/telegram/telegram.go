package telegram

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"path"
	"strconv"

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

func (c *Client) SendMessage(log *slog.Logger, chatID int64, message string, inlineKeyboardMarkup InlineKeyboardMarkup) error {
	op := "client.SendMessage"
	log = log.With(
		slog.String("op", op))
	reqBody := Message{
		ChatID:               chatID,
		Text:                 message,
		InlineKeyboardMarkUp: inlineKeyboardMarkup,
	}

	//q := url.Values{}
	//q.Add("chat_id", strconv.FormatInt(chatID, 10))
	//q.Add("text", message)
	jsondata, err := json.Marshal(reqBody)
	if err != nil {
		return e.Wrap(op, err)
	}
	_, err := c.doRequest(log, sendMessageMethod)
	if err != nil {
		log.Error("can't send message", "chat_id", chatID)
		return e.Wrap(op, err)
	}
	return nil
}

func (c *Client) Updates(log *slog.Logger, offset int64, limit int64) (updates []Update, err error) {
	defer func() {
		err = e.WrapIfErr("can't get updates", err)
		if err != nil {
			log.Error(err.Error())
		}
	}()
	op := "client.telegram.Updates"
	log = log.With(
		slog.String("op", op),
	)
	q := url.Values{}
	q.Add("offset", strconv.FormatInt(offset, 10))
	q.Add("limit", strconv.FormatInt(limit, 10))
	data, err := c.doRequest(log, getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}
	var res UpdatesResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (c *Client) doRequest(log *slog.Logger, method string) (data []byte, err error) {
	defer func() {
		err = e.WrapIfErr("can't do request", err)
		if err != nil {
			log.Error(err.Error())
		}
	}()
	op := "client.telegram.DoRequest"
	log = log.With(
		slog.String("op", op),
	)
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath + method),
	}
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = query.Encode()
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
	log.Info("Request completed", "status", resp.StatusCode)
	return body, nil
}
