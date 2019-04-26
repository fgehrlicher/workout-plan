package main

import (
	"fmt"
	"os"

	"workout-plan/config"
	"workout-plan/plan"
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
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
