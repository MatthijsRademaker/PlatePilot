package handler

import (
	"context"

	"github.com/google/uuid"

	"github.com/platepilot/backend/internal/mealplanner/domain"
)

// MealPlanner defines the planning operations needed by the handler
type MealPlanner interface {
	SuggestMeals(ctx context.Context, req domain.SuggestionRequest) ([]uuid.UUID, error)
}
