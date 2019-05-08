package main

import (
	"time"

	"workout-plan/config"
	"workout-plan/db"
	"workout-plan/plan"

	log "github.com/sirupsen/logrus"
)

func main() {

	conf, err := config.GetConfig()
	handleError(err)

	err = plan.InitializeExerciseDefinitions(conf.Plans.DefinitionsFile)
	handleError(err)

	err = plan.InitializePlans(conf.Plans.Directory)
	handleError(err)

	_, err = db.GetDatabase(
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.User,
		conf.Database.Password,
		conf.Database.Database,
		time.Duration(conf.Database.Timeout.Startup)*time.Second,
	)
	handleError(err)

	log.Info("Configuration Valid!")
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
