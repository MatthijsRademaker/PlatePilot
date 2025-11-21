# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**PlatePilot** is an intelligent meal planning and recipe management application built with microservices architecture. The application uses AI-powered recommendations (planned) to suggest personalized meal combinations based on user preferences, dietary restrictions, and recipe similarity using vector search.

**Current Status**: POC phase focused on architecture - LLM/embedding functionality is planned but not yet active. The system uses simple hash-based vectors as a placeholder for future Azure OpenAI embeddings.

**Frontend**: Currently implemented in Flutter for iOS. Plan to migrate to Quasar/Vue for cross-platform web and native iOS support.

## Technology Stack

- **Backend**: .NET 9.0 (C# 12+) with ASP.NET Core
- **Orchestration**: .NET Aspire for local development
- **Databases**: PostgreSQL with pgvector extension for vector similarity search
- **Messaging**: Azure Service Bus (emulated locally)
- **Inter-service Communication**: gRPC (proto3)
- **Mobile API Gateway**: REST API via MobileBFF
- **ORM**: Entity Framework Core 9.0
- **CQRS**: MediatR 12.4.1
- **Frontend**: Flutter ^3.6.0 (current), Quasar/Vue (planned)

## Build & Development Commands

### Backend Development

```bash
# Run all services with .NET Aspire orchestration
make run-backend

# Run with hot reload (watch mode)
make run-backend-watch

# Or run directly
dotnet run --project src/backend/Hosting/Hosting.csproj
dotnet watch --project src/backend/Hosting/Hosting.csproj
```

### Database Migrations

```bash
# Create new migration for RecipeApi
cd src/backend/RecipeApi/Infrastructure
dotnet ef migrations add <MigrationName> --startup-project ../Application --context RecipeContext

# Create new migration for MealPlannerApi
cd src/backend/MealPlannerApi/Infrastructure
dotnet ef migrations add <MigrationName> --startup-project ../Application --context MealPlannerContext
```

### Testing

**Note**: No test projects currently exist in the codebase. Tests should be added when implementing new features.

## Architecture Overview

### Microservices Structure

The backend follows **CQRS pattern** with event-driven architecture:

1. **RecipeApi** (Write Service)
   - Handles recipe creation, updates, and command operations
   - PostgreSQL database with pgvector for embeddings
   - Publishes events to Azure Service Bus on recipe changes
   - gRPC service for inter-service communication
   - Located: `src/backend/RecipeApi/`

2. **MealPlannerApi** (Read Service)
   - Handles meal planning queries and recipe suggestions
   - Uses materialized view for denormalized read model
   - Vector similarity search using pgvector
   - Consumes recipe events to update read model
   - Located: `src/backend/MealPlannerApi/`

3. **MobileBFF** (Backend-for-Frontend)
   - REST API gateway for mobile/web clients
   - Aggregates calls to RecipeApi and MealPlannerApi via gRPC
   - OpenAPI/Swagger documentation
   - Located: `src/backend/MobileBFF/`

4. **Hosting** (Aspire AppHost)
   - Orchestrates all services for local development
   - Configures PostgreSQL, Azure Service Bus emulator, Redis
   - Located: `src/backend/Hosting/`

### Project Organization

Each API follows **Clean Architecture** with three layers:

- **Application**: gRPC services, API endpoints, MediatR handlers
- **Domain**: Business logic, aggregates, domain models
- **Infrastructure**: EF Core contexts, migrations, repositories, event handlers

**Common**: Shared DTOs, events, and cross-cutting concerns (`src/backend/Common/`)

### Event-Driven Communication

- **Azure Service Bus** topic: `recipe-events`
- **Subscriptions**:
  - `recipe-api`: Recipe service event processing
  - `meal-planner-api`: Updates materialized view on recipe changes
  - `bff-web-api`: Cache invalidation
- **Events**: `RecipeCreatedEvent`, `RecipeUpdatedEvent` (in `Common/Events/`)

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

### When Making Changes

- **Read before modifying**: Always read existing files before making changes
- **Maintain layer separation**: Respect Application/Domain/Infrastructure boundaries
- **Event publishing**: Publish domain events for state changes in RecipeApi
- **Event handling**: Update read models in MealPlannerApi when consuming events
- **gRPC contracts**: Update proto files when changing service contracts
- **Database changes**: Create EF migrations for schema changes

### Code Quality Considerations

- **Security**: Watch for command injection, XSS, SQL injection
- **Vector operations**: Current hash-based vectors are POC only - replace with real embeddings for production
- **Error handling**: Use MediatR pipeline behaviors for cross-cutting concerns
- **API versioning**: Use Asp.Versioning for REST endpoints
- **gRPC deadlines**: Set appropriate timeouts for inter-service calls

### .NET Aspire Local Development

All services are orchestrated through `src/backend/Hosting/Hosting.csproj`:
- PostgreSQL instances auto-provisioned with pgvector extension
- Azure Service Bus emulator with pre-configured topics/subscriptions
- Service discovery and health checks
- OpenTelemetry instrumentation for observability

### Known Limitations (POC Phase)

- No test coverage currently implemented
- Hash-based vectors instead of proper embeddings
- No Redis caching layer active yet
- MealPlanner REST endpoints not fully implemented
- Frontend is minimal POC structure
- No authentication/authorization implemented
