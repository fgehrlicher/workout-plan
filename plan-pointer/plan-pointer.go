package plan_pointer

import (
	"workout-plan/plan"
)

func CreatePlanPointer(plan *plan.Plan, userId string) (PlanPointer, error) {
	planPointer := PlanPointer{
		PlanId: plan.ID,
		PlanVersion: plan.Version,
		UserId: userId,
	}

	planPointer.Position.Unit = &plan.Units[0]
	planPointer.Position.ExerciseKey = 0

	return planPointer, nil
}

type PlanPointer struct {
	PlanId      string
	PlanVersion string
	UserId      string
	Position    struct {
		Unit        *plan.Unit
		ExerciseKey int
	}
}
