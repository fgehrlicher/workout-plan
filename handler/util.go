package handler

import (
	"time"

	"workout-plan/config"
	"workout-plan/db"
	"workout-plan/plan"
	"workout-plan/plan-pointer"
)

func NewPlanPointerRepository() (*plan_pointer.PlanPointerRepository, error) {
	conf, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	database, err := db.GetDatabase(
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.User,
		conf.Database.Password,
		conf.Database.Database,
		time.Duration(conf.Database.Timeout.Startup)*time.Second,
	)

	if err != nil {
		return nil, err
	}

	planPointerRepository := plan_pointer.NewPlanPointerRepository(
		database,
		time.Duration(conf.Database.Timeout.Request)*time.Second,
	)

	return planPointerRepository, nil
}

func hasPlanEnded(pointer plan_pointer.PlanPointer, userPlan *plan.Plan) bool {
	unitKey := pointer.Position.Unit
	exerciseKey := pointer.Position.Exercise

	return unitKey > len(userPlan.Units) ||
		(unitKey == len(userPlan.Units) && exerciseKey >= len(userPlan.Units[unitKey -1].Exercises))
}

func hasPlanUnitsLeft(pointer plan_pointer.PlanPointer, userPlan *plan.Plan) bool {
	return len(userPlan.Units) > pointer.Position.Unit
}
