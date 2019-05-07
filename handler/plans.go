package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"workout-plan/plan"
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
