# PlatePilot: .NET to Go Migration Plan

## Executive Summary

This document outlines the migration strategy for PlatePilot from .NET 9.0 to Go, moving from a microservices architecture with Entity Framework Core, MediatR, and .NET Aspire to a Go-based stack with explicit patterns and simpler tooling.

**Current State:** ~3,500 LOC across 52 C# files
**Target State:** ~2,500-3,000 LOC across Go packages
**Estimated Duration:** 6-10 weeks (incremental, service-by-service)

### Key Decisions

- **Hobby Project**: No backwards compatibility needed, delete freely, keep it simple
- **RabbitMQ**: Replacing Azure Service Bus (simpler, better local dev)
- **No ORM**: Use pgx with raw SQL for explicit, maintainable queries
- **Delete .NET**: Remove legacy code once Go equivalent works

---

## 1. Technology Stack Decisions

### 1.1 Go Stack Selection

| Concern | .NET Current | Go Target | Rationale |
|---------|--------------|-----------|-----------|
| **Web Framework** | ASP.NET Core Minimal APIs | `chi` router | Lightweight, idiomatic, middleware support |
| **gRPC** | Grpc.AspNetCore | `google.golang.org/grpc` | Standard, mature ecosystem |
| **Database** | Entity Framework Core | `pgx` + `sqlc` | Type-safe queries, no ORM magic |
| **Migrations** | EF Core Migrations | `golang-migrate` | SQL-based, version controlled |
| **DI Container** | Microsoft.Extensions.DI | Manual / `wire` | Explicit wiring, compile-time safety |
| **CQRS** | MediatR | Custom dispatcher | Simple interface-based handlers |
| **Message Bus** | Azure Service Bus | `rabbitmq/amqp091-go` | Simple, great local dev |
| **Vector Search** | Pgvector.EntityFrameworkCore | `pgvector-go` | Native pgvector support |
| **Config** | appsettings.json + Aspire | `viper` + env vars | 12-factor app compliance |
| **Logging** | ILogger | `slog` (stdlib) | Structured logging, zero deps |
| **Testing** | (none currently) | `testify` + `testcontainers` | Table-driven tests, integration |
| **API Docs** | Swagger/OpenAPI | `swaggo/swag` | Generate from annotations |
| **Observability** | OpenTelemetry | `go.opentelemetry.io/otel` | Same standard, Go SDK |

### 1.2 Project Structure (Target)

```
src/backend-go/
├── cmd/
│   ├── recipe-api/
│   │   └── main.go
│   ├── mealplanner-api/
│   │   └── main.go
│   └── mobile-bff/
│       └── main.go
├── internal/
│   ├── recipe/
│   │   ├── domain/
│   │   │   ├── recipe.go
│   │   │   ├── ingredient.go
│   │   │   └── cuisine.go
│   │   ├── handler/
│   │   │   ├── grpc.go
│   │   │   ├── commands.go
│   │   │   └── queries.go
│   │   ├── repository/
│   │   │   ├── postgres.go
│   │   │   └── queries.sql
│   │   └── events/
│   │       ├── publisher.go
│   │       └── types.go
│   ├── mealplanner/
│   │   ├── domain/
│   │   │   ├── planner.go
│   │   │   └── suggestion.go
│   │   ├── handler/
│   │   │   └── grpc.go
│   │   ├── repository/
│   │   │   └── postgres.go
│   │   └── events/
│   │       └── consumer.go
│   ├── bff/
│   │   ├── handler/
│   │   │   ├── recipe.go
│   │   │   └── mealplan.go
│   │   └── client/
│   │       ├── recipe_client.go
│   │       └── mealplanner_client.go
│   └── common/
│       ├── events/
│       │   ├── bus.go
│       │   └── types.go
│       ├── vector/
│       │   └── extensions.go
│       └── config/
│           └── config.go
├── api/
│   └── proto/
│       ├── recipe/v1/
│       │   └── recipe.proto
│       └── mealplanner/v1/
│           └── mealplanner.proto
├── migrations/
│   ├── recipe/
│   │   ├── 000001_init.up.sql
│   │   └── 000001_init.down.sql
│   └── mealplanner/
│       ├── 000001_init.up.sql
│       └── 000001_init.down.sql
├── deployments/
│   ├── docker-compose.yml
│   ├── docker-compose.dev.yml
│   └── Dockerfile
├── scripts/
│   ├── generate-proto.sh
│   └── migrate.sh
├── go.mod
├── go.sum
└── Makefile
```

---

## 2. Migration Phases

### Phase Overview

```
Phase 0: Foundation Setup          [Week 1]     ✅ COMPLETE
    ↓
Phase 1: Common Layer              [Week 1-2]   ✅ COMPLETE
    ↓
Phase 2: Mobile BFF                [Week 2-3]   ✅ COMPLETE
    ↓
Phase 3: MealPlanner API           [Week 3-5]
    ↓
Phase 4: Recipe API                [Week 5-8]
    ↓
Phase 5: Integration & Cleanup     [Week 8-10]
```

---

## 3. Phase 0: Foundation Setup

**Duration:** 1 week
**Goal:** Set up Go project structure, tooling, and infrastructure

### 3.1 Tasks

- [ ] **P0-1: Initialize Go module**
  - Create `src/backend-go/` directory
  - Run `go mod init github.com/user/platepilot`
  - Set up `.gitignore` for Go

- [ ] **P0-2: Set up development tooling**
  - Install `golangci-lint` for linting
  - Configure `air` for hot reload
  - Set up VS Code / GoLand settings
  - Create `Makefile` with common commands

- [ ] **P0-3: Create Docker Compose for local dev**
  - PostgreSQL with pgvector extension
  - Azure Service Bus emulator (or switch to RabbitMQ for simplicity)
  - Redis (optional, for caching)
  - Adminer/pgAdmin for DB management

- [ ] **P0-4: Set up proto generation**
  - Install `protoc` and Go plugins
  - Create `scripts/generate-proto.sh`
  - Copy existing `.proto` files to `api/proto/`
  - Generate Go code from protos

- [ ] **P0-5: Set up database migrations**
  - Install `golang-migrate`
  - Convert EF Core migrations to SQL
  - Create `migrations/recipe/000001_init.up.sql`
  - Create `migrations/mealplanner/000001_init.up.sql`

- [ ] **P0-6: Create base configuration**
  - Set up `viper` for config loading
  - Create `config.yaml` template
  - Environment variable overrides
  - Document all config options

### 3.2 Deliverables

```
src/backend-go/
├── go.mod
├── go.sum
├── Makefile
├── .golangci.yml
├── docker-compose.yml
├── api/proto/...
├── migrations/...
└── internal/common/config/config.go
```

### 3.3 Makefile Template

```makefile
.PHONY: all build run test lint proto migrate

# Build all services
build:
	go build -o bin/recipe-api ./cmd/recipe-api
	go build -o bin/mealplanner-api ./cmd/mealplanner-api
	go build -o bin/mobile-bff ./cmd/mobile-bff

# Run with hot reload
dev-recipe:
	air -c .air.recipe.toml

dev-mealplanner:
	air -c .air.mealplanner.toml

dev-bff:
	air -c .air.bff.toml

# Run all services
run:
	docker-compose up -d postgres servicebus
	@sleep 2
	make migrate-up
	@$(MAKE) -j3 dev-recipe dev-mealplanner dev-bff

# Generate protobuf code
proto:
	./scripts/generate-proto.sh

# Run migrations
migrate-up:
	migrate -path migrations/recipe -database "postgres://..." up
	migrate -path migrations/mealplanner -database "postgres://..." up

migrate-down:
	migrate -path migrations/recipe -database "postgres://..." down 1
	migrate -path migrations/mealplanner -database "postgres://..." down 1

# Linting
lint:
	golangci-lint run ./...

# Testing
test:
	go test -v -race ./...

test-integration:
	go test -v -tags=integration ./...

# Generate sqlc queries
sqlc:
	sqlc generate
```

---

## 4. Phase 1: Common Layer

**Duration:** 1 week
**Goal:** Migrate shared domain models, events, and utilities

### 4.1 Tasks

- [x] **P1-1: Migrate domain models**
  - `internal/common/domain/recipe.go` - Recipe, Ingredient, Cuisine, Allergy structs
  - `internal/common/domain/metadata.go` - Metadata, NutritionalInfo
  - Add JSON tags and validation tags

- [x] **P1-2: Migrate event types**
  - `internal/common/events/types.go` - Event interfaces and concrete types
  - `RecipeCreatedEvent`, `RecipeUpdatedEvent`
  - JSON serialization for Service Bus

- [x] **P1-3: Create event bus abstraction**
  - `internal/common/events/bus.go` - Publisher/Subscriber interfaces
  - `internal/common/events/servicebus.go` - Azure implementation
  - Consider adding in-memory implementation for testing

- [x] **P1-4: Migrate vector utilities**
  - `internal/common/vector/generator.go` - Vector generation
  - Port hash-based POC implementation
  - Add interface for future Azure OpenAI integration

- [x] **P1-5: Create DTO types**
  - `internal/common/dto/recipe.go` - RecipeDTO for events
  - Conversion functions: `ToDTO()`, `FromDTO()`

- [x] **P1-6: Set up sqlc for type-safe queries**
  - Create `sqlc.yaml` configuration
  - Define query files per service
  - Generate Go code from SQL

### 4.2 Code Examples

**Domain Model (recipe.go):**
```go
package domain

import (
    "time"

    "github.com/google/uuid"
    "github.com/pgvector/pgvector-go"
)

type Recipe struct {
    ID             uuid.UUID      `json:"id" db:"id"`
    Name           string         `json:"name" db:"name"`
    Description    string         `json:"description" db:"description"`
    PrepTime       string         `json:"prepTime" db:"prep_time"`
    CookTime       string         `json:"cookTime" db:"cook_time"`
    MainIngredient *Ingredient    `json:"mainIngredient,omitempty"`
    Cuisine        *Cuisine       `json:"cuisine,omitempty"`
    Ingredients    []Ingredient   `json:"ingredients"`
    Directions     []string       `json:"directions" db:"directions"`
    Metadata       Metadata       `json:"metadata"`
    CreatedAt      time.Time      `json:"createdAt" db:"created_at"`
    UpdatedAt      time.Time      `json:"updatedAt" db:"updated_at"`
}

type Metadata struct {
    SearchVector  pgvector.Vector `json:"-" db:"search_vector"`
    ImageURL      string          `json:"imageUrl" db:"image_url"`
    Tags          []string        `json:"tags" db:"tags"`
    PublishedDate *time.Time      `json:"publishedDate" db:"published_date"`
}

type Ingredient struct {
    ID        uuid.UUID  `json:"id" db:"id"`
    Name      string     `json:"name" db:"name"`
    Allergies []Allergy  `json:"allergies,omitempty"`
}

type Cuisine struct {
    ID   uuid.UUID `json:"id" db:"id"`
    Name string    `json:"name" db:"name"`
}

type Allergy struct {
    ID   uuid.UUID `json:"id" db:"id"`
    Name string    `json:"name" db:"name"`
}
```

**Event Types (types.go):**
```go
package events

import (
    "time"

    "github.com/google/uuid"
)

// Event is the base interface for all domain events
type Event interface {
    EventID() uuid.UUID
    EventType() string
    OccurredAt() time.Time
    AggregateID() uuid.UUID
}

// BaseEvent provides common event fields
type BaseEvent struct {
    ID         uuid.UUID `json:"id"`
    Type       string    `json:"type"`
    OccurredOn time.Time `json:"occurredOn"`
    AggregateId uuid.UUID `json:"aggregateId"`
}

func (e BaseEvent) EventID() uuid.UUID     { return e.ID }
func (e BaseEvent) EventType() string      { return e.Type }
func (e BaseEvent) OccurredAt() time.Time  { return e.OccurredOn }
func (e BaseEvent) AggregateID() uuid.UUID { return e.AggregateId }

// RecipeCreatedEvent is published when a new recipe is created
type RecipeCreatedEvent struct {
    BaseEvent
    Recipe RecipeDTO `json:"recipe"`
}

// RecipeUpdatedEvent is published when a recipe is updated
type RecipeUpdatedEvent struct {
    BaseEvent
}

// NewRecipeCreatedEvent creates a new RecipeCreatedEvent
func NewRecipeCreatedEvent(recipe RecipeDTO) RecipeCreatedEvent {
    return RecipeCreatedEvent{
        BaseEvent: BaseEvent{
            ID:          uuid.New(),
            Type:        "RecipeCreatedEvent",
            OccurredOn:  time.Now().UTC(),
            AggregateId: recipe.ID,
        },
        Recipe: recipe,
    }
}
```

**Event Bus Interface (bus.go):**
```go
package events

import "context"

// Publisher publishes events to the message bus
type Publisher interface {
    Publish(ctx context.Context, event Event) error
    Close() error
}

// Subscriber subscribes to events from the message bus
type Subscriber interface {
    Subscribe(ctx context.Context, handler EventHandler) error
    Close() error
}

// EventHandler processes incoming events
type EventHandler interface {
    Handle(ctx context.Context, event Event) error
    EventTypes() []string // Returns event types this handler can process
}

// EventBus combines publishing and subscribing capabilities
type EventBus interface {
    Publisher
    Subscriber
}
```

**Vector Extensions (extensions.go):**
```go
package vector

import (
    "hash/fnv"
    "strings"

    "github.com/pgvector/pgvector-go"
)

const VectorDimensions = 128

// Generator creates vector embeddings for recipes
type Generator interface {
    Generate(text string) pgvector.Vector
    GenerateForRecipe(name, description, mainIngredient string) pgvector.Vector
}

// HashGenerator is a POC implementation using hash-based vectors
// TODO: Replace with Azure OpenAI embeddings
type HashGenerator struct{}

func NewHashGenerator() *HashGenerator {
    return &HashGenerator{}
}

func (g *HashGenerator) Generate(text string) pgvector.Vector {
    words := strings.Fields(strings.ToLower(text))
    vector := make([]float32, VectorDimensions)

    for _, word := range words {
        h := fnv.New32a()
        h.Write([]byte(word))
        idx := h.Sum32() % uint32(VectorDimensions)
        vector[idx] += 1.0
    }

    // Normalize
    var sum float32
    for _, v := range vector {
        sum += v * v
    }
    if sum > 0 {
        norm := float32(1.0 / sqrt64(float64(sum)))
        for i := range vector {
            vector[i] *= norm
        }
    }

    return pgvector.NewVector(vector)
}

func (g *HashGenerator) GenerateForRecipe(name, description, mainIngredient string) pgvector.Vector {
    combined := name + " " + description + " " + mainIngredient
    return g.Generate(combined)
}

func sqrt64(x float64) float64 {
    // Simple Newton-Raphson
    z := x / 2
    for i := 0; i < 10; i++ {
        z = z - (z*z-x)/(2*z)
    }
    return z
}
```

### 4.3 Deliverables

```
internal/common/
├── domain/
│   ├── recipe.go
│   ├── ingredient.go
│   ├── cuisine.go
│   └── allergy.go
├── dto/
│   └── recipe.go
├── events/
│   ├── types.go
│   ├── bus.go
│   └── servicebus.go
├── vector/
│   └── extensions.go
└── config/
    └── config.go
```

---

## 5. Phase 2: Mobile BFF

**Duration:** 1 week
**Goal:** Migrate the REST API gateway

### 5.1 Tasks

- [x] **P2-1: Set up HTTP server**
  - `cmd/mobile-bff/main.go` - Entry point
  - Chi router setup with middleware
  - Graceful shutdown handling

- [x] **P2-2: Create gRPC clients**
  - `internal/bff/client/recipe.go`
  - `internal/bff/client/mealplanner.go`
  - Connection pooling and retry logic

- [x] **P2-3: Implement REST handlers**
  - `internal/bff/handler/recipe.go` - All recipe endpoints
  - `internal/bff/handler/mealplan.go` - Meal planning endpoints
  - Request validation and error handling

- [x] **P2-4: Add middleware**
  - Logging middleware (request/response)
  - Recovery middleware (panic handling)
  - CORS middleware
  - Request ID middleware

- [ ] **P2-5: OpenAPI documentation** (deferred)
  - Add swaggo annotations to handlers
  - Generate OpenAPI spec
  - Serve Swagger UI

- [ ] **P2-6: Write tests** (deferred)
  - Unit tests for handlers (mock gRPC clients)
  - Integration tests with testcontainers

### 5.2 Code Examples

**Main Entry Point (main.go):**
```go
package main

import (
    "context"
    "log/slog"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"

    "github.com/user/platepilot/internal/bff/client"
    "github.com/user/platepilot/internal/bff/handler"
    "github.com/user/platepilot/internal/common/config"
)

func main() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        slog.Error("failed to load config", "error", err)
        os.Exit(1)
    }

    // Set up structured logging
    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
    }))
    slog.SetDefault(logger)

    // Create gRPC clients
    recipeClient, err := client.NewRecipeClient(cfg.RecipeAPI.Address)
    if err != nil {
        slog.Error("failed to create recipe client", "error", err)
        os.Exit(1)
    }
    defer recipeClient.Close()

    mealPlannerClient, err := client.NewMealPlannerClient(cfg.MealPlannerAPI.Address)
    if err != nil {
        slog.Error("failed to create meal planner client", "error", err)
        os.Exit(1)
    }
    defer mealPlannerClient.Close()

    // Create handlers
    recipeHandler := handler.NewRecipeHandler(recipeClient)
    mealPlanHandler := handler.NewMealPlanHandler(mealPlannerClient)

    // Set up router
    r := chi.NewRouter()

    // Middleware
    r.Use(middleware.RequestID)
    r.Use(middleware.RealIP)
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Use(middleware.Timeout(30 * time.Second))

    // Health checks
    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })

    // API v1 routes
    r.Route("/v1", func(r chi.Router) {
        r.Route("/recipe", func(r chi.Router) {
            r.Get("/{id}", recipeHandler.GetByID)
            r.Get("/all", recipeHandler.GetAll)
            r.Get("/similar", recipeHandler.GetSimilar)
            r.Get("/cuisine/{id}", recipeHandler.GetByCuisine)
            r.Get("/ingredient/{id}", recipeHandler.GetByIngredient)
            r.Get("/allergy/{id}", recipeHandler.GetByAllergy)
            r.Post("/create", recipeHandler.Create)
        })
        r.Route("/mealplan", func(r chi.Router) {
            r.Post("/suggest", mealPlanHandler.Suggest)
        })
    })

    // Create server
    srv := &http.Server{
        Addr:         cfg.BFF.Address,
        Handler:      r,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    // Start server in goroutine
    go func() {
        slog.Info("starting server", "address", cfg.BFF.Address)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            slog.Error("server error", "error", err)
            os.Exit(1)
        }
    }()

    // Wait for interrupt signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    slog.Info("shutting down server...")

    // Graceful shutdown with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        slog.Error("server forced to shutdown", "error", err)
    }

    slog.Info("server stopped")
}
```

**Recipe Handler (recipe.go):**
```go
package handler

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/go-chi/chi/v5"
    "github.com/google/uuid"

    "github.com/user/platepilot/internal/bff/client"
)

type RecipeHandler struct {
    client *client.RecipeClient
}

func NewRecipeHandler(client *client.RecipeClient) *RecipeHandler {
    return &RecipeHandler{client: client}
}

// GetByID godoc
// @Summary Get a recipe by ID
// @Tags recipes
// @Produce json
// @Param id path string true "Recipe ID"
// @Success 200 {object} domain.Recipe
// @Failure 404 {object} ErrorResponse
// @Router /v1/recipe/{id} [get]
func (h *RecipeHandler) GetByID(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        writeError(w, http.StatusBadRequest, "invalid recipe ID")
        return
    }

    recipe, err := h.client.GetByID(r.Context(), id)
    if err != nil {
        writeError(w, http.StatusNotFound, "recipe not found")
        return
    }

    writeJSON(w, http.StatusOK, recipe)
}

// GetAll godoc
// @Summary Get all recipes with pagination
// @Tags recipes
// @Produce json
// @Param pageIndex query int false "Page index" default(1)
// @Param pageSize query int false "Page size" default(20)
// @Success 200 {array} domain.Recipe
// @Router /v1/recipe/all [get]
func (h *RecipeHandler) GetAll(w http.ResponseWriter, r *http.Request) {
    pageIndex, _ := strconv.Atoi(r.URL.Query().Get("pageIndex"))
    if pageIndex < 1 {
        pageIndex = 1
    }

    pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
    if pageSize < 1 || pageSize > 100 {
        pageSize = 20
    }

    recipes, err := h.client.GetAll(r.Context(), pageIndex, pageSize)
    if err != nil {
        writeError(w, http.StatusInternalServerError, "failed to fetch recipes")
        return
    }

    writeJSON(w, http.StatusOK, recipes)
}

// GetSimilar godoc
// @Summary Get similar recipes
// @Tags recipes
// @Produce json
// @Param recipe query string true "Recipe ID to find similar recipes for"
// @Param amount query int false "Number of results" default(5)
// @Success 200 {array} domain.Recipe
// @Router /v1/recipe/similar [get]
func (h *RecipeHandler) GetSimilar(w http.ResponseWriter, r *http.Request) {
    recipeID := r.URL.Query().Get("recipe")
    if recipeID == "" {
        writeError(w, http.StatusBadRequest, "recipe parameter required")
        return
    }

    id, err := uuid.Parse(recipeID)
    if err != nil {
        writeError(w, http.StatusBadRequest, "invalid recipe ID")
        return
    }

    amount, _ := strconv.Atoi(r.URL.Query().Get("amount"))
    if amount < 1 || amount > 50 {
        amount = 5
    }

    recipes, err := h.client.GetSimilar(r.Context(), id, amount)
    if err != nil {
        writeError(w, http.StatusInternalServerError, "failed to fetch similar recipes")
        return
    }

    writeJSON(w, http.StatusOK, recipes)
}

// Create godoc
// @Summary Create a new recipe
// @Tags recipes
// @Accept json
// @Produce json
// @Param recipe body CreateRecipeRequest true "Recipe data"
// @Success 201 {object} domain.Recipe
// @Failure 400 {object} ErrorResponse
// @Router /v1/recipe/create [post]
func (h *RecipeHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req CreateRecipeRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "invalid request body")
        return
    }

    if err := req.Validate(); err != nil {
        writeError(w, http.StatusBadRequest, err.Error())
        return
    }

    recipe, err := h.client.Create(r.Context(), req.ToCommand())
    if err != nil {
        writeError(w, http.StatusInternalServerError, "failed to create recipe")
        return
    }

    writeJSON(w, http.StatusCreated, recipe)
}

// GetByCuisine, GetByIngredient, GetByAllergy follow same pattern...

// Helper functions
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
    writeJSON(w, status, ErrorResponse{Error: message})
}

type ErrorResponse struct {
    Error string `json:"error"`
}

type CreateRecipeRequest struct {
    Name             string   `json:"name"`
    Description      string   `json:"description"`
    PrepTime         string   `json:"prepTime"`
    CookTime         string   `json:"cookTime"`
    MainIngredientID string   `json:"mainIngredientId"`
    CuisineID        string   `json:"cuisineId"`
    IngredientIDs    []string `json:"ingredientIds"`
    Directions       []string `json:"directions"`
}

func (r *CreateRecipeRequest) Validate() error {
    // Add validation logic
    return nil
}
```

### 5.3 Deliverables

```
cmd/mobile-bff/
└── main.go

internal/bff/
├── handler/
│   ├── recipe.go
│   ├── mealplan.go
│   └── middleware.go
├── client/
│   ├── recipe_client.go
│   └── mealplanner_client.go
└── dto/
    └── requests.go
```

---

## 6. Phase 3: MealPlanner API

**Duration:** 2 weeks
**Goal:** Migrate the read model service with event consumption

### 6.1 Tasks

- [ ] **P3-1: Set up gRPC server**
  - `cmd/mealplanner-api/main.go`
  - gRPC server configuration
  - Health check service

- [ ] **P3-2: Create database repository**
  - `internal/mealplanner/repository/postgres.go`
  - sqlc queries for read model
  - Vector similarity queries

- [ ] **P3-3: Implement meal planner domain logic**
  - `internal/mealplanner/domain/planner.go`
  - Suggestion algorithm with diversity scoring
  - Constraint filtering logic

- [ ] **P3-4: Implement gRPC handlers**
  - `internal/mealplanner/handler/grpc.go`
  - SuggestRecipes implementation
  - Error handling and logging

- [ ] **P3-5: Implement event consumer**
  - `internal/mealplanner/events/consumer.go`
  - Subscribe to recipe-events topic
  - Handle RecipeCreatedEvent
  - Update read model

- [ ] **P3-6: Write SQL migrations**
  - `migrations/mealplanner/000001_init.up.sql`
  - Read model table schema
  - Vector index for similarity search

- [ ] **P3-7: Write tests**
  - Unit tests for planner algorithm
  - Integration tests for repository
  - Event handler tests

### 6.2 Code Examples

**Meal Planner Domain (planner.go):**
```go
package domain

import (
    "context"
    "math"
    "sort"

    "github.com/google/uuid"
    "github.com/pgvector/pgvector-go"
)

type MealPlanner interface {
    SuggestMeals(ctx context.Context, req SuggestionRequest) ([]uuid.UUID, error)
}

type SuggestionRequest struct {
    DailyConstraints       []DailyConstraints
    AlreadySelectedRecipes []uuid.UUID
    Amount                 int
}

type DailyConstraints struct {
    IngredientConstraints []uuid.UUID
    CuisineConstraints    []uuid.UUID
}

type RecipeRepository interface {
    GetAll(ctx context.Context) ([]Recipe, error)
    GetByID(ctx context.Context, id uuid.UUID) (*Recipe, error)
}

type Recipe struct {
    ID            uuid.UUID
    Name          string
    CuisineID     uuid.UUID
    IngredientIDs []uuid.UUID
    SearchVector  pgvector.Vector
}

type mealPlanner struct {
    repo RecipeRepository
}

func NewMealPlanner(repo RecipeRepository) MealPlanner {
    return &mealPlanner{repo: repo}
}

func (p *mealPlanner) SuggestMeals(ctx context.Context, req SuggestionRequest) ([]uuid.UUID, error) {
    recipes, err := p.repo.GetAll(ctx)
    if err != nil {
        return nil, err
    }

    // Filter by constraints
    filtered := p.filterByConstraints(recipes, req.DailyConstraints)

    // Remove already selected
    filtered = p.removeSelected(filtered, req.AlreadySelectedRecipes)

    // Score and rank
    scored := p.scoreRecipes(filtered, req.AlreadySelectedRecipes, recipes)

    // Sort by score descending
    sort.Slice(scored, func(i, j int) bool {
        return scored[i].score > scored[j].score
    })

    // Take top N
    result := make([]uuid.UUID, 0, req.Amount)
    for i := 0; i < len(scored) && i < req.Amount; i++ {
        result = append(result, scored[i].id)
    }

    return result, nil
}

type scoredRecipe struct {
    id    uuid.UUID
    score float64
}

func (p *mealPlanner) filterByConstraints(recipes []Recipe, constraints []DailyConstraints) []Recipe {
    if len(constraints) == 0 {
        return recipes
    }

    var filtered []Recipe
    for _, recipe := range recipes {
        if p.matchesConstraints(recipe, constraints) {
            filtered = append(filtered, recipe)
        }
    }
    return filtered
}

func (p *mealPlanner) matchesConstraints(recipe Recipe, constraints []DailyConstraints) bool {
    for _, daily := range constraints {
        // Check cuisine constraints
        if len(daily.CuisineConstraints) > 0 {
            found := false
            for _, cuisineID := range daily.CuisineConstraints {
                if recipe.CuisineID == cuisineID {
                    found = true
                    break
                }
            }
            if !found {
                return false
            }
        }

        // Check ingredient constraints
        if len(daily.IngredientConstraints) > 0 {
            found := false
            for _, ingredientID := range daily.IngredientConstraints {
                for _, recipeIngredientID := range recipe.IngredientIDs {
                    if recipeIngredientID == ingredientID {
                        found = true
                        break
                    }
                }
                if found {
                    break
                }
            }
            if !found {
                return false
            }
        }
    }
    return true
}

func (p *mealPlanner) removeSelected(recipes []Recipe, selected []uuid.UUID) []Recipe {
    selectedSet := make(map[uuid.UUID]bool)
    for _, id := range selected {
        selectedSet[id] = true
    }

    var filtered []Recipe
    for _, recipe := range recipes {
        if !selectedSet[recipe.ID] {
            filtered = append(filtered, recipe)
        }
    }
    return filtered
}

func (p *mealPlanner) scoreRecipes(candidates []Recipe, selected []uuid.UUID, allRecipes []Recipe) []scoredRecipe {
    // Get selected recipes for diversity calculation
    selectedRecipes := make([]Recipe, 0)
    for _, id := range selected {
        for _, r := range allRecipes {
            if r.ID == id {
                selectedRecipes = append(selectedRecipes, r)
                break
            }
        }
    }

    scored := make([]scoredRecipe, len(candidates))
    for i, candidate := range candidates {
        diversityScore := p.calculateDiversityScore(candidate, selectedRecipes)
        scored[i] = scoredRecipe{
            id:    candidate.ID,
            score: diversityScore,
        }
    }
    return scored
}

func (p *mealPlanner) calculateDiversityScore(candidate Recipe, selected []Recipe) float64 {
    if len(selected) == 0 {
        return 1.0
    }

    var totalSimilarity float64
    for _, s := range selected {
        similarity := p.calculateSimilarity(candidate, s)
        totalSimilarity += similarity
    }

    avgSimilarity := totalSimilarity / float64(len(selected))
    return 1.0 - avgSimilarity // Higher score for more diverse recipes
}

func (p *mealPlanner) calculateSimilarity(a, b Recipe) float64 {
    // Cosine similarity between vectors
    if len(a.SearchVector.Slice()) == 0 || len(b.SearchVector.Slice()) == 0 {
        return 0
    }

    vecA := a.SearchVector.Slice()
    vecB := b.SearchVector.Slice()

    var dotProduct, normA, normB float64
    for i := 0; i < len(vecA) && i < len(vecB); i++ {
        dotProduct += float64(vecA[i]) * float64(vecB[i])
        normA += float64(vecA[i]) * float64(vecA[i])
        normB += float64(vecB[i]) * float64(vecB[i])
    }

    if normA == 0 || normB == 0 {
        return 0
    }

    return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}
```

**Event Consumer (consumer.go):**
```go
package events

import (
    "context"
    "encoding/json"
    "log/slog"

    "github.com/user/platepilot/internal/common/events"
    "github.com/user/platepilot/internal/mealplanner/repository"
)

type RecipeEventHandler struct {
    repo   *repository.PostgresRepository
    logger *slog.Logger
}

func NewRecipeEventHandler(repo *repository.PostgresRepository, logger *slog.Logger) *RecipeEventHandler {
    return &RecipeEventHandler{
        repo:   repo,
        logger: logger,
    }
}

func (h *RecipeEventHandler) Handle(ctx context.Context, event events.Event) error {
    switch e := event.(type) {
    case *events.RecipeCreatedEvent:
        return h.handleRecipeCreated(ctx, e)
    case *events.RecipeUpdatedEvent:
        return h.handleRecipeUpdated(ctx, e)
    default:
        h.logger.Warn("unknown event type", "type", event.EventType())
        return nil
    }
}

func (h *RecipeEventHandler) EventTypes() []string {
    return []string{"RecipeCreatedEvent", "RecipeUpdatedEvent"}
}

func (h *RecipeEventHandler) handleRecipeCreated(ctx context.Context, event *events.RecipeCreatedEvent) error {
    h.logger.Info("handling recipe created event",
        "eventId", event.EventID(),
        "recipeId", event.Recipe.ID,
    )

    // Convert DTO to domain model and save to read model
    recipe := event.Recipe.ToDomain()

    if err := h.repo.Upsert(ctx, recipe); err != nil {
        h.logger.Error("failed to upsert recipe",
            "error", err,
            "recipeId", event.Recipe.ID,
        )
        return err
    }

    h.logger.Info("recipe upserted to read model", "recipeId", event.Recipe.ID)
    return nil
}

func (h *RecipeEventHandler) handleRecipeUpdated(ctx context.Context, event *events.RecipeUpdatedEvent) error {
    h.logger.Info("handling recipe updated event",
        "eventId", event.EventID(),
        "aggregateId", event.AggregateID(),
    )

    // For updates, we might need to fetch the full recipe from the source
    // or the event could contain the full updated data
    // Implementation depends on event design decisions

    return nil
}
```

### 6.3 SQL Migration

```sql
-- migrations/mealplanner/000001_init.up.sql

CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE recipes (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    prep_time VARCHAR(50),
    cook_time VARCHAR(50),
    cuisine_id UUID,
    cuisine_name VARCHAR(100),
    main_ingredient_id UUID,
    main_ingredient_name VARCHAR(100),
    ingredient_ids UUID[] DEFAULT '{}',
    directions TEXT[] DEFAULT '{}',
    search_vector vector(128),
    image_url TEXT,
    tags TEXT[] DEFAULT '{}',
    published_date TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Index for vector similarity search
CREATE INDEX recipes_search_vector_idx ON recipes
USING ivfflat (search_vector vector_cosine_ops)
WITH (lists = 100);

-- Index for cuisine filtering
CREATE INDEX recipes_cuisine_id_idx ON recipes (cuisine_id);

-- Index for ingredient filtering (GIN for array)
CREATE INDEX recipes_ingredient_ids_idx ON recipes USING GIN (ingredient_ids);
```

### 6.4 Deliverables

```
cmd/mealplanner-api/
└── main.go

internal/mealplanner/
├── domain/
│   ├── planner.go
│   └── types.go
├── handler/
│   └── grpc.go
├── repository/
│   ├── postgres.go
│   └── queries.sql
└── events/
    └── consumer.go

migrations/mealplanner/
├── 000001_init.up.sql
└── 000001_init.down.sql
```

---

## 7. Phase 4: Recipe API

**Duration:** 3 weeks
**Goal:** Migrate the write model service (most complex)

### 7.1 Tasks

- [ ] **P4-1: Set up gRPC server**
  - `cmd/recipe-api/main.go`
  - Server configuration
  - Health and reflection services

- [ ] **P4-2: Create database repository**
  - `internal/recipe/repository/postgres.go`
  - Full CRUD operations
  - Transaction support
  - Relationship handling (ingredients, cuisines, allergies)

- [ ] **P4-3: Implement command handlers**
  - `internal/recipe/handler/commands.go`
  - CreateRecipe with validation
  - UpdateRecipe
  - DeleteRecipe

- [ ] **P4-4: Implement query handlers**
  - `internal/recipe/handler/queries.go`
  - GetByID with eager loading
  - GetAll with pagination
  - SearchSimilar with vector search
  - Filter by cuisine/ingredient/allergy

- [ ] **P4-5: Implement gRPC service**
  - `internal/recipe/handler/grpc.go`
  - Map gRPC requests to commands/queries
  - Error translation

- [ ] **P4-6: Implement event publisher**
  - `internal/recipe/events/publisher.go`
  - Publish after successful database operations
  - Transactional outbox pattern (optional)

- [ ] **P4-7: Database seeding**
  - `internal/recipe/seed/seeder.go`
  - Load from recipes.json
  - Deduplication logic

- [ ] **P4-8: Write SQL migrations**
  - Full schema with all tables
  - Relationships and indexes
  - pgvector setup

- [ ] **P4-9: Write tests**
  - Repository integration tests
  - Handler unit tests
  - gRPC service tests
  - Seeder tests

### 7.2 SQL Migration

```sql
-- migrations/recipe/000001_init.up.sql

CREATE EXTENSION IF NOT EXISTS vector;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Cuisines table
CREATE TABLE cuisines (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Allergies table
CREATE TABLE allergies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Ingredients table
CREATE TABLE ingredients (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Ingredient-Allergy relationship
CREATE TABLE ingredient_allergies (
    ingredient_id UUID REFERENCES ingredients(id) ON DELETE CASCADE,
    allergy_id UUID REFERENCES allergies(id) ON DELETE CASCADE,
    PRIMARY KEY (ingredient_id, allergy_id)
);

-- Recipes table
CREATE TABLE recipes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    prep_time VARCHAR(50),
    cook_time VARCHAR(50),
    main_ingredient_id UUID REFERENCES ingredients(id),
    cuisine_id UUID REFERENCES cuisines(id),
    directions TEXT[] DEFAULT '{}',
    search_vector vector(128),
    image_url TEXT,
    tags TEXT[] DEFAULT '{}',
    published_date TIMESTAMPTZ,
    calories INTEGER,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Recipe-Ingredient relationship
CREATE TABLE recipe_ingredients (
    recipe_id UUID REFERENCES recipes(id) ON DELETE CASCADE,
    ingredient_id UUID REFERENCES ingredients(id) ON DELETE CASCADE,
    PRIMARY KEY (recipe_id, ingredient_id)
);

-- Indexes
CREATE INDEX recipes_cuisine_id_idx ON recipes (cuisine_id);
CREATE INDEX recipes_main_ingredient_id_idx ON recipes (main_ingredient_id);
CREATE INDEX recipes_search_vector_idx ON recipes
    USING ivfflat (search_vector vector_cosine_ops)
    WITH (lists = 100);
CREATE INDEX recipe_ingredients_recipe_idx ON recipe_ingredients (recipe_id);
CREATE INDEX recipe_ingredients_ingredient_idx ON recipe_ingredients (ingredient_id);
```

### 7.3 Repository Pattern

```go
package repository

import (
    "context"
    "errors"
    "fmt"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/pgvector/pgvector-go"

    "github.com/user/platepilot/internal/common/domain"
)

var (
    ErrRecipeNotFound = errors.New("recipe not found")
    ErrDuplicateKey   = errors.New("duplicate key")
)

type PostgresRepository struct {
    pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
    return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Recipe, error) {
    query := `
        SELECT
            r.id, r.name, r.description, r.prep_time, r.cook_time,
            r.directions, r.search_vector, r.image_url, r.tags,
            r.published_date, r.calories, r.created_at, r.updated_at,
            c.id as cuisine_id, c.name as cuisine_name,
            mi.id as main_ingredient_id, mi.name as main_ingredient_name
        FROM recipes r
        LEFT JOIN cuisines c ON r.cuisine_id = c.id
        LEFT JOIN ingredients mi ON r.main_ingredient_id = mi.id
        WHERE r.id = $1
    `

    var recipe domain.Recipe
    var cuisine domain.Cuisine
    var mainIngredient domain.Ingredient
    var searchVector pgvector.Vector

    err := r.pool.QueryRow(ctx, query, id).Scan(
        &recipe.ID, &recipe.Name, &recipe.Description,
        &recipe.PrepTime, &recipe.CookTime, &recipe.Directions,
        &searchVector, &recipe.Metadata.ImageURL, &recipe.Metadata.Tags,
        &recipe.Metadata.PublishedDate, &recipe.NutritionalInfo.Calories,
        &recipe.CreatedAt, &recipe.UpdatedAt,
        &cuisine.ID, &cuisine.Name,
        &mainIngredient.ID, &mainIngredient.Name,
    )

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, ErrRecipeNotFound
        }
        return nil, fmt.Errorf("query recipe: %w", err)
    }

    recipe.Cuisine = &cuisine
    recipe.MainIngredient = &mainIngredient
    recipe.Metadata.SearchVector = searchVector

    // Load ingredients
    ingredients, err := r.getRecipeIngredients(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("load ingredients: %w", err)
    }
    recipe.Ingredients = ingredients

    return &recipe, nil
}

func (r *PostgresRepository) Create(ctx context.Context, recipe *domain.Recipe) error {
    tx, err := r.pool.Begin(ctx)
    if err != nil {
        return fmt.Errorf("begin transaction: %w", err)
    }
    defer tx.Rollback(ctx)

    // Insert recipe
    query := `
        INSERT INTO recipes (
            id, name, description, prep_time, cook_time,
            main_ingredient_id, cuisine_id, directions, search_vector,
            image_url, tags, published_date, calories
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
    `

    _, err = tx.Exec(ctx, query,
        recipe.ID, recipe.Name, recipe.Description,
        recipe.PrepTime, recipe.CookTime,
        recipe.MainIngredient.ID, recipe.Cuisine.ID,
        recipe.Directions, recipe.Metadata.SearchVector,
        recipe.Metadata.ImageURL, recipe.Metadata.Tags,
        recipe.Metadata.PublishedDate, recipe.NutritionalInfo.Calories,
    )
    if err != nil {
        return fmt.Errorf("insert recipe: %w", err)
    }

    // Insert recipe-ingredient relationships
    for _, ingredient := range recipe.Ingredients {
        _, err = tx.Exec(ctx,
            `INSERT INTO recipe_ingredients (recipe_id, ingredient_id) VALUES ($1, $2)`,
            recipe.ID, ingredient.ID,
        )
        if err != nil {
            return fmt.Errorf("insert recipe ingredient: %w", err)
        }
    }

    if err := tx.Commit(ctx); err != nil {
        return fmt.Errorf("commit transaction: %w", err)
    }

    return nil
}

func (r *PostgresRepository) SearchSimilar(ctx context.Context, recipeID uuid.UUID, limit int) ([]domain.Recipe, error) {
    // First get the vector for the target recipe
    var targetVector pgvector.Vector
    err := r.pool.QueryRow(ctx,
        `SELECT search_vector FROM recipes WHERE id = $1`,
        recipeID,
    ).Scan(&targetVector)
    if err != nil {
        return nil, fmt.Errorf("get target vector: %w", err)
    }

    // Find similar recipes using cosine distance
    query := `
        SELECT
            r.id, r.name, r.description, r.prep_time, r.cook_time,
            c.id as cuisine_id, c.name as cuisine_name,
            1 - (r.search_vector <=> $1) as similarity
        FROM recipes r
        LEFT JOIN cuisines c ON r.cuisine_id = c.id
        WHERE r.id != $2
        ORDER BY r.search_vector <=> $1
        LIMIT $3
    `

    rows, err := r.pool.Query(ctx, query, targetVector, recipeID, limit)
    if err != nil {
        return nil, fmt.Errorf("query similar recipes: %w", err)
    }
    defer rows.Close()

    var recipes []domain.Recipe
    for rows.Next() {
        var recipe domain.Recipe
        var cuisine domain.Cuisine
        var similarity float64

        if err := rows.Scan(
            &recipe.ID, &recipe.Name, &recipe.Description,
            &recipe.PrepTime, &recipe.CookTime,
            &cuisine.ID, &cuisine.Name,
            &similarity,
        ); err != nil {
            return nil, fmt.Errorf("scan recipe: %w", err)
        }

        recipe.Cuisine = &cuisine
        recipes = append(recipes, recipe)
    }

    return recipes, nil
}

func (r *PostgresRepository) getRecipeIngredients(ctx context.Context, recipeID uuid.UUID) ([]domain.Ingredient, error) {
    query := `
        SELECT i.id, i.name
        FROM ingredients i
        JOIN recipe_ingredients ri ON i.id = ri.ingredient_id
        WHERE ri.recipe_id = $1
    `

    rows, err := r.pool.Query(ctx, query, recipeID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var ingredients []domain.Ingredient
    for rows.Next() {
        var ingredient domain.Ingredient
        if err := rows.Scan(&ingredient.ID, &ingredient.Name); err != nil {
            return nil, err
        }
        ingredients = append(ingredients, ingredient)
    }

    return ingredients, nil
}

// Additional methods: GetAll, GetByCuisine, GetByIngredient,
// GetByAllergy, Update, Delete...
```

### 7.4 Deliverables

```
cmd/recipe-api/
└── main.go

internal/recipe/
├── domain/
│   └── errors.go
├── handler/
│   ├── grpc.go
│   ├── commands.go
│   └── queries.go
├── repository/
│   ├── postgres.go
│   └── queries.sql
├── events/
│   └── publisher.go
└── seed/
    └── seeder.go

migrations/recipe/
├── 000001_init.up.sql
└── 000001_init.down.sql
```

---

## 8. Phase 5: Integration & Cleanup

**Duration:** 2 weeks
**Goal:** Full system integration, testing, and .NET decommission

### 8.1 Tasks

- [ ] **P5-1: End-to-end integration testing**
  - Full flow: BFF → Recipe API → Event → MealPlanner
  - Load testing with k6 or similar
  - Error scenario testing

- [ ] **P5-2: Docker Compose production setup**
  - Multi-stage Dockerfile for each service
  - Production docker-compose.yml
  - Health checks and restart policies

- [ ] **P5-3: Observability setup**
  - OpenTelemetry tracing across services
  - Prometheus metrics endpoints
  - Structured logging correlation

- [ ] **P5-4: Documentation**
  - Update CLAUDE.md for Go stack
  - API documentation (OpenAPI specs)
  - Deployment guide
  - Development setup guide

- [ ] **P5-5: CI/CD pipeline**
  - GitHub Actions for build/test
  - Container image builds
  - Automated migrations

- [ ] **P5-6: Frontend integration**
  - Update Flutter app (or start Vue migration)
  - Verify all endpoints work
  - Mobile testing

- [ ] **P5-7: .NET decommission**
  - Archive .NET code (or keep in separate branch)
  - Update repository structure
  - Clean up old dependencies

- [ ] **P5-8: Performance comparison**
  - Benchmark against .NET version
  - Memory usage comparison
  - Cold start time comparison

### 8.2 Docker Setup

**Dockerfile (multi-stage):**
```dockerfile
# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build all services
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/recipe-api ./cmd/recipe-api
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/mealplanner-api ./cmd/mealplanner-api
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/mobile-bff ./cmd/mobile-bff

# Recipe API image
FROM alpine:3.19 AS recipe-api
RUN apk add --no-cache ca-certificates
COPY --from=builder /bin/recipe-api /usr/local/bin/
COPY --from=builder /app/migrations/recipe /migrations
EXPOSE 8080 9090
CMD ["recipe-api"]

# MealPlanner API image
FROM alpine:3.19 AS mealplanner-api
RUN apk add --no-cache ca-certificates
COPY --from=builder /bin/mealplanner-api /usr/local/bin/
COPY --from=builder /app/migrations/mealplanner /migrations
EXPOSE 8080 9090
CMD ["mealplanner-api"]

# Mobile BFF image
FROM alpine:3.19 AS mobile-bff
RUN apk add --no-cache ca-certificates
COPY --from=builder /bin/mobile-bff /usr/local/bin/
EXPOSE 8080
CMD ["mobile-bff"]
```

**docker-compose.yml:**
```yaml
version: '3.8'

services:
  postgres:
    image: ankane/pgvector:latest
    environment:
      POSTGRES_USER: platepilot
      POSTGRES_PASSWORD: platepilot
      POSTGRES_DB: platepilot
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U platepilot"]
      interval: 5s
      timeout: 5s
      retries: 5

  servicebus:
    image: mcr.microsoft.com/azure-messaging/servicebus-emulator:latest
    environment:
      ACCEPT_EULA: "Y"
      SQL_SERVER: mssql
    depends_on:
      mssql:
        condition: service_healthy
    ports:
      - "5672:5672"

  mssql:
    image: mcr.microsoft.com/mssql/server:2022-latest
    environment:
      ACCEPT_EULA: "Y"
      SA_PASSWORD: "YourStrong!Passw0rd"
    healthcheck:
      test: /opt/mssql-tools/bin/sqlcmd -S localhost -U sa -P "YourStrong!Passw0rd" -Q "SELECT 1"
      interval: 10s
      timeout: 5s
      retries: 5

  recipe-api:
    build:
      context: .
      target: recipe-api
    environment:
      DATABASE_URL: postgres://platepilot:platepilot@postgres:5432/recipedb?sslmode=disable
      SERVICEBUS_CONNECTION: Endpoint=sb://servicebus;SharedAccessKeyName=...
    depends_on:
      postgres:
        condition: service_healthy
      servicebus:
        condition: service_started
    ports:
      - "8081:8080"
      - "9091:9090"

  mealplanner-api:
    build:
      context: .
      target: mealplanner-api
    environment:
      DATABASE_URL: postgres://platepilot:platepilot@postgres:5432/mealplannerdb?sslmode=disable
      SERVICEBUS_CONNECTION: Endpoint=sb://servicebus;SharedAccessKeyName=...
    depends_on:
      postgres:
        condition: service_healthy
      servicebus:
        condition: service_started
    ports:
      - "8082:8080"
      - "9092:9090"

  mobile-bff:
    build:
      context: .
      target: mobile-bff
    environment:
      RECIPE_API_ADDRESS: recipe-api:9090
      MEALPLANNER_API_ADDRESS: mealplanner-api:9090
    depends_on:
      - recipe-api
      - mealplanner-api
    ports:
      - "8080:8080"

volumes:
  postgres_data:
```

### 8.3 Deliverables

```
deployments/
├── Dockerfile
├── docker-compose.yml
├── docker-compose.dev.yml
└── docker-compose.prod.yml

.github/workflows/
├── build.yml
├── test.yml
└── deploy.yml

docs/
├── setup.md
├── deployment.md
└── api.md
```

---

## 9. Risk Mitigation

### 9.1 Identified Risks

| Risk | Impact | Likelihood | Mitigation |
|------|--------|------------|------------|
| Event bus incompatibility | HIGH | LOW | Test Azure SDK early, have RabbitMQ fallback |
| pgvector performance | MEDIUM | LOW | Benchmark early, tune ivfflat lists parameter |
| Missing EF Core features | MEDIUM | MEDIUM | Use sqlc for type safety, manual migrations |
| Learning curve delays | MEDIUM | MEDIUM | Start with simpler BFF service |
| Frontend breaking changes | HIGH | LOW | API contract testing, versioned endpoints |

### 9.2 Rollback Strategy

1. Keep .NET code in `legacy/dotnet` branch
2. Both stacks can run in parallel during migration
3. Feature flags in BFF to route to either backend
4. Database schemas are compatible (same PostgreSQL)

---

## 10. Success Criteria

### 10.1 Functional Requirements

- [ ] All existing REST endpoints work identically
- [ ] gRPC services maintain same contracts
- [ ] Event-driven sync between services works
- [ ] Vector similarity search returns comparable results
- [ ] Database seeding produces same data

### 10.2 Non-Functional Requirements

- [ ] Cold start time < 500ms (vs ~2-3s for .NET)
- [ ] Memory usage < 50MB per service (vs ~150-200MB)
- [ ] Container image size < 20MB (vs ~200MB)
- [ ] P95 latency comparable or better
- [ ] Test coverage > 70%

### 10.3 Developer Experience

- [ ] Single `make run` starts all services
- [ ] Hot reload works for development
- [ ] Clear error messages and logging
- [ ] Comprehensive API documentation

---

## 11. Timeline Summary

| Phase | Duration | Start | End | Dependencies |
|-------|----------|-------|-----|--------------|
| Phase 0: Foundation | 1 week | Week 1 | Week 1 | None |
| Phase 1: Common | 1 week | Week 1 | Week 2 | Phase 0 |
| Phase 2: Mobile BFF | 1 week | Week 2 | Week 3 | Phase 1 |
| Phase 3: MealPlanner | 2 weeks | Week 3 | Week 5 | Phase 1, 2 |
| Phase 4: Recipe API | 3 weeks | Week 5 | Week 8 | Phase 1, 3 |
| Phase 5: Integration | 2 weeks | Week 8 | Week 10 | All |

**Total: 10 weeks**

---

## 12. Next Steps

1. **Review this plan** - Validate assumptions and adjust scope
2. **Set up development environment** - Go 1.23+, Docker, protoc
3. **Create Phase 0 PR** - Foundation setup
4. **Begin Phase 1** - Common layer migration

---

*Last Updated: 2025-12-28*
*Version: 1.0*
