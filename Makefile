.PHONY: run build test clean

include .env

run:
	@go run main.go
build:
	@go build -o $(APP_NAME) main.go
test:
	@go test ./test -v
