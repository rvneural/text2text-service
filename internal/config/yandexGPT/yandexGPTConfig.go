package yandexgpt

import "os"

const (
	URI        = "https://llm.api.cloud.yandex.net/foundationModels/v1/completion"
	MAX_TOKENS = "32000"

	DEFAULT_TEMPERATURE = "0.3"
	DEFAULT_MODEL       = "pro"
)

var (
	GPT_API_KEY = os.Getenv("GPT_API_KEY")
	STORAGE_ID  = os.Getenv("STORAGE_ID")
)
