package anothersiteparser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog"
)

type Parser struct {
	logger *zerolog.Logger
	url    string
}

func New(logger *zerolog.Logger) *Parser {
	return &Parser{
		logger: logger,
		url:    os.Getenv("WEB_PARSER"),
	}
}

func (p *Parser) Parse(url string) (string, error) {
	body := map[string]string{
		"url": url,
	}
	byteBody, err := json.Marshal(body)
	if err != nil {
		p.logger.Error().Msg("Error marshaling body: " + err.Error())
		return "", err
	}
	reader := bytes.NewReader(byteBody)
	response, err := http.Post(p.url, "application/json", reader)
	if err != nil {
		p.logger.Error().Msg("Error sending request: " + err.Error())
		return "", err
	}
	type Response struct {
		Text  string `json:"text"`
		Error string `json:"error"`
	}
	resp := Response{}
	err = json.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		p.logger.Error().Msg("Error decoding response: " + err.Error())
		return "", err
	} else if resp.Error != "" {
		p.logger.Error().Msg("Error from parser: " + resp.Error)
		return "", fmt.Errorf(resp.Error)
	} else {
		return resp.Text, nil
	}
}
