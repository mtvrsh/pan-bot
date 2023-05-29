package crappers //(s)crappers

import (
	"html"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
)

const (
	sjpURL     = "https://sjp.pl/"
	emptyQuery = "Puste zapytanie, pusta odpowiedź :)"
)

// QuerySjp returns meaning(s) of word found in html from sjpURL.
// Doesn't return error if there is no description or word doesn't exist.
func QuerySjp(word string) (string, error) {
	word = strings.TrimSpace(word)
	if word == "" {
		return emptyQuery, nil
	}

	resp, err := http.Get(sjpURL + word)
	if err != nil {
		log.Println(err)
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer resp.Body.Close()

	if strings.Contains(string(body), "<p>nie występuje w słowniku</p>") {
		return "Nie występuje w słowniku", nil
	}

	parAll := regexp.MustCompile("<p style=\"margin: .5em 0; font: medium/1.4 sans-serif; max-width: 34em; \">.*</p>")
	parStart := regexp.MustCompile("<p style=\"margin: .5em 0; font: medium/1.4 sans-serif; max-width: 34em; \">")
	parEnd := regexp.MustCompile("</p>")
	linebreak := regexp.MustCompile("<br />")
	trailing := regexp.MustCompile("[;,\\. ]+\n")

	text := parAll.FindAllString(string(body), -1)

	var result string
	for _, x := range text {
		s := linebreak.ReplaceAllString(x, "\n")
		s = parStart.ReplaceAllString(s, "")
		s = parEnd.ReplaceAllString(s, "")
		s = trailing.ReplaceAllString(s, "\n")
		result += s + "\n"
	}

	result = strings.Trim(result, " \t\n")
	if result == "" {
		return word + " nie ma opisu", nil
	}

	return html.UnescapeString(result), nil
}
