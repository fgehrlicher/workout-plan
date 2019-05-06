package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

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

	router := mux.NewRouter()
	server := &http.Server{
		Addr:         "127.0.0.1:8080",
		WriteTimeout: time.Second * 60,
		ReadTimeout:  time.Second * 360,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	err = server.ListenAndServe()
	if err == http.ErrServerClosed {
		os.Exit(0)
	} else {
		log.Fatal(err.Error())
	}
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
