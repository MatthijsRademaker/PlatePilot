# Multi-stage Dockerfile for PlatePilot Go services
# Builds small, secure Alpine-based images

# =============================================================================
# Build stage - compiles all services
# =============================================================================
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build all services with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /bin/recipe-api ./cmd/recipe-api

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /bin/mealplanner-api ./cmd/mealplanner-api

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /bin/mobile-bff ./cmd/mobile-bff

# =============================================================================
# Recipe API image
# =============================================================================
FROM alpine:3.20 AS recipe-api

# Add non-root user for security
RUN addgroup -g 1000 app && adduser -u 1000 -G app -s /bin/sh -D app

# Install CA certificates for HTTPS calls
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy binary
COPY --from=builder /bin/recipe-api /usr/local/bin/recipe-api

# Copy migrations for the recipe service
COPY --from=builder /app/migrations/recipe /app/migrations

# Copy seed data
COPY --from=builder /app/data /app/data

# Set ownership
RUN chown -R app:app /app

USER app

# gRPC port
EXPOSE 50051

# Health check using grpc_health_probe (install separately if needed)
# For now, we rely on the gRPC health service

ENTRYPOINT ["recipe-api"]

# =============================================================================
# MealPlanner API image
# =============================================================================
FROM alpine:3.20 AS mealplanner-api

# Add non-root user for security
RUN addgroup -g 1000 app && adduser -u 1000 -G app -s /bin/sh -D app

# Install CA certificates for HTTPS calls
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy binary
COPY --from=builder /bin/mealplanner-api /usr/local/bin/mealplanner-api

# Copy migrations for the mealplanner service
COPY --from=builder /app/migrations/mealplanner /app/migrations

# Set ownership
RUN chown -R app:app /app

USER app

# gRPC port
EXPOSE 50052

ENTRYPOINT ["mealplanner-api"]

# =============================================================================
# Mobile BFF image
# =============================================================================
FROM alpine:3.20 AS mobile-bff

# Add non-root user for security
RUN addgroup -g 1000 app && adduser -u 1000 -G app -s /bin/sh -D app

# Install CA certificates for HTTPS calls
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy binary
COPY --from=builder /bin/mobile-bff /usr/local/bin/mobile-bff

# Set ownership
RUN chown -R app:app /app

USER app

# HTTP port
EXPOSE 8080

# HTTP health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

ENTRYPOINT ["mobile-bff"]
