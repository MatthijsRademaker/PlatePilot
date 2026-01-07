# Changelog

## User auth setup

### Architecture
- Added user ownership for recipes and mealplanner read model; all recipe queries are now user-scoped.
- Added auth subsystem in Mobile BFF with JWT access tokens and refresh tokens stored in Postgres.
- Added `recipe_shares` table for future sharing between users.

### Auth Flow
- Register/login returns `accessToken` (JWT), `refreshToken`, `expiresIn`, `tokenType`.
- Clients send `Authorization: Bearer <accessToken>` to all recipe and mealplan endpoints.
- Refresh rotates refresh tokens; logout revokes refresh token.

### New/Updated REST Endpoints (Mobile BFF)
- `POST /v1/auth/register`
- `POST /v1/auth/login`
- `POST /v1/auth/refresh`
- `POST /v1/auth/logout`
- All `/v1/recipe/*` and `/v1/mealplan/*` endpoints now require `Authorization`.

### gRPC Changes
- Added `user_id` to recipe and mealplanner request messages to enforce scoping.

### Seed User (for local/dev)
- Default seed user: `seed@platepilot.local` / `platepilot`
- Override via `PLATEPILOT_SEED_USER_EMAIL` and `PLATEPILOT_SEED_USER_PASSWORD`.

### Notes
- Schema changes are breaking; wipe local volumes before bringing up the stack.
