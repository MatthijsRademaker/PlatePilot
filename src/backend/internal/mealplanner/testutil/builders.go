package testutil

import (
	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"

	"github.com/platepilot/backend/internal/mealplanner/repository"
)

// RecipeBuilder helps construct repository.Recipe instances for testing
type RecipeBuilder struct {
	recipe repository.Recipe
}

// NewRecipeBuilder creates a new RecipeBuilder with default values
func NewRecipeBuilder() *RecipeBuilder {
	dims := make([]float32, 1536)
	dims[0] = 1.0

	return &RecipeBuilder{
		recipe: repository.Recipe{
			ID:                 uuid.New(),
			UserID:             uuid.New(),
			Name:               "Test Recipe",
			Description:        "A test recipe description",
			PrepTimeMinutes:    15,
			CookTimeMinutes:    30,
			TotalTimeMinutes:   45,
			Servings:           2,
			SearchVector:       pgvector.NewVector(dims),
			CuisineID:          uuid.New(),
			CuisineName:        "Test Cuisine",
			MainIngredientID:   uuid.New(),
			MainIngredientName: "Test Ingredient",
			IngredientIDs:      []uuid.UUID{},
			AllergyIDs:         []uuid.UUID{},
			ImageURL:           "",
			Tags:               []string{},
			CaloriesTotal:      500,
			CaloriesPerServing: 250,
		},
	}
}

// WithID sets the recipe ID
func (b *RecipeBuilder) WithID(id uuid.UUID) *RecipeBuilder {
	b.recipe.ID = id
	return b
}

// WithUserID sets the recipe owner
func (b *RecipeBuilder) WithUserID(id uuid.UUID) *RecipeBuilder {
	b.recipe.UserID = id
	return b
}

// WithName sets the recipe name
func (b *RecipeBuilder) WithName(name string) *RecipeBuilder {
	b.recipe.Name = name
	return b
}

// WithDescription sets the recipe description
func (b *RecipeBuilder) WithDescription(desc string) *RecipeBuilder {
	b.recipe.Description = desc
	return b
}

// WithCuisineID sets the cuisine ID
func (b *RecipeBuilder) WithCuisineID(id uuid.UUID) *RecipeBuilder {
	b.recipe.CuisineID = id
	return b
}

// WithCuisineName sets the cuisine name
func (b *RecipeBuilder) WithCuisineName(name string) *RecipeBuilder {
	b.recipe.CuisineName = name
	return b
}

// WithMainIngredientID sets the main ingredient ID
func (b *RecipeBuilder) WithMainIngredientID(id uuid.UUID) *RecipeBuilder {
	b.recipe.MainIngredientID = id
	return b
}

// WithMainIngredientName sets the main ingredient name
func (b *RecipeBuilder) WithMainIngredientName(name string) *RecipeBuilder {
	b.recipe.MainIngredientName = name
	return b
}

// WithIngredientIDs sets the ingredient IDs
func (b *RecipeBuilder) WithIngredientIDs(ids []uuid.UUID) *RecipeBuilder {
	b.recipe.IngredientIDs = ids
	return b
}

// WithAllergyIDs sets the allergy IDs
func (b *RecipeBuilder) WithAllergyIDs(ids []uuid.UUID) *RecipeBuilder {
	b.recipe.AllergyIDs = ids
	return b
}

// WithSearchVector sets the search vector
func (b *RecipeBuilder) WithSearchVector(vector pgvector.Vector) *RecipeBuilder {
	b.recipe.SearchVector = vector
	return b
}

// Build returns the constructed Recipe
func (b *RecipeBuilder) Build() repository.Recipe {
	return b.recipe
}

// CreateTestVector creates a simple test vector for diversity testing
// The index parameter allows creating vectors that have different similarity scores
func CreateTestVector(index int) pgvector.Vector {
	dims := make([]float32, 1536)
	// Set different values based on index to create distinguishable vectors
	for i := 0; i < 128 && i < len(dims); i++ {
		if i == index%128 {
			dims[i] = 1.0
		}
	}
	return pgvector.NewVector(dims)
}

// CreateSimilarVector creates a vector similar to the given base vector
func CreateSimilarVector(base pgvector.Vector) pgvector.Vector {
	slice := base.Slice()
	dims := make([]float32, len(slice))
	copy(dims, slice)
	// Add small perturbation to first few elements
	for i := 0; i < 10 && i < len(dims); i++ {
		dims[i] += 0.01
	}
	return pgvector.NewVector(dims)
}

// CreateDifferentVector creates a vector different from the given base vector
func CreateDifferentVector(base pgvector.Vector) pgvector.Vector {
	slice := base.Slice()
	dims := make([]float32, len(slice))
	// Create orthogonal-ish vector by using different indices
	for i := 0; i < len(dims); i++ {
		if slice[i] == 0 && i < 128 {
			dims[i] = 1.0
		}
	}
	return pgvector.NewVector(dims)
}
