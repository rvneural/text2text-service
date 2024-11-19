package app

import "os"

const (
	ADDR = ":8081"
)

var (
	DB_URL = os.Getenv("DB_URL")
)
