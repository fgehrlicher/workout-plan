package handler

import (
	"encoding/json"
	"net/http"

	"workout-plan/plan"
)

func GetAllExercises(response http.ResponseWriter, request *http.Request) {
	err := json.NewEncoder(response).Encode(
		plan.GetExerciseDefinitionsInstance().GetAll(),
	)
	if err != nil {
		internalServerErrorHandler(response, request, err)
	}
}
