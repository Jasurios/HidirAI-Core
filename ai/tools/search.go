package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Result struct {
	Results []struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		URL     string `json:"url"`
	} `json:"results"`
}

func SearchInWeb(query string) (string, error) {
	API := os.Getenv("TAVILY_API")
	body := fmt.Sprintf(`{"api_key": "%s", "query": "%s", "max_results": 3}`, API, query)

	resp, err := http.Post(
		"https://api.tavily.com/search",
		"application/json",
		strings.NewReader(body),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)

	var result Result
	json.Unmarshal(data, &result)

	var out strings.Builder
	for _, r := range result.Results {
		out.WriteString(fmt.Sprintf("**%s**\n%s\n%s\n\n", r.Title, r.Content, r.URL))
	}

	return out.String(), nil
}
