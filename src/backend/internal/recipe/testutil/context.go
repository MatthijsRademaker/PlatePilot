package testutil

import (
	"context"
	"io"
	"log/slog"

	"github.com/google/uuid"

	"github.com/platepilot/backend/internal/recipe/handler"
)

// TestContext contains all test dependencies for handler tests
type TestContext struct {
	Ctx       context.Context
	UserID    uuid.UUID
	Repo      *FakeRecipeRepository
	Publisher *FakeEventPublisher
	VectorGen *FakeVectorGenerator
	Handler   *handler.GRPCHandler
	Logger    *slog.Logger
}

// NewTestContext creates a new test context with all dependencies wired up
func NewTestContext() *TestContext {
	ctx := context.Background()
	userID := uuid.New()
	repo := NewFakeRecipeRepository()
	publisher := NewFakeEventPublisher()
	vectorGen := NewFakeVectorGenerator()

	// Create a silent logger for tests (writes to io.Discard)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	h := handler.NewGRPCHandler(repo, vectorGen, publisher, logger)

	return &TestContext{
		Ctx:       ctx,
		UserID:    userID,
		Repo:      repo,
		Publisher: publisher,
		VectorGen: vectorGen,
		Handler:   h,
		Logger:    logger,
	}
}

// NewTestContextWithoutPublisher creates a test context without an event publisher
// Use this to test behavior when event publishing is disabled
func NewTestContextWithoutPublisher() *TestContext {
	ctx := context.Background()
	userID := uuid.New()
	repo := NewFakeRecipeRepository()
	vectorGen := NewFakeVectorGenerator()

	// Create a silent logger for tests
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	h := handler.NewGRPCHandler(repo, vectorGen, nil, logger)

	return &TestContext{
		Ctx:       ctx,
		UserID:    userID,
		Repo:      repo,
		Publisher: nil, // No publisher
		VectorGen: vectorGen,
		Handler:   h,
		Logger:    logger,
	}
}
