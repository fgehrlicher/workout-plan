package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"workout-plan/plan"
	"workout-plan/plan-pointer"
)

func GetCurrentUnit(response http.ResponseWriter, request *http.Request) {
	queryParameter := request.URL.Query()
	userId := queryParameter.Get("user")
	planId := mux.Vars(request)["planId"]

	plans := plan.GetPlansInstance()
	planPointerRepository, err := NewPlanPointerRepository()
	if err != nil {
		internalServerErrorHandler(response, request, err)
		return
	}

	planPointer, err := planPointerRepository.GetByPlan(userId, planId)
	if err != nil {
		if err == plan_pointer.NoPlanFoundError {
			notFoundErrorHandler(response, request, err)
		} else {
			internalServerErrorHandler(response, request, err)
		}
		return
	}

	userPlan, err := plans.Get(planPointer.PlanId, planPointer.PlanVersion)
	if err != nil {
		internalServerErrorHandler(response, request, err)
	}

	if hasPlanEnded(planPointer, userPlan) {
		err := planPointerRepository.Delete(planPointer)
		if err != nil {
			internalServerErrorHandler(response, request, err)
			return
		}

		notFoundErrorHandler(
			response,
			request,
			errors.New("the requested plan has already been finished"),
		)
		return
	}

	currentUnit := userPlan.Units[planPointer.Position.UnitKey]
	err = json.NewEncoder(response).Encode(currentUnit)
	if err != nil {
		internalServerErrorHandler(response, request, err)
	}
}
