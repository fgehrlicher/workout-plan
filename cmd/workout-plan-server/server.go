package main

import (
	"fmt"
	"os"
	"workout-plan/config"
	"workout-plan/plans"
)

const ConfigFilePath = "./config.yml"

func main() {

	conf, err := config.LoadConfig(ConfigFilePath)
	handleError(err)

	err = plans.InitializePlans(conf.Plans.Directory)
	handleError(err)

	plan := plans.GetInstance()

	fmt.Printf("%v", plan)
}

func handleError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}