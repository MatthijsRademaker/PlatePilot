package domain

import "github.com/google/uuid"

// RecipeFilter defines optional filters for listing recipes.
type RecipeFilter struct {
	CuisineID    *uuid.UUID
	IngredientID *uuid.UUID
	AllergyID    *uuid.UUID
	Tags         []string
}
