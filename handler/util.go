package handler

import (
	"time"

	"workout-plan/config"
	"workout-plan/db"
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

type ReturnMessage struct {
	Message string `json:"message"`
}