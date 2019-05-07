package handler

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func NotFoundErrorHandler(responseWriter http.ResponseWriter, request *http.Request, err error) {
	handleError(responseWriter, request, http.StatusNotFound, err, log.WarnLevel)
}

func BadRequestErrorHandler(responseWriter http.ResponseWriter, request *http.Request, err error) {
	handleError(responseWriter, request, http.StatusBadRequest, err, log.WarnLevel)
}

func ForbiddenErrorHandler(responseWriter http.ResponseWriter, request *http.Request, err error) {
	handleError(responseWriter, request, http.StatusForbidden, err, log.WarnLevel)
}

func MethodNotAllowedErrorHandler(responseWriter http.ResponseWriter, request *http.Request, err error) {
	handleError(responseWriter, request, http.StatusMethodNotAllowed, err, log.WarnLevel)
}

func InternalServerErrorHandler(responseWriter http.ResponseWriter, request *http.Request, err error) {
	handleError(responseWriter, request, http.StatusInternalServerError, err, log.ErrorLevel)
}

func handleError(response http.ResponseWriter, request *http.Request, errorCode int, err error, level log.Level) {

	logEntry := log.WithFields(log.Fields{
		"Remote Adress":  request.RemoteAddr,
		"Request Uri":    request.RequestURI,
		"Request Method": request.Method,
	})

	logEntry.Log(level, err.Error())

	response.WriteHeader(errorCode)
	err = json.NewEncoder(response).Encode(struct {
		error string
	}{
		error: err.Error(),
	})

}
