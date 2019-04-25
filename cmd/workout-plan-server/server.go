package main

import (
	"fmt"
	"os"

	"workout-plan/config"
	"workout-plan/exercise-definitions"
	"workout-plan/plans"
)

const ConfigFilePath = "./config.yml"

func main() {

	conf, err := config.LoadConfig(ConfigFilePath)
	handleError(err)

	err = exercise_definitions.InitializeExerciseDefinitions(conf.Plans.ExerciseDefinition)
	handleError(err)

	err = plans.InitializePlans(conf.Plans.Directory)
	handleError(err)

	plan := plans.GetInstance()
	exerciseDefinition := exercise_definitions.GetInstance()

	fmt.Printf("%v \n %v", plan, exerciseDefinition)
}

func handleError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
