package wwf

import (
	"HidirAi-Core/structs"
	"encoding/json"
	"errors"
	"os"
)

func DefaultMCP(userid string) error {
	os.MkdirAll("users/"+userid, 0755)
	file, err := os.OpenFile("./users/"+userid+"/MCP.jsonl", os.O_APPEND|os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return nil
		}
		return err
	}
	defer file.Close()

	mcps := []structs.MCPServer{
		{
			Url: "https://mcp.tavily.com/mcp/?tavilyApiKey=" + os.Getenv("TAVILY_API"),
			Key: "",
		},
		{
			Url: "https://api.githubcopilot.com/mcp/",
			Key: os.Getenv("GHAPI"),
		},
	}

	for _, m := range mcps {
		data, _ := json.Marshal(m)
		data = append(data, '\n')

		file.Write(data)
	}

	return nil
}
