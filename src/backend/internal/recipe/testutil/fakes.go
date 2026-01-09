package testutil

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"

	"github.com/platepilot/backend/internal/common/domain"
	"github.com/platepilot/backend/internal/recipe/repository"
)

// FakeRecipeRepository is an in-memory implementation of RecipeRepository for testing
type FakeRecipeRepository struct {
	Recipes     map[uuid.UUID]*domain.Recipe
	Ingredients map[uuid.UUID]*domain.Ingredient
	Cuisines    map[uuid.UUID]*domain.Cuisine
	Units       map[uuid.UUID]*domain.Unit

	// Failure modes for testing error paths
	FailOnGetByID               bool
	FailOnGetAll                bool
	FailOnCreate                bool
	FailOnGetSimilar            bool
	FailOnGetByCuisine          bool
	FailOnGetByIngredient       bool
	FailOnGetExcludingAllergy   bool
	FailOnGetIngredientByID     bool
	FailOnGetCuisineByID        bool
	FailOnGetOrCreateIngredient bool
	FailOnGetOrCreateCuisine    bool
	FailOnGetUnits              bool
	FailOnGetUnitByName         bool
	FailOnCreateUnit            bool

	// Call tracking for assertions
	CreateCalls  []CreateCall
	GetByIDCalls []uuid.UUID
}

// CreateCall records a call to Create
type CreateCall struct {
	Recipe *domain.Recipe
}

// NewFakeRecipeRepository creates a new fake repository
func NewFakeRecipeRepository() *FakeRecipeRepository {
	return &FakeRecipeRepository{
		Recipes:      make(map[uuid.UUID]*domain.Recipe),
		Ingredients:  make(map[uuid.UUID]*domain.Ingredient),
		Cuisines:     make(map[uuid.UUID]*domain.Cuisine),
		Units:        make(map[uuid.UUID]*domain.Unit),
		CreateCalls:  []CreateCall{},
		GetByIDCalls: []uuid.UUID{},
	}
}

// GetByID retrieves a recipe by ID
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

// GetAll retrieves all recipes with pagination
func (r *FakeRecipeRepository) GetAll(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.Recipe, error) {
	if r.FailOnGetAll {
		return nil, errors.New("fake repository error")
	}

	recipes := make([]domain.Recipe, 0, len(r.Recipes))
	for _, recipe := range r.Recipes {
		if recipe.UserID == userID {
			recipes = append(recipes, *recipe)
		}
	}

	// Apply pagination
	if offset >= len(recipes) {
		return []domain.Recipe{}, nil
	}
	end := offset + limit
	if end > len(recipes) {
		end = len(recipes)
	}
	return recipes[offset:end], nil
}

// Count returns the total number of recipes
func (r *FakeRecipeRepository) Count(ctx context.Context, userID uuid.UUID) (int64, error) {
	if r.FailOnGetAll {
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

// Create creates a new recipe
func (r *FakeRecipeRepository) Create(ctx context.Context, recipe *domain.Recipe) error {
	r.CreateCalls = append(r.CreateCalls, CreateCall{Recipe: recipe})

	if r.FailOnCreate {
		return errors.New("fake repository error")
	}

	// Generate ID if not set
	if recipe.ID == uuid.Nil {
		recipe.ID = uuid.New()
	}

	r.Recipes[recipe.ID] = recipe
	return nil
}

// GetSimilar retrieves similar recipes
func (r *FakeRecipeRepository) GetSimilar(ctx context.Context, userID, recipeID uuid.UUID, limit int) ([]domain.Recipe, error) {
	if r.FailOnGetSimilar {
		return nil, errors.New("fake repository error")
	}

	// Check if source recipe exists
	if recipe, ok := r.Recipes[recipeID]; !ok || recipe.UserID != userID {
		return nil, repository.ErrRecipeNotFound
	}

	// Return all other recipes (simplified similarity)
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

// GetByCuisine retrieves recipes by cuisine
func (r *FakeRecipeRepository) GetByCuisine(ctx context.Context, userID, cuisineID uuid.UUID, limit, offset int) ([]domain.Recipe, error) {
	if r.FailOnGetByCuisine {
		return nil, errors.New("fake repository error")
	}

	recipes := make([]domain.Recipe, 0)
	for _, recipe := range r.Recipes {
		if recipe.UserID == userID && recipe.Cuisine != nil && recipe.Cuisine.ID == cuisineID {
			recipes = append(recipes, *recipe)
		}
	}
	return recipes, nil
}

// GetByIngredient retrieves recipes containing a specific ingredient
func (r *FakeRecipeRepository) GetByIngredient(ctx context.Context, userID, ingredientID uuid.UUID, limit, offset int) ([]domain.Recipe, error) {
	if r.FailOnGetByIngredient {
		return nil, errors.New("fake repository error")
	}

	recipes := make([]domain.Recipe, 0)
	for _, recipe := range r.Recipes {
		if recipe.UserID != userID {
			continue
		}
		// Check main ingredient
		if recipe.MainIngredient != nil && recipe.MainIngredient.ID == ingredientID {
			recipes = append(recipes, *recipe)
			continue
		}
		// Check ingredients list
		for _, ing := range recipe.Ingredients {
			if ing.ID == ingredientID {
				recipes = append(recipes, *recipe)
				break
			}
		}
	}
	return recipes, nil
}

// GetExcludingAllergy retrieves recipes excluding a specific allergy
func (r *FakeRecipeRepository) GetExcludingAllergy(ctx context.Context, userID, allergyID uuid.UUID, limit, offset int) ([]domain.Recipe, error) {
	if r.FailOnGetExcludingAllergy {
		return nil, errors.New("fake repository error")
	}

	// For simplicity, return all recipes (real implementation would filter)
	recipes := make([]domain.Recipe, 0, len(r.Recipes))
	for _, recipe := range r.Recipes {
		if recipe.UserID == userID {
			recipes = append(recipes, *recipe)
		}
	}
	return recipes, nil
}

// GetIngredientByID retrieves an ingredient by ID
func (r *FakeRecipeRepository) GetIngredientByID(ctx context.Context, id uuid.UUID) (*domain.Ingredient, error) {
	if r.FailOnGetIngredientByID {
		return nil, errors.New("fake repository error")
	}

	ingredient, ok := r.Ingredients[id]
	if !ok {
		return nil, repository.ErrIngredientNotFound
	}
	return ingredient, nil
}

// GetOrCreateIngredient retrieves an ingredient by name or creates it.
func (r *FakeRecipeRepository) GetOrCreateIngredient(ctx context.Context, name string, quantity string) (*domain.Ingredient, error) {
	if r.FailOnGetOrCreateIngredient {
		return nil, errors.New("fake repository error")
	}

	for _, ingredient := range r.Ingredients {
		if ingredient.Name == name {
			return ingredient, nil
		}
	}

	ingredient := &domain.Ingredient{
		ID:       uuid.New(),
		Name:     name,
		Quantity: quantity,
	}
	r.Ingredients[ingredient.ID] = ingredient
	return ingredient, nil
}

// GetCuisineByID retrieves a cuisine by ID
func (r *FakeRecipeRepository) GetCuisineByID(ctx context.Context, id uuid.UUID) (*domain.Cuisine, error) {
	if r.FailOnGetCuisineByID {
		return nil, errors.New("fake repository error")
	}

	cuisine, ok := r.Cuisines[id]
	if !ok {
		return nil, repository.ErrCuisineNotFound
	}
	return cuisine, nil
}

// GetOrCreateCuisine retrieves a cuisine by name or creates it.
func (r *FakeRecipeRepository) GetOrCreateCuisine(ctx context.Context, name string) (*domain.Cuisine, error) {
	if r.FailOnGetOrCreateCuisine {
		return nil, errors.New("fake repository error")
	}

	for _, cuisine := range r.Cuisines {
		if cuisine.Name == name {
			return cuisine, nil
		}
	}

	cuisine := &domain.Cuisine{
		ID:   uuid.New(),
		Name: name,
	}
	r.Cuisines[cuisine.ID] = cuisine
	return cuisine, nil
}

// GetUnits retrieves all units for a user.
func (r *FakeRecipeRepository) GetUnits(ctx context.Context, userID uuid.UUID) ([]domain.Unit, error) {
	if r.FailOnGetUnits {
		return nil, errors.New("fake repository error")
	}

	units := make([]domain.Unit, 0, len(r.Units))
	for _, unit := range r.Units {
		if unit.UserID == userID {
			units = append(units, *unit)
		}
	}
	return units, nil
}

// GetUnitByName retrieves a unit by name for a user.
func (r *FakeRecipeRepository) GetUnitByName(ctx context.Context, userID uuid.UUID, name string) (*domain.Unit, error) {
	if r.FailOnGetUnitByName {
		return nil, errors.New("fake repository error")
	}

	for _, unit := range r.Units {
		if unit.UserID == userID && unit.Name == name {
			return unit, nil
		}
	}
	return nil, repository.ErrUnitNotFound
}

// CreateUnit creates a new unit.
func (r *FakeRecipeRepository) CreateUnit(ctx context.Context, unit *domain.Unit) error {
	if r.FailOnCreateUnit {
		return errors.New("fake repository error")
	}
	if unit.ID == uuid.Nil {
		unit.ID = uuid.New()
	}
	r.Units[unit.ID] = unit
	return nil
}

// AddRecipe adds a recipe to the fake repository for test setup
func (r *FakeRecipeRepository) AddRecipe(recipe *domain.Recipe) {
	r.Recipes[recipe.ID] = recipe
}

// AddIngredient adds an ingredient to the fake repository for test setup
func (r *FakeRecipeRepository) AddIngredient(ingredient *domain.Ingredient) {
	r.Ingredients[ingredient.ID] = ingredient
}

// AddCuisine adds a cuisine to the fake repository for test setup
func (r *FakeRecipeRepository) AddCuisine(cuisine *domain.Cuisine) {
	r.Cuisines[cuisine.ID] = cuisine
}

// AddUnit adds a unit to the fake repository for test setup
func (r *FakeRecipeRepository) AddUnit(unit *domain.Unit) {
	r.Units[unit.ID] = unit
}

// FakeEventPublisher is an in-memory implementation of EventPublisher for testing
type FakeEventPublisher struct {
	RecipeCreatedEvents []*domain.Recipe
	RecipeUpdatedEvents []*domain.Recipe

	// Failure modes
	FailOnPublishCreated bool
	FailOnPublishUpdated bool
}

// NewFakeEventPublisher creates a new fake event publisher
func NewFakeEventPublisher() *FakeEventPublisher {
	return &FakeEventPublisher{
		RecipeCreatedEvents: []*domain.Recipe{},
		RecipeUpdatedEvents: []*domain.Recipe{},
	}
}

// PublishRecipeCreated records a RecipeCreatedEvent
func (p *FakeEventPublisher) PublishRecipeCreated(ctx context.Context, recipe *domain.Recipe) error {
	if p.FailOnPublishCreated {
		return errors.New("fake publisher error")
	}
	p.RecipeCreatedEvents = append(p.RecipeCreatedEvents, recipe)
	return nil
}

// PublishRecipeUpdated records a RecipeUpdatedEvent
func (p *FakeEventPublisher) PublishRecipeUpdated(ctx context.Context, recipe *domain.Recipe) error {
	if p.FailOnPublishUpdated {
		return errors.New("fake publisher error")
	}
	p.RecipeUpdatedEvents = append(p.RecipeUpdatedEvents, recipe)
	return nil
}

// CreatedEventCount returns the number of RecipeCreatedEvents published
func (p *FakeEventPublisher) CreatedEventCount() int {
	return len(p.RecipeCreatedEvents)
}

// UpdatedEventCount returns the number of RecipeUpdatedEvents published
func (p *FakeEventPublisher) UpdatedEventCount() int {
	return len(p.RecipeUpdatedEvents)
}

// FakeVectorGenerator is a fake implementation of vector.Generator for testing
type FakeVectorGenerator struct {
	FixedVector pgvector.Vector
}

// NewFakeVectorGenerator creates a new fake vector generator
func NewFakeVectorGenerator() *FakeVectorGenerator {
	// Create a simple fixed vector for testing
	dims := make([]float32, 1536)
	dims[0] = 1.0 // Simple non-zero vector
	return &FakeVectorGenerator{
		FixedVector: pgvector.NewVector(dims),
	}
}

// Generate returns a fixed vector for testing
func (g *FakeVectorGenerator) Generate(text string) pgvector.Vector {
	return g.FixedVector
}

// GenerateForRecipe returns a fixed vector for testing
func (g *FakeVectorGenerator) GenerateForRecipe(recipe *domain.Recipe) pgvector.Vector {
	return g.FixedVector
}
