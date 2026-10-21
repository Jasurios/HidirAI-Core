package structs

type Newkey struct {
	Type string `json:"type"`
	User string `json:"user"`
	API  string `json:"API"`
}

type ChatService struct {
	User      string `json:"user"`
	Sysprompt string `json:"sysprompt"`
	Text      string `json:"text"`
	Image     string `json:"image"`
	Audio     []byte `json:"audio"`
	File      []byte `json:"file"`
}

type ChatUser struct {
	Sysprompt string `json:"sysprompt"`
	Text      string `json:"text"`
	Image     string `json:"image"`
	Audio     []byte `json:"audio"`
	File      []byte `json:"file"`
}

type File struct {
	Name string `json:"name"`
	Data []byte `json:"data"`
}

type ChatResponse struct {
	Text  string `json:"text"`
	Image string `json:"image"`
	Audio []byte `json:"audio"`
	File  *File  `json:"file"`
}

type MCPServerService struct {
	User string `json:"user"`
	Url  string `json:"url"`
	Key string `json:"key"`
}

type MCPServer struct{
	Url string `json:"url"`
	Key string `json:"key"`
}
