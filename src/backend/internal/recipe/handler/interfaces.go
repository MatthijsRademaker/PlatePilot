package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/platepilot/backend/internal/common/domain"
)

// RecipeRepository defines the repository operations needed by the handler
type RecipeRepository interface {
	// Recipe operations
	GetByID(ctx context.Context, userID, id uuid.UUID) (*domain.Recipe, error)
	GetAll(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.Recipe, error)
	Count(ctx context.Context, userID uuid.UUID) (int64, error)
	Create(ctx context.Context, recipe *domain.Recipe) error
	GetSimilar(ctx context.Context, userID, recipeID uuid.UUID, limit int) ([]domain.Recipe, error)
	GetByCuisine(ctx context.Context, userID, cuisineID uuid.UUID, limit, offset int) ([]domain.Recipe, error)
	GetByIngredient(ctx context.Context, userID, ingredientID uuid.UUID, limit, offset int) ([]domain.Recipe, error)
	GetExcludingAllergy(ctx context.Context, userID, allergyID uuid.UUID, limit, offset int) ([]domain.Recipe, error)

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
