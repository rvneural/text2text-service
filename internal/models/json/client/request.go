package client

import "encoding/xml"

// Request - структура для запроса от клиента
type Request struct {
	XMLName      xml.Name `xml:"request"`
	Operation_ID string   `json:"operation_id" xml:"operation_id" form:"operation_id"`
	Model        string   `json:"model" xml:"model" form:"model"`
	Prompt       string   `json:"prompt" xml:"prompt" form:"prompt"`
	Text         string   `json:"text" xml:"text" form:"text"`
	Temperature  string   `json:"temperature" xml:"temperature" form:"temperature"`
}
