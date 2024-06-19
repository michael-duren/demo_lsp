package main

import (
	"demo_lsp/lsp"
	"demo_lsp/rpc"
	"demo_lsp/thesaurus"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {
	logger := getLogger(nil)
	logger.Println("Hey buddy, I started")

	// Setup scanner
	scanner := rpc.CreateScanner()
	state := thesaurus.NewState()

	for scanner.Scan() {
		// Listener loop
		writer := os.Stdout
		msg := scanner.Bytes()
		handleMessage(msg, state, writer, logger)
	}
}

func handleMessage(msg []byte, state thesaurus.State, writer *os.File, logger *log.Logger) {
	logger.Printf("Message received: %s\n", msg)
	method, content, err := rpc.DecodeMsg(msg)
	if err != nil {
		logger.Printf("In \"handleMessage\" Error decoding message: %s \n", err)
		return
	}

	switch method {
	case "initialize":
		var initializeRequest lsp.Initialize

		if err := json.Unmarshal(content, &initializeRequest); err != nil {
			logger.Printf("Error unmarshalling initialize request: %s\n", err)
			return
		}

		logger.Printf("Initialize Request by client: %s, version: %s",
			initializeRequest.Params.ClientInfo.Name,
			initializeRequest.Params.ClientInfo.Version)

		initializeResponse := lsp.NewInitializeResponse(*initializeRequest.Id)
		writeResponse(initializeResponse, writer, logger)

	case "textDocument/didOpen":
		var request lsp.DidOpenTextDocumentNotification

		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("textDocument/didOpen: %s\n", err)
			return
		}

		logger.Printf("Connected to: %s", request.Params.TextDocument.Uri)
		logger.Printf("Contents: %s", request.Params.TextDocument.Text)
		state.OpenDocument(request.Params.TextDocument.Uri, request.Params.TextDocument.Text) // open document
	case "textDocument/didChange":
		var request lsp.TextDocumentDidChangeNotification

		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("textDocument/didChange: %s\n", err)
			return
		}

		logger.Printf("Changed document: %s", request.Params.TextDocument.Uri)
		logger.Printf("Changes: %s", request.Params.ContentChanges[0].Text)

		for _, change := range request.Params.ContentChanges {
			// should only be one change
			state.UpdateDocument(request.Params.TextDocument.Uri, change.Text)
		}

	case "textDocument/hover":
		logger.Println("Hover Request")
		var hoverRequest lsp.HoverRequest

		if err := json.Unmarshal(content, &hoverRequest); err != nil {
			logger.Printf("Error unmarshalling hover request: %s\n", err)
			return
		}

		hoveredWord := state.GetWordFromRange(hoverRequest.Params.TextDocument.Uri, hoverRequest.Params.TextDocumentPositionParams)

		// reply with the word
		hoverResponse := lsp.NewHoverResponse(*hoverRequest.Id, hoveredWord)
		writeResponse(hoverResponse, writer, logger)
	}
}

func writeResponse(msg any, writer io.Writer, logger *log.Logger) {
	encodedMsg := rpc.EncodeMsg(msg)
	_, err := writer.Write([]byte(encodedMsg))
	if err != nil {
		logger.Printf("Error writing response: %s\n", err)
	} else {
		logger.Printf("Response sent: %s\n", encodedMsg)
	}
}

func getLogger(filename *string) *log.Logger {
	cwd, err := os.Getwd()
	if err != nil {
		panic("Could not get the current working directory")
	}
	if filename == nil {
		name := fmt.Sprintf("log_%s.log", time.Now().Format("2006-01-02"))
		filename = &name
	}
	log_path := filepath.Join(cwd, "logs", *filename)
	// make logs directory if it doesn't exist
	// check if the directory exists
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		// create the directory
		err := os.Mkdir("logs", os.ModePerm)
		if err != nil {
			panic("Could not create the logs directory")
		}
	}

	logFile, err := os.OpenFile(log_path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o666)
	if err != nil {
		panic("hey, you didn't give a good file bozzo")
	}

	return log.New(logFile, "[demo_lsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
