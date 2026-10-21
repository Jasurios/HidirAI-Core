package wwf

import (
	"HidirAi-Core/structs"
	"encoding/json"
	"os"
)

func Newmcp(userid, url, key string) error {
	os.MkdirAll("users/"+userid, 0755)
	file, err := os.OpenFile("./users/"+userid+"/MCP.jsonl", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	mcps := []structs.MCPServer{
		{
			Url: url,
			Key: key,
		},
	}

	for _, m := range mcps {
		data, _ := json.Marshal(m)
		data = append(data, '\n')

		file.Write(data)
	}

	return nil
}
