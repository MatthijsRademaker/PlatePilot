.PHONY: all build run test lint proto sqlc migrate-up migrate-down clean help
.PHONY: dev dev-recipe dev-mealplanner dev-bff
.PHONY: docker-up docker-down docker-build docker-run docker-run-all generate
.PHONY: test-e2e test-integration test-coverage

# Variables
GO := go
GOFLAGS := -v
BIN_DIR := bin
PROTO_DIR := api/proto
MIGRATIONS_DIR := migrations

# Database URLs (override with environment variables)
RECIPE_DB_URL ?= postgres://platepilot:platepilot@localhost:5432/recipedb?sslmode=disable
MEALPLANNER_DB_URL ?= postgres://platepilot:platepilot@localhost:5432/mealplannerdb?sslmode=disable

# Colors for output
GREEN := \033[0;32m
NC := \033[0m # No Color

## help: Show this help message
help:
	@echo "PlatePilot Go Backend - Available commands:"
	@echo ""
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'

## all: Build all services
all: build

## build: Build all service binaries
build:
	@echo "$(GREEN)Building services...$(NC)"
	@mkdir -p $(BIN_DIR)
	$(GO) build $(GOFLAGS) -o $(BIN_DIR)/recipe-api ./cmd/recipe-api
	$(GO) build $(GOFLAGS) -o $(BIN_DIR)/mealplanner-api ./cmd/mealplanner-api
	$(GO) build $(GOFLAGS) -o $(BIN_DIR)/mobile-bff ./cmd/mobile-bff
	@echo "$(GREEN)Build complete!$(NC)"

## build-recipe: Build recipe-api service
build-recipe:
	@mkdir -p $(BIN_DIR)
	$(GO) build $(GOFLAGS) -o $(BIN_DIR)/recipe-api ./cmd/recipe-api

## build-mealplanner: Build mealplanner-api service
build-mealplanner:
	@mkdir -p $(BIN_DIR)
	$(GO) build $(GOFLAGS) -o $(BIN_DIR)/mealplanner-api ./cmd/mealplanner-api

## build-bff: Build mobile-bff service
build-bff:
	@mkdir -p $(BIN_DIR)
	$(GO) build $(GOFLAGS) -o $(BIN_DIR)/mobile-bff ./cmd/mobile-bff

## run: Run all services (requires infrastructure)
run: docker-up migrate-up
	@echo "$(GREEN)Starting services...$(NC)"
	@$(MAKE) -j3 run-recipe run-mealplanner run-bff

run-recipe:
	$(GO) run ./cmd/recipe-api

run-mealplanner:
	$(GO) run ./cmd/mealplanner-api

run-bff:
	$(GO) run ./cmd/mobile-bff

## dev: Run all services with hot reload
dev: docker-up migrate-up
	@echo "$(GREEN)Starting services with hot reload...$(NC)"
	@$(MAKE) -j3 dev-recipe dev-mealplanner dev-bff

## dev-recipe: Run recipe-api with hot reload
dev-recipe:
	air -c .air.recipe.toml

## dev-mealplanner: Run mealplanner-api with hot reload
dev-mealplanner:
	air -c .air.mealplanner.toml

## dev-bff: Run mobile-bff with hot reload
dev-bff:
	air -c .air.bff.toml

## test: Run all tests
test:
	@echo "$(GREEN)Running tests...$(NC)"
	$(GO) test -v -race ./...

## test-coverage: Run tests with coverage
test-coverage:
	@echo "$(GREEN)Running tests with coverage...$(NC)"
	$(GO) test -v -race -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report: coverage.html$(NC)"

## test-integration: Run integration tests
test-integration:
	@echo "$(GREEN)Running integration tests...$(NC)"
	$(GO) test -v -tags=integration ./...

## lint: Run linter
lint:
	@echo "$(GREEN)Running linter...$(NC)"
	golangci-lint run ./...

## lint-fix: Run linter and fix issues
lint-fix:
	@echo "$(GREEN)Running linter with fixes...$(NC)"
	golangci-lint run --fix ./...

## proto: Generate protobuf code
proto:
	@echo "$(GREEN)Generating protobuf code...$(NC)"
	./scripts/generate-proto.sh

## sqlc: Generate type-safe SQL code
sqlc:
	@echo "$(GREEN)Generating sqlc code...$(NC)"
	sqlc generate

## generate: Generate all code (proto + sqlc)
generate: proto sqlc

## migrate-up: Run all database migrations
migrate-up:
	@echo "$(GREEN)Running migrations...$(NC)"
	migrate -path $(MIGRATIONS_DIR)/recipe -database "$(RECIPE_DB_URL)" up
	migrate -path $(MIGRATIONS_DIR)/mealplanner -database "$(MEALPLANNER_DB_URL)" up

## migrate-down: Rollback last migration
migrate-down:
	@echo "$(GREEN)Rolling back migrations...$(NC)"
	migrate -path $(MIGRATIONS_DIR)/recipe -database "$(RECIPE_DB_URL)" down 1
	migrate -path $(MIGRATIONS_DIR)/mealplanner -database "$(MEALPLANNER_DB_URL)" down 1

## migrate-create: Create a new migration (usage: make migrate-create name=migration_name service=recipe)
migrate-create:
	@if [ -z "$(name)" ] || [ -z "$(service)" ]; then \
		echo "Usage: make migrate-create name=migration_name service=recipe|mealplanner"; \
		exit 1; \
	fi
	migrate create -ext sql -dir $(MIGRATIONS_DIR)/$(service) -seq $(name)

## docker-up: Start infrastructure containers
docker-up:
	@echo "$(GREEN)Starting infrastructure...$(NC)"
	docker-compose -f deployments/docker-compose.yml up -d postgres rabbitmq
	@echo "Waiting for services to be ready..."
	@sleep 5
	@echo "$(GREEN)Infrastructure ready!$(NC)"
	@echo "  PostgreSQL: localhost:5432"
	@echo "  RabbitMQ:   localhost:5672 (AMQP), localhost:15672 (Management UI)"

## docker-down: Stop infrastructure containers
docker-down:
	@echo "$(GREEN)Stopping infrastructure...$(NC)"
	docker-compose -f deployments/docker-compose.yml down

## docker-build: Build Docker images for all services
docker-build:
	@echo "$(GREEN)Building Docker images...$(NC)"
	docker build -f deployments/Dockerfile --target recipe-api -t platepilot/recipe-api:latest .
	docker build -f deployments/Dockerfile --target mealplanner-api -t platepilot/mealplanner-api:latest .
	docker build -f deployments/Dockerfile --target mobile-bff -t platepilot/mobile-bff:latest .

## docker-run: Run all services in Docker (infrastructure only)
docker-run:
	docker-compose -f deployments/docker-compose.yml up postgres rabbitmq

## docker-run-all: Build and run the complete stack
docker-run-all:
	@echo "$(GREEN)Building and running complete stack...$(NC)"
	docker-compose -f deployments/docker-compose.yml up --build

## docker-run-detached: Run complete stack in background
docker-run-detached:
	@echo "$(GREEN)Starting complete stack in background...$(NC)"
	docker-compose -f deployments/docker-compose.yml up --build -d
	@echo "$(GREEN)Stack started! Services:$(NC)"
	@echo "  BFF (REST):     http://localhost:8080"
	@echo "  Recipe API:     localhost:50051 (gRPC)"
	@echo "  MealPlanner:    localhost:50052 (gRPC)"
	@echo "  PostgreSQL:     localhost:5432"
	@echo "  RabbitMQ:       localhost:5672 (AMQP), http://localhost:15672 (UI)"

## clean: Clean build artifacts
clean:
	@echo "$(GREEN)Cleaning...$(NC)"
	rm -rf $(BIN_DIR)
	rm -rf tmp/
	rm -f coverage.out coverage.html

## deps: Download dependencies
deps:
	@echo "$(GREEN)Downloading dependencies...$(NC)"
	$(GO) mod download
	$(GO) mod tidy

## tools: Install development tools
tools:
	@echo "$(GREEN)Installing development tools...$(NC)"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/air-verse/air@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@echo "$(GREEN)Don't forget to install:$(NC)"
	@echo "  - protoc (Protocol Buffers compiler)"
	@echo "  - golang-migrate (database migrations)"

## verify: Verify project setup
verify:
	@echo "$(GREEN)Verifying project setup...$(NC)"
	@echo "Go version:"
	@go version
	@echo ""
	@echo "Checking tools..."
	@which golangci-lint > /dev/null 2>&1 && echo "✓ golangci-lint" || echo "✗ golangci-lint (run: make tools)"
	@which air > /dev/null 2>&1 && echo "✓ air" || echo "✗ air (run: make tools)"
	@which protoc > /dev/null 2>&1 && echo "✓ protoc" || echo "✗ protoc (install manually)"
	@which migrate > /dev/null 2>&1 && echo "✓ migrate" || echo "✗ migrate (install manually)"
	@which sqlc > /dev/null 2>&1 && echo "✓ sqlc" || echo "✗ sqlc (run: make tools)"
	@which docker > /dev/null 2>&1 && echo "✓ docker" || echo "✗ docker"
	@echo ""
	@echo "$(GREEN)Verification complete!$(NC)"

## test-e2e: Run end-to-end integration tests
test-e2e:
	@echo "$(GREEN)Running E2E tests...$(NC)"
	./scripts/e2e-test.sh

## seed: Seed the database with sample data (requires running services)
seed:
	@echo "$(GREEN)Seeding database...$(NC)"
	$(GO) run ./cmd/recipe-api -seed data/recipes.json

## docker-logs: Show logs from all services
docker-logs:
	docker-compose -f deployments/docker-compose.yml logs -f

## docker-ps: Show running containers
docker-ps:
	docker-compose -f deployments/docker-compose.yml ps
