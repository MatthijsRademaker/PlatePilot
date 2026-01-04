package testutil

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/platepilot/backend/internal/mealplanner/domain"
	"github.com/platepilot/backend/internal/mealplanner/repository"
)

// FakeRecipeRepository is an in-memory implementation of RecipeRepository for testing
type FakeRecipeRepository struct {
	Recipes []repository.Recipe

	// Failure modes for testing error paths
	FailOnGetAll bool

	// Call tracking for assertions
	GetAllCalls []GetAllCall
}

// GetAllCall records a call to GetAll
type GetAllCall struct {
	Limit  int
	Offset int
}

// NewFakeRecipeRepository creates a new fake repository
func NewFakeRecipeRepository() *FakeRecipeRepository {
	return &FakeRecipeRepository{
		Recipes:     []repository.Recipe{},
		GetAllCalls: []GetAllCall{},
	}
}

// GetAll retrieves all recipes with pagination
func (r *FakeRecipeRepository) GetAll(ctx context.Context, limit, offset int) ([]repository.Recipe, error) {
	r.GetAllCalls = append(r.GetAllCalls, GetAllCall{Limit: limit, Offset: offset})

	if r.FailOnGetAll {
		return nil, errors.New("fake repository error")
	}

	// Apply pagination
	if offset >= len(r.Recipes) {
		return []repository.Recipe{}, nil
	}
	end := offset + limit
	if end > len(r.Recipes) {
		end = len(r.Recipes)
	}
	return r.Recipes[offset:end], nil
}

// AddRecipe adds a recipe to the fake repository for test setup
func (r *FakeRecipeRepository) AddRecipe(recipe repository.Recipe) {
	r.Recipes = append(r.Recipes, recipe)
}

// FakeMealPlanner is a fake implementation of MealPlanner for handler testing
type FakeMealPlanner struct {
	SuggestedRecipes []uuid.UUID

	// Failure modes
	FailOnSuggestMeals bool

	// Call tracking
	SuggestMealsCalls []domain.SuggestionRequest
}

// NewFakeMealPlanner creates a new fake meal planner
func NewFakeMealPlanner() *FakeMealPlanner {
	return &FakeMealPlanner{
		SuggestedRecipes:  []uuid.UUID{},
		SuggestMealsCalls: []domain.SuggestionRequest{},
	}
}

// SuggestMeals returns configured suggestions or an error
func (p *FakeMealPlanner) SuggestMeals(ctx context.Context, req domain.SuggestionRequest) ([]uuid.UUID, error) {
	p.SuggestMealsCalls = append(p.SuggestMealsCalls, req)

	if p.FailOnSuggestMeals {
		return nil, errors.New("fake planner error")
	}

	return p.SuggestedRecipes, nil
}

// SetSuggestedRecipes configures the recipes to return
func (p *FakeMealPlanner) SetSuggestedRecipes(ids ...uuid.UUID) {
	p.SuggestedRecipes = ids
}

// SuggestMealsCallCount returns the number of SuggestMeals calls
func (p *FakeMealPlanner) SuggestMealsCallCount() int {
	return len(p.SuggestMealsCalls)
}
