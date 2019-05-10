package plan_pointer

import (
	"workout-plan/plan"
)

func CreatePlanPointer(plan *plan.Plan, userId string) PlanPointer {
	planPointer := PlanPointer{
		PlanId:      plan.ID,
		PlanVersion: plan.Version,
		UserId:      userId,
	}

	planPointer.Position.Unit = 1
	planPointer.Position.Exercise = 1

	return planPointer
}

type PlanPointer struct {
	PlanId      string `bson:"plan_id"`
	PlanVersion string `bson:"plan_version"`
	UserId      string `bson:"user_id"`
	Position    struct {
		Unit     int `bson:"unit"`
		Exercise int `bson:"exercise"`
	} `bson:"position"`
	Data map[string]string `bson:"data"`
}
