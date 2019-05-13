package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"workout-plan/plan"
	"workout-plan/plan-pointer"
	"workout-plan/template"
)

const (
	PlanFinishedState = "PLAN_FINISHED"
	PlanOkState       = "OK"
)

func GetCurrentUnit(response http.ResponseWriter, request *http.Request) {
	userId := request.URL.Query().Get(UserQuerySegment)
	planId := mux.Vars(request)[PlanIdQuerySegment]

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

	currentUnit := userPlan.Units[planPointer.Position.Unit-1]

	err = template.EvaluateTemplate(&currentUnit, planPointer.Data)
	if err != nil {
		internalServerErrorHandler(response, request, err)
	}

	err = json.NewEncoder(response).Encode(currentUnit)
	if err != nil {
		internalServerErrorHandler(response, request, err)
	}
}

func FinishCurrentUnit(response http.ResponseWriter, request *http.Request) {
	queryParameter := request.URL.Query()
	userId := queryParameter.Get(UserQuerySegment)
	planId := mux.Vars(request)[PlanIdQuerySegment]

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

	currentUnit := userPlan.Units[planPointer.Position.Unit-1]
	requiredVariables := currentUnit.GetRequiredVariables()

	for _, requiredVariable := range requiredVariables {
		variableInForm := request.FormValue(requiredVariable)
		if variableInForm == "" {
			badRequestErrorHandler(
				response, request,
				errors.New(
					fmt.Sprintf("the variable `%v` must be sent when finishing this unit", requiredVariable),
				),
			)
			return
		}
		intValue, err := strconv.Atoi(variableInForm)
		if err != nil {
			badRequestErrorHandler(response, request, err)
		}

		planPointer.Data[requiredVariable] = intValue
	}

	if hasPlanUnitsLeft(planPointer, userPlan) {
		planPointer.Position.Unit ++
		err = planPointerRepository.Update(planPointer)
		if err != nil {
			internalServerErrorHandler(response, request, err)
			return
		}

		err = WriteStateMessage(response, "unit finished", PlanOkState)
		if err != nil {
			internalServerErrorHandler(response, request, err)
			return
		}
	} else {
		err := planPointerRepository.Delete(planPointer)
		if err != nil {
			internalServerErrorHandler(response, request, err)
			return
		}
		err = WriteStateMessage(response, "plan finished", PlanFinishedState)
		if err != nil {
			internalServerErrorHandler(response, request, err)
			return
		}
	}
}
