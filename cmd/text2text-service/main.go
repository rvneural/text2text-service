package main

import (
	"Text2TextService/cmd/log"
	"Text2TextService/internal/pkg/app"
)

// Main — входная точка в приложение
func main() {

	// Инициализация логгера
	logger := log.NewLogger()

	// Создание приложения
	app := app.New(&logger)

	// Запуск приложения и логирование завершения работы
	logger.Fatal().Msg(app.Run().Error())
}
