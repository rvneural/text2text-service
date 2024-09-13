package getTemplates

import (
	"Text2TextService/internal/models/json/client"
	"Text2TextService/internal/models/templates"
	"encoding/xml"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"net/http"
	"os"
)

type Handler struct {
	logger    *zerolog.Logger
	templates templates.Templates
}

func New(logger *zerolog.Logger) *Handler {
	templateList := templates.Templates{}
	templateFile, err := os.Open("../../internal/models/templates/templates.xml")
	if err != nil {
		logger.Error().Msg("Error opening template file in getTemplates: " + err.Error())
	}
	err = xml.NewDecoder(templateFile).Decode(&templateList)
	if err != nil {
		logger.Error().Msg("Error decoding template file in getTemplates: " + err.Error())
	}
	return &Handler{
		logger:    logger,
		templates: templateList,
	}
}

func (h *Handler) HandleRequest(c echo.Context) error {
	h.logger.Info().Msg("Get request received for getting templates from: " + c.RealIP())
	if c.Request().Header.Get("Content-Type") == "application/json" || c.Request().Header.Get("Content-Type") == "" {
		return c.JSON(http.StatusOK, h.templates)
	} else if c.Request().Header.Get("Content-Type") == "application/xml" {
		return c.XML(http.StatusOK, h.templates)
	} else {
		return c.JSON(http.StatusBadRequest, client.Error{Error: "Invalid content type header",
			Details: "Content type header should be application/json or application/xml"})
	}
}
