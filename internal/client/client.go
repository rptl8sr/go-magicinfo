package client

import "net/http"

type Client struct {
	HTTPClient http.Client
	Token      string
}

func New(c *http.Client, t string) *Client {
	return &Client{
		HTTPClient: *c,
		Token:      t,
	}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("api_key", c.Token)
	return c.HTTPClient.Do(req)
}
