package domain

import (
	"context"

	"github.com/platepilot/backend/internal/mealplanner/repository"
)

// RecipeRepository defines the repository operations needed by the planner
type RecipeRepository interface {
	GetAll(ctx context.Context, limit, offset int) ([]repository.Recipe, error)
}
