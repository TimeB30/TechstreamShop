package telegram

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

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

func (c *Client) Updates(offset int64, limit int) ([]Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.FormatInt(offset, 10))
	// TODO  do request get updates
}
func (c *Client) DoRequest(log *slog.Logger, method string, query url.Values) ([]byte, error) {
	errMsg := "Failed to do request"
	op := "client.telegram.DoRequest"
	log := log.With(
		slog.String("op", op),
	)
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath + method),
	}
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		log.Error(errMsg, err)
		return nil, fmt.Errorf("%s %s %w", errMsg, op, err)
	}
	req.URL.RawQuery = query.Encode()
	resp, err := c.client.Do(req)
	if err != nil {
		log.Error(errMsg, err)
		return nil, fmt.Errorf("%s %s %w", errMsg, op, err)
	}
	return resp
}

func (c *Client) SendMessage() {

}
