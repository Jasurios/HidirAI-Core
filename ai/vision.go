package ai

import (
	"context"
	"regexp"
	"strings"

	"os"

	"github.com/sashabaranov/go-openai"
)

var visionStandartPrompt = `Опиши это изображение максимально подробно и структурированно, по пунктам:

1. Общая сцена — что в целом происходит на фото, где это (помещение, улица, природа и т.д.)
2. Объекты — перечисли все заметные предметы, людей, животных, их расположение относительно друг друга
3. Текст — если на изображении есть любой читаемый текст (вывески, надписи, экран, документ), приведи его дословно, без пропусков
4. Цвета и освещение — основная цветовая гамма, время суток, источник света
5. Детали и контекст — эмоции людей (если есть), состояние предметов, любые примечательные мелочи, которые могут быть важны

Пиши простым текстом, без markdown, без нумерованных списков в ответе — просто связный текст по каждому пункту.`

var visionStickerPrompt = `Это стикер из чата — почти наверняка мем, шутка или эмоциональная реакция. ` +
	`Не описывай композицию, цвета, фон и качество изображения как техническую съёмку. ` +
	`Кратко опиши суть шутки/мема: что смешного или абсурдного происходит, какая эмоция передаётся, ` +
	`и если узнаёшь конкретный известный мем — назови его. 2-4 предложения, без длинных разборов.`

var thinkTags = regexp.MustCompile(`(?s)<think>.*?</think>`)

func delThinkTags(text string) string {
	cleaned := thinkTags.ReplaceAllString(text, "")
	return strings.TrimSpace(cleaned)
}
func visionText(prompt, base64Img string) []openai.ChatCompletionMessage {
	return []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "Не используй markdown форматирование, таблицы и звёздочки. Пиши простым текстом.",
		},
		{
			Role: openai.ChatMessageRoleUser,
			MultiContent: []openai.ChatMessagePart{
				{
					Type: openai.ChatMessagePartTypeImageURL,
					ImageURL: &openai.ChatMessageImageURL{
						URL: "data:image/jpeg;base64," + base64Img,
					},
				},
				{
					Type: openai.ChatMessagePartTypeText,
					Text: prompt,
				},
			},
		},
	}
}

func ViewImage(userid, base64Img, typeImg string) (string, error) {
	config := openai.DefaultConfig(os.Getenv("GROQ_API"))
	config.BaseURL = os.Getenv("GROQ_URL")
	client := openai.NewClientWithConfig(config)

	prompt := visionStandartPrompt
	if typeImg == "стикер" {
		prompt = visionStickerPrompt
	}

	resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model:     os.Getenv("MODEL_VISION"),
		Messages:  visionText(prompt, base64Img),
		MaxTokens: 995,
	})
	if err != nil {
		return "", err
	}

	switch typeImg {
	case "sticker":
		return "[Пользователь прислал стикер]" + delThinkTags(resp.Choices[0].Message.Content), nil
	case "photo":
		return "[Пользователь прислал фото]" + delThinkTags(resp.Choices[0].Message.Content), nil
	}

	return "", nil
}
