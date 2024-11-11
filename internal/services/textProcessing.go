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

type DigestText struct {
	Text string
	URL  string
}

type Parser interface {
	Parse(*string) string
}

type RVParser interface {
	ParseRV(string) (string, error)
}

type Service struct {
	queueLen uint8
	mu       sync.Mutex
	parser   Parser
	logger   *zerolog.Logger
	rvParser RVParser
}

func New(parser Parser, rvParser RVParser, logger *zerolog.Logger) *Service {
	return &Service{
		parser:   parser,
		logger:   logger,
		rvParser: rvParser,
		queueLen: 0,
	}
}

func (s *Service) parseDigest(model, text string) (string, error) {

	const prompt = "{{ short }}"

	links := make([]string, 0, 30)
	s.logger.Info().Msg("New request for DIGEST: " + text)
	// Split text by '\n' and ' ', and put it into links
	for _, line := range strings.Fields(text) {
		links = append(links, strings.Fields(line)...)
	}

	texts := make([]DigestText, 0, len(links))

	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}

	for _, link := range links {
		wg.Add(1)
		go func(link string) {
			defer wg.Done()
			text, err := s.rvParser.ParseRV(link)
			if err != nil {
				s.logger.Error().Msg("Error while parsing url: " + err.Error())
				mutex.Lock()
				defer mutex.Unlock()
				texts = append(texts, DigestText{Text: "НЕ УДАЛОСЬ ОБРАБОТАТЬ: " + link, URL: link})
			}
			if text != "" {
				mutex.Lock()
				defer mutex.Unlock()
				s_link := strings.ReplaceAll(link, "https://realnoevremya.ru", "")
				s_link = strings.ReplaceAll(s_link, "https://m.realnoevremya.ru", "")
				texts = append(texts, DigestText{Text: text, URL: s_link})
			} else {
				s.logger.Error().Msg("Invalid url: " + link)
			}
		}(link)
	}
	var resultTexts = ""
	resultTextList := make([]string, 0, len(texts))
	wg.Wait()

	for _, text := range texts {
		wg.Add(1)
		go func(text DigestText) {
			defer wg.Done()
			var resultText string
			var err error
			if strings.HasPrefix(text.Text, "НЕ УДАЛОСЬ ОБРАБОТАТЬ: ") {
				resultText = text.Text
			} else {
				resultText, err = s.ProcessText(model, prompt, text.Text, "0.1")
			}

			if err != nil {
				s.logger.Error().Msg("Error while processing text: " + err.Error())
				return
			}
			url := "<a href=\"" + text.URL + "\" target=\"_blank\">"

			var id int = 0

			for i, s := range resultText {
				if s == '.' || s == '!' || s == '?' || s == ',' || s == ';' || s == ':' || s == '—' || s == '-' {
					if (i+1 < len(resultText) && resultText[i+1] == ' ') || (i+1 == len(resultText)) {
						id = i
						break
					}
				}
			}
			resultText = strings.TrimSpace(strings.ReplaceAll(url+resultText[:id]+"</a>"+resultText[id:], "\n", " "))
			mutex.Lock()
			resultTextList = append(resultTextList, resultText)
			mutex.Unlock()
		}(text)
	}
	wg.Wait()
	for _, text := range resultTextList {
		resultTexts += "<p>" + text + "</p>\n\n"
	}
	resultTexts = strings.TrimSpace(resultTexts)
	return strings.TrimSpace(resultTexts), nil
}

// ProcessText - обработка текста с использованием Yandex GPT
func (s *Service) ProcessText(model, prompt, text, temperature string) (string, error) {

	if prompt == "{{ digest }}" {
		return s.parseDigest(model, text)
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
		}
	}

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
