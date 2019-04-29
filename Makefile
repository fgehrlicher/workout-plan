build-all: build-workout-plan-server build-workout-plan-validator

build-workout-plan-server:
	go build -o bin/server cmd/server/server.go

build-workout-plan-validator:
	go build -o bin/validator cmd/validator/validator.go
