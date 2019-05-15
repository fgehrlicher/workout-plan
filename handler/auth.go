package handler

import (
	"fmt"
	"net/http"

	"workout-plan/auth"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		authorizationHeader := request.Header.Get(auth.AuthorizationHeader)
		if authorizationHeader == "" {
			badRequestErrorHandler(
				responseWriter,
				request,
				fmt.Errorf("authorization header `%v` missing", auth.AuthorizationHeader),
			)
			return
		}

		err := auth.ParseAuth(authorizationHeader)
		if err != nil {
			_, ok := err.(*auth.BadRequestError)
			if ok {
				badRequestErrorHandler(
					responseWriter,
					request,
					err,
				)
			} else {
				forbiddenErrorHandler(
					responseWriter,
					request,
					err,
				)
			}
		}

		next.ServeHTTP(responseWriter, request)
	})
}
