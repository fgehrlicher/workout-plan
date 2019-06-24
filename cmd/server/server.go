package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"workout-plan/config"
	"workout-plan/db"
	"workout-plan/handler"
	"workout-plan/plan"
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

	database, err := db.GetDatabase(
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.User,
		conf.Database.Password,
		conf.Database.Database,
		time.Duration(conf.Database.Timeout.Startup)*time.Second,
	)
	handleError(err)

	planPointerRepository := db.NewPlanPointerRepository(
		database,
		time.Duration(conf.Database.Timeout.Request)*time.Second,
	)

	err = planPointerRepository.InitIndices()
	handleError(err)

	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(handler.NotFound)
	router.MethodNotAllowedHandler = http.HandlerFunc(handler.MethodNotAllowed)
	router.Use(handler.ConfigMiddleware)
	router.Use(handler.AuthMiddleware)
	router.Use(handler.HeaderMiddleware)

	router.HandleFunc(
		"/plans/",
		handler.GetAllPlans,
	).Methods("GET")
	router.HandleFunc(
		"/plans/active/",
		handler.GetActivePlans,
	).Methods("GET")
	router.HandleFunc(
		fmt.Sprintf("/plans/{%v}/", handler.PlanIdQuerySegment),
		handler.GetPlan,
	).Methods("GET")
	router.HandleFunc(
		fmt.Sprintf("/plans/{%v}/start/", handler.PlanIdQuerySegment),
		handler.StartPlan,
	).Methods("POST")
	router.HandleFunc(
		fmt.Sprintf("/plans/{%v}/stop/", handler.PlanIdQuerySegment),
		handler.StopPlan,
	).Methods("POST")
	router.HandleFunc(
		fmt.Sprintf("/plans/{%v}/stats/", handler.PlanIdQuerySegment),
		handler.GetStats,
	).Methods("GET")
	router.HandleFunc(
		fmt.Sprintf("/plans/{%v}/units/current/", handler.PlanIdQuerySegment),
		handler.GetCurrentUnit,
	).Methods("GET")
	router.HandleFunc(
		fmt.Sprintf("/plans/{%v}/units/current/finish/", handler.PlanIdQuerySegment),
		handler.FinishCurrentUnit,
	).Methods("POST")
	router.HandleFunc(
		fmt.Sprintf("/plans/{%v}/units/{%v}/", handler.PlanIdQuerySegment, handler.UnitIdQuerySegment),
		handler.GetUnit,
	).Methods("GET")
	router.HandleFunc(
		fmt.Sprintf("/plans/{%v}/units/{%v}/finish/", handler.PlanIdQuerySegment, handler.UnitIdQuerySegment),
		handler.FinishUnit,
	).Methods("POST")
	router.HandleFunc(
		"/exercises/",
		handler.GetAllExercises,
	).Methods("GET")
	router.HandleFunc(
		fmt.Sprintf("/exercises/{%v}/", handler.ExerciseIdQuerySegment),
		handler.GetExercise,
	).Methods("GET")

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
