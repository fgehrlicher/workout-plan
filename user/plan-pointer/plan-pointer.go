package plan_pointer

import (
	"workout-plan/plan"
)

var plans *plan.Plans

func initializePlans() {
	if plans == nil {
		plans = plan.GetPlansInstance()
	}
}

func CreatePlanPointer(planId string) (*PlanPointer, error) {
	initializePlans()

	retrievedPlan, err := plans.Get(planId)
	if err != nil {
		return nil, err
	}

	planPointer := &PlanPointer{
		PlanId: planId,
	}

	planPointer.Position.Unit = &retrievedPlan.Units[0]
	planPointer.Position.ExercisesKey = 0

	return planPointer, nil
}

type PlanPointer struct {
	PlanId   string
	Position struct {
		Unit         *plan.Unit
		ExercisesKey int
	}
}
