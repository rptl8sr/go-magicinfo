package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-magicinfo/internal/client"
	"io"
	"log/slog"
	"net/http"
)

var (
	SETUP = "/restapi/v1.0/rms/devices/%s/setup"
)

func ChangeUrl(ctx context.Context, c *client.Client, id, url, newUrl string) {
	slog.Info(fmt.Sprintf("Change device %s URL from %s to %s", id, url, newUrl))

	payload := map[string]string{
		"magicinfoServerUrl": newUrl,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		slog.Error("Unable to marshal payload", "err", err)
		return
	}

	fullUrl := url + fmt.Sprintf(SETUP, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, fullUrl, bytes.NewBuffer(jsonData))
	req.Header.Add("Content-Type", "application/json")

	res, err := c.Do(req)
	if err != nil {
		slog.Error("Unable to make request", "err", err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		slog.Error("Request failed", "status", res.StatusCode)
		return
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("Unable to read response body", "err", err)
		return
	}

	slog.Info("URL change successful", "msg", string(body))
}
