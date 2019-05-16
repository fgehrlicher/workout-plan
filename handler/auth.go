package handler

import (
	"context"
	"fmt"
	"net/http"

	"workout-plan/auth"
	"workout-plan/config"
)

const UserGrantCtxKey = "usergrant"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		authorizationHeader := request.Header.Get(auth.Header)
		if authorizationHeader == "" {
			badRequestErrorHandler(
				responseWriter,
				request,
				fmt.Errorf("authorization header `%v` missing", auth.Header),
			)
			return
		}

		conf, _ := config.GetConfig()

		userGrant, err := auth.ParseAuth(authorizationHeader, conf.Auth)
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

		ctx := context.WithValue(request.Context(), UserGrantCtxKey, userGrant)
		next.ServeHTTP(responseWriter, request.WithContext(ctx))
	})
}
