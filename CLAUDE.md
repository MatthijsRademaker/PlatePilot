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
| Frontend | `src/frontend/` | Vue/Quasar app with vertical slice architecture |

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

### Frontend
- **Framework**: Vue.js 3 + Quasar 2
- **State**: Pinia stores
- **Styling**: UnoCSS with Wind preset (Tailwind-compatible, `tw-` prefix)
- **Mobile**: Capacitor (planned)
- **Architecture**: Feature-based vertical slices
- **Package Manager**: bun (NEVER use npm)

## Docker Compose Convention

| File | Location | Purpose |
|------|----------|---------|
| `docker-compose.yml` | Project root | **Local development** - hot reload, debug logging |
| `docker-compose.prod.yml` | Project root | Production builds - optimized images |
| `deployments/docker-compose.yml` | `src/backend-go/` | Backend-only production (legacy, prefer root files) |

**For local development**, always run from the project root:
```bash
docker compose up              # Start all services with hot reload
docker compose up -d           # Start in background
docker compose logs -f         # Follow logs
docker compose down            # Stop all services
docker compose down -v         # Stop and remove volumes
```

**For production builds**:
```bash
docker compose -f docker-compose.prod.yml up --build
```

## Build & Development Commands

```bash
cd src/backend-go

# First time setup
make tools              # Install dev tools (golangci-lint, air, protoc plugins)
make docker-up          # Start PostgreSQL + RabbitMQ (infrastructure only)
make migrate-up         # Run database migrations

# Development (local with hot reload - prefer root docker-compose)
make dev                # Run all services with hot reload (without Docker)
make dev-recipe         # Run recipe-api only
make dev-mealplanner    # Run mealplanner-api only
make dev-bff            # Run mobile-bff only

# Docker (full stack) - prefer root docker-compose.yml
make docker-run-all     # Build and run complete stack
make docker-run-detached # Run in background
make docker-logs        # View logs
make docker-down        # Stop everything

# Build & Test
make build              # Build all binaries
make test               # Run unit tests
make test-e2e           # Run E2E integration tests
make lint               # Run linter
make verify             # Verify setup + generated code is up-to-date (run in CI)

# Database
make migrate-up         # Apply all migrations
make migrate-down       # Rollback last migration
make seed               # Seed with sample recipes

# Code generation
make proto              # Generate gRPC code from protos
make verify-proto       # Check if proto code is up-to-date (fails if stale)
```

**Important**: Generated protobuf code is committed to the repo. After modifying `.proto` files, run `make proto` and commit the regenerated `*.pb.go` files. CI runs `make verify` to catch forgotten regenerations.

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

## Frontend

### Technology Stack
- **Framework**: Vue.js 3 with Quasar 2
- **State Management**: Pinia
- **Language**: TypeScript
- **Mobile**: Capacitor for iOS/Android (planned)
- **Styling**: UnoCSS with Wind preset (Tailwind-compatible)
- **Package Manager**: bun (NEVER use npm)

### Frontend Commands

```bash
cd src/frontend

# Development
bun install              # Install dependencies
bun run dev              # Start dev server (hot reload)

# Build
bun run build            # Production build
bun run lint             # Run ESLint
bun run format           # Format with Prettier
```

### UnoCSS / Tailwind CSS Integration

This project uses **UnoCSS** with the **Wind preset** for Tailwind-compatible utility classes. To avoid conflicts with Quasar's built-in classes, all Tailwind utilities use the `tw-` prefix.

#### Configuration Files

| File | Purpose |
|------|---------|
| `uno.config.ts` | UnoCSS configuration with Wind preset and `tw-` prefix |
| `src/boot/unocss.ts` | Boot file that imports UnoCSS virtual CSS |
| `quasar.config.ts` | Registers UnoCSS Vite plugin and boot file |

#### Usage Examples

```vue
<template>
  <!-- Use tw- prefix for all Tailwind utilities -->
  <div class="tw-flex tw-items-center tw-gap-4 tw-p-4">
    <span class="tw-text-lg tw-font-bold">Title</span>
  </div>

  <!-- Modifiers: Place tw- AFTER the modifier -->
  <button class="hover:tw-underline focus:tw-ring-2 md:tw-flex">
    Hover me
  </button>

  <!-- Combining with Quasar (no prefix for Quasar classes) -->
  <q-btn class="tw-mt-4" color="primary">Submit</q-btn>
</template>
```

#### Common Patterns

| Tailwind Class | UnoCSS Equivalent | Notes |
|----------------|-------------------|-------|
| `flex` | `tw-flex` | Always use prefix |
| `hover:underline` | `hover:tw-underline` | Modifier BEFORE tw- |
| `md:hidden` | `md:tw-hidden` | Responsive modifier BEFORE tw- |
| `dark:bg-gray-800` | `dark:tw-bg-gray-800` | Dark mode modifier BEFORE tw- |

#### When to Use Quasar vs UnoCSS

- **Quasar components**: Complex UI elements (buttons, dialogs, forms, tables)
- **UnoCSS/Tailwind**: Layout, spacing, typography, custom styling
- **Combine them**: Apply Tailwind utilities to Quasar components for fine-tuning

```vue
<template>
  <!-- Quasar component with Tailwind spacing/layout -->
  <q-card class="tw-max-w-md tw-mx-auto tw-shadow-lg">
    <q-card-section class="tw-space-y-4">
      <h2 class="tw-text-xl tw-font-semibold">Card Title</h2>
    </q-card-section>
  </q-card>
</template>
```

### Vertical Slice Architecture

The frontend follows a **feature-based vertical slice** pattern where each feature is self-contained with its own types, API, store, composables, components, and pages.

```
src/frontend/src/
├── features/                    # Feature modules (vertical slices)
│   ├── recipe/                  # Recipe feature
│   │   ├── api/                 # API calls (recipeApi.ts)
│   │   ├── components/          # Feature components (RecipeCard, RecipeList)
│   │   ├── composables/         # Vue composables (useRecipeList, useRecipeDetail)
│   │   ├── pages/               # Route pages (RecipeListPage, RecipeDetailPage)
│   │   ├── store/               # Pinia store (recipeStore.ts)
│   │   ├── types/               # TypeScript types (Recipe, Ingredient)
│   │   ├── routes.ts            # Feature routes
│   │   └── index.ts             # Barrel export
│   ├── mealplan/                # Meal planning feature
│   │   ├── api/
│   │   ├── components/          # WeekView, MealSlotCard
│   │   ├── composables/
│   │   ├── pages/
│   │   ├── store/
│   │   ├── types/
│   │   ├── routes.ts
│   │   └── index.ts
│   ├── search/                  # Search feature
│   │   ├── pages/
│   │   ├── types/
│   │   ├── routes.ts
│   │   └── index.ts
│   └── home/                    # Home/dashboard feature
│       ├── pages/
│       ├── routes.ts
│       └── index.ts
├── shared/                      # Shared/common code
│   ├── api/                     # HTTP client (apiClient)
│   ├── components/              # Shared UI components
│   ├── composables/             # Shared composables
│   ├── types/                   # Shared types (pagination)
│   └── utils/                   # Utility functions
├── layouts/                     # App layouts (MainLayout)
├── router/                      # Vue Router setup
├── stores/                      # Pinia setup
├── boot/                        # Quasar boot files
└── i18n/                        # Internationalization
```

### Feature Slice Structure

Each feature follows this pattern:

```
feature/
├── types/           # Domain types and DTOs
│   └── index.ts     # Barrel export
├── api/             # API layer - calls to backend
│   └── featureApi.ts
├── store/           # Pinia store - state management
│   └── featureStore.ts
├── composables/     # Vue composables - reusable logic
│   └── useFeature.ts
├── components/      # Feature-specific components
│   └── FeatureCard.vue
├── pages/           # Route pages
│   └── FeaturePage.vue
├── routes.ts        # Feature route definitions
└── index.ts         # Public API (barrel export)
```

### Key Frontend Patterns

```typescript
// Feature barrel export (index.ts)
export * from './types';
export { featureApi } from './api';
export { useFeatureStore } from './store';
export { useFeature } from './composables';
export { FeatureCard } from './components';
export { featureRoutes } from './routes';

// Composable pattern
export function useRecipeList() {
  const store = useRecipeStore();
  const { recipes, loading, error } = storeToRefs(store);

  onMounted(() => store.fetchRecipes());

  return { recipes, loading, error };
}

// Store pattern (Pinia composition API)
export const useRecipeStore = defineStore('recipe', () => {
  const recipes = ref<Recipe[]>([]);
  const loading = ref(false);

  async function fetchRecipes() { /* ... */ }

  return { recipes, loading, fetchRecipes };
});
```

### Frontend Best Practices

- **Import from feature index**: `import { Recipe, useRecipeStore } from '@/features/recipe'`
- **Keep features isolated**: Features should not import from other features' internal modules
- **Shared code in shared/**: Cross-feature utilities go in `shared/`
- **Composables for logic**: Extract reusable logic into composables
- **Types first**: Define types before implementing API/store
- **Barrel exports**: Use index.ts to control public API

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
- **Update protos**: Modify `api/proto/` files when changing gRPC contracts
- **Run migrations**: Use `make migrate-create` for schema changes
- **Error handling**: Wrap errors with context using `fmt.Errorf("context: %w", err)`
- **Timeouts**: Set context timeouts for all external calls
- **Graceful shutdown**: Handle SIGTERM, drain connections

### Known Limitations

- **Hash-based vectors**: POC only, real embeddings planned
- **No authentication**: Auth/authz not yet implemented
- **Mobile apps**: Capacitor integration not yet configured

### Hobby Project Philosophy

- **Delete freely**: No deprecation cycles, just remove

## Pommel - Semantic Code Search

This project uses Pommel for semantic code search. Pommel indexes your codebase into semantic chunks (files, classes, methods) and enables natural language search to find relevant code quickly.

**Supported platforms:** macOS, Linux, Windows
**Supported languages** (full AST-aware chunking): Go, Java, C#, Python, JavaScript, TypeScript, JSX, TSX

### Code Search Priority

**IMPORTANT: Use `pm search` BEFORE using Grep/Glob for code exploration.**

When looking for:
- How something is implemented → `pm search "authentication flow"`
- Where a pattern is used → `pm search "error handling"`
- Related code/concepts → `pm search "database connection"`
- Code that does X → `pm search "validate user input"`

Only fall back to Grep/Glob when:
- Searching for an exact string literal (e.g., a specific error message)
- Looking for a specific identifier name you already know
- Pommel daemon is not running

### Quick Search Examples
```bash
# Find code by semantic meaning (not just keywords)
pm search "authentication logic"
pm search "error handling patterns"
pm search "database connection setup"

# Search with JSON output for programmatic use
pm search "user validation" --json

# Limit results
pm search "API endpoints" --limit 5

# Search specific chunk levels
pm search "class definitions" --level class
pm search "function implementations" --level method
```

### Available Commands
- `pm search <query>` - Semantic search across the codebase
- `pm status` - Check daemon status and index statistics
- `pm reindex` - Force a full reindex of the codebase
- `pm start` / `pm stop` - Control the background daemon

### Tips
- Use natural language queries - Pommel understands semantic meaning
- Keep the daemon running (`pm start`) for always-current search results
- Use `--json` flag when you need structured output for processing
- Chunk levels: file (entire files), class (structs/interfaces/classes), method (functions/methods)
