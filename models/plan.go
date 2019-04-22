package models

import (
	log "github.com/sirupsen/logrus"
)

type Plan struct {
	ID      string `yaml:"id"`
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Units   []Unit `yaml:"units"`
}

type Plans struct {
	underlyingSlice []*Plan
}

func (plans *Plans) Add(plan Plan) {
	logEntry := log.WithFields(log.Fields{
		"Id":      plan.ID,
		"Version": plan.Version,
	})

	for _, existingPlan := range plans.underlyingSlice {
		if existingPlan.ID == plan.ID && existingPlan.Version == plan.Version {
			logEntry.Error("Plan already exists. skipping...")
			return
		}
	}

	plans.underlyingSlice = append(plans.underlyingSlice, &plan)
	logEntry.Info("Plan added")
}
