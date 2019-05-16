package main

import (
	"time"

	"workout-plan/config"
	"workout-plan/db"
	"workout-plan/plan"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func main() {
	conf, err := config.GetConfig()
	handleError(err)

	log.SetLevel(logrus.InfoLevel)

	err = plan.InitializeExerciseDefinitions(conf.Plans.DefinitionsFile, log)
	handleError(err)

	err = plan.InitializePlans(conf.Plans.Directory, log)
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
