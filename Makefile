.PHONY: run-server run-client migrate migrate-down migrate-status

lint:
	golangci-lint run ./...

run-server:
	go run ./cmd/server

run-client:
	go run ./cmd/client


DB_URL ?= $(shell grep DB_URL .env | cut -d '=' -f2-)

migrate:
	goose -dir ./cmd/server/migrations postgres "$(DB_URL)" up

migrate-down:
	goose -dir ./cmd/server/migrations postgres "$(DB_URL)" down

migrate-status:
	goose -dir ./cmd/server/migrations postgres "$(DB_URL)" status

