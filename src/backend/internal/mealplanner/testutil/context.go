package testutil

import (
	"context"
	"io"
	"log/slog"

	"github.com/platepilot/backend/internal/mealplanner/domain"
	"github.com/platepilot/backend/internal/mealplanner/handler"
)

// HandlerTestContext contains all test dependencies for handler tests
type HandlerTestContext struct {
	Ctx     context.Context
	Planner *FakeMealPlanner
	Handler *handler.GRPCHandler
	Logger  *slog.Logger
}

// NewHandlerTestContext creates a new test context for handler testing
func NewHandlerTestContext() *HandlerTestContext {
	ctx := context.Background()
	planner := NewFakeMealPlanner()

	// Create a silent logger for tests (writes to io.Discard)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	h := handler.NewGRPCHandler(planner, logger)

	return &HandlerTestContext{
		Ctx:     ctx,
		Planner: planner,
		Handler: h,
		Logger:  logger,
	}
}

// PlannerTestContext contains all test dependencies for planner/domain tests
type PlannerTestContext struct {
	Ctx     context.Context
	Repo    *FakeRecipeRepository
	Planner *domain.Planner
	Logger  *slog.Logger
}

// NewPlannerTestContext creates a new test context for planner testing
func NewPlannerTestContext() *PlannerTestContext {
	ctx := context.Background()
	repo := NewFakeRecipeRepository()

	// Create a silent logger for tests
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	planner := domain.NewPlanner(repo)

	return &PlannerTestContext{
		Ctx:     ctx,
		Repo:    repo,
		Planner: planner,
		Logger:  logger,
	}
}
