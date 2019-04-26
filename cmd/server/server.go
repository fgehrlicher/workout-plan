package main

import (
	"fmt"
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

	plans := plan.GetPlansInstance()
	exerciseDefinitions := plan.GetExerciseDefinitionsInstance()

	fmt.Printf("%v \n %v", plans, exerciseDefinitions)
}

func handleError(err error) {
	if err != nil {
		log.Error(err.Error())
		log.Exit(1)
	}
}
