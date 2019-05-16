build-all: build-workout-plan-server build-workout-plan-validator build-workout-plan-claim-generator

build-workout-plan-server:
	go build -o bin/server cmd/server/server.go

build-workout-plan-validator:
	go build -o bin/validator cmd/validator/validator.go

build-workout-plan-claim-generator:
	go build -o bin/claim-generator cmd/claim-generator/generator.go
