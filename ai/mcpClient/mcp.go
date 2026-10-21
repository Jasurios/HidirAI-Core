package mcpclient

import (
	"HidirAi-Core/structs"
	"HidirAi-Core/wwf"
	"bufio"
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/sashabaranov/go-openai"
	"golang.org/x/oauth2"
)

var tools []openai.Tool
var servers = map[string]*mcp.ClientSession{}

func GetTools(session *mcp.ClientSession) error {
	list, err := session.ListTools(context.Background(), nil)
	if err != nil {
		log.Println(err)
	}

	for _, tool := range list.Tools {
		servers[tool.Name] = session
		tools = append(tools, openai.Tool{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        tool.Name,
				Description: tool.Description,
				Parameters:  tool.InputSchema,
			},
		})
	}

	return nil
}

func Load(userid string) ([]openai.Tool, map[string]*mcp.ClientSession) {
	tools = []openai.Tool{}
	servers = map[string]*mcp.ClientSession{}

	wwf.DefaultMCP(userid)

	file, err := os.Open("./users/" + userid + "/MCP.jsonl")
	if err != nil {
		log.Println(err)
		return tools, servers
	}

	client := mcp.NewClient(&mcp.Implementation{
		Name:    "HidirAi-Core",
		Version: "4.0",
	}, nil)

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 10*1024*1024)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var mcpserver structs.MCPServer

		err := json.Unmarshal(line, &mcpserver)
		if err != nil {
			log.Println(mcpserver)
		}

		if mcpserver.Key == "" {
			session, err := client.Connect(context.Background(), &mcp.StreamableClientTransport{Endpoint: mcpserver.Url}, nil)
			if err != nil {
				log.Println(err)
				return tools, servers
			}

			err = GetTools(session)
			if err != nil {
				log.Println(err)
			}
		} else if mcpserver.Key != "" {
			session, err := client.Connect(context.Background(), &mcp.StreamableClientTransport{
				Endpoint: mcpserver.Url,
				HTTPClient: oauth2.NewClient(
					context.Background(),
					oauth2.StaticTokenSource(&oauth2.Token{
						AccessToken: mcpserver.Key,
						TokenType:   "Bearer",
					}),
				),
			}, nil)
			if err != nil {
				log.Println(err)
				return tools, servers
			}

			err = GetTools(session)
			if err != nil {
				log.Println(err)
			}
		}
	}

	return tools, servers
}

func Call(server *mcp.ClientSession, name, argsJSON string) string {
	var args map[string]any
	json.Unmarshal([]byte(argsJSON), &args)

	res, err := server.CallTool(context.Background(), &mcp.CallToolParams{Name: name, Arguments: args})
	if err != nil {
		return "ошибка: " + err.Error()
	}

	text := ""
	for _, c := range res.Content {
		if t, ok := c.(*mcp.TextContent); ok {
			text += t.Text
		}
	}
	return text
}

func Close(servers map[string]*mcp.ClientSession) {
	for _, s := range servers {
		s.Close()
	}
}
