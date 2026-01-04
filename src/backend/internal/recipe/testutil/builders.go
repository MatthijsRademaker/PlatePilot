package testutil

import (
	"time"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"

	"github.com/platepilot/backend/internal/common/domain"
)

// RecipeBuilder helps construct Recipe instances for testing
type RecipeBuilder struct {
	recipe domain.Recipe
}

// NewRecipeBuilder creates a new RecipeBuilder with default values
func NewRecipeBuilder() *RecipeBuilder {
	dims := make([]float32, 1536)
	dims[0] = 1.0

	return &RecipeBuilder{
		recipe: domain.Recipe{
			ID:          uuid.New(),
			Name:        "Test Recipe",
			Description: "A test recipe description",
			PrepTime:    "15 mins",
			CookTime:    "30 mins",
			Directions:  []string{"Step 1", "Step 2"},
			Metadata: domain.Metadata{
				SearchVector:  pgvector.NewVector(dims),
				Tags:          []string{},
				PublishedDate: time.Now().UTC(),
			},
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
	}
}

// WithID sets the recipe ID
func (b *RecipeBuilder) WithID(id uuid.UUID) *RecipeBuilder {
	b.recipe.ID = id
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

// WithPrepTime sets the prep time
func (b *RecipeBuilder) WithPrepTime(prepTime string) *RecipeBuilder {
	b.recipe.PrepTime = prepTime
	return b
}

// WithCookTime sets the cook time
func (b *RecipeBuilder) WithCookTime(cookTime string) *RecipeBuilder {
	b.recipe.CookTime = cookTime
	return b
}

// WithMainIngredient sets the main ingredient
func (b *RecipeBuilder) WithMainIngredient(ingredient *domain.Ingredient) *RecipeBuilder {
	b.recipe.MainIngredient = ingredient
	return b
}

// WithCuisine sets the cuisine
func (b *RecipeBuilder) WithCuisine(cuisine *domain.Cuisine) *RecipeBuilder {
	b.recipe.Cuisine = cuisine
	return b
}

// WithIngredients sets the ingredients list
func (b *RecipeBuilder) WithIngredients(ingredients []domain.Ingredient) *RecipeBuilder {
	b.recipe.Ingredients = ingredients
	return b
}

// WithDirections sets the directions
func (b *RecipeBuilder) WithDirections(directions []string) *RecipeBuilder {
	b.recipe.Directions = directions
	return b
}

// Build returns the constructed Recipe
func (b *RecipeBuilder) Build() *domain.Recipe {
	return &b.recipe
}

// IngredientBuilder helps construct Ingredient instances for testing
type IngredientBuilder struct {
	ingredient domain.Ingredient
}

// NewIngredientBuilder creates a new IngredientBuilder with default values
func NewIngredientBuilder() *IngredientBuilder {
	return &IngredientBuilder{
		ingredient: domain.Ingredient{
			ID:        uuid.New(),
			Name:      "Test Ingredient",
			Quantity:  "1 cup",
			Allergies: []domain.Allergy{},
			CreatedAt: time.Now().UTC(),
		},
	}
}

// WithID sets the ingredient ID
func (b *IngredientBuilder) WithID(id uuid.UUID) *IngredientBuilder {
	b.ingredient.ID = id
	return b
}

// WithName sets the ingredient name
func (b *IngredientBuilder) WithName(name string) *IngredientBuilder {
	b.ingredient.Name = name
	return b
}

// WithQuantity sets the ingredient quantity
func (b *IngredientBuilder) WithQuantity(quantity string) *IngredientBuilder {
	b.ingredient.Quantity = quantity
	return b
}

// WithAllergies sets the allergies
func (b *IngredientBuilder) WithAllergies(allergies []domain.Allergy) *IngredientBuilder {
	b.ingredient.Allergies = allergies
	return b
}

// Build returns the constructed Ingredient
func (b *IngredientBuilder) Build() *domain.Ingredient {
	return &b.ingredient
}

// CuisineBuilder helps construct Cuisine instances for testing
type CuisineBuilder struct {
	cuisine domain.Cuisine
}

// NewCuisineBuilder creates a new CuisineBuilder with default values
func NewCuisineBuilder() *CuisineBuilder {
	return &CuisineBuilder{
		cuisine: domain.Cuisine{
			ID:        uuid.New(),
			Name:      "Test Cuisine",
			CreatedAt: time.Now().UTC(),
		},
	}
}

// WithID sets the cuisine ID
func (b *CuisineBuilder) WithID(id uuid.UUID) *CuisineBuilder {
	b.cuisine.ID = id
	return b
}

// WithName sets the cuisine name
func (b *CuisineBuilder) WithName(name string) *CuisineBuilder {
	b.cuisine.Name = name
	return b
}

// Build returns the constructed Cuisine
func (b *CuisineBuilder) Build() *domain.Cuisine {
	return &b.cuisine
}

// AllergyBuilder helps construct Allergy instances for testing
type AllergyBuilder struct {
	allergy domain.Allergy
}

// NewAllergyBuilder creates a new AllergyBuilder with default values
func NewAllergyBuilder() *AllergyBuilder {
	return &AllergyBuilder{
		allergy: domain.Allergy{
			ID:        uuid.New(),
			Name:      "Test Allergy",
			CreatedAt: time.Now().UTC(),
		},
	}
}

// WithID sets the allergy ID
func (b *AllergyBuilder) WithID(id uuid.UUID) *AllergyBuilder {
	b.allergy.ID = id
	return b
}

// WithName sets the allergy name
func (b *AllergyBuilder) WithName(name string) *AllergyBuilder {
	b.allergy.Name = name
	return b
}

// Build returns the constructed Allergy
func (b *AllergyBuilder) Build() *domain.Allergy {
	return &b.allergy
}
