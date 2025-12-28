# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**PlatePilot** is an intelligent meal planning and recipe management application built with microservices architecture. The application uses AI-powered recommendations (planned) to suggest personalized meal combinations based on user preferences, dietary restrictions, and recipe similarity using vector search.

**This is a hobby project** - no backwards compatibility requirements, no legacy constraints. We can make breaking changes freely and choose the simplest solutions.

## Migration Status: COMPLETE ✅

The backend migration from .NET to Go is **complete**. The .NET code has been removed and the project is now Go-only.

See `MIGRATION_PLAN.md` for the detailed migration history.

### Why Go?
- **AI Agent Friendly**: Explicit control flow, no magic/reflection, errors as values
- **Simpler Deployment**: Single binaries, smaller containers (<20MB vs ~200MB)
- **Lower Resource Usage**: ~50MB memory vs ~150-200MB per service
- **Faster Cold Starts**: <500ms vs 2-3s for .NET

### Architecture

| Component | Location | Description |
|-----------|----------|-------------|
| Recipe API | `src/backend-go/cmd/recipe-api/` | Write service, gRPC, event publishing |
| MealPlanner API | `src/backend-go/cmd/mealplanner-api/` | Read service, vector search, event consuming |
| Mobile BFF | `src/backend-go/cmd/mobile-bff/` | REST gateway for clients |
| Common | `src/backend-go/internal/common/` | Shared domain, events, config |
| Frontend | `src/frontend/` | Vue/Quasar (planned) |

### Completed Migration Phases
1. **Phase 0: Foundation** ✅ - Go project setup, tooling, Docker, migrations
2. **Phase 1: Common Layer** ✅ - Domain models, events, vector utilities
3. **Phase 2: Mobile BFF** ✅ - REST gateway, gRPC clients, handlers
4. **Phase 3: MealPlanner API** ✅ - Read model, gRPC server, RabbitMQ consumer
5. **Phase 4: Recipe API** ✅ - Write model, gRPC server, RabbitMQ publisher, seeder
6. **Phase 5: Integration** ✅ - Docker, CI/CD, E2E testing, .NET cleanup

## Technology Stack

### Backend (Go)
- **Language**: Go 1.23+
- **Web Framework**: `chi` router
- **gRPC**: `google.golang.org/grpc`
- **Database**: `pgx` + raw SQL (no ORM magic)
- **Migrations**: `golang-migrate`
- **Configuration**: `viper`
- **Logging**: `slog` (stdlib)
- **Vector Search**: `pgvector-go`
- **Messaging**: `rabbitmq/amqp091-go`

### Infrastructure
- **Database**: PostgreSQL with pgvector extension
- **Message Broker**: RabbitMQ
- **Inter-service**: gRPC (proto3)
- **Containers**: Docker, docker-compose
- **CI/CD**: GitHub Actions

### Frontend (Planned)
- **Framework**: Vue.js 3 + Quasar
- **Mobile**: Capacitor for iOS/Android
- **Status**: Not yet started

## Build & Development Commands

```bash
cd src/backend-go

# First time setup
make tools              # Install dev tools (golangci-lint, air, protoc plugins)
make docker-up          # Start PostgreSQL + RabbitMQ
make migrate-up         # Run database migrations

# Development (local with hot reload)
make dev                # Run all services with hot reload
make dev-recipe         # Run recipe-api only
make dev-mealplanner    # Run mealplanner-api only
make dev-bff            # Run mobile-bff only

# Docker (full stack)
make docker-run-all     # Build and run complete stack
make docker-run-detached # Run in background
make docker-logs        # View logs
make docker-down        # Stop everything

# Build & Test
make build              # Build all binaries
make test               # Run unit tests
make test-e2e           # Run E2E integration tests
make lint               # Run linter

# Database
make migrate-up         # Apply all migrations
make migrate-down       # Rollback last migration
make seed               # Seed with sample recipes

# Code generation
make proto              # Generate gRPC code from protos
```

## Architecture Overview

### Microservices Structure

The backend follows **CQRS pattern** with event-driven architecture:

1. **Recipe API** (Write Service) - `cmd/recipe-api/` + `internal/recipe/`
   - Handles recipe creation, updates, and command operations
   - PostgreSQL database with pgvector for embeddings
   - Publishes events to RabbitMQ on recipe changes
   - gRPC service for inter-service communication

2. **MealPlanner API** (Read Service) - `cmd/mealplanner-api/` + `internal/mealplanner/`
   - Handles meal planning queries and recipe suggestions
   - Denormalized read model for query performance
   - Vector similarity search using pgvector
   - Consumes recipe events to update read model

3. **Mobile BFF** (Backend-for-Frontend) - `cmd/mobile-bff/` + `internal/bff/`
   - REST API gateway for mobile/web clients
   - Aggregates calls to RecipeApi and MealPlannerApi via gRPC
   - Health/ready endpoints for orchestration

### Project Structure

```
src/backend-go/
├── cmd/                          # Service entry points
│   ├── recipe-api/main.go
│   ├── mealplanner-api/main.go
│   └── mobile-bff/main.go
├── internal/                     # Private application code
│   ├── recipe/
│   │   ├── domain/               # Business entities
│   │   ├── handler/              # gRPC + HTTP handlers
│   │   ├── repository/           # Database access
│   │   └── events/               # Event publishing
│   ├── mealplanner/
│   │   ├── domain/               # Planner logic
│   │   ├── handler/              # gRPC handlers
│   │   ├── repository/           # Read model access
│   │   └── events/               # Event consumption
│   ├── bff/
│   │   ├── handler/              # REST handlers
│   │   └── client/               # gRPC clients
│   └── common/
│       ├── config/               # Viper configuration
│       ├── domain/               # Shared domain types
│       ├── dto/                  # Data transfer objects
│       ├── events/               # Event bus abstraction
│       └── vector/               # Vector utilities
├── api/proto/                    # Protobuf definitions
├── migrations/                   # SQL migrations
├── deployments/                  # Docker configs
└── Makefile
```

### Key Go Patterns

- **No ORM**: Use `pgx` with raw SQL for explicit queries
- **Explicit errors**: Return errors, don't panic
- **Interface-based DI**: Define interfaces, inject implementations
- **Table-driven tests**: Use subtests with `t.Run()`
- **Structured logging**: Use `slog` with context

### Event-Driven Communication

- **RabbitMQ** exchange: `recipe-events` (topic exchange)
- **Queues**:
  - `mealplanner.recipe-events`: Updates read model on recipe changes
- **Events**: `RecipeCreatedEvent`, `RecipeUpdatedEvent`
- **Go library**: `github.com/rabbitmq/amqp091-go`

### API Endpoints

**gRPC Services** (inter-service):
- `RecipeService`: GetRecipeById, GetAllRecipes, CreateRecipe
- `MealPlannerService`: SuggestRecipes

**REST API** (client-facing via MobileBFF):
- `GET /v1/recipe/{id}` - Get recipe by ID
- `GET /v1/recipe/all?pageIndex=1&pageSize=20` - Paginated recipes
- `GET /v1/recipe/similar?recipe=...&amount=5` - Similar recipes (vector search)
- `GET /v1/recipe/cuisine/{id}` - Recipes by cuisine
- `GET /v1/recipe/ingredient/{id}` - Recipes by ingredient
- `GET /v1/recipe/allergy/{id}` - Recipes avoiding allergen
- `POST /v1/recipe/create` - Create new recipe

## AI & Vector Search (Planned Enhancement)

### Current State

**Vector embeddings**: Currently using placeholder hash-based vectors in `internal/common/vector/generator.go`
- Simple TF-IDF implementation using word hashing (128-dimensional)
- **POC only** - intended for architecture validation

### Future: Azure OpenAI Integration

**Implementation Steps**:
1. Replace `HashGenerator` with Azure OpenAI embedding client
2. Use `text-embedding-ada-002` or newer model (1536 dimensions)
3. Update database schema for larger vectors
4. Implement batch embedding for seeding
5. Add embedding caching strategy

**Planned Use Cases**:
- Semantic recipe similarity search
- Natural language recipe search
- Personalized meal suggestions
- Ingredient substitution recommendations

## Database Schema

### RecipeApi Database (`recipedb`)
- **Recipes**: Main recipe table with vector embeddings
- **Ingredients**: Ingredient master list
- **Cuisines**: Cuisine types
- **Allergies**: Allergen tracking
- **Migrations**: `migrations/recipe/`

### MealPlannerApi Database (`mealplannerdb`)
- **Recipes**: Denormalized recipe data for queries
- **Migrations**: `migrations/mealplanner/`

### Seeding Data
- JSON seed file: `data/recipes.json`
- Run with: `make seed` or `recipe-api -seed data/recipes.json`
- Handles deduplication for ingredients and cuisines

## Frontend (Planned)

### Technology Choice: Vue.js + Quasar
- **Framework**: Vue.js 3 with Quasar Framework
- **Mobile**: Capacitor for native iOS/Android apps
- **Web**: PWA support out of the box
- **Status**: Not yet started

### Why Quasar?
- Single codebase for web, iOS, and Android
- 70+ Material Design components built-in
- Capacitor integration for native device APIs
- TypeScript support matches Go backend contracts
- Smaller bundle size than Flutter web

## Key Architectural Patterns

1. **CQRS**: Separate read (MealPlanner) and write (Recipe) models
2. **Event-Driven**: RabbitMQ for service communication
3. **Repository Pattern**: Data access abstraction with raw SQL
4. **Backend-for-Frontend**: BFF gateway for mobile clients
5. **Explicit over Magic**: No ORM, no DI containers, no reflection

## Development Guidelines

### Code Style

```go
// Explicit error handling
result, err := repo.GetByID(ctx, id)
if err != nil {
    return nil, fmt.Errorf("get recipe: %w", err)
}

// Interface-based dependencies
type RecipeRepository interface {
    GetByID(ctx context.Context, id uuid.UUID) (*Recipe, error)
}

// Structured logging
slog.Info("recipe created", "id", recipe.ID, "name", recipe.Name)

// Functional naming (what it does, not how)
type MealPlanner interface {           // NOT: VectorSimilaritySearcher
    SuggestRecipes(ctx, constraints)   // NOT: FindByCosineSimilarity
}
```

### Best Practices

- **Read before modifying**: Always read existing code first
- **Use explicit patterns**: No magic, no reflection, explicit error handling
- **Write tests**: Add table-driven tests for new functionality
- **SQL over ORM**: Write explicit SQL queries, avoid abstractions
- **Update protos**: Modify `api/proto/` files when changing gRPC contracts
- **Run migrations**: Use `make migrate-create` for schema changes
- **Error handling**: Wrap errors with context using `fmt.Errorf("context: %w", err)`
- **Timeouts**: Set context timeouts for all external calls
- **Graceful shutdown**: Handle SIGTERM, drain connections

### Known Limitations

- **Hash-based vectors**: POC only, real embeddings planned
- **No authentication**: Auth/authz not yet implemented
- **Frontend**: Not yet implemented (Vue/Quasar planned)

### Hobby Project Philosophy

- **Keep it simple**: Choose boring, proven solutions
- **No enterprise patterns**: Skip abstractions until needed
- **Delete freely**: No deprecation cycles, just remove
