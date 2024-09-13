package client

// Response - структура для ответа клиенту
type Response struct {
	NewText string `json:"newText"`
	OldText string `json:"oldText"`
}

type Error struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}
