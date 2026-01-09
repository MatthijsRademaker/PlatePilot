package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
)

// Recipe represents a recipe in the system
type Recipe struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	Name             string
	Description      string
	PrepTimeMinutes  int
	CookTimeMinutes  int
	TotalTimeMinutes int
	Servings         int
	YieldQuantity    *float64
	YieldUnit        string
	MainIngredient   *Ingredient
	Cuisine          *Cuisine
	IngredientLines  []RecipeIngredientLine
	Steps            []RecipeStep
	Tags             []string
	ImageURL         string
	Nutrition        RecipeNutrition
	SearchVector     pgvector.Vector
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time
}

// RecipeIngredientLine represents a per-recipe ingredient line item.
type RecipeIngredientLine struct {
	ID            uuid.UUID
	Ingredient    Ingredient
	QuantityValue *float64
	QuantityText  string
	Unit          string
	IsOptional    bool
	Note          string
	SortOrder     int
}

// RecipeStep represents a structured instruction step.
type RecipeStep struct {
	ID               uuid.UUID
	StepIndex        int
	Instruction      string
	DurationSeconds  *int
	TemperatureValue *float64
	TemperatureUnit  string
	MediaURL         string
}

// RecipeNutrition contains aggregated nutritional information for a recipe.
type RecipeNutrition struct {
	CaloriesTotal      int
	CaloriesPerServing int
	ProteinG           float64
	CarbsG             float64
	FatG               float64
	FiberG             float64
	SugarG             float64
	SodiumMg           float64
}

// Allergies returns all unique allergies from the recipe's ingredients.
func (r *Recipe) Allergies() []Allergy {
	seen := make(map[uuid.UUID]bool)
	var allergies []Allergy

	for _, line := range r.IngredientLines {
		for _, allergy := range line.Ingredient.Allergies {
			if !seen[allergy.ID] {
				seen[allergy.ID] = true
				allergies = append(allergies, allergy)
			}
		}
	}

	return allergies
}

// AllergyIDs returns all unique allergy IDs from the recipe's ingredients
func (r *Recipe) AllergyIDs() []uuid.UUID {
	allergies := r.Allergies()
	ids := make([]uuid.UUID, len(allergies))
	for i, a := range allergies {
		ids[i] = a.ID
	}
	return ids
}

// IngredientIDs returns all ingredient IDs
func (r *Recipe) IngredientIDs() []uuid.UUID {
	ids := make([]uuid.UUID, len(r.IngredientLines))
	for i, line := range r.IngredientLines {
		ids[i] = line.Ingredient.ID
	}
	return ids
}
