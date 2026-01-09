package testutil

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"

	"github.com/platepilot/backend/internal/common/domain"
	"github.com/platepilot/backend/internal/recipe/repository"
)

// FakeRecipeRepository is an in-memory implementation of RecipeRepository for testing.
type FakeRecipeRepository struct {
	Recipes     map[uuid.UUID]*domain.Recipe
	Ingredients map[uuid.UUID]*domain.Ingredient
	Cuisines    map[uuid.UUID]*domain.Cuisine

	// Failure modes for testing error paths
	FailOnGetByID               bool
	FailOnList                  bool
	FailOnCount                 bool
	FailOnCreate                bool
	FailOnUpdate                bool
	FailOnDelete                bool
	FailOnGetSimilar            bool
	FailOnGetIngredientByID     bool
	FailOnGetCuisineByID        bool
	FailOnGetOrCreateIngredient bool
	FailOnGetOrCreateCuisine    bool
	FailOnGetCuisines           bool

	// Call tracking for assertions
	CreateCalls  []CreateCall
	UpdateCalls  []UpdateCall
	DeleteCalls  []uuid.UUID
	GetByIDCalls []uuid.UUID
}

// CreateCall records a call to Create.
type CreateCall struct {
	Recipe *domain.Recipe
}

// UpdateCall records a call to Update.
type UpdateCall struct {
	Recipe *domain.Recipe
}

// NewFakeRecipeRepository creates a new fake repository.
func NewFakeRecipeRepository() *FakeRecipeRepository {
	return &FakeRecipeRepository{
		Recipes:      make(map[uuid.UUID]*domain.Recipe),
		Ingredients:  make(map[uuid.UUID]*domain.Ingredient),
		Cuisines:     make(map[uuid.UUID]*domain.Cuisine),
		CreateCalls:  []CreateCall{},
		UpdateCalls:  []UpdateCall{},
		DeleteCalls:  []uuid.UUID{},
		GetByIDCalls: []uuid.UUID{},
	}
}

// GetByID retrieves a recipe by ID.
func (r *FakeRecipeRepository) GetByID(ctx context.Context, userID, id uuid.UUID) (*domain.Recipe, error) {
	r.GetByIDCalls = append(r.GetByIDCalls, id)

	if r.FailOnGetByID {
		return nil, errors.New("fake repository error")
	}

	recipe, ok := r.Recipes[id]
	if !ok || recipe.UserID != userID {
		return nil, repository.ErrRecipeNotFound
	}
	return recipe, nil
}

// List retrieves recipes with pagination.
func (r *FakeRecipeRepository) List(ctx context.Context, userID uuid.UUID, filter domain.RecipeFilter, limit, offset int) ([]domain.Recipe, error) {
	if r.FailOnList {
		return nil, errors.New("fake repository error")
	}

	recipes := make([]domain.Recipe, 0, len(r.Recipes))
	for _, recipe := range r.Recipes {
		if recipe.UserID == userID {
			recipes = append(recipes, *recipe)
		}
	}

	if offset >= len(recipes) {
		return []domain.Recipe{}, nil
	}
	end := offset + limit
	if end > len(recipes) {
		end = len(recipes)
	}
	return recipes[offset:end], nil
}

// Count returns the total number of recipes.
func (r *FakeRecipeRepository) Count(ctx context.Context, userID uuid.UUID, filter domain.RecipeFilter) (int64, error) {
	if r.FailOnCount {
		return 0, errors.New("fake repository error")
	}
	count := int64(0)
	for _, recipe := range r.Recipes {
		if recipe.UserID == userID {
			count++
		}
	}
	return count, nil
}

// Create creates a new recipe.
func (r *FakeRecipeRepository) Create(ctx context.Context, recipe *domain.Recipe) error {
	r.CreateCalls = append(r.CreateCalls, CreateCall{Recipe: recipe})

	if r.FailOnCreate {
		return errors.New("fake repository error")
	}

	if recipe.ID == uuid.Nil {
		recipe.ID = uuid.New()
	}

	r.Recipes[recipe.ID] = recipe
	return nil
}

// Update updates an existing recipe.
func (r *FakeRecipeRepository) Update(ctx context.Context, recipe *domain.Recipe) error {
	r.UpdateCalls = append(r.UpdateCalls, UpdateCall{Recipe: recipe})

	if r.FailOnUpdate {
		return errors.New("fake repository error")
	}

	if _, ok := r.Recipes[recipe.ID]; !ok {
		return repository.ErrRecipeNotFound
	}

	r.Recipes[recipe.ID] = recipe
	return nil
}

// Delete removes a recipe.
func (r *FakeRecipeRepository) Delete(ctx context.Context, userID, id uuid.UUID) error {
	r.DeleteCalls = append(r.DeleteCalls, id)

	if r.FailOnDelete {
		return errors.New("fake repository error")
	}

	recipe, ok := r.Recipes[id]
	if !ok || recipe.UserID != userID {
		return repository.ErrRecipeNotFound
	}

	delete(r.Recipes, id)
	return nil
}

// GetSimilar retrieves similar recipes.
func (r *FakeRecipeRepository) GetSimilar(ctx context.Context, userID, recipeID uuid.UUID, limit int) ([]domain.Recipe, error) {
	if r.FailOnGetSimilar {
		return nil, errors.New("fake repository error")
	}

	if recipe, ok := r.Recipes[recipeID]; !ok || recipe.UserID != userID {
		return nil, repository.ErrRecipeNotFound
	}

	recipes := make([]domain.Recipe, 0)
	for id, recipe := range r.Recipes {
		if id != recipeID && recipe.UserID == userID {
			recipes = append(recipes, *recipe)
			if len(recipes) >= limit {
				break
			}
		}
	}
	return recipes, nil
}

// GetIngredientByID retrieves an ingredient by ID.
func (r *FakeRecipeRepository) GetIngredientByID(ctx context.Context, userID, id uuid.UUID) (*domain.Ingredient, error) {
	if r.FailOnGetIngredientByID {
		return nil, errors.New("fake repository error")
	}

	ingredient, ok := r.Ingredients[id]
	if !ok || ingredient.UserID != userID {
		return nil, repository.ErrIngredientNotFound
	}
	return ingredient, nil
}

// GetOrCreateIngredient retrieves an ingredient by name or creates it.
func (r *FakeRecipeRepository) GetOrCreateIngredient(ctx context.Context, userID uuid.UUID, name string) (*domain.Ingredient, error) {
	if r.FailOnGetOrCreateIngredient {
		return nil, errors.New("fake repository error")
	}

	for _, ingredient := range r.Ingredients {
		if ingredient.UserID == userID && ingredient.Name == name {
			return ingredient, nil
		}
	}

	ingredient := &domain.Ingredient{
		ID:     uuid.New(),
		UserID: userID,
		Name:   name,
	}
	r.Ingredients[ingredient.ID] = ingredient
	return ingredient, nil
}

// GetCuisineByID retrieves a cuisine by ID for a user.
func (r *FakeRecipeRepository) GetCuisineByID(ctx context.Context, userID, id uuid.UUID) (*domain.Cuisine, error) {
	if r.FailOnGetCuisineByID {
		return nil, errors.New("fake repository error")
	}

	cuisine, ok := r.Cuisines[id]
	if !ok || cuisine.UserID != userID {
		return nil, repository.ErrCuisineNotFound
	}
	return cuisine, nil
}

// GetOrCreateCuisine retrieves a cuisine by name or creates it.
func (r *FakeRecipeRepository) GetOrCreateCuisine(ctx context.Context, userID uuid.UUID, name string) (*domain.Cuisine, error) {
	if r.FailOnGetOrCreateCuisine {
		return nil, errors.New("fake repository error")
	}

	for _, cuisine := range r.Cuisines {
		if cuisine.UserID == userID && cuisine.Name == name {
			return cuisine, nil
		}
	}

	cuisine := &domain.Cuisine{
		ID:     uuid.New(),
		UserID: userID,
		Name:   name,
	}
	r.Cuisines[cuisine.ID] = cuisine
	return cuisine, nil
}

// GetCuisines retrieves all cuisines.
func (r *FakeRecipeRepository) GetCuisines(ctx context.Context, userID uuid.UUID) ([]domain.Cuisine, error) {
	if r.FailOnGetCuisines {
		return nil, errors.New("fake repository error")
	}

	cuisines := make([]domain.Cuisine, 0, len(r.Cuisines))
	for _, cuisine := range r.Cuisines {
		if cuisine.UserID == userID {
			cuisines = append(cuisines, *cuisine)
		}
	}
	return cuisines, nil
}

// AddRecipe adds a recipe to the fake repository for test setup.
func (r *FakeRecipeRepository) AddRecipe(recipe *domain.Recipe) {
	r.Recipes[recipe.ID] = recipe
}

// AddIngredient adds an ingredient to the fake repository for test setup.
func (r *FakeRecipeRepository) AddIngredient(ingredient *domain.Ingredient) {
	r.Ingredients[ingredient.ID] = ingredient
}

// AddCuisine adds a cuisine to the fake repository for test setup.
func (r *FakeRecipeRepository) AddCuisine(cuisine *domain.Cuisine) {
	r.Cuisines[cuisine.ID] = cuisine
}

// FakeEventPublisher is an in-memory implementation of EventPublisher for testing.
type FakeEventPublisher struct {
	RecipeUpsertedEvents []*domain.Recipe
	RecipeDeletedEvents  []DeletedEvent

	FailOnPublishUpserted bool
	FailOnPublishDeleted  bool
}

// DeletedEvent represents a deleted recipe event in tests.
type DeletedEvent struct {
	RecipeID uuid.UUID
	UserID   uuid.UUID
}

// NewFakeEventPublisher creates a new fake event publisher.
func NewFakeEventPublisher() *FakeEventPublisher {
	return &FakeEventPublisher{
		RecipeUpsertedEvents: []*domain.Recipe{},
		RecipeDeletedEvents:  []DeletedEvent{},
	}
}

// PublishRecipeUpserted records a RecipeUpsertedEvent.
func (p *FakeEventPublisher) PublishRecipeUpserted(ctx context.Context, recipe *domain.Recipe) error {
	if p.FailOnPublishUpserted {
		return errors.New("fake publisher error")
	}
	p.RecipeUpsertedEvents = append(p.RecipeUpsertedEvents, recipe)
	return nil
}

// PublishRecipeDeleted records a RecipeDeletedEvent.
func (p *FakeEventPublisher) PublishRecipeDeleted(ctx context.Context, recipeID, userID uuid.UUID) error {
	if p.FailOnPublishDeleted {
		return errors.New("fake publisher error")
	}
	p.RecipeDeletedEvents = append(p.RecipeDeletedEvents, DeletedEvent{RecipeID: recipeID, UserID: userID})
	return nil
}

// UpsertedEventCount returns the number of RecipeUpsertedEvents published.
func (p *FakeEventPublisher) UpsertedEventCount() int {
	return len(p.RecipeUpsertedEvents)
}

// DeletedEventCount returns the number of RecipeDeletedEvents published.
func (p *FakeEventPublisher) DeletedEventCount() int {
	return len(p.RecipeDeletedEvents)
}

// FakeVectorGenerator is a fake implementation of vector.Generator for testing.
type FakeVectorGenerator struct {
	FixedVector pgvector.Vector
}

// NewFakeVectorGenerator creates a new fake vector generator.
func NewFakeVectorGenerator() *FakeVectorGenerator {
	dims := make([]float32, 1536)
	dims[0] = 1.0
	return &FakeVectorGenerator{
		FixedVector: pgvector.NewVector(dims),
	}
}

// Generate returns a fixed vector for testing.
func (g *FakeVectorGenerator) Generate(text string) pgvector.Vector {
	return g.FixedVector
}

// GenerateForRecipe returns a fixed vector for testing.
func (g *FakeVectorGenerator) GenerateForRecipe(recipe *domain.Recipe) pgvector.Vector {
	return g.FixedVector
}
