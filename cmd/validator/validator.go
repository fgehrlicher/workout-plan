package main

import (
	"fmt"
	"os"

	"workout-plan/config"
	"workout-plan/plans"

	log "github.com/sirupsen/logrus"
)

const ConfigFilePath = "./config.yml"

func main() {

	conf, err := config.LoadConfig(ConfigFilePath)
	handleError(err)

	err = plans.InitializePlans(conf.Plans.Directory)
	handleError(err)

	log.Info("Configuration Valid!")
}

func handleError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}