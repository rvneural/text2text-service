package rest

import (
	"Text2TextService/internal/models/json/client"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type DBWorker interface {
	RegisterOperation(uniqID string, operation_type string, user_id int) error
	SetResult(uniqID string, data []byte) error
}

type Text2TextHandler struct {
	service  Service
	dbWorker DBWorker
	logger   *zerolog.Logger
}

func New(service Service, dbWorker DBWorker, logger *zerolog.Logger) *Text2TextHandler {
	return &Text2TextHandler{service: service, logger: logger, dbWorker: dbWorker}
}

// HandleRequest - обработка HTTP запроса от клиента
func (handler *Text2TextHandler) HandleRequest(c echo.Context) error {

	handler.logger.Info().Msg("Processing request from: " + c.RealIP())

	if c.Request().Header.Get("Content-Type") == "" {
		handler.logger.Error().Msg("Missing content type header")
		return c.JSON(http.StatusBadRequest, client.Error{Error: "Missing content type header",
			Details: "Content type header is required\nUser application/json or application/x-www-form-urlencoded or application/xml"})
	}

	request := new(client.Request)

	// Привязка запроса к структуре request
	err := c.Bind(request)

	if err != nil {
		handler.logger.Error().Err(err).Msg("Error while binding request")
		return c.JSON(http.StatusBadRequest, client.Error{Error: "Invalid request body", Details: err.Error()})
	}

	// Проверка полей запроса.
	// Если поля отсутствуют, возвращается ошибка HTTP 400 и сообщение об ошибке.
	// Если все поля заполнены, обрабатывается текст и возвращается результат в виде JSON.
	if request.Prompt == "" || request.Text == "" {
		handler.logger.Error().Msg("Request body is missing required fields")
		return c.JSON(http.StatusBadRequest, client.Error{Error: "Missing required fields in request body",
			Details: "You could use incorrect Content-Type in header"})
	}

	if request.Operation_ID != "" {
		handler.logger.Info().Msg("Saving operation ID: " + request.Operation_ID)
		user_id, err := strconv.Atoi(request.UserID)
		if err != nil {
			user_id = 0
		}
		handler.dbWorker.RegisterOperation(request.Operation_ID, "text", user_id)
	}

	// Обработка текста
	result, err := handler.service.ProcessText(request.Model, request.Prompt, request.Text, request.Temperature)

	if err != nil {
		handler.logger.Error().Err(err).Msg("Error while reading body")
		return c.JSON(http.StatusInternalServerError, client.Error{Error: "Internal server error", Details: err.Error()})
	}

	var response client.Response

	response.NewText = result
	response.OldText = request.Text

	if request.Operation_ID != "" {
		dbResult := client.DBResult{
			NewText: response.NewText,
			OldText: response.OldText,
			Prompt:  request.Prompt,
		}
		byteResponse, _ := json.Marshal(dbResult)
		handler.dbWorker.SetResult(request.Operation_ID, byteResponse)
	}

	return c.JSON(http.StatusOK, response)
}
