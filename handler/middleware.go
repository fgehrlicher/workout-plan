package handler

import (
	"errors"
	"fmt"
	"net/http"
)

const UserQuerySegment = "user"

func UserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		userId := request.URL.Query().Get(UserQuerySegment)
		if userId == "" {
			badRequestErrorHandler(
				responseWriter,
				request,
				errors.New(fmt.Sprintf("`%v` parameter is required ", UserQuerySegment)),
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
