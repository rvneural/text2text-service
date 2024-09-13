package rest

import (
	client2 "Text2TextService/internal/models/json/client"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"net/http"
)

type Text2TextHandler struct {
	service Service
	logger  *zerolog.Logger
}

func New(service Service, logger *zerolog.Logger) *Text2TextHandler {
	return &Text2TextHandler{service: service, logger: logger}
}

// HandleRequest - обработка HTTP запроса от клиента
func (handler *Text2TextHandler) HandleRequest(c echo.Context) error {

	handler.logger.Info().Msg("Processing request from: " + c.RealIP())

	if c.Request().Header.Get("Content-Type") == "" {
		handler.logger.Error().Msg("Missing content type header")
		return c.JSON(http.StatusBadRequest, client2.Error{Error: "Missing content type header",
			Details: "Content type header is required\nUser application/json or application/x-www-form-urlencoded or application/xml"})
	}

	request := new(client2.Request)

	// Привязка запроса к структуре request
	err := c.Bind(request)

	if err != nil {
		handler.logger.Error().Msg("Error while binding request: " + err.Error())
		return c.JSON(http.StatusBadRequest, client2.Error{Error: "Invalid request body", Details: err.Error()})
	}

	// Проверка полей запроса.
	// Если поля отсутствуют, возвращается ошибка HTTP 400 и сообщение об ошибке.
	// Если все поля заполнены, обрабатывается текст и возвращается результат в виде JSON.
	if request.Prompt == "" || request.Text == "" {
		handler.logger.Error().Msg("Request body is missing required fields")
		return c.JSON(http.StatusBadRequest, client2.Error{Error: "Missing required fields in request body",
			Details: "You could use incorrect Content-Type in header"})
	}

	// Обработка текста
	result, err := handler.service.ProcessText(request.Model, request.Prompt, request.Text, request.Temperature)

	if err != nil {
		handler.logger.Error().Msg("Error while reading body: " + err.Error())
		return c.JSON(http.StatusInternalServerError, client2.Error{Error: "Internal server error", Details: err.Error()})
	}

	var response client2.Response

	response.NewText = result
	response.OldText = request.Text

	return c.JSON(http.StatusOK, response)
}
