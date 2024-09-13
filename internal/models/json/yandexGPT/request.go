package yandexgpt

type Message struct {
	Role string `json:"role"`
	Text string `json:"text"`
}

// Request - структура для запроса к Yandex GPT
type Request struct {
	ModelURI          string `json:"modelUri"`
	CompletionOptions struct {
		Stream      bool   `json:"stream"`
		Temperature string `json:"temperature"`
		MaxTokens   string `json:"maxTokens"`
	} `json:"completionOptions"`
	Messages []Message `json:"messages"`
}
