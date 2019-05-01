package user

import (
	"workout-plan/plan"
)

type PlanPointer struct {
	PlanId   string
	Position *plan.Exercise
}
