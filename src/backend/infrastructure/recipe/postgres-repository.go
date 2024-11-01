package recipeRepository

import (
	"PlatePilot/domain/recipe"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresRecipeRepository struct {
	db *gorm.DB
}

func NewPostgresRecipeRepository(db *gorm.DB) recipe.RecipeRepository {
	return &PostgresRecipeRepository{db: db}
}

func (r *PostgresRecipeRepository) Save(ctx context.Context, recipe *recipe.Recipe) error {
	return r.db.WithContext(ctx).Create(recipe).Error
}

func (r *PostgresRecipeRepository) FindById(ctx context.Context, id uuid.UUID) (*recipe.Recipe, error) {
	var recipe recipe.Recipe
	r.db.WithContext(ctx).First(&recipe, id)

	// TODO return error?
	return &recipe, nil
}

func (r *PostgresRecipeRepository) FindByName(ctx context.Context, name string) ([]recipe.Recipe, error) {
	var recipes []recipe.Recipe
	r.db.WithContext(ctx).Where(&recipe.Recipe{Name: name}).Find(&recipes)

	// TODO return error?
	return recipes, nil
}

func (r *PostgresRecipeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(id).Error
}
