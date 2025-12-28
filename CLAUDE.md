# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**PlatePilot** is an intelligent meal planning and recipe management application built with microservices architecture. The application uses AI-powered recommendations (planned) to suggest personalized meal combinations based on user preferences, dietary restrictions, and recipe similarity using vector search.

**This is a hobby project** - no backwards compatibility requirements, no legacy constraints. We can make breaking changes freely and choose the simplest solutions.

## Migration Status: .NET → Go

**IMPORTANT**: This project is actively being migrated from .NET to Go. See `MIGRATION_PLAN.md` for the detailed migration plan, phases, and progress.

### Why Go?
- **AI Agent Friendly**: Explicit control flow, no magic/reflection, errors as values
- **Simpler Deployment**: Single binaries, smaller containers (<20MB vs ~200MB)
- **Lower Resource Usage**: ~50MB memory vs ~150-200MB per service
- **Faster Cold Starts**: <500ms vs 2-3s for .NET

### Key Migration Decisions
- **No backwards compatibility**: Break anything, delete freely, no migration paths needed
- **RabbitMQ over Azure Service Bus**: Simpler, better local dev experience, no emulator issues
- **Delete .NET code when Go equivalent works**: No need to maintain both
- **Functional naming over technical naming**: Name types/packages by what they do, not how they do it (e.g., `mealplanner` not `vectordbsearch`, `SuggestRecipes` not `CosineSimilaritySearch`). This keeps code readable for LLMs and humans alike.

### Current State

| Component | .NET (Legacy) | Go (Target) | Status |
|-----------|---------------|-------------|--------|
| Recipe API | `src/backend/RecipeApi/` | `src/backend-go/cmd/recipe-api/` | Phase 0 Complete |
| MealPlanner API | `src/backend/MealPlannerApi/` | `src/backend-go/cmd/mealplanner-api/` | Phase 3 Complete |
| Mobile BFF | `src/backend/MobileBFF/` | `src/backend-go/cmd/mobile-bff/` | Phase 2 Complete |
| Common | `src/backend/Common/` | `src/backend-go/internal/common/` | Phase 1 Complete |
| Frontend | Flutter (iOS) | Vue/Quasar | Not Started |

### Migration Phases (from MIGRATION_PLAN.md)
1. **Phase 0: Foundation** ✅ - Go project setup, tooling, Docker, migrations
2. **Phase 1: Common Layer** ✅ - Domain models, events, vector utilities, sqlc
3. **Phase 2: Mobile BFF** ✅ - REST gateway, gRPC clients, handlers
4. **Phase 3: MealPlanner API** ✅ - Read model, gRPC server, RabbitMQ consumer
5. **Phase 4: Recipe API** - Write model (most complex)
6. **Phase 5: Integration** - E2E testing, cleanup

## Technology Stack

### Target Stack (Go) - PREFERRED FOR NEW WORK
- **Language**: Go 1.21+
- **Web Framework**: `chi` router
- **gRPC**: `google.golang.org/grpc`
- **Database**: `pgx` + raw SQL (no ORM magic)
- **Migrations**: `golang-migrate`
- **Configuration**: `viper`
- **Logging**: `slog` (stdlib)
- **Vector Search**: `pgvector-go`

### Legacy Stack (.NET) - REFERENCE ONLY
- **Backend**: .NET 9.0 (C# 12+) with ASP.NET Core
- **Orchestration**: .NET Aspire for local development
- **ORM**: Entity Framework Core 9.0
- **CQRS**: MediatR 12.4.1

### Shared Infrastructure
- **Databases**: PostgreSQL with pgvector extension
- **Messaging**: RabbitMQ (replacing Azure Service Bus)
- **Inter-service Communication**: gRPC (proto3)
- **Frontend**: Vue/Quasar (planned)

## Build & Development Commands

### Go Backend (Target) - Use This

```bash
cd src/backend-go

# First time setup
make tools              # Install dev tools (golangci-lint, air, protoc plugins)
make docker-up          # Start PostgreSQL with pgvector
make migrate-up         # Run database migrations

# Development
make dev                # Run all services with hot reload
make build              # Build all service binaries
make test               # Run tests
make lint               # Run linter
make proto              # Generate gRPC code from protos

# Individual services
make dev-recipe         # Run recipe-api with hot reload
make dev-mealplanner    # Run mealplanner-api with hot reload
make dev-bff            # Run mobile-bff with hot reload

# Database
make migrate-up         # Apply all migrations
make migrate-down       # Rollback last migration
make migrate-create name=add_users service=recipe  # Create new migration
```

### .NET Backend (Legacy) - Reference Only

```bash
# Run all services with .NET Aspire orchestration
dotnet run --project src/backend/Hosting/Hosting.csproj

# Database migrations (EF Core)
cd src/backend/RecipeApi/Infrastructure
dotnet ef migrations add <MigrationName> --startup-project ../Application --context RecipeContext
```

### Testing

```bash
# Go tests
cd src/backend-go
make test               # Unit tests
make test-coverage      # With coverage report
make test-integration   # Integration tests (requires Docker)
```

## Architecture Overview

### Microservices Structure

The backend follows **CQRS pattern** with event-driven architecture:

1. **Recipe API** (Write Service)
   - Handles recipe creation, updates, and command operations
   - PostgreSQL database with pgvector for embeddings
   - Publishes events to Azure Service Bus on recipe changes
   - gRPC service for inter-service communication
   - Go: `src/backend-go/cmd/recipe-api/` + `internal/recipe/`
   - .NET (legacy): `src/backend/RecipeApi/`

2. **MealPlanner API** (Read Service)
   - Handles meal planning queries and recipe suggestions
   - Denormalized read model for query performance
   - Vector similarity search using pgvector
   - Consumes recipe events to update read model
   - Go: `src/backend-go/cmd/mealplanner-api/` + `internal/mealplanner/`
   - .NET (legacy): `src/backend/MealPlannerApi/`

3. **Mobile BFF** (Backend-for-Frontend)
   - REST API gateway for mobile/web clients
   - Aggregates calls to RecipeApi and MealPlannerApi via gRPC
   - OpenAPI/Swagger documentation
   - Go: `src/backend-go/cmd/mobile-bff/` + `internal/bff/`
   - .NET (legacy): `src/backend/MobileBFF/`

### Go Project Structure (Target)

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

## Azure OpenAI & LLM Integration (Planned)

### Current State

**Vector embeddings**: Currently using placeholder hash-based vectors in `src/backend/Common/Recipe/VectorExtensions.cs:14-27`
- Simple TF-IDF implementation using word hashing
- 128-dimensional vectors
- **Not production-ready** - intended for POC/architecture validation

### Future Integration Plan

**Commented package references** in `src/backend/RecipeApi/Application/RecipeApplication.csproj:11-13`:
```xml
<!-- <PackageReference Include="Aspire.Azure.AI.OpenAI" />
<PackageReference Include="Microsoft.Extensions.AI" />
<PackageReference Include="Microsoft.Extensions.AI.OpenAI" /> -->
```

**Implementation Steps** (when ready):
1. Uncomment Azure OpenAI package references
2. Configure Azure OpenAI service in Aspire host
3. Replace `VectorExtensions.GenerateVector()` with Azure OpenAI embedding API calls
4. Use `text-embedding-ada-002` or newer model for recipe embeddings
5. Update vector dimensions in database schema if needed (currently 128)
6. Consider batch embedding for seed data and new recipes
7. Implement embedding caching strategy

**Planned Use Cases**:
- Recipe similarity search based on semantic understanding
- Ingredient substitution suggestions
- Personalized meal plan generation based on user history
- Natural language recipe search
- Dietary restriction analysis and suggestions

**Configuration Notes**:
- Use Azure OpenAI service endpoints (not OpenAI direct)
- Store API keys in Azure Key Vault or Aspire secrets
- Consider embedding caching to reduce API costs
- Implement retry policies for API resilience

## Database Schema

### RecipeApi Database (`recipedb`)
- **Recipes**: Main recipe table with vector embeddings
- **Ingredients**: Ingredient master list
- **Cuisines**: Cuisine types
- **Allergies**: Allergen tracking
- **Migrations**: `src/backend/RecipeApi/Infrastructure/Migrations/`

### MealPlannerApi Database (`mealplannerdb`)
- **Materialized View**: Denormalized recipe data for queries
- **Migrations**: `src/backend/MealPlannerApi/Infrastructure/Migrations/`

### Seeding Data
- JSON seed file: `src/backend/RecipeApi/Application/DatabaseSeed/recipes.json`
- Auto-seeded on application startup
- Handles deduplication for ingredients and cuisines

## Frontend

### Current Implementation (Flutter)
- **Location**: `src/frontend/IosApp/plate_pilot/`
- **Platform**: iOS/macOS
- **Status**: Early POC with basic UI structure
- **Features**:
  - 4-tab navigation: Home, Mealplan, All Recipes, Search
  - Basic recipe display components
  - HTTP client for backend communication

### Planned Migration
- **Target**: Quasar/Vue.js framework
- **Goal**: Cross-platform web app with native iOS support
- **Status**: Not yet started

## Key Architectural Patterns

1. **CQRS**: Separate read (MealPlanner) and write (Recipe) models
2. **Event Sourcing Foundation**: Event-driven updates via Service Bus
3. **Clean Architecture**: Application/Domain/Infrastructure separation
4. **Repository Pattern**: Data access abstraction with EF Core
5. **Mediator Pattern**: MediatR for command/query handling
6. **Backend-for-Frontend**: BFF gateway for mobile clients
7. **Domain-Driven Design**: Rich domain models with aggregates

## Important Development Notes

### When Working on Go Code (Preferred)

- **Consult MIGRATION_PLAN.md**: Check the migration plan for context and current phase
- **Read before modifying**: Always read existing Go and .NET files for reference
- **Use explicit patterns**: No magic, no reflection, explicit error handling
- **Write tests**: Add table-driven tests for new functionality
- **SQL over ORM**: Write explicit SQL queries, avoid abstractions
- **Update protos**: Modify `api/proto/` files when changing gRPC contracts
- **Run migrations**: Use `make migrate-create` for schema changes

### Go Code Style

```go
// Good: Explicit error handling
result, err := repo.GetByID(ctx, id)
if err != nil {
    return nil, fmt.Errorf("get recipe: %w", err)
}

// Good: Interface-based dependencies
type RecipeRepository interface {
    GetByID(ctx context.Context, id uuid.UUID) (*Recipe, error)
}

// Good: Structured logging
slog.Info("recipe created", "id", recipe.ID, "name", recipe.Name)

// Good: Functional naming (what it does, not how)
type MealPlanner interface {           // NOT: VectorSimilaritySearcher
    SuggestRecipes(ctx, constraints)   // NOT: FindByCosineSimilarity
}
func (p *Planner) SuggestRecipes(...)  // NOT: QueryPgvectorIndex
```

### When Referencing .NET Code (Legacy)

- **.NET is reference only**: Use to understand business logic, then reimplement in Go
- **Don't modify .NET**: Focus all new development on Go codebase
- **Pattern translation**: MediatR handlers → Go handler functions, EF Core → raw SQL

### Code Quality Considerations

- **Security**: Parameterized queries, input validation at boundaries
- **Vector operations**: Hash-based vectors are POC - plan for Azure OpenAI embeddings
- **Error handling**: Wrap errors with context using `fmt.Errorf("context: %w", err)`
- **Timeouts**: Set context timeouts for all external calls
- **Graceful shutdown**: Handle SIGTERM, drain connections

### Local Development

```bash
# Start infrastructure
cd src/backend-go
make docker-up          # PostgreSQL with pgvector
make migrate-up         # Apply migrations

# Run services
make dev                # All services with hot reload

# Or run individually
make dev-bff            # Just the BFF on :8080
```

### Known Limitations

- **Go migration in progress**: Not all functionality ported yet
- **Hash-based vectors**: POC only, real embeddings planned
- **No authentication**: Auth/authz not yet implemented
- **Frontend**: Flutter POC exists, Vue migration not started

### Hobby Project Philosophy

- **Keep it simple**: Choose boring, proven solutions
- **No enterprise patterns**: Skip abstractions until needed
- **Delete freely**: No deprecation cycles, just remove
- **Learn by doing**: Experiment with Go patterns
