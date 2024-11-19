package app

import (
	endpoint "Text2TextService/internal/endpoint/app"
	service "Text2TextService/internal/services"
	rvParser "Text2TextService/internal/services/rvparser"
	promptParser "Text2TextService/internal/services/templates"
	textHandler "Text2TextService/internal/transport/rest"
	templatesHandler "Text2TextService/internal/transport/rest/getTemplates"

	config "Text2TextService/internal/config/app"
	"Text2TextService/internal/services/db"

	"github.com/rs/zerolog"
)

type App struct {
	endpoint *endpoint.App
	service  *service.Service
	handler  *textHandler.Text2TextHandler
	logger   *zerolog.Logger
}

// New - конструктор приложения
func New(logger *zerolog.Logger) *App {
	logger.Debug().Msg("Initializing app...")
	parser := promptParser.New(logger)
	rvpars := rvParser.New(logger)
	service := service.New(parser, rvpars, logger)

	dbWorker := db.New(config.DB_URL)
	logger.Info().Msg("DBURL:" + config.DB_URL)

	textHandler := textHandler.New(service, dbWorker, logger)
	templatesHandler := templatesHandler.New(logger)

	endpoint := endpoint.New(textHandler, templatesHandler, logger)
	return &App{endpoint: endpoint, service: service, handler: textHandler, logger: logger}
}

// Run - запуск приложения на сервере
func (app *App) Run() error {
	app.logger.Info().Msg("Starting server...")
	return app.endpoint.Start()
}
