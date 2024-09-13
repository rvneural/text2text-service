package rest

type Service interface {
	ProcessText(model, prompt, text, temperature string) (string, error)
}
