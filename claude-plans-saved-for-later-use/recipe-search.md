# Recipe Search Feature Implementation Plan

## Overview
Implement a hybrid search feature combining text search (PostgreSQL full-text) and semantic search (Ollama embeddings) for recipes.

## API Design

**Endpoint:** `GET /v1/recipe/search`

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `q` | string | Yes | Search query |
| `cuisineId` | string | No | Filter by cuisine UUID |
| `maxPrepTime` | int | No | Max prep time in minutes |
| `maxCookTime` | int | No | Max cook time in minutes |
| `pageIndex` | int | No | Page number (default: 1) |
| `pageSize` | int | No | Results per page (default: 20) |

**Response:** Paginated results with relevance scores. Uses hybrid search (text + vector) automatically, falls back to text-only if LLM unavailable.

---

## Implementation Phases

### Phase 1: Database Migration
**File:** `src/backend/migrations/recipe/000003_add_text_search.up.sql`

- Add `search_tsv` generated tsvector column for full-text search
- Create GIN index for fast text search
- Combine recipe name (weight A) + description (weight B)

### Phase 2: Repository Layer
**File:** `src/backend/internal/recipe/repository/postgres.go`

Add `Search(ctx, params)` method with hybrid query:
- Combines text search (`ts_rank`) with vector similarity (`<=>`)
- Uses Reciprocal Rank Fusion (RRF) to merge rankings
- Falls back to text-only if no query vector provided

### Phase 3: gRPC Service
**Files:**
- `src/backend/api/proto/recipe/v1/recipe.proto` - Add `SearchRecipes` RPC
- `src/backend/internal/recipe/handler/grpc.go` - Handler with LLM fallback

Key logic:
1. If `semantic` or `hybrid` mode, try to generate embedding via LLM
2. If LLM fails, fall back to `text` mode and indicate in response
3. Execute appropriate search query

### Phase 4: Recipe API Integration
**File:** `src/backend/cmd/recipe-api/main.go`

- Initialize `EmbeddingGenerator` when `cfg.LLM.IsConfigured()` is true
- Pass to gRPC handler for search embedding generation
- Graceful degradation if LLM unavailable

### Phase 5: BFF REST Endpoint
**Files:**
- `src/backend/internal/bff/handler/recipe.go` - Add `Search` handler
- `src/backend/internal/bff/client/recipe.go` - Add `Search` client method
- `src/backend/cmd/mobile-bff/main.go` - Add route `/recipe/search`

### Phase 6: OpenAPI & Code Generation
- Regenerate `swagger.yaml` with `make swagger`
- Run `bun run orval` in frontend to generate new API hooks

### Phase 7: Frontend
**Files:**
- `src/frontend/src/features/search/composables/useSearch.ts` - New search composable
- `src/frontend/src/features/search/pages/SearchPage.vue` - Replace client-side filtering

---

## Key Files to Modify

| File | Changes |
|------|---------|
| `src/backend/migrations/recipe/000003_add_text_search.up.sql` | NEW - FTS column + index |
| `src/backend/api/proto/recipe/v1/recipe.proto` | Add SearchRecipes RPC |
| `src/backend/internal/recipe/repository/postgres.go` | Add Search() method |
| `src/backend/internal/recipe/handler/grpc.go` | Add SearchRecipes handler |
| `src/backend/cmd/recipe-api/main.go` | Wire EmbeddingGenerator |
| `src/backend/internal/bff/handler/recipe.go` | Add Search REST handler |
| `src/backend/internal/bff/client/recipe.go` | Add Search client method |
| `src/backend/cmd/mobile-bff/main.go` | Add /recipe/search route |
| `src/backend/api/openapi/swagger.yaml` | Add search endpoint spec |
| `src/frontend/src/features/search/composables/useSearch.ts` | NEW - search composable |
| `src/frontend/src/features/search/pages/SearchPage.vue` | Server-side search |

---

## Design Decisions

1. **Single Smart Mode:** Hybrid-only API - simpler for users, combines text + vector automatically
2. **Graceful LLM Degradation:** If Ollama unavailable, automatically fall back to text-only search
3. **Hybrid Search (RRF):** Combine text + vector results using Reciprocal Rank Fusion for best quality
4. **Performance:** GIN index for text, IVFFlat for vectors, pagination limits

---

## Testing Plan
- Add E2E test for `/v1/recipe/search` endpoint
- Test all three search modes
- Test LLM fallback behavior
- Update frontend SearchPage tests

