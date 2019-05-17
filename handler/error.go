package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	"workout-plan/auth"
)

func NotFound(responseWriter http.ResponseWriter, request *http.Request) {
	notFoundErrorHandler(
		responseWriter,
		request,
		errors.New(
			http.StatusText(http.StatusNotFound),
		),
	)
}

func MethodNotAllowed(responseWriter http.ResponseWriter, request *http.Request) {
	methodNotAllowedErrorHandler(
		responseWriter,
		request,
		errors.New(
			http.StatusText(http.StatusMethodNotAllowed),
		),
	)
}

func notFoundErrorHandler(responseWriter http.ResponseWriter, request *http.Request, err error) {
	handleError(responseWriter, request, http.StatusNotFound, err, log.WarnLevel)
}

func badRequestErrorHandler(responseWriter http.ResponseWriter, request *http.Request, err error) {
	handleError(responseWriter, request, http.StatusBadRequest, err, log.WarnLevel)
}

func forbiddenErrorHandler(responseWriter http.ResponseWriter, request *http.Request, err error) {
	conf, configErr := GetConfig(request)
	if configErr != nil {
		internalServerErrorHandler(
			responseWriter,
			request,
			configErr,
		)
		return
	}
	responseWriter.Header().Set(
		auth.GetTokenAuthenticateHeader(conf.Auth.Token),
	)
	handleError(responseWriter, request, http.StatusForbidden, err, log.WarnLevel)
}

func methodNotAllowedErrorHandler(responseWriter http.ResponseWriter, request *http.Request, err error) {
	handleError(responseWriter, request, http.StatusMethodNotAllowed, err, log.WarnLevel)
}

func internalServerErrorHandler(responseWriter http.ResponseWriter, request *http.Request, err error) {
	handleError(responseWriter, request, http.StatusInternalServerError, err, log.ErrorLevel)
}

func handleError(response http.ResponseWriter, request *http.Request, errorCode int, err error, level log.Level) {
	logEntry := log.WithFields(log.Fields{
		"remote_address": request.RemoteAddr,
		"uri":            request.RequestURI,
		"http_method":    request.Method,
		"status_code":    errorCode,
	})

	logEntry.Log(level, err.Error())

	response.WriteHeader(errorCode)
	err = json.NewEncoder(response).Encode(struct {
		Error      string `json:"error"`
		StatusCode string `json:"status_code"`
	}{
		Error:      err.Error(),
		StatusCode: strconv.Itoa(errorCode),
	})

}
