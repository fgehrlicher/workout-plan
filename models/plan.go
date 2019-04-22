package models

import (
	log "github.com/sirupsen/logrus"
)

type Plan struct {
	ID      string `yaml:"id" json:"id"`
	Name    string `yaml:"name" json:"name"`
	Version string `yaml:"version" json:"version"`
	Units   []Unit `yaml:"units" json:"units"`
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
			logEntry.Error("Plan id - version combination already exists")
			return
		}
	}

	plans.underlyingSlice = append(plans.underlyingSlice, &plan)
	logEntry.Info("Plan added")
}
