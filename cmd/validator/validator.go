package main

import (
	"workout-plan/config"
	"workout-plan/db"
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

	_, err = db.GetDatabase(
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.User,
		conf.Database.Password,
		conf.Database.Database,
	)

	log.Info("Configuration Valid!")
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
