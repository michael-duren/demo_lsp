package lsp

type BaseMessage struct {
	Method string `json:"method"`
}

type Request struct {
	Id      *int   `json:"id,omitempty"`
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	// we will add the params when we implement
}

type Response struct {
	Id      *int   `json:"id,omitempty"`
	Jsonrpc string `json:"jsonrpc"`
	// we will either have a result or an error when we implement
}

type Notification struct {
	Method string `json:"method"`
}
