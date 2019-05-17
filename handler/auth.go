package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"workout-plan/auth"
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

		conf, err := GetConfig(request)
		if err != nil {
			internalServerErrorHandler(
				responseWriter,
				request,
				err,
			)
			return
		}

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
			return
		}

		planId := mux.Vars(request)[PlanIdQuerySegment]
		if planId != "" && !userGrant.IsAuthorizedForPlan(planId) {
			forbiddenErrorHandler(responseWriter, request, fmt.Errorf("not authorized for '%v'", planId))
			return
		}

		ctx := context.WithValue(request.Context(), UserGrantCtxKey, userGrant)
		next.ServeHTTP(responseWriter, request.WithContext(ctx))
	})
}
