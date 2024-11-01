package recipe

import (
	"context"
)

type RecipeRepository interface {
	Save(ctx context.Context, recipe *Recipe) error
	FindById(ctx context.Context, id int) (*Recipe, error)
	FindByName(ctx context.Context, name string) ([]Recipe, error)
	Delete(ctx context.Context, id int) error
}
