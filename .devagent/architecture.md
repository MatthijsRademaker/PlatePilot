# PlatePilot Architecture Documentation

## Overview

PlatePilot is an intelligent meal planning and recipe management application built with a microservices architecture. The system enables users to discover, organize, and plan meals with AI-powered recommendations based on preferences, dietary restrictions, and recipe similarity using vector search.

The backend is implemented in Go following CQRS (Command Query Responsibility Segregation) with event-driven communication. The frontend is a Vue.js/Quasar application with vertical slice architecture.

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                                  CLIENTS                                         │
│                    ┌─────────────────────────────────────┐                       │
│                    │   Vue.js/Quasar Frontend (SPA)      │                       │
│                    │   - Recipe Management               │                       │
│                    │   - Meal Planning UI                │                       │
│                    │   - Search Interface                │                       │
│                    └──────────────┬──────────────────────┘                       │
│                                   │ REST/HTTP                                    │
└───────────────────────────────────┼─────────────────────────────────────────────┘
                                    │
┌───────────────────────────────────┼─────────────────────────────────────────────┐
│                                   ▼                                              │
│                    ┌─────────────────────────────────────┐                       │
│                    │      Mobile BFF (Gateway)           │                       │
│                    │      HTTP :8080                     │                       │
│                    │      - REST API Aggregation         │                       │
│                    │      - Protocol Translation         │                       │
│                    │      - CORS Handling                │                       │
│                    └───────────┬─────────────┬───────────┘                       │
│                                │             │                                   │
│                      gRPC      │             │     gRPC                          │
│                                ▼             ▼                                   │
│     ┌────────────────────────────────┐ ┌────────────────────────────────────┐   │
│     │    Recipe API (Write)          │ │    MealPlanner API (Read)          │   │
│     │    gRPC :50051                 │ │    gRPC :50052                     │   │
│     │    ─────────────────────────   │ │    ─────────────────────────       │   │
│     │    • Recipe CRUD Operations    │ │    • Recipe Queries                │   │
│     │    • Vector Embedding Gen      │ │    • Vector Similarity Search      │   │
│     │    • Event Publishing          │ │    • Meal Plan Suggestions         │   │
│     │    • Data Seeding              │ │    • Event Consumption             │   │
│     └─────────┬──────────┬───────────┘ └─────────┬──────────┬────────────────┘   │
│               │          │                       │          │                    │
│               │          │                       │          │                    │
│               ▼          │                       ▼          │                    │
│     ┌─────────────────┐  │             ┌─────────────────┐  │                    │
│     │   PostgreSQL    │  │             │   PostgreSQL    │  │                    │
│     │   (recipedb)    │  │             │ (mealplannerdb) │  │                    │
│     │   + pgvector    │  │             │   + pgvector    │  │                    │
│     └─────────────────┘  │             └─────────────────┘  │                    │
│                          │                                  │                    │
│                          │      ┌─────────────────────┐     │                    │
│                          │      │     RabbitMQ        │     │                    │
│                          └─────►│  Exchange:          │─────┘                    │
│                         publish │  recipe-events      │ consume                  │
│                                 │  (topic)            │                          │
│                                 └─────────────────────┘                          │
│                                                                                  │
└──────────────────────────────────────────────────────────────────────────────────┘
```

## Key Components

### Backend Services

| Service | Port | Responsibility |
|---------|------|----------------|
| **Mobile BFF** | HTTP 8080 | REST gateway for clients, aggregates gRPC services, handles CORS |
| **Recipe API** | gRPC 50051 | Write model - recipe CRUD, event publishing, vector generation |
| **MealPlanner API** | gRPC 50052 | Read model - queries, vector search, event consumption |

### Backend Directory Structure

```
src/backend/
├── cmd/                              # Service entry points
│   ├── recipe-api/main.go            # Recipe write service
│   ├── mealplanner-api/main.go       # MealPlanner read service
│   └── mobile-bff/main.go            # REST gateway
├── internal/
│   ├── recipe/                       # Recipe service internals
│   │   ├── domain/                   # Business entities
│   │   ├── handler/                  # gRPC handlers
│   │   ├── repository/               # PostgreSQL access
│   │   └── events/                   # RabbitMQ publisher
│   ├── mealplanner/                  # MealPlanner service internals
│   │   ├── domain/                   # Planner logic
│   │   ├── handler/                  # gRPC handlers
│   │   ├── repository/               # Read model access
│   │   └── events/                   # RabbitMQ consumer
│   ├── bff/                          # BFF service internals
│   │   ├── handler/                  # REST handlers
│   │   └── client/                   # gRPC clients
│   └── common/                       # Shared code
│       ├── config/                   # Viper configuration
│       ├── domain/                   # Shared domain types
│       ├── dto/                      # Data transfer objects
│       ├── events/                   # Event definitions
│       └── vector/                   # Vector utilities
├── api/proto/                        # Protobuf definitions
│   ├── recipe/v1/recipe.proto
│   └── mealplanner/v1/mealplanner.proto
├── migrations/                       # Database migrations
│   ├── recipe/                       # Recipe DB schema
│   └── mealplanner/                  # MealPlanner DB schema
└── deployments/                      # Docker configurations
```

### Frontend Structure (Vertical Slices)

```
src/frontend/src/
├── features/                         # Feature modules
│   ├── recipe/                       # Recipe feature
│   │   ├── types/                    # TypeScript interfaces
│   │   ├── api/                      # API calls
│   │   ├── store/                    # Pinia store
│   │   ├── composables/              # Vue composables
│   │   ├── components/               # Feature components
│   │   ├── pages/                    # Route pages
│   │   └── routes.ts                 # Feature routes
│   ├── mealplan/                     # Meal planning feature
│   ├── search/                       # Search feature
│   └── home/                         # Dashboard feature
├── shared/                           # Cross-feature code
│   ├── api/                          # HTTP client
│   ├── components/                   # Shared components
│   ├── composables/                  # Shared composables
│   └── types/                        # Shared types
├── layouts/                          # App layouts
├── router/                           # Vue Router setup
└── stores/                           # Pinia setup
```

## Data Flow

### Write Flow (Recipe Creation)

```
┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐
│  Client  │───►│   BFF    │───►│  Recipe  │───►│ Recipe   │───►│ RabbitMQ │
│          │    │  (REST)  │    │   API    │    │   DB     │    │          │
└──────────┘    └──────────┘    │ (gRPC)   │    └──────────┘    └────┬─────┘
                                │          │                         │
                                │  • Generate vector                 │
                                │  • Persist recipe                  │
                                │  • Publish event                   │
                                └──────────┘                         │
                                                                     ▼
                                ┌──────────┐    ┌──────────┐    ┌──────────┐
                                │MealPlan  │◄───│  Event   │◄───│  Queue   │
                                │   DB     │    │ Consumer │    │          │
                                └──────────┘    └──────────┘    └──────────┘
```

### Read Flow (Query/Search)

```
┌──────────┐    ┌──────────┐    ┌────────────┐    ┌──────────┐
│  Client  │───►│   BFF    │───►│ MealPlanner│───►│MealPlan  │
│          │    │  (REST)  │    │    API     │    │   DB     │
└──────────┘    └──────────┘    │   (gRPC)   │    └──────────┘
                                │            │
                                │ • Query recipes
                                │ • Vector similarity
                                │ • Apply constraints
                                │ • Diversity scoring
                                └────────────┘
```

### Event-Driven Synchronization

```
Recipe API                     RabbitMQ                      MealPlanner API
    │                              │                              │
    │  RecipeCreatedEvent          │                              │
    │─────────────────────────────►│                              │
    │  routing: recipe.created     │  Deliver to queue            │
    │                              │─────────────────────────────►│
    │                              │  mealplanner.recipe-events   │
    │                              │                              │
    │  RecipeUpdatedEvent          │                              │ Upsert to
    │─────────────────────────────►│                              │ read model
    │  routing: recipe.updated     │─────────────────────────────►│
    │                              │                              │
```

## Key Design Decisions

### 1. CQRS Pattern
- **Rationale**: Separates read and write concerns for optimized query performance
- **Implementation**: Recipe API handles writes, MealPlanner API handles reads
- **Trade-off**: Eventual consistency between models (acceptable for this use case)

### 2. Event-Driven Architecture
- **Rationale**: Loose coupling between services, async data synchronization
- **Implementation**: RabbitMQ with topic exchange, durable queues
- **Events**: `RecipeCreatedEvent`, `RecipeUpdatedEvent`

### 3. Vector Search with pgvector
- **Rationale**: Enable AI-powered recipe similarity without external vector DB
- **Implementation**: 128-dimensional vectors with IVFFlat indexing
- **Current**: Hash-based POC vectors (Azure OpenAI planned)

### 4. Backend-for-Frontend (BFF)
- **Rationale**: Single entry point for clients, protocol translation
- **Implementation**: Chi router, aggregates gRPC calls to REST
- **Benefits**: Simplified client logic, centralized CORS/auth

### 5. No ORM (Raw SQL)
- **Rationale**: Explicit queries, better performance control, AI-agent friendly
- **Implementation**: pgx driver with hand-written SQL
- **Benefits**: No magic, clear data access patterns

### 6. Vertical Slice Frontend Architecture
- **Rationale**: Feature isolation, maintainability, clear ownership
- **Implementation**: Each feature has its own types, API, store, components
- **Rule**: No cross-feature imports (use shared/ for common code)

### 7. Go Over .NET Migration
- **Rationale**: Simpler deployment, lower resources, faster cold starts
- **Benefits**: ~20MB containers vs ~200MB, ~50MB RAM vs ~150-200MB

## Dependencies

### Infrastructure

| Component | Technology | Purpose |
|-----------|------------|---------|
| Database | PostgreSQL + pgvector | Recipe/MealPlanner data, vector search |
| Message Broker | RabbitMQ | Event-driven communication |
| Container Runtime | Docker | Service deployment |

### Backend Libraries

| Library | Purpose |
|---------|---------|
| `chi` | HTTP router |
| `google.golang.org/grpc` | Inter-service communication |
| `pgx` | PostgreSQL driver |
| `golang-migrate` | Database migrations |
| `viper` | Configuration management |
| `slog` | Structured logging (stdlib) |
| `pgvector-go` | Vector operations |
| `amqp091-go` | RabbitMQ client |

### Frontend Libraries

| Library | Purpose |
|---------|---------|
| Vue.js 3 | UI framework |
| Quasar 2 | Component library |
| Pinia | State management |
| Vue Router | Routing |
| TypeScript | Type safety |

## Development Guidelines

### Go Patterns

```go
// Explicit error handling - always wrap with context
result, err := repo.GetByID(ctx, id)
if err != nil {
    return nil, fmt.Errorf("get recipe by id: %w", err)
}

// Interface-based dependencies
type RecipeRepository interface {
    GetByID(ctx context.Context, id uuid.UUID) (*Recipe, error)
    Create(ctx context.Context, recipe *Recipe) error
}

// Structured logging with context
slog.Info("recipe created",
    "id", recipe.ID,
    "name", recipe.Name,
    "cuisine", recipe.Cuisine.Name)

// Functional naming (what it does, not how)
type MealPlanner interface {
    SuggestRecipes(ctx context.Context, constraints Constraints) ([]Recipe, error)
}
```

### Frontend Patterns

```typescript
// Import from feature barrel exports
import { Recipe, useRecipeStore } from '@/features/recipe';

// Composable pattern for reusable logic
export function useRecipeList() {
  const store = useRecipeStore();
  const { recipes, loading, error } = storeToRefs(store);

  onMounted(() => store.fetchRecipes());

  return { recipes, loading, error };
}

// Pinia composition API store
export const useRecipeStore = defineStore('recipe', () => {
  const recipes = ref<Recipe[]>([]);
  const loading = ref(false);

  async function fetchRecipes() { /* ... */ }

  return { recipes, loading, fetchRecipes };
});
```

### API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v1/recipe/{id}` | Get recipe by ID |
| GET | `/v1/recipe/all` | Paginated recipe list |
| GET | `/v1/recipe/similar` | Vector similarity search |
| GET | `/v1/recipe/cuisine/{id}` | Filter by cuisine |
| GET | `/v1/recipe/ingredient/{id}` | Filter by ingredient |
| GET | `/v1/recipe/allergy/{id}` | Filter avoiding allergen |
| POST | `/v1/recipe/create` | Create new recipe |
| POST | `/v1/mealplan/suggest` | Get meal suggestions |

### Database Schemas

**Recipe DB (Write Model)**
```
recipes           - Main recipe table with vector embeddings
cuisines          - Cuisine types
ingredients       - Ingredient master list
allergies         - Allergen tracking
recipe_ingredients - M:N recipe-ingredient relationship
ingredient_allergies - M:N ingredient-allergy relationship
```

**MealPlanner DB (Read Model)**
```
recipes           - Denormalized recipe data
                   (arrays: ingredient_ids[], allergy_ids[], directions[], tags[])
                   GIN indexes for array queries
                   IVFFlat index for vector search
```

### Key Commands

```bash
# Development
make dev              # Run all services with hot reload
make docker-up        # Start PostgreSQL + RabbitMQ
make migrate-up       # Apply database migrations

# Build & Test
make build            # Build all binaries
make test             # Run unit tests
make lint             # Run linter

# Production
make docker-run-all   # Build and run complete stack
make seed             # Seed with sample recipes
```
