package auth

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

type Payload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Token      string `json:"token"`
	DateFormat string `json:"dateFormat"`
	TimeFormat string `json:"timeFormat"`
}

func Token(c *http.Client, url, user, pass string) (t string, err error) {
	payload := &Payload{
		Username: user,
		Password: pass,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		slog.Error("Unable to marshal payload", "err", err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, url+"/auth", bytes.NewBuffer(jsonData))
	req.Header.Add("Content-Type", "application/json")

	res, err := c.Do(req)
	if err != nil {
		slog.Error("Unable to make request", "err", err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("Unable to read response body", "err", err)
		return
	}
	if res.StatusCode != http.StatusOK {
		slog.Error("Request failed", "err", err, "status", res.StatusCode)
		return
	}

	data := &Response{}
	resJson := json.Unmarshal(body, data)
	if resJson != nil {
		slog.Error("Unable to unmarshal response", "err", resJson)
		return
	}

	t = data.Token
	return
}
