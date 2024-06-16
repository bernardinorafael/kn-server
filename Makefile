.PHONY: run build test clean pg

include .env

pg:
	@docker exec -it db_kn_server psql -U $(DB_USER) -d $(DB_NAME)
run:
	@go run cmd/api/main.go
build:
	@go build -o $(APP_NAME) main.go
test:
	@go test ./...
