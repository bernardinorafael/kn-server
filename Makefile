.PHONY: run build test clean

include .env

APP_NAME=kn.co

run:
	@go run ./cmd/main.go
build:
	@go build -o $(APP_NAME) main.go
test:
	@go test ./ ...
clean:
	@rm -f $(APP_NAME)
	@rm -rf ./docs
