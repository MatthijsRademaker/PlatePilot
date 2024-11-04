package entities

import (
	"time"

	"gorm.io/gorm"
)

// Recipe Entity
type RecipeEntity struct {
	gorm.Model
	Name string
	// Many-to-many with ingredients through RecipeIngredientEntity as join table
	Ingredients  []RecipeIngredientEntity `gorm:"foreignKey:RecipeEntityID"`
	Instructions []string                 `gorm:"serializer:json"`
	// Many-to-many with cuisines
	Cuisines    []CuisineEntity `gorm:"many2many:recipe_cuisines;"`
	KCalories   uint
	CookingTime time.Duration
}

// Cuisine Entity
type CuisineEntity struct {
	gorm.Model
	Name string
	// Remove RecipeEntityID - not needed for many-to-many
	Recipes []RecipeEntity `gorm:"many2many:recipe_cuisines;"`
}

// Base Ingredient Entity
type IngredientEntity struct {
	gorm.Model
	Name string
	// Optional: Add reverse relation if you need to query ingredients' recipes
	Recipes []RecipeIngredientEntity `gorm:"foreignKey:IngredientEntityID"`
}

// Join table for Recipe-Ingredient relationship
type RecipeIngredientEntity struct {
	gorm.Model
	Quantity int
	Unit     string
	// // Foreign keys for the relationships
	// IngredientEntityID uint
	// RecipeEntityID     uint
	// Add the actual relations for easier querying
	Ingredient IngredientEntity //`gorm:"foreignKey:IngredientEntityID"`
	Recipe     RecipeEntity     //`gorm:"foreignKey:RecipeEntityID"`
}
