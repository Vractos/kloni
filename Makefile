#!make
include .env

BINARY_NAME=kloni

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## test: run tests
test:
	@echo "Running tests..."
	@go test ./... -v | sed "/PASS/s//$$(printf "\033[32m✅ PASS\033[0m")/" | sed "/FAIL/s//$$(printf "\033[31m❌ FAIL\033[0m")/" | sed "/RUN/s//$$(printf "\033[33m🏃 RUN\033[0m")/"

## coverage: run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	@go test ./... -coverprofile=coverage.out -coverpkg=./... -covermode=atomic
	@go tool cover -html=coverage.out


## create_migration: create a new migration using migrate cli, e.g. make create_migration NAME=create_users_table
create_migration:
	@echo "Creating migration..."
	@$(eval NAME := $(shell read -p "Enter migration name: " && echo $$REPLY)) \
	migrate create -ext sql -dir ops/db/migrations -seq $(NAME)

## backup_db: backup database running in docker container
backup_db:
	@echo "Backing up database..."
	@docker exec -i postgres pg_dump -U ${POSTGRES_USER} ${POSTGRES_DB_NAME} > ./postgres-backup-`date +%Y%m%d%H%M%S`.sql

## migrate_up: run all migrations
migrate_up:
	@echo "Running migrations..."
	@migrate -path ops/db/migrations -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:5432/${POSTGRES_DB_NAME}?sslmode=disable" --verbose up