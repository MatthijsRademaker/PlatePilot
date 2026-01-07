package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
)

// Recipe represents a recipe in the system
type Recipe struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	Name            string
	Description     string
	PrepTime        string
	CookTime        string
	MainIngredient  *Ingredient
	Cuisine         *Cuisine
	Ingredients     []Ingredient
	Directions      []string
	NutritionalInfo NutritionalInfo
	Metadata        Metadata
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Allergies returns all unique allergies from the recipe's ingredients
func (r *Recipe) Allergies() []Allergy {
	seen := make(map[uuid.UUID]bool)
	var allergies []Allergy

	for _, ingredient := range r.Ingredients {
		for _, allergy := range ingredient.Allergies {
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
	ids := make([]uuid.UUID, len(r.Ingredients))
	for i, ing := range r.Ingredients {
		ids[i] = ing.ID
	}
	return ids
}

// Metadata contains recipe metadata including search vector
type Metadata struct {
	SearchVector  pgvector.Vector
	ImageURL      string
	Tags          []string
	PublishedDate time.Time
}

// NutritionalInfo contains nutritional information for a recipe
type NutritionalInfo struct {
	Calories int
}
