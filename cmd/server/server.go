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
	"workout-plan/db"
	"workout-plan/handler"
	"workout-plan/plan"
	"workout-plan/plan-pointer"
)

func main() {

	conf, err := config.GetConfig()
	handleError(err)

	err = plan.InitializeExerciseDefinitions(conf.Plans.ExerciseDefinition)
	handleError(err)

	err = plan.InitializePlans(conf.Plans.Directory)
	handleError(err)

	database, err := db.GetDatabase(
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.User,
		conf.Database.Password,
		conf.Database.Database,
	)
	handleError(err)

	planPointerRepository := plan_pointer.NewPlanPointerRepository(
		database,
		time.Duration(conf.Database.Timeout.Request)*time.Second,
	)

	err = planPointerRepository.InitIndices()
	handleError(err)

	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(handler.NotFound)
	router.MethodNotAllowedHandler = http.HandlerFunc(handler.MethodNotAllowed)

	router.HandleFunc("/plans/", handler.GetAllPlans).Methods("GET")
	router.HandleFunc("/plans/active/", handler.GetActivePlans).Methods("GET")
	router.HandleFunc("/plans/{planId}/", handler.GetPlan).Methods("GET")
	router.HandleFunc("/plans/{planId}/start", handler.StartPlan).Methods("POST")
	router.HandleFunc("/plans/{planId}/stop", handler.StopPlan).Methods("POST")

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

	log.Info(fmt.Sprintf("Server started at %v:%v", conf.Server.Ip, conf.Server.Port))
	log.Info("Startup finished!")
	fmt.Println()
	err = server.ListenAndServe()
	if err == http.ErrServerClosed {
		os.Exit(0)
	} else {
		handleError(err)
	}
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
