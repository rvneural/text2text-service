package services

import (
	config "Text2TextService/internal/config/yandexGPT"
	yandexgpt "Text2TextService/internal/models/json/yandexGPT"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type PromptParser interface {
	Parse(*string) string
}

type RVParser interface {
	ParseRV(string) (string, error)
}

type AnotherParser interface {
	Parse(string) (string, error)
}

type Service struct {
	queueLen      uint8
	mu            sync.Mutex
	promptParser  PromptParser
	logger        *zerolog.Logger
	rvParser      RVParser
	anotherParser AnotherParser
}

func New(promptParser PromptParser, rvParser RVParser, anotherParser AnotherParser, logger *zerolog.Logger) *Service {
	return &Service{
		promptParser:  promptParser,
		logger:        logger,
		rvParser:      rvParser,
		queueLen:      0,
		anotherParser: anotherParser,
	}
}

func (s *Service) ProcessText(model, prompt, text, temperature string) (string, error) {

	if prompt == "{{ digest }}" {
		return s.parseDigest(model, text), nil
	}

	if s.getCurrentQueueLength() >= s.getMaxQueueLength() {
		time.Sleep(time.Second)
		return s.ProcessText(model, prompt, text, temperature)
	}

	s.incQueueLength()
	defer s.decQueueLength()

	var err error

	if len(strings.Fields(text)) == 1 {
		if strings.HasPrefix(text, "https://realnoevremya.ru/") || strings.HasPrefix(text, "https://m.realnoevremya.ru") {
			text, err = s.rvParser.ParseRV(strings.TrimSpace(text))
			if err != nil {
				return "", err
			}
		} else if strings.HasPrefix(text, "https://") || strings.HasPrefix(text, "http://") {
			text, err = s.anotherParser.Parse(text)
			if err != nil {
				return "", err
			}
		}
	}

	temp := s.promptParser.Parse(&prompt)

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

	s.logger.Info().Msg("Response text:" + Res.Result.Alternatives[0].Message.Text)
	return Res.Result.Alternatives[0].Message.Text, nil
}

func (s *Service) getMaxQueueLength() uint8 {
	if config.ERR == nil && config.MAX_PARALLEL_STR != "" {
		return uint8(config.MAX_PARALLEL)
	} else {
		if config.ERR != nil {
			s.logger.Error().Msg("Error while parsing MAX_PARALLEL_STR: " + config.ERR.Error())
		}
		return 2
	}
}

func (s *Service) getCurrentQueueLength() uint8 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.queueLen
}

func (s *Service) incQueueLength() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.queueLen++
}

func (s *Service) decQueueLength() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.queueLen--
}
