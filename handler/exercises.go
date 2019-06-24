package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"workout-plan/plan"
)

const ExerciseIdQuerySegment = "exerciseId"

func GetAllExercises(response http.ResponseWriter, request *http.Request) {
	err := json.NewEncoder(response).Encode(
		plan.GetExerciseDefinitionsInstance().GetAll(),
	)
	if err != nil {
		internalServerErrorHandler(response, request, err)
	}
}

func GetExercise(response http.ResponseWriter, request *http.Request) {
	exerciseDefinitions := plan.GetExerciseDefinitionsInstance()
	exerciseId := mux.Vars(request)[ExerciseIdQuerySegment]

	exerciseDefinition, err := exerciseDefinitions.Get(exerciseId)
	if err != nil {
		notFoundErrorHandler(response, request, err)
		return
	}

	err = json.NewEncoder(response).Encode(exerciseDefinition)
	if err != nil {
		internalServerErrorHandler(response, request, err)
	}
}
