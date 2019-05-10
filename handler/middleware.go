package handler

import (
	"errors"
	"net/http"
)

func UserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		userId := request.URL.Query().Get("user")
		if userId == "" {
			badRequestErrorHandler(
				responseWriter,
				request,
				errors.New("`user` parameter is required "),
			)
			return
		}

		next.ServeHTTP(responseWriter, request)
	})
}


func HeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(responseWriter, request)
	})
}
