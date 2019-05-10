package handler

import (
	"encoding/json"
	"net/http"
)

type ReturnMessage struct {
	Message string `json:"message"`
	State   string `json:"state,omitempty"`
}

func WriteMessage(response http.ResponseWriter, message string) error {
	return writeMessage(response, message, "")
}

func WriteStateMessage(response http.ResponseWriter, message, state string) error {
	return writeMessage(response, state, message)
}

func writeMessage(response http.ResponseWriter, message, state string) error {
	return json.NewEncoder(response).Encode(ReturnMessage{
		Message: message,
		State:   state,
	})
}
