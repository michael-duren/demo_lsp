package main

import (
	"demo_lsp/lsp"
	"demo_lsp/rpc"
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

	for scanner.Scan() {
		// Listener loop
		writer := os.Stdout
		msg := scanner.Bytes()
		handleMessage(msg, writer, logger)
	}
}

func handleMessage(msg []byte, writer *os.File, logger *log.Logger) {
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

	logFile, err := os.OpenFile(log_path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("hey, you didn't give a good file bozzo")
	}

	return log.New(logFile, "[demo_lsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
