package yandexgpt

type Alternative struct {
	Message Message `json:"message"`
	Status  string  `json:"status"`
}

// Response - структура для ответа Yandex GPT
type Response struct {
	Result struct {
		Alternatives []Alternative `json:"alternatives"`
	} `json:"result"`
}
