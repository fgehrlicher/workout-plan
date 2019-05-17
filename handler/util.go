package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"workout-plan/auth"

	"workout-plan/config"
	"workout-plan/db"
	"workout-plan/plan"
	"workout-plan/plan-pointer"
)

const ConfigCtxKey = "usergrant"

func NewPlanPointerRepository() (*plan_pointer.PlanPointerRepository, error) {
	conf, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	database, err := db.GetDatabase(
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.User,
		conf.Database.Password,
		conf.Database.Database,
		time.Duration(conf.Database.Timeout.Startup)*time.Second,
	)

	if err != nil {
		return nil, err
	}

	planPointerRepository := plan_pointer.NewPlanPointerRepository(
		database,
		time.Duration(conf.Database.Timeout.Request)*time.Second,
	)

	return planPointerRepository, nil
}

func hasPlanEnded(pointer plan_pointer.PlanPointer, userPlan *plan.Plan) bool {
	unitKey := pointer.Position.Unit
	exerciseKey := pointer.Position.Exercise

	return unitKey > len(userPlan.Units) ||
		(unitKey == len(userPlan.Units) && exerciseKey >= len(userPlan.Units[unitKey-1].Exercises))
}

func hasPlanUnitsLeft(pointer plan_pointer.PlanPointer, userPlan *plan.Plan) bool {
	return len(userPlan.Units) > pointer.Position.Unit
}

func HeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(responseWriter, request)
	})
}

func ConfigMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		conf, err := config.GetConfig()
		if err != nil {
			internalServerErrorHandler(
				responseWriter,
				request,
				err,
			)
			return
		}

		ctx := context.WithValue(request.Context(), ConfigCtxKey, conf)
		next.ServeHTTP(responseWriter, request.WithContext(ctx))
	})
}

func GetUserGrant(request *http.Request) (*auth.Grant, error) {
	rawUserGrant := request.Context().Value(UserGrantCtxKey)
	userGrant, ok := rawUserGrant.(*auth.Grant)
	if !ok {
		return nil, fmt.Errorf("invalid user grant")
	}
	return userGrant, nil
}

func GetConfig(request *http.Request) (*config.Config, error) {
	rawConfig := request.Context().Value(ConfigCtxKey)
	conf, ok := rawConfig.(*config.Config)
	if !ok {
		return nil, fmt.Errorf("invalid ctx config")
	}
	return conf, nil
}
