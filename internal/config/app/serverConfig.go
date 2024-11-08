package app

import "os"

const (
	ADDR = ":80"
)

var (
	BEARER_KEY = os.Getenv("KEY")
)
