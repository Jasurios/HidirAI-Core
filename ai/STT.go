package ai

import (
	"bytes"
	"context"
	"log"
	"os"

	"github.com/sashabaranov/go-openai"
)

func STT(audioBytes []byte)(string, error){
	config := openai.DefaultConfig(os.Getenv("GROQ_API"))
	config.BaseURL = os.Getenv("GROQ_URL")
	client := openai.NewClientWithConfig(config)

	resp, err := client.CreateTranscription(context.Background(), openai.AudioRequest{
		Model:    os.Getenv("MODEL_SPEECH"),
		Reader:   bytes.NewReader(audioBytes),
		FilePath: "voice.ogg",
	})
	if err != nil {
		return "", err
	}
	log.Println(resp.Text)

	return resp.Text, nil
}