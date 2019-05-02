package plan_pointer

import (
	"workout-plan/plan"
)

type PlanPointer struct {
	PlanId   string
	Position struct {
		Unit         *plan.Unit
		ExercisesKey int
	}
}
