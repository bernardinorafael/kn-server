include .env

.DEFAULT_GOAL := run

vet:
	@go vet ./...
.PHONY: vet

pg:
	@docker exec -it postgres-kn psql -U $(DB_USER) -d $(DB_NAME)
.PHONY: pg

run:
	@go run cmd/server/main.go
.PHONY: run

test:
	@go test ./...
.PHONY: test
