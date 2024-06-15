package rpc

import (
	"bytes"
	"demo_lsp/lsp"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

func EncodeMsg(msg any) string {
	content, err := json.Marshal(msg)
	if err != nil {
		panic("UNABLE TO ENCODE MSG")
	}

	return fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content)
}

// Goal of this function is to just determine the method of the message,
// if it's a valid message and return the content of the message
// This way we can unmarshall the content of the message into the correct struct
func DecodeMsg(msg []byte) (string, []byte, error) {
	// Takes a message and decodes to a struct
	header, content, found := bytes.Cut(msg, []byte{'\r', '\n', '\r', '\n'})

	if !found {
		return "", nil, errors.New("header not found")
	}

	// Parse the header
	contentLengthStr := header[len("Content-Length: "):]
	contentLength, err := strconv.Atoi(string(contentLengthStr))
	if err != nil {
		return "", nil, err
	}

	// Determine the method
	var baseMessage lsp.BaseMessage

	if err := json.Unmarshal(content[:contentLength], &baseMessage); err != nil {
		return "", nil, err
	}

	return baseMessage.Method, content[:contentLength], nil
}
