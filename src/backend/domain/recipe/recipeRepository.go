package recipe

import (
	"context"

	"github.com/google/uuid"
)

type RecipeRepository interface {
	Save(ctx context.Context, recipe *Recipe) error
	FindById(ctx context.Context, id uuid.UUID) (*Recipe, error)
	FindByName(ctx context.Context, name string) ([]Recipe, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
