package handler

import (
	"encoding/json"
	"net/http"
)

type ReturnMessage struct {
	Message string `json:"message"`
	State   string `json:"state"`
}

func WriteMessage(response http.ResponseWriter, message string) error {
	return json.NewEncoder(response).Encode(ReturnMessage{
		Message: message,
	})
}
