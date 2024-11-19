package app

import (
	config "Text2TextService/internal/config/app"

	"github.com/rs/zerolog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type App struct {
	Text2TextHandler    HandleFunc
	GetTemplatesHandler HandleFunc
	logger              *zerolog.Logger
}

func New(text2textHandler HandleFunc, getTemplates HandleFunc, logger *zerolog.Logger) *App {
	return &App{
		Text2TextHandler:    text2textHandler,
		GetTemplatesHandler: getTemplates,
		logger:              logger,
	}
}

// Start - запуск сервера на указанном порте и адресе
func (app *App) Start() error {

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CSRF())
	e.Use(middleware.Recover())

	e.POST("/", app.Text2TextHandler.HandleRequest)
	e.GET("/", app.GetTemplatesHandler.HandleRequest)

	app.logger.Info().Msg("Server started at " + config.ADDR)
	return e.Start(config.ADDR)
}
