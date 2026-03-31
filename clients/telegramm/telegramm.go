package telegramm

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/Aidajy111/Read-adviser-bot/lib/e"
)

type Client struct {
	host    string
	baseUrl string
	client  http.Client
}

const (
	getUpdatesMethod  = "getUpdates"
	errMsg            = "failed to get updates:"
	sendMessageMethod = "sendMessage"
)

func NewClient(host, token string) Client {
	return Client{
		host:    host,
		baseUrl: newBasePath(token),
		client:  http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) Updates(offset, limit int) ([]Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	// go request

	body, err := c.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, e.Wrap(errMsg, err)
	}

	var resp UpdateResponse

	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, e.Wrap(errMsg, err)
	}

	return resp.Result, nil
}

func (c *Client) SendMessage(chatID int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	_, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return e.Wrap(errMsg, err)
	}

	return nil
}

func (c *Client) doRequest(method string, query url.Values) ([]byte, error) {
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.baseUrl, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, e.Wrap(errMsg, err)
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, e.Wrap(errMsg, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, e.Wrap(errMsg, err)
	}

	return body, nil
}
