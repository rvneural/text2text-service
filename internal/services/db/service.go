package db

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type DBService struct {
	url string
}

func New(url string) *DBService {
	return &DBService{
		url: url,
	}
}

func (w *DBService) RegisterOperation(uniqID string, operation_type string) error {
	uri := w.url

	type Request struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	}

	var request Request
	request.ID = uniqID
	request.Type = operation_type

	data, err := json.Marshal(request)
	if err != nil {
		return err
	}

	response, err := http.Post(uri, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Errorf("Error")
	}
	return nil
}

func (w *DBService) SetResult(uniqID string, data []byte) error {
	uri := w.url + "operation/" + uniqID
	response, err := http.Post(uri, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return fmt.Errorf("Error")
	}
	return nil
}
