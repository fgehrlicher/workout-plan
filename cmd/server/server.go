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

func main() {

	conf, err := config.LoadConfig()
	handleError(err)

	err = plan.InitializeExerciseDefinitions(conf.Plans.ExerciseDefinition)
	handleError(err)

	err = plan.InitializePlans(conf.Plans.Directory)
	handleError(err)

	router := mux.NewRouter()
	server := &http.Server{
		Addr:         fmt.Sprintf("%v:%v", conf.Server.Ip, conf.Server.Port),
		WriteTimeout: time.Second * time.Duration(conf.Server.Timeout.Write),
		ReadTimeout:  time.Second * time.Duration(conf.Server.Timeout.Read),
		IdleTimeout:  time.Second * time.Duration(conf.Server.Timeout.Idle),
		Handler:      router,
	}

	go func() {
		osSignalChannel := make(chan os.Signal, 1)
		signal.Notify(osSignalChannel, os.Interrupt)
		signal.Notify(osSignalChannel, os.Kill)

		<-osSignalChannel

		fmt.Println("shutting down...")
		ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
		err := server.Shutdown(ctx)
		if err != nil {
			log.Fatal(err.Error())
		}
	}()

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
