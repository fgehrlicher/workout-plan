package stats

import (
	"time"

	"workout-plan/plan"
)

type Stats struct {
	TotalUnitCount     int            `json:"total_unit_count"`
	UnitsDone          int            `json:"units_done"`
	TotalExerciseCount int            `json:"total_exercise_count"`
	ExercisesUsed      int            `json:"exercises_used"`
	PlanStarted        string         `json:"plan_started"`
	LastWorkout        string         `json:"last_workout"`
	Variables          map[string]int `json:"variables,omitempty"`
}

func RetrieveStats(plan plan.Plan, pointer plan.Pointer) (Stats, error) {
	currentUnit := pointer.Position.Unit - 1
	stats := Stats{
		TotalUnitCount: len(plan.Units),
		UnitsDone:      currentUnit,
		Variables:      pointer.Data,
		PlanStarted:    pointer.Started.Format(time.RFC1123),
	}

	if !pointer.Moved.IsZero() {
		stats.LastWorkout = pointer.Moved.Format(time.RFC1123)
	}

	usedExercises := getExercises(plan.Units[0:currentUnit])
	stats.ExercisesUsed = len(usedExercises)

	unusedExercises := getExercises(plan.Units[currentUnit:stats.TotalUnitCount])

	allExercises := mergeUniqueExercises(usedExercises, unusedExercises)
	stats.TotalExerciseCount = len(allExercises)

	return stats, nil
}

func getExercises(units []plan.Unit) map[string]bool {
	exerciseMap := make(map[string]bool, 0)

	for unitIndex := 0; unitIndex < len(units); unitIndex ++ {
		currentUnit := units[unitIndex]
		for _, exercise := range currentUnit.Exercises {
			exerciseId := exercise.Definition.Id
			if !exerciseMap[exerciseId] {
				exerciseMap[exerciseId] = true
			}
		}
	}

	return exerciseMap
}

func mergeUniqueExercises(exerciseMaps ...map[string]bool) map[string]bool {
	returnMap := make(map[string]bool, 0)
	for _, exerciseMap := range exerciseMaps {
		for exercise := range exerciseMap {
			if !returnMap[exercise] {
				returnMap[exercise] = true
			}
		}
	}

	return returnMap
}
