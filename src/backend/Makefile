.PHONY: all build run test lint proto sqlc migrate-up migrate-down clean help
.PHONY: dev dev-recipe dev-mealplanner dev-bff
.PHONY: docker-up docker-down docker-build docker-run docker-run-all generate
.PHONY: test-e2e test-integration test-coverage
.PHONY: verify-proto verify-generated verify-swagger swagger
.PHONY: ollama-start ollama-stop ollama-pull ollama-status
.PHONY: up down up-no-llm logs restart

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
RED := \033[0;31m
NC := \033[0m # No Color

## help: Show this help message
help:
	@echo "PlatePilot Go Backend - Available commands:"
	@echo ""
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'

# =============================================================================
# Quick Start (recommended for local development)
# =============================================================================

## up: Start complete dev environment (Ollama + all services with hot reload)
up:
	@echo "$(GREEN)Starting PlatePilot development environment...$(NC)"
	@echo ""
	@echo "Step 1/4: Starting Ollama..."
	@cd ../.. && docker compose --profile llm up -d ollama
	@echo "Waiting for Ollama to be ready..."
	@sleep 5
	@echo ""
	@echo "Step 2/4: Pulling LLM models (first run may take a few minutes)..."
	@docker exec platepilot-ollama ollama pull llama3.1 2>/dev/null || echo "  llama3.1 already available or pulling..."
	@docker exec platepilot-ollama ollama pull nomic-embed-text 2>/dev/null || echo "  nomic-embed-text already available or pulling..."
	@echo ""
	@echo "Step 3/4: Starting infrastructure and services..."
	@cd ../.. && docker compose up -d
	@echo ""
	@echo "Step 4/4: Waiting for services to be healthy..."
	@sleep 10
	@echo ""
	@echo "$(GREEN)==============================================$(NC)"
	@echo "$(GREEN)PlatePilot is ready!$(NC)"
	@echo "$(GREEN)==============================================$(NC)"
	@echo ""
	@echo "  Frontend:     http://localhost:9000"
	@echo "  API (BFF):    http://localhost:8080"
	@echo "  RabbitMQ UI:  http://localhost:15672 (platepilot/platepilot)"
	@echo "  Ollama:       http://localhost:11434"
	@echo ""
	@echo "Logs: docker compose logs -f"
	@echo "Stop: make down"
	@echo ""

## down: Stop all services including Ollama
down:
	@echo "$(GREEN)Stopping PlatePilot...$(NC)"
	@cd ../.. && docker compose --profile llm down
	@echo "$(GREEN)All services stopped.$(NC)"

## up-no-llm: Start dev environment without Ollama (for when using external LLM)
up-no-llm:
	@echo "$(GREEN)Starting PlatePilot (without local LLM)...$(NC)"
	@cd ../.. && docker compose up -d
	@echo ""
	@echo "$(GREEN)PlatePilot is ready (LLM features disabled)$(NC)"
	@echo ""
	@echo "  Frontend:     http://localhost:9000"
	@echo "  API (BFF):    http://localhost:8080"
	@echo ""

## logs: Follow logs from all services
logs:
	@cd ../.. && docker compose --profile llm logs -f

## restart: Restart all services (keeps Ollama running)
restart:
	@echo "$(GREEN)Restarting services...$(NC)"
	@cd ../.. && docker compose restart recipe-api mealplanner-api mobile-bff frontend
	@echo "$(GREEN)Services restarted.$(NC)"

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

## generate: Generate all code (proto + sqlc + swagger)
generate: proto sqlc swagger

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
	go install github.com/swaggo/swag/cmd/swag@latest
	@echo "$(GREEN)Don't forget to install:$(NC)"
	@echo "  - protoc (Protocol Buffers compiler)"
	@echo "  - golang-migrate (database migrations)"

## verify: Verify project setup (tools + generated code)
verify: verify-proto verify-swagger
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
	@which swag > /dev/null 2>&1 && echo "✓ swag" || echo "✗ swag (run: make tools)"
	@which docker > /dev/null 2>&1 && echo "✓ docker" || echo "✗ docker"
	@echo ""
	@echo "$(GREEN)Verification complete!$(NC)"

## verify-proto: Verify generated protobuf code is up-to-date
verify-proto:
	@echo "$(GREEN)Verifying protobuf code is up-to-date...$(NC)"
	@# Create temp directory for comparison
	@mkdir -p /tmp/proto-verify
	@# Copy current generated files
	@cp internal/recipe/pb/*.go /tmp/proto-verify/ 2>/dev/null || true
	@cp internal/mealplanner/pb/*.go /tmp/proto-verify/ 2>/dev/null || true
	@# Regenerate
	@./scripts/generate-proto.sh > /dev/null 2>&1
	@# Compare
	@if ! diff -q internal/recipe/pb/recipe.pb.go /tmp/proto-verify/recipe.pb.go > /dev/null 2>&1; then \
		echo "$(RED)ERROR: recipe.pb.go is out of date. Run 'make proto' to regenerate.$(NC)"; \
		rm -rf /tmp/proto-verify; \
		exit 1; \
	fi
	@if ! diff -q internal/mealplanner/pb/mealplanner.pb.go /tmp/proto-verify/mealplanner.pb.go > /dev/null 2>&1; then \
		echo "$(RED)ERROR: mealplanner.pb.go is out of date. Run 'make proto' to regenerate.$(NC)"; \
		rm -rf /tmp/proto-verify; \
		exit 1; \
	fi
	@rm -rf /tmp/proto-verify
	@echo "$(GREEN)✓ Protobuf code is up-to-date$(NC)"

## swagger: Generate OpenAPI spec from swagger annotations
swagger:
	@echo "$(GREEN)Generating OpenAPI spec...$(NC)"
	@mkdir -p api/openapi
	swag init -g cmd/mobile-bff/main.go -o api/openapi --outputTypes yaml
	@echo "$(GREEN)OpenAPI spec generated at api/openapi/swagger.yaml$(NC)"

## verify-swagger: Verify OpenAPI spec is up-to-date
verify-swagger:
	@echo "$(GREEN)Verifying OpenAPI spec is up-to-date...$(NC)"
	@mkdir -p /tmp/swagger-verify
	@cp api/openapi/swagger.yaml /tmp/swagger-verify/ 2>/dev/null || true
	@swag init -g cmd/mobile-bff/main.go -o api/openapi --outputTypes yaml > /dev/null 2>&1
	@if ! diff -q api/openapi/swagger.yaml /tmp/swagger-verify/swagger.yaml > /dev/null 2>&1; then \
		echo "$(RED)ERROR: swagger.yaml is out of date. Run 'make swagger' to regenerate.$(NC)"; \
		rm -rf /tmp/swagger-verify; \
		exit 1; \
	fi
	@rm -rf /tmp/swagger-verify
	@echo "$(GREEN)✓ OpenAPI spec is up-to-date$(NC)"

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

# =============================================================================
# Ollama (Local LLM)
# =============================================================================

## ollama-start: Start Ollama service in Docker
ollama-start:
	@echo "$(GREEN)Starting Ollama...$(NC)"
	cd ../.. && docker compose --profile llm up -d ollama
	@echo "Waiting for Ollama to be ready..."
	@sleep 10
	@echo "$(GREEN)Ollama ready at http://localhost:11434$(NC)"

## ollama-stop: Stop Ollama service
ollama-stop:
	@echo "$(GREEN)Stopping Ollama...$(NC)"
	cd ../.. && docker compose --profile llm stop ollama

## ollama-pull: Pull required LLM models
ollama-pull:
	@echo "$(GREEN)Pulling LLM models...$(NC)"
	@echo "Pulling llama3.1 (chat model)..."
	docker exec platepilot-ollama ollama pull llama3.1
	@echo "Pulling nomic-embed-text (embedding model)..."
	docker exec platepilot-ollama ollama pull nomic-embed-text
	@echo "$(GREEN)Models ready!$(NC)"

## ollama-status: Show Ollama status and available models
ollama-status:
	@echo "$(GREEN)Ollama Status:$(NC)"
	@docker exec platepilot-ollama ollama list 2>/dev/null || echo "Ollama is not running. Use 'make ollama-start' to start it."

## ollama-logs: Show Ollama logs
ollama-logs:
	docker logs -f platepilot-ollama
