package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/platepilot/backend/internal/common/domain"
)

// RecipeRepository defines the repository operations needed by the handler
type RecipeRepository interface {
	// Recipe operations
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Recipe, error)
	GetAll(ctx context.Context, limit, offset int) ([]domain.Recipe, error)
	Create(ctx context.Context, recipe *domain.Recipe) error
	GetSimilar(ctx context.Context, recipeID uuid.UUID, limit int) ([]domain.Recipe, error)
	GetByCuisine(ctx context.Context, cuisineID uuid.UUID, limit, offset int) ([]domain.Recipe, error)
	GetByIngredient(ctx context.Context, ingredientID uuid.UUID, limit, offset int) ([]domain.Recipe, error)
	GetExcludingAllergy(ctx context.Context, allergyID uuid.UUID, limit, offset int) ([]domain.Recipe, error)

	// Ingredient operations
	GetIngredientByID(ctx context.Context, id uuid.UUID) (*domain.Ingredient, error)

	// Cuisine operations
	GetCuisineByID(ctx context.Context, id uuid.UUID) (*domain.Cuisine, error)
}

// EventPublisher defines the event publishing operations needed by the handler
type EventPublisher interface {
	PublishRecipeCreated(ctx context.Context, recipe *domain.Recipe) error
	PublishRecipeUpdated(ctx context.Context, recipe *domain.Recipe) error
}
