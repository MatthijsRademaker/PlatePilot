package entities

import (
	"time"

	"gorm.io/gorm"
)

type RecipeEntity struct {
	gorm.Model
	Name         string
	Ingredients  []IngredientEntity
	Instructions []string
	CookingTime  time.Duration
}

type IngredientEntity struct {
	gorm.Model
	Name string
}
