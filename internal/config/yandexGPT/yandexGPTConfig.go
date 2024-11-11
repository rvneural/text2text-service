package yandexgpt

import (
	"os"
	"strconv"
)

const (
	URI        = "https://llm.api.cloud.yandex.net/foundationModels/v1/completion"
	MAX_TOKENS = "32000"

	DEFAULT_TEMPERATURE = "0.3"
	DEFAULT_MODEL       = "pro"
)

var (
	GPT_API_KEY       = os.Getenv("API_KEY")
	STORAGE_ID        = os.Getenv("STORAGE_ID")
	MAX_PARALLEL_STR  = os.Getenv("MAX_PARALLEL")
	MAX_PARALLEL, ERR = strconv.Atoi(MAX_PARALLEL_STR)
)
