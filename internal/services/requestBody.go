package services

import (
	config "Text2TextService/internal/config/yandexGPT"
	yandexgpt "Text2TextService/internal/models/json/yandexGPT"
	"encoding/json"
	"errors"
	"os"
	"strings"
)

// getRequestBody - формирование тела запроса к Yandex GPT
func (s *Service) getRequestBody(model, prompt, text, temperature string) ([]byte, error) {
	var Req yandexgpt.Request

	model = strings.TrimSpace(strings.ToLower(model))

	model_type := os.Getenv("MODEL_TYPE")
	if model_type == "" {
		model_type = "latest"
	}

	if model == "lite" {
		Req.ModelURI = "gpt://" + config.STORAGE_ID + "/yandexgpt-lite/" + model_type
	} else if model == "pro" {
		Req.ModelURI = "gpt://" + config.STORAGE_ID + "/yandexgpt/" + model_type
	} else {
		return nil, errors.New("unsupported model")
	}

	Req.CompletionOptions.MaxTokens = config.MAX_TOKENS
	Req.CompletionOptions.Stream = false
	Req.CompletionOptions.Temperature = temperature

	var systemMessage yandexgpt.Message
	systemMessage.Role = "system"
	systemMessage.Text = prompt

	var userMessage yandexgpt.Message
	userMessage.Role = "user"
	userMessage.Text = text

	Req.Messages = []yandexgpt.Message{systemMessage, userMessage}
	dataReq, err := json.Marshal(&Req)

	if err != nil {
		s.logger.Error().Msg("Error while marshalling request body >>>>>> " + err.Error())
		return []byte{}, err
	}
	return dataReq, nil
}
