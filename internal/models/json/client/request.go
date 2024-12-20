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
	UserID       int      `json:"user_id" xml:"user_id" form:"user_id"`
}

type DBResult struct {
	OldText string `json:"old_text"`
	Prompt  string `json:"prompt"`
	NewText string `json:"new_text"`
}
