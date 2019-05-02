package user

import (
	"errors"
	"fmt"

	"workout-plan/user/plan-pointer"
)

type User struct {
	ID              string
	AuthorizedPlans []string
	ActivePlans     []*plan_pointer.PlanPointer
}

func (user *User) GetPlanPointer(planId string) (*plan_pointer.PlanPointer, error) {
	for _, planPointer := range user.ActivePlans {
		if planPointer.PlanId == planId {
			return planPointer, nil
		}
	}

	return nil, errors.New(
		fmt.Sprintf(
			"no plan pointer with plan id `%v` found.",
			planId,
		),
	)
}

func (user *User) CreatePlanPointer(planId string) error {
	_, err := user.GetPlanPointer(planId)
	if err == nil {
		return errors.New(
			fmt.Sprintf(
				"a plan pointer with plan id `%v` already exists.",
				planId,
			),
		)
	}

	planPointer, err := plan_pointer.CreatePlanPointer(planId)
	if err != nil {
		return err
	}

	user.ActivePlans = append(user.ActivePlans, planPointer)
	return nil
}
