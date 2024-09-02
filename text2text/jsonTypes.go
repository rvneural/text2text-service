package text2text

type Message struct {
	Role string `json:"role"`
	Text string `json:"text"`
}

type Request struct {
	ModelURI          string `json:"modelUri"`
	CompletionOptions struct {
		Stream      bool   `json:"stream"`
		Temperature string `json:"temperature"`
		MaxTokens   string `json:"maxTokens"`
	} `json:"completionOptions"`
	Messages []Message `json:"messages"`
}

type Alternative struct {
	Message Message `json:"message"`
	Status  string  `json:"status"`
}

type Response struct {
	Result struct {
		Alternatives []Alternative `json:"alternatives"`
	} `json:"result"`
}
