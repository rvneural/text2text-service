package templates

import (
	"Text2TextService/internal/models/templates"
	"encoding/xml"
	"github.com/rs/zerolog"
	"os"
	"strconv"
	"strings"
)

type Parser struct {
	logger    *zerolog.Logger
	templates templates.Templates
}

// New creates a new Parser instance with the provided logger and reads the template file.
// If there is an error opening or decoding the template file, it logs the error using the provided logger.
//
// Parameters:
//   - logger: A pointer to a zerolog.Logger instance for logging errors.
//
// Returns:
//   - A pointer to a Parser instance with the initialized templates.
//     If there is an error opening or decoding the template file, it returns nil.
func New(logger *zerolog.Logger) *Parser {
	templateList := templates.Templates{}
	templateFile, err := os.Open("../../internal/models/templates/templates.xml")
	if err != nil {
		logger.Error().Msg("Error opening template file: " + err.Error())
		return nil
	}
	err = xml.NewDecoder(templateFile).Decode(&templateList)
	if err != nil {
		logger.Error().Msg("Error decoding template file: " + err.Error())
		return nil
	}
	return &Parser{logger: logger, templates: templateList}
}

// Parse parses the given content by replacing template placeholders with their corresponding values.
// It also determines the lowest temperature value among the replaced templates and returns it as a string.
// If no template is replaced, an empty string is returned.
//
// Parameters:
//   - content: A pointer to a string containing the content to be parsed.
//
// Returns:
//   - temperature: A string representing the lowest temperature value among the replaced templates.
//     If no template is replaced, an empty string is returned.
func (p *Parser) Parse(content *string) (temperature string) {
	p.logger.Info().Msg("Parsing content... " + *content)
	pTemperature := 1.0
	isTemperatureEdited := false
	for _, template := range p.templates.Templates {
		name := template.Name
		value := template.Value
		tTemperature, _ := strconv.ParseFloat(template.Temperature, 64)
		if strings.Contains(*content, "{{ "+name+" }}") {
			*content = strings.TrimSpace(strings.ReplaceAll(*content, "{{ "+name+" }}", value+" "))
			if tTemperature <= pTemperature {
				isTemperatureEdited = true
				pTemperature = tTemperature
			}
		}
	}
	p.logger.Info().Msg("Parsed content: " + *content)
	if isTemperatureEdited {
		temperature = strconv.FormatFloat(pTemperature, 'f', -1, 64)
		p.logger.Info().Msg("Lowest temperature value: " + temperature)
		return temperature
	}
	return ""
}
