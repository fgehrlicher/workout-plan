package user

import (
	"workout-plan/user/plan-pointer"
)

type User struct {
	ID              string
	AuthorizedPlans []string
	ActivePlans     []*plan_pointer.PlanPointer
}
