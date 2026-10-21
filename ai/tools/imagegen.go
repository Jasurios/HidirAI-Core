package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type cfReq struct {
	Prompt string `json:"prompt"`
	Steps  int    `json:"steps,omitempty"`
}

type cfResp struct {
	Success bool `json:"success"`
	Result  struct {
		Image string `json:"image"`
	} `json:"result"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

var client = &http.Client{Timeout: 60 * time.Second}

func GenImage(prompt string) (string, error) {
	id, token := os.Getenv("CF_ACCOUNT_ID"), os.Getenv("CF_API_TOKEN")
	if id == "" || token == "" {
		return "", fmt.Errorf("нет CF_ACCOUNT_ID / CF_API_TOKEN")
	}

	body, err := json.Marshal(cfReq{Prompt: prompt, Steps: 8})
	if err != nil {
		return "", err
	}

	url := "https://api.cloudflare.com/client/v4/accounts/" + id +
		"/ai/run/@cf/black-forest-labs/flux-1-schnell"

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var out cfResp
	if err := json.Unmarshal(raw, &out); err != nil {
		return "", fmt.Errorf("разбор: %w, тело: %s", err, raw)
	}
	if resp.StatusCode != http.StatusOK || !out.Success {
		return "", fmt.Errorf("CF статус %d, ошибки: %+v", resp.StatusCode, out.Errors)
	}
	return out.Result.Image, nil
}
