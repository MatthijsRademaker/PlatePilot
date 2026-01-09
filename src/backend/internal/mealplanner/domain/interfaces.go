package domain

import (
	"context"

	"github.com/google/uuid"
)

// RecipeRepository defines the repository operations needed by the planner
type RecipeRepository interface {
	GetAll(ctx context.Context, userID uuid.UUID, limit, offset int) ([]Recipe, error)
}
