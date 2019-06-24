package stats

import (
	"time"
)

type Stats struct {
	TotalUnitCount     int `json:"total_unit_count"`
	UnitsDone          int `json:"units_done"`
	TotalExerciseCount int `json:"total_exercise_count"`
	ExercisesUsed      int `json:"exercises_used"`
	PlanStarted        time.Time `json:"plan_started"`
	LastWorkout        time.Time `json:"last_workout"`
	Variables map[string]string `json:"variables"`
}
