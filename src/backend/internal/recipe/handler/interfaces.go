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
	List(ctx context.Context, userID uuid.UUID, filter domain.RecipeFilter, limit, offset int) ([]domain.Recipe, error)
	Count(ctx context.Context, userID uuid.UUID, filter domain.RecipeFilter) (int64, error)
	Create(ctx context.Context, recipe *domain.Recipe) error
	Update(ctx context.Context, recipe *domain.Recipe) error
	Delete(ctx context.Context, userID, id uuid.UUID) error
	GetSimilar(ctx context.Context, userID, recipeID uuid.UUID, limit int) ([]domain.Recipe, error)

	// Ingredient operations
	GetIngredientByID(ctx context.Context, userID, id uuid.UUID) (*domain.Ingredient, error)
	GetOrCreateIngredient(ctx context.Context, userID uuid.UUID, name string) (*domain.Ingredient, error)

	// Cuisine operations
	GetCuisineByID(ctx context.Context, userID, id uuid.UUID) (*domain.Cuisine, error)
	GetOrCreateCuisine(ctx context.Context, userID uuid.UUID, name string) (*domain.Cuisine, error)
	GetCuisines(ctx context.Context, userID uuid.UUID) ([]domain.Cuisine, error)
}

// EventPublisher defines the event publishing operations needed by the handler
type EventPublisher interface {
	PublishRecipeUpserted(ctx context.Context, recipe *domain.Recipe) error
	PublishRecipeDeleted(ctx context.Context, recipeID, userID uuid.UUID) error
}
