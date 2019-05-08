package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"workout-plan/config"
	"workout-plan/db"
	"workout-plan/plan"
	"workout-plan/plan-pointer"
)

func GetAllPlans(response http.ResponseWriter, request *http.Request) {
	plans := plan.GetPlansInstance()

	var sanitizedPlans []plan.Plan
	latestPlans, err := plans.GetAllLatest()

	for _, rawPlan := range latestPlans {
		sanitizedPlans = append(sanitizedPlans, rawPlan.GetSanitizedCopy())
	}

	err = json.NewEncoder(response).Encode(sanitizedPlans)
	if err != nil {
		InternalServerErrorHandler(response, request, err)
	}
}

func GetPlan(response http.ResponseWriter, request *http.Request) {
	plans := plan.GetPlansInstance()
	vars := mux.Vars(request)
	planId := vars["planId"]

	latestPlan, err := plans.GetLatest(planId)
	if err != nil {
		NotFoundErrorHandler(response, request, err)
		return
	}

	err = json.NewEncoder(response).Encode(latestPlan.GetSanitizedCopy())
	if err != nil {
		InternalServerErrorHandler(response, request, err)
		return
	}
}

func GetActivePlans(response http.ResponseWriter, request *http.Request) {
	userId, ok := mux.Vars(request)["user"]
	if !ok {
		BadRequestErrorHandler(response, request, errors.New("`user` parameter is required for this endpoint"))
		return
	}

	plans := plan.GetPlansInstance()
	conf, _ := config.GetConfig()
	database, err := db.GetDatabase(
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.User,
		conf.Database.Password,
		conf.Database.Database,
	)

	if err != nil {
		InternalServerErrorHandler(response, request, err)
		return
	}

	planPointerRepository := plan_pointer.NewPlanPointerRepository(
		database,
		time.Duration(conf.Database.Timeout.Request)*time.Second,
	)

	planPointers, err := planPointerRepository.GetAll(userId)
	if err != nil {
		InternalServerErrorHandler(response, request, err)
		return
	}

	var returnPlans []plan.Plan

	for _, planPointer := range planPointers {
		existingPlan, err := plans.Get(planPointer.PlanId, planPointer.PlanVersion)
		if err != nil {
			InternalServerErrorHandler(response, request, err)
			return
		}
		returnPlans = append(returnPlans, existingPlan.GetSanitizedCopy())
	}

	err = json.NewEncoder(response).Encode(returnPlans)
	if err != nil {
		InternalServerErrorHandler(response, request, err)
	}
}
