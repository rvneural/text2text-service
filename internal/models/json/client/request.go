package client

import "encoding/xml"

// Request - структура для запроса от клиента
type Request struct {
	XMLName     xml.Name `xml:"request"`
	Model       string   `json:"model" xml:"model" form:"model"`
	Prompt      string   `json:"prompt" xml:"prompt" form:"prompt"`
	Text        string   `json:"text" xml:"text" form:"text"`
	Temperature string   `json:"temperature" xml:"temperature" form:"temperature"`
}
