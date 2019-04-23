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
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	plans.InitializePlans(conf.Plans.Directory)

	plan := plans.GetInstance()

	fmt.Printf("%v", plan)
}
