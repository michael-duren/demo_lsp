package rpc

import (
	"bufio"
	"bytes"
	"os"
	"strconv"
)

func CreateScanner() *bufio.Scanner {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(SplitMessage)
	return scanner
}

// type SplitFunc func(data []byte, atEOF bool) (advance int, token []byte, err error)

func SplitMessage(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// look for the seperator between header and content
	header, content, found := bytes.Cut(data, []byte{'\r', '\n', '\r', '\n'})

	if !found {
		return 0, nil, nil
	}

	// get the content-length
	contentLengthBytes := header[len("Content-Length: "):]
	// parse to int
	contentLength, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return 0, nil, err
	}

	if len(content) < contentLength {
		return 0, nil, nil // missing content
	}

	totalLength := len(header) + contentLength + 4

	return totalLength, data[:totalLength], nil
}
