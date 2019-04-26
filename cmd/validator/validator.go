package main

import (
	"workout-plan/config"
	"workout-plan/plan"

	log "github.com/sirupsen/logrus"
)

const ConfigFilePath = "./config.yml"

func main() {

	conf, err := config.LoadConfig(ConfigFilePath)
	handleError(err)

	err = plan.InitializeExerciseDefinitions(conf.Plans.ExerciseDefinition)
	handleError(err)

	err = plan.InitializePlans(conf.Plans.Directory)
	handleError(err)

	log.Info("Configuration Valid!")
}

func handleError(err error) {
	if err != nil {
		log.Error(err.Error())
		log.Exit(1)
	}
}