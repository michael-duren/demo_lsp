package lsp

type Initialize struct {
	Request
	Params InitializeParams `json:"params"`
}

type InitializeParams struct {
	ClientInfo ClientInfo `json:"clientInfo"`
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InitializeResponse struct {
	// this is the first message that the server sends to the client
	Response
	Result InitializeResult `json:"result"`
}

type InitializeResult struct {
	ServerInfo   ServerInfo         `json:"serverInfo"`
	Capabilities ServerCapabilities `json:"capabilities"`
}

type ServerCapabilities struct {
	TextDocumentSync int `json:"textDocumentSync"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func NewInitializeResponse(id int) InitializeResponse {
	return InitializeResponse{
		Response: Response{
			Id:      &id,
			Jsonrpc: "2.0",
		},
		Result: InitializeResult{
			Capabilities: ServerCapabilities{
				TextDocumentSync: 1,
			}, // what we're able to do
			ServerInfo: ServerInfo{
				Name:    "demo_lsp", // name of lsp
				Version: "0.0.0.0",  // version of lsp
			},
		},
	}
}
