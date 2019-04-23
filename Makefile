build-all: build-workout-plan-server build-workout-plan-validator

build-workout-plan-server:
	go build -o bin/workout-plan-server cmd/workout-plan-server/server.go

build-workout-plan-validator:
	go build -o bin/workout-plan-validator cmd/workout-plan-validator/validator.go
