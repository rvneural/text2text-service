package text2text

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

const (
	GPT_API_KEY = "AQVNzU5pE0yY9za2gCmN9dGxE0B0nb6hKm8aWzdX"
	STORAGE_ID  = "b1gjtlqofdt5mu5io6a9"
	URI         = "https://llm.api.cloud.yandex.net/foundationModels/v1/completion"
	MAX_TOKENS  = "2000"
)

func getRequestBody(model, promt, text, temperature string) ([]byte, error) {
	var Req Request

	if strings.ToLower(model) == "lite" {
		Req.ModelURI = "gpt://" + STORAGE_ID + "/yandexgpt-lite/rc"
	} else {
		Req.ModelURI = "gpt://" + STORAGE_ID + "/yandexgpt/latest"
	}

	Req.CompletionOptions.MaxTokens = MAX_TOKENS
	Req.CompletionOptions.Stream = false
	Req.CompletionOptions.Temperature = temperature

	var systemMessage Message
	systemMessage.Role = "system"
	systemMessage.Text = promt

	var userMessage Message
	userMessage.Role = "user"
	userMessage.Text = text

	Req.Messages = []Message{systemMessage, userMessage}
	dataReq, err := json.Marshal(&Req)

	if err != nil {
		log.Println("Error while marshalling request body >>>>>>", err)
		return []byte{}, err
	}
	return dataReq, nil
}

func ProccessText(model, promt, text, temperature string) (string, error) {

	log.Println("Promt: ", promt)
	log.Println("Text:", strings.ReplaceAll(text, "\n", " "))
	log.Println("Temperature:", temperature)

	// Создание тела запроса
	data, err := getRequestBody(model, promt, text, temperature)

	if err != nil {
		log.Println("Error while creating request body >>>>>>", err)
		return "", err
	}

	// Создание запроса
	reader := bytes.NewReader(data)
	request, err := http.NewRequest("POST", URI, reader)
	request.Header.Set("Authorization", "Api-Key "+GPT_API_KEY)

	if err != nil {
		log.Println("Error while creating request >>>>>>", err)
		return "", err
	}

	log.Println("Request:", request.Host)

	// Создание HTTP клиента
	client := &http.Client{}

	// Отправка запроса
	response, err := client.Do(request)

	if err != nil {
		log.Println("Error while sending request >>>>>>", err)
		return "", err
	}
	defer response.Body.Close()

	log.Println("Request sent")

	if response.StatusCode != 200 {
		log.Println("Status Code:", response.StatusCode)
		return response.Status, nil
	}

	// Чтение результата
	responseData, err := io.ReadAll(response.Body)

	if err != nil {
		log.Println("Error while reading response data >>>>>>", err)
		return "", err
	}

	// Превращение JSON в структуру
	var Res Response
	err = json.Unmarshal(responseData, &Res)

	if err != nil {
		log.Println("Error while unmarshalling response data >>>>>>", err)
		return "", err
	}

	return Res.Result.Alternatives[0].Message.Text, nil
}
