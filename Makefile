build-all: build-workout-plan-server

build-workout-plan-server:
	go build -o bin/workout-plan-server cmd/workout-plan-server/server.go