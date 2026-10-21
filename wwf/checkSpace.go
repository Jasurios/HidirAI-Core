package wwf

import (
	"bufio"
	"encoding/json"
	"os"
	"strconv"

	"github.com/sashabaranov/go-openai"
)

func readJSONL(path string) ([]openai.ChatCompletionMessage, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var messages []openai.ChatCompletionMessage

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 10*1024*1024)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var m openai.ChatCompletionMessage
		if err := json.Unmarshal(line, &m); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}

	return messages, scanner.Err()
}

func writeJSONL(path string, messages []openai.ChatCompletionMessage) error {
	tmp := path + ".tmp"

	f, err := os.Create(tmp)
	if err != nil {
		return err
	}

	w := bufio.NewWriter(f)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)

	for _, m := range messages {
		if err := enc.Encode(m); err != nil {
			f.Close()
			os.Remove(tmp)
			return err
		}
	}

	if err := w.Flush(); err != nil {
		f.Close()
		os.Remove(tmp)
		return err
	}
	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return err
	}

	return os.Rename(tmp, path)
}

func CheckSpace(userid string) error {
	limit, _ := strconv.Atoi(os.Getenv("MESSAGE"))
	file := "./users/" + userid + "/" + userid + ".jsonl"

	messages, err := readJSONL(file)
	if err != nil {
		return err
	}

	if len(messages) <= limit {
		return nil
	}

	cut := 2
	for cut < len(messages) && messages[cut].Role != "user" {
		cut++
	}

	messages = append(messages[:1], messages[cut:]...)

	return writeJSONL(file, messages)
}