package app

import "os"

const (
	ADDR = ":8081"
)

var (
	BEARER_KEY = os.Getenv("KEY")
)
