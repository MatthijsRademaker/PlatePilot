package domain

import (
	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
)

// Recipe represents a recipe in the meal planner read model.
type Recipe struct {
	ID                 uuid.UUID
	UserID             uuid.UUID
	Name               string
	Description        string
	PrepTimeMinutes    int
	CookTimeMinutes    int
	TotalTimeMinutes   int
	Servings           int
	YieldQuantity      *float64
	YieldUnit          string
	SearchVector       pgvector.Vector
	CuisineID          uuid.UUID
	CuisineName        string
	MainIngredientID   uuid.UUID
	MainIngredientName string
	IngredientIDs      []uuid.UUID
	AllergyIDs         []uuid.UUID
	Tags               []string
	ImageURL           string
	CaloriesTotal      int
	CaloriesPerServing int
	ProteinG           float64
	CarbsG             float64
	FatG               float64
	FiberG             float64
	SugarG             float64
	SodiumMg           float64
}
