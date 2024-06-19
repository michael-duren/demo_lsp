package thesaurus

import (
	"demo_lsp/lsp"
	"demo_lsp/util"
	"strings"
)

type State struct {
	// Filenames to contents
	Documents map[string]string
}

func NewState() State {
	return State{
		Documents: make(map[string]string),
	}
}

func (s *State) OpenDocument(uri, text string) {
	s.Documents[uri] = text
}

func (s *State) UpdateDocument(uri, text string) {
	s.Documents[uri] = text
}

func (s *State) GetWordFromRange(uri string, textDocumentPositionParams lsp.TextDocumentPositionParams) string {
	doc := s.Documents[uri]

	lines := strings.Split(doc, "\n")
	currentLine := lines[textDocumentPositionParams.Position.Line]

	selectedWord := ""

	// prepend letter
	for i := textDocumentPositionParams.Position.Character; i >= 0; i-- {
		currentByte := currentLine[i]
		if util.IsWhitespace(currentByte) {
			break
		}
		selectedWord = string(currentByte) + selectedWord
	}

	// append letters
	for i := textDocumentPositionParams.Position.Character + 1; i < len(currentLine); i++ {
		currentByte := currentLine[i]

		if util.IsWhitespace(currentByte) {
			break
		}

		selectedWord = selectedWord + string(currentByte)
	}

	return selectedWord
}
