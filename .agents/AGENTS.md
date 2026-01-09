# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**PlatePilot** is an intelligent meal planning and recipe management application built with microservices architecture. The application uses AI-powered recommendations (planned) to suggest personalized meal combinations based on user preferences, dietary restrictions, and recipe similarity using vector search.

**This is a hobby project** - no backwards compatibility requirements, no legacy constraints. We can make breaking changes freely and choose the simplest solutions.

## Mobile Development Strategy

### iOS-Native Approach

PlatePilot is now a **native iOS application** built with SwiftUI:

| Platform | Technology | Status |
|----------|------------|--------|
| iOS/iPadOS | SwiftUI | **Primary platform** (iOS 26+) |
| watchOS | SwiftUI | Planned companion app |
| Web (Vue.js) | Vue.js 3 + Quasar | **DEPRECATED** - No longer maintained |

**Why iOS-native:**
- **Liquid Glass design** - Apple's new design language is SwiftUI-native
- **Better UX** - Native performance, gestures, and platform integration
- **watchOS future** - Companion app requires native Swift anyway
- **Focused development** - Single codebase for mobile-first experience

### Development Workflow

```
1. Implement features directly in SwiftUI
2. Test on iOS Simulator and real devices
3. All features consume Mobile BFF REST endpoints (/v1/*)
```

### Project Structure

```
src/
├── ios/                   # Native iOS app (PRIMARY)
│   └── PlatePilot/        # iOS/iPadOS SwiftUI app
│       ├── App/           # App entry, navigation, routing
│       ├── Features/      # Feature modules
│       │   ├── Home/      # Dashboard with daily plan
│       │   ├── Recipes/   # Recipe browsing and detail
│       │   ├── MealPlan/  # Weekly meal planning
│       │   ├── Search/    # Search interface
│       │   ├── Insights/  # Analytics and stats
│       │   └── Auth/      # Authentication flow
│       ├── Shared/        # Shared code
│       │   ├── API/       # REST client for BFF
│       │   ├── Models/    # Swift models (matches Go domain)
│       │   ├── Components/# Reusable UI components
│       │   ├── Extensions/# Swift extensions
│       │   └── Utils/     # Utility functions
│       └── Resources/     # Assets, colors, fonts
├── frontend/              # DEPRECATED - Vue.js (no longer maintained)
└── backend/               # Go services (unchanged)
```

### iOS Tech Stack
- **Language**: Swift 6+
- **UI Framework**: SwiftUI (Liquid Glass on iOS 26+)
- **Networking**: URLSession + async/await
- **State**: @Observable (modern Swift observation)
- **Architecture**: Feature-based modules
- **Project Generation**: XcodeGen (never edit .xcodeproj directly)
- **Storage**: Keychain for secure credential storage

### Current iOS Features (Implemented)

The iOS app currently includes:
1. ✅ Recipe browsing with cards and list views
2. ✅ Recipe detail view with ingredients and directions
3. ✅ Meal planning with weekly view
4. ✅ Home dashboard with today's plan and calorie tracker
5. ✅ Search interface
6. ✅ Authentication flow (sign in/sign up)
7. ✅ Shopping list (recent feature)

### Future Features
- watchOS companion app: Quick recipe reference, cooking timers
- Recipe creation with image upload
- Enhanced calorie tracking integration

### Notes
- `src-capacitor/` and `src/frontend/` are deprecated
- iOS deployment requires Xcode, Apple Developer account ($99/year)
- Use XcodeGen to generate the Xcode project from `project.yml`
- Claude maintains the Swift codebase; user runs Xcode builds

### Why Go?
- **AI Agent Friendly**: Explicit control flow, no magic/reflection, errors as values
- **Simpler Deployment**: Single binaries, smaller containers (<20MB vs ~200MB)
- **Lower Resource Usage**: ~50MB memory vs ~150-200MB per service
- **Faster Cold Starts**: <500ms vs 2-3s for .NET

### Architecture

| Component | Location | Description |
|-----------|----------|-------------|
| Recipe API | `src/backend/cmd/recipe-api/` | Write service, gRPC, event publishing |
| MealPlanner API | `src/backend/cmd/mealplanner-api/` | Read service, vector search, event consuming |
| Mobile BFF | `src/backend/cmd/mobile-bff/` | REST gateway for clients |
| Common | `src/backend/internal/common/` | Shared domain, events, config |
| **Frontend (iOS)** | `src/ios/PlatePilot/` | **Native SwiftUI app for iOS/iPadOS (PRIMARY)** |
| ~~Frontend (Web)~~ | ~~`src/frontend/`~~ | ~~DEPRECATED - Vue/Quasar app (no longer maintained)~~ |

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

### Frontend (iOS Native - PRIMARY)
- **Language**: Swift 6+
- **UI Framework**: SwiftUI with Liquid Glass design
- **State Management**: @Observable macro
- **Networking**: URLSession with async/await
- **Storage**: Keychain for credentials
- **Architecture**: Feature-based modules
- **Project Generation**: XcodeGen

### ~~Frontend (Vue.js - DEPRECATED)~~
- ~~**Framework**: Vue.js 3 + Quasar 2~~
- ~~**State**: Pinia stores~~
- ~~**Styling**: UnoCSS with Wind preset~~
- ~~No longer maintained - iOS is the primary platform~~

## Docker Compose Convention

| File | Location | Purpose |
|------|----------|---------|
| `docker-compose.yml` | Project root | **Local development** - hot reload, debug logging |
| `docker-compose.prod.yml` | Project root | Production builds - optimized images |
| `deployments/docker-compose.yml` | `src/backend/` | Backend-only production (legacy, prefer root files) |

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
cd src/backend

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
src/backend/
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

## iOS Frontend

### Technology Stack
- **Language**: Swift 6+
- **UI Framework**: SwiftUI with Liquid Glass design
- **State Management**: @Observable macro (modern Swift observation)
- **Networking**: URLSession with async/await
- **Storage**: Keychain for secure credential storage
- **Project Generation**: XcodeGen (never edit .xcodeproj directly)
- **Architecture**: Feature-based modules

### iOS Development Commands

```bash
cd src/ios/PlatePilot

# Generate Xcode project (required after adding/removing files)
xcodegen generate

# Open project in Xcode
open PlatePilot.xcodeproj

# Build and run from Xcode
# - Select PlatePilot scheme
# - Choose iOS 26+ simulator or device
# - Cmd+R to build and run
```

### Feature-Based Architecture

The iOS app follows a **feature-based architecture** where each feature is self-contained:

```
src/ios/PlatePilot/
├── App/                         # App entry and navigation
│   ├── PlatePilotApp.swift      # App entry point (@main)
│   ├── AppState.swift           # Global observable app state
│   ├── RootView.swift           # Root navigation container
│   ├── AppView.swift            # Main tab view
│   ├── Router.swift             # Navigation router
│   └── AppTab.swift             # Tab definitions
├── Features/                    # Feature modules
│   ├── Home/                    # Dashboard
│   ├── Recipes/                 # Recipe browsing and detail
│   ├── MealPlan/                # Meal planning (week view)
│   ├── Search/                  # Search interface
│   ├── Insights/                # Analytics and stats
│   └── Auth/                    # Authentication flow
├── Shared/                      # Shared code
│   ├── API/                     # REST API client
│   ├── Models/                  # Data models (matches Go domain)
│   ├── Components/              # Reusable UI components
│   ├── Extensions/              # Swift extensions
│   └── Utils/                   # Utility functions
└── Resources/                   # Assets, colors, fonts
```

### Key iOS/SwiftUI Patterns

```swift
// Observable view model pattern
@Observable
class RecipeListViewModel {
    var recipes: [Recipe] = []
    var isLoading = false
    var errorMessage: String?

    func fetchRecipes() async {
        isLoading = true
        defer { isLoading = false }

        do {
            recipes = try await APIClient.shared.getRecipes()
        } catch {
            errorMessage = error.localizedDescription
        }
    }
}

// SwiftUI view with async data loading
struct RecipeListView: View {
    @State private var viewModel = RecipeListViewModel()

    var body: some View {
        List(viewModel.recipes) { recipe in
            RecipeCardView(recipe: recipe)
        }
        .task {
            await viewModel.fetchRecipes()
        }
        .overlay {
            if viewModel.isLoading {
                ProgressView()
            }
        }
    }
}

// API client with async/await
actor APIClient {
    static let shared = APIClient()

    func getRecipes() async throws -> [Recipe] {
        let url = APIConfig.defaultBaseURL
            .appendingPathComponent("recipe/all")
        let (data, _) = try await URLSession.shared.data(from: url)
        return try JSONDecoder().decode([Recipe].self, from: data)
    }
}
```

### iOS Best Practices

- **Use XcodeGen**: Never edit `.xcodeproj` directly, always use `project.yml`
- **Feature isolation**: Features should be self-contained modules
- **Shared code in Shared/**: Cross-feature utilities go in `Shared/`
- **Observable for state**: Use `@Observable` macro for view models
- **Async/await for networking**: Always use modern concurrency
- **Models match backend**: Swift models mirror Go domain structs
- **Test on real devices**: Validate networking with physical iPhone when possible

### ~~Vue.js Frontend (DEPRECATED)~~

The `src/frontend/` directory contains a deprecated Vue.js/Quasar application that is **no longer maintained**. All frontend development is now focused on the native iOS app.

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

- **Hash-based vectors**: POC only, Azure OpenAI embeddings planned
- **E2E testing**: Legacy Playwright tests for deprecated Vue frontend need rewriting for iOS (XCUITest)
- **watchOS app**: Not yet implemented (planned for future)

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
