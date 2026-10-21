package ai

import (
	"HidirAi-Core/ai/tools"
	"HidirAi-Core/structs"
	"HidirAi-Core/wwf"
	"encoding/json"
	"log"

	"github.com/sashabaranov/go-openai"
)

type ResultType string

const (
	ResultText  ResultType = "text"
	ResultPhoto ResultType = "photo"
)

type SkillResult struct {
	Type  ResultType
	Text  string
	Photo string
	Audio []byte
	File structs.File
}

var MyTools = []openai.Tool{
	{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        "create_qrcode",
			Description: "Создать qr код",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"text": {"type": "string", "description": "Ссылка или текст для QR кода. Если пользователь не дал текст то попроси его что добавить в QR код"}
				},
				"required": ["text"]
			}`),
		},
	},
	{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        "new_mcp",
			Description: "Позволяет пользователю подключить внешний MCP сервер для расширения функций. НО покачто не реализован OAuth2.0 поэтому апи должен быть в самой ссылке",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"user": {"type": "string", "description": "userid пользователя"},
					"url": {"type": "string", "description": "Ссылка на MCP сервер"},
					"key": {"type": "string", "description": "API к MCP серверу если есть. Если нету по оставь строку пустой ''"}
				},
				"required": ["user","url","key"]
			}`),
		},
	},
	{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        "image_gen",
			Description: "Генерация картинки по описанию. Prompt ВСЕГДА переводи на английский язык перед передачей, даже если пользователь написал на русском или таджикском.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"prompt": {"type": "string", "description": "То что пользователь хочет видить в изрбражении"}
				},
				"required": ["prompt"]
			}`),
		},
	},
}

var ToolHandlers = map[string]func(args string) SkillResult{
	"create_qrcode": func(args string) SkillResult {
		var a struct {
			Text string `json:"text"`
		}
		json.Unmarshal([]byte(args), &a)
		return SkillResult{Type: ResultPhoto, Photo: tools.MakeQrcode(a.Text)}
	},
	"new_mcp": func(args string) SkillResult {
		var a structs.MCPServerService
		json.Unmarshal([]byte(args), &a)

		if a.User == ""{
			return SkillResult{Type: ResultText, Text: "Отсутствует userid"}
		}
		
		err := wwf.Newmcp(a.User, a.Url, a.Key)
		log.Println(a.User, a.Url, a.Key)
		if err != nil{
			return SkillResult{Type: ResultText, Text: err.Error()}
		}

		return SkillResult{Type: ResultText, Text: "Success"}
	},
	"image_gen": func(args string) SkillResult {
		var a struct {
			Prompt string `json:"prompt"`
		}
		json.Unmarshal([]byte(args), &a)
		photo, err := tools.GenImage(a.Prompt)
		if photo == "" {
			return SkillResult{Type: ResultText, Text: "Не удалось сгенерировать картинку"+err.Error()}
		}
		return SkillResult{Type: ResultPhoto, Photo: photo}
	},
}