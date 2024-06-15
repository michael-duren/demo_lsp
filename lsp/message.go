package lsp

type Message struct {
	Jsonrpc string `json:"jsonrpc"`
}

type BaseMessage struct {
	Method string `json:"method"`
}

type Request struct {
	Id *int `json:"id"`
	BaseMessage
	// we will add the params when we implement
}

type Response struct {
	Id *int `json:"id"`
	Message
	// we will either have a result or an error when we implement
}

type Notification struct {
	Method string `json:"method"`
}
