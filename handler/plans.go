package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"workout-plan/auth"
	"workout-plan/plan"
	"workout-plan/plan-pointer"
)

const PlanIdQuerySegment = "planId"

func GetAllPlans(response http.ResponseWriter, request *http.Request) {
	plans := plan.GetPlansInstance()

	var sanitizedPlans []plan.Plan
	latestPlans, err := plans.GetAllLatest()

	for _, rawPlan := range latestPlans {
		sanitizedPlans = append(sanitizedPlans, rawPlan.GetSanitizedCopy())
	}

	err = json.NewEncoder(response).Encode(sanitizedPlans)
	if err != nil {
		internalServerErrorHandler(response, request, err)
	}
}

func GetPlan(response http.ResponseWriter, request *http.Request) {
	plans := plan.GetPlansInstance()
	planId := mux.Vars(request)[PlanIdQuerySegment]

	latestPlan, err := plans.GetLatest(planId)
	if err != nil {
		notFoundErrorHandler(response, request, err)
		return
	}

	err = json.NewEncoder(response).Encode(latestPlan.GetSanitizedCopy())
	if err != nil {
		internalServerErrorHandler(response, request, err)
		return
	}
}

func GetActivePlans(response http.ResponseWriter, request *http.Request) {
	rawUserGrant := request.Context().Value(UserGrantCtxKey)
	userGrant, ok := rawUserGrant.(*auth.UserAccessClaim)
	if !ok {
		internalServerErrorHandler(response, request, nil)
	}
	userId := userGrant.UserName

	plans := plan.GetPlansInstance()
	planPointerRepository, err := NewPlanPointerRepository()

	if err != nil {
		internalServerErrorHandler(response, request, err)
		return
	}

	planPointers, err := planPointerRepository.GetAll(userId)
	if err != nil {
		internalServerErrorHandler(response, request, err)
		return
	}

	returnPlans := make([]plan.Plan, 0)

	for _, planPointer := range planPointers {
		existingPlan, err := plans.Get(planPointer.PlanId, planPointer.PlanVersion)
		if err != nil {
			internalServerErrorHandler(response, request, err)
			return
		}
		returnPlans = append(returnPlans, existingPlan.GetSanitizedCopy())
	}

	err = json.NewEncoder(response).Encode(returnPlans)
	if err != nil {
		internalServerErrorHandler(response, request, err)
	}
}

func StartPlan(response http.ResponseWriter, request *http.Request) {
	rawUserGrant := request.Context().Value(UserGrantCtxKey)
	userGrant, ok := rawUserGrant.(*auth.UserAccessClaim)
	if !ok {
		internalServerErrorHandler(response, request, nil)
	}
	userId := userGrant.UserName
	plans := plan.GetPlansInstance()
	planId := mux.Vars(request)[PlanIdQuerySegment]

	requestedPlan, err := plans.GetLatest(planId)
	if err != nil {
		notFoundErrorHandler(response, request, err)
		return
	}

	planPointerRepository, err := NewPlanPointerRepository()
	if err != nil {
		internalServerErrorHandler(response, request, err)
		return
	}

	planPointers, err := planPointerRepository.GetAll(userId)
	if err != nil {
		internalServerErrorHandler(response, request, err)
		return
	}

	if len(planPointers) > 0 {
		for _, planPointer := range planPointers {
			if planPointer.PlanId == planId {
				badRequestErrorHandler(
					response,
					request,
					errors.New("the plan has already been started by this user"),
				)
				return
			}
		}
	}

	userPlanPointer := plan_pointer.CreatePlanPointer(requestedPlan, userId)
	_, err = planPointerRepository.Insert(userPlanPointer)
	if err != nil {
		internalServerErrorHandler(response, request, err)
		return
	}

	err = WriteMessage(response, "plan started")
	if err != nil {
		internalServerErrorHandler(response, request, err)
	}
}

func StopPlan(response http.ResponseWriter, request *http.Request) {
	rawUserGrant := request.Context().Value(UserGrantCtxKey)
	userGrant, ok := rawUserGrant.(*auth.UserAccessClaim)
	if !ok {
		internalServerErrorHandler(response, request, nil)
	}
	userId := userGrant.UserName
	planId := mux.Vars(request)[PlanIdQuerySegment]

	planPointerRepository, err := NewPlanPointerRepository()
	if err != nil {
		internalServerErrorHandler(response, request, err)
		return
	}

	planPointers, err := planPointerRepository.GetAll(userId)
	if err != nil {
		internalServerErrorHandler(response, request, err)
		return
	}

	for _, planPointer := range planPointers {
		if planPointer.PlanId == planId {
			err = planPointerRepository.Delete(planPointer)
			if err != nil {
				internalServerErrorHandler(response, request, err)
				return
			}

			err = WriteMessage(response, "plan deleted")
			if err != nil {
				internalServerErrorHandler(response, request, err)
			}
			return
		}
	}
	notFoundErrorHandler(response, request, errors.New("no active plan with that name found"))
}
