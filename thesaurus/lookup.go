package thesaurus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Synonyms struct {
	Word string `json:"word"`
}

const baseUrl = "https://api.datamuse.com"

func Lookup(word string) ([]string, error) {
	url := fmt.Sprintf("%s/words?ml=%s", baseUrl, word)

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var synonymResult []Synonyms

	err = json.Unmarshal(body, &synonymResult)
	if err != nil {
		return nil, err
	}

	var synonyms []string

	for _, synonym := range synonymResult {
		synonyms = append(synonyms, synonym.Word)
	}

	return synonyms, nil
}
