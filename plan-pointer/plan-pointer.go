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

	planPointer.Position.UnitId = plan.Units[0].Id
	planPointer.Position.ExerciseKey = 0

	return planPointer, nil
}

type PlanPointer struct {
	PlanId      string `bson:"plan_id"`
	PlanVersion string `bson:"plan_version"`
	UserId      string `bson:"user_id"`
	Position    struct {
		UnitId      string `bson:"unit_id"`
		ExerciseKey int    `bson:"exercise_key"`
	} `bson:"position"`
}
