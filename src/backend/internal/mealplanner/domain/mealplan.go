package domain

import (
	"time"

	"github.com/google/uuid"
)

// WeekPlan represents a persisted weekly meal plan.
type WeekPlan struct {
	UserID    uuid.UUID
	StartDate time.Time
	EndDate   time.Time
	Slots     []MealSlot
}

// MealSlot represents a planned meal slot tied to a recipe.
type MealSlot struct {
	Date              time.Time
	MealType          string
	RecipeID          uuid.UUID
	RecipeName        string
	RecipeDescription string
}
