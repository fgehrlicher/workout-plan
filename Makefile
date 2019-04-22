build-all: build-workout-plan-server

build-workout-plan-server:
	go build -o bin/test cmd/workout-plan-server/server.go 