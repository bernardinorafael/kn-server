.PHONY: run build test clean pg vet

include .env

vet:
	@go vet ./...
pg:
	@docker exec -it postgres-kn psql -U $(DB_USER) -d $(DB_NAME)
run:
	@go run cmd/server/main.go
build:
	@go build -o $(APP_NAME) main.go
test:
	@go test ./...
