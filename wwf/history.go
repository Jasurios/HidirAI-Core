package wwf

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/sashabaranov/go-openai"
)

func Write(userid string, msgs ...openai.ChatCompletionMessage)error  {
	os.MkdirAll("users/"+userid, 0755)

	file, err := os.OpenFile("./users/" + userid + "/" + userid + ".jsonl", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil{
		return err
	}
	defer file.Close()

	for _, m := range(msgs){
		data, _ := json.Marshal(m)
		data = append(data, '\n')

		file.Write(data)
	}

	return nil
}

func Read(userid string) []openai.ChatCompletionMessage {
	file, err := os.Open("./users/" + userid + "/" + userid + ".jsonl")
	if err != nil{
		return nil
	}
	defer file.Close()

	var msgs []openai.ChatCompletionMessage

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 10*1024*1024)

	for scanner.Scan(){
		line := scanner.Bytes()
		if len(line)==0{
			continue
		}
		var m openai.ChatCompletionMessage
		err := json.Unmarshal(line, &m)
		if err != nil{
			continue
		}

		msgs = append(msgs, m)
	}

	return msgs
}