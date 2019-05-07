package handler

import (
	"encoding/json"
	"net/http"

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
