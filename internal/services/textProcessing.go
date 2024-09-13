package services

import (
	config "Text2TextService/internal/config/yandexGPT"
	"Text2TextService/internal/models/json/yandexGPT"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
)

type Parser interface {
	Parse(*string) string
}

type Service struct {
	parser Parser
	logger *zerolog.Logger
}

func New(parser Parser, logger *zerolog.Logger) *Service {
	return &Service{
		parser: parser,
		logger: logger,
	}
}

// ProcessText - обработка текста с использованием Yandex GPT
func (s *Service) ProcessText(model, prompt, text, temperature string) (string, error) {

	temp := s.parser.Parse(&prompt)

	if model == "" {
		model = config.DEFAULT_MODEL
	}

	if temperature == "" {
		if temp != "" {
			temperature = temp
		} else {
			temperature = config.DEFAULT_TEMPERATURE
		}
	}

	s.logger.Info().Msg("Promt: " + prompt)
	s.logger.Info().Msg("Text: " + strings.ReplaceAll(text, "\n", " "))
	s.logger.Info().Msg("Temperature: " + temperature)

	// Создание тела запроса к Yandex GPT
	data, err := s.getRequestBody(model, prompt, text, temperature)

	if err != nil {
		s.logger.Error().Msg("Error while creating request body >>>>>> " + err.Error())
		return "", err
	}

	// Создание http запроса к Yandex GPT
	reader := bytes.NewReader(data)
	request, err := http.NewRequest("POST", config.URI, reader)
	request.Header.Set("Authorization", "Api-Key "+config.GPT_API_KEY)

	if err != nil {
		s.logger.Error().Msg("Error while creating request >>>>>> " + err.Error())
		return "", err
	}

	s.logger.Debug().Msg("Request: " + request.Host)

	// Создание HTTP клиента
	client := &http.Client{}

	// Отправка запроса
	response, err := client.Do(request)

	if err != nil {
		s.logger.Error().Msg("Error while sending request >>>>>> " + err.Error())
		return "", err
	}

	// Закрытие соединения с ресурсом
	defer response.Body.Close()

	s.logger.Debug().Msg("Request sent")

	// Проверка кода ответа и получение текста из результата
	if response.StatusCode != 200 {
		s.logger.Error().Msg("Status Code:" + strconv.Itoa(response.StatusCode))
		body, _ := io.ReadAll(response.Body)
		s.logger.Debug().Msg(string(body))
		return response.Status, nil
	}

	// Чтение результата
	responseData, err := io.ReadAll(response.Body)

	if err != nil {
		s.logger.Error().Msg("Error while reading response data >>>>>>" + err.Error())
		return "", err
	}

	// Превращение JSON в структуру
	var Res yandexgpt.Response
	err = json.Unmarshal(responseData, &Res)

	if err != nil {
		s.logger.Error().Msg("Error while unmarshalling response data >>>>>>" + err.Error())
		return "", err
	}

	return Res.Result.Alternatives[0].Message.Text, nil
}
