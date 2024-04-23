package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

type Payload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Params struct {
	Cmd      string `json:"cmd"`
	Id       string `json:"id"`
	Password string `json:"pw"`
}

type Response struct {
	Token      string `json:"token"`
	DateFormat string `json:"dateFormat"`
	TimeFormat string `json:"timeFormat"`
}

var (
	RestAPIAuth = "/auth"
	OpenAPIAuth = "/openapi/auth"
)

func RestAPIToken(ctx context.Context, c *http.Client, u, user, pass string) (t string, err error) {
	fullUrl := u + RestAPIAuth
	slog.Info(fmt.Sprintf("Getting RestAPIToken from %s", fullUrl))

	payload := &Payload{
		Username: user,
		Password: pass,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		slog.Error("Unable to marshal payload", "err", err)
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullUrl, bytes.NewBuffer(jsonData))
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

func OpenAPIToken(c *http.Client, u, user, pass string) (t string, err error) {
	fullUrl := u + OpenAPIAuth
	slog.Info(fmt.Sprintf("Getting OpenAPIToken from %s", fullUrl))

	params := &Params{
		Cmd:      "getAuthToken",
		Id:       user,
		Password: pass,
	}

	queryParams := url.Values{
		"cmd": {params.Cmd},
		"id":  {params.Id},
		"pw":  {params.Password},
	}

	req, err := http.NewRequest(http.MethodGet, fullUrl, nil)
	req.URL.RawQuery = queryParams.Encode()
	slog.Info("Full path", "path", req.URL.String())
	if err != nil {
		slog.Error("Unable to make request", "err", err)
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("Unable to read response body", "err", err)
		return
	}
	if res.StatusCode != http.StatusOK {
		slog.Error("Request failed", "err", err, "status", res.StatusCode)
		return
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("Unable to read response body", "err", err)
		return
	}

	t = string(body)
	return
}
