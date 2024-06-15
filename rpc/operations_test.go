package rpc_test

import (
	"demo_lsp/rpc"
	"fmt"
	"testing"
)

type SomeMessage struct {
	MessageName string `json:"messageName"`
}

func TestEncodeMsg(t *testing.T) {
	expected := "Content-Length: 22\r\n\r\n{\"messageName\":\"test\"}"
	actual := rpc.EncodeMsg(SomeMessage{MessageName: "test"})

	if expected != actual {
		t.Fatalf("TestEncodeMsg did not properly encode expected: %s actual: %s", expected, actual)
	}
}

func TestDecodeMsg(t *testing.T) {
	// arrange
	testContent := "{\"method\":\"someMethod\"}"
	testMessage := fmt.Sprintf("Content-Length: 23\r\n\r\n%s", testContent)

	// act
	method, content, err := rpc.DecodeMsg([]byte(testMessage))
	// assert
	if err != nil {
		t.Fatal(err)
	}

	if method != "someMethod" {
		t.Fatalf("Method: %s did not equal \"someMethod\"", method)
	}

	if string(content) != testContent {
		t.Fatalf("Content: %s did not equal testContent: %s", content, testContent)
	}
}
