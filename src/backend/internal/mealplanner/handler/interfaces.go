package handler

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/platepilot/backend/internal/mealplanner/domain"
)

// MealPlanner defines the planning operations needed by the handler
type MealPlanner interface {
	SuggestMeals(ctx context.Context, req domain.SuggestionRequest) ([]uuid.UUID, error)
}

// MealPlanStore defines persistence operations for week plans.
type MealPlanStore interface {
	GetWeekPlan(ctx context.Context, userID uuid.UUID, startDate time.Time) (*domain.WeekPlan, error)
	UpsertWeekPlan(ctx context.Context, plan domain.WeekPlan) (*domain.WeekPlan, error)
}
