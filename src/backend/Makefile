# Plate Pilot Makefile

# Variables
# Include environment variables from .env file
include .env
export

APP_NAME=plate-pilot
BUILD_DIR=bin
DOCKER_COMPOSE=docker compose
RECIPE_DB_DIR=infrastructure/recipe/db
# Help
help:
	@echo "Plate Pilot - Available Commands:"
	@echo "==============================="
	@echo "make dev          - Start development environment"
	@echo "make build        - Build the application"
	@echo "make test         - Run tests"
	@echo "make docker-up    - Start Docker containers"
	@echo "make docker-down  - Stop Docker containers"
	@echo "make migrate-up   - Run database migrations"
	@echo "make migrate-down - Revert database migrations"
	@echo "make clean        - Clean built files"

# Development
dev: docker-up migrate-up

# Docker commands
docker-up:
	$(DOCKER_COMPOSE) -f docker-compose.yml -f docker-compose.backend-only.yml  up --watch

docker-down:
	$(DOCKER_COMPOSE) down

docker-logs:
	$(DOCKER_COMPOSE) logs -f

migrate-create:
	migrate create -ext sql -dir ${RECIPE_DB_DIR}/migrations -seq ${NAME}

# Database commands
migrate-up:
	@echo "Running database migrations..."
	migrate -path ${RECIPE_DB_DIR}/migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@localhost:5432/${DB_NAME}?sslmode=disable" up

migrate-down:
	@echo "Reverting database migrations..."
	migrate -path ${RECIPE_DB_DIR}/migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@localhost:5432/${DB_NAME}?sslmode=disable" down

# Build commands
build:
	@echo "Building $(APP_NAME)..."
	go build -o $(BUILD_DIR)/$(APP_NAME) cmd/api/main.go


# Testing and code quality
test:
	go test ./... -v

fmt:
	go fmt ./...

lint:
	golangci-lint run

# Cleanup
clean:
	rm -rf $(BUILD_DIR)



.DEFAULT_GOAL := help