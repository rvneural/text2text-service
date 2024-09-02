package server

type Request struct {
	Model       string `json:"model"`
	Promt       string `json:"promt"`
	Text        string `json:"text"`
	Temperature string `json:"temperature"`
}

type Response struct {
	NewText string `json:"newText"`
	OldText string `json:"oldText"`
}
