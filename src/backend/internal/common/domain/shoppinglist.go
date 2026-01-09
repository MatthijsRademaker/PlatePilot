package domain

import (
	"strconv"
	"time"

	"github.com/google/uuid"
)

// IngredientCategory represents a category for grouping ingredients
type IngredientCategory struct {
	ID           uuid.UUID
	Name         string
	DisplayOrder int
	CreatedAt    time.Time
}

// ShoppingList represents a user's shopping list
type ShoppingList struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	Name          string
	WeekStartDate *time.Time
	Items         []ShoppingListItem
	Recipes       []Recipe
	CreatedAt     time.Time
	UpdatedAt     time.Time
	CompletedAt   *time.Time
}

// ShoppingListItem represents a single item in a shopping list
type ShoppingListItem struct {
	ID             uuid.UUID
	ShoppingListID uuid.UUID
	IngredientID   *uuid.UUID
	Ingredient     *Ingredient
	Category       *IngredientCategory
	CustomName     *string
	Quantity       *float64
	Unit           *string
	Checked        bool
	Notes          *string
	IsCustom       bool
	Sources        []ShoppingListItemSource
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// ShoppingListItemSource tracks which recipe contributed to an item
type ShoppingListItemSource struct {
	RecipeID   uuid.UUID
	RecipeName string
	Quantity   *float64
	Unit       *string
}

// ItemIDs returns all item IDs from the shopping list
func (s *ShoppingList) ItemIDs() []uuid.UUID {
	ids := make([]uuid.UUID, len(s.Items))
	for i, item := range s.Items {
		ids[i] = item.ID
	}
	return ids
}

// RecipeIDs returns all recipe IDs associated with the shopping list
func (s *ShoppingList) RecipeIDs() []uuid.UUID {
	ids := make([]uuid.UUID, len(s.Recipes))
	for i, recipe := range s.Recipes {
		ids[i] = recipe.ID
	}
	return ids
}

// UncheckedItems returns all unchecked items
func (s *ShoppingList) UncheckedItems() []ShoppingListItem {
	var items []ShoppingListItem
	for _, item := range s.Items {
		if !item.Checked {
			items = append(items, item)
		}
	}
	return items
}

// CheckedItems returns all checked items
func (s *ShoppingList) CheckedItems() []ShoppingListItem {
	var items []ShoppingListItem
	for _, item := range s.Items {
		if item.Checked {
			items = append(items, item)
		}
	}
	return items
}

// IsCompleted returns true if all items are checked
func (s *ShoppingList) IsCompleted() bool {
	if len(s.Items) == 0 {
		return false
	}
	for _, item := range s.Items {
		if !item.Checked {
			return false
		}
	}
	return true
}

// DisplayName returns the item's display name (ingredient name or custom name)
func (i *ShoppingListItem) DisplayName() string {
	if i.IsCustom && i.CustomName != nil {
		return *i.CustomName
	}
	if i.Ingredient != nil {
		return i.Ingredient.Name
	}
	return ""
}

// DisplayQuantity returns a formatted quantity string
func (i *ShoppingListItem) DisplayQuantity() string {
	if i.Quantity == nil || *i.Quantity == 0 {
		return "as needed"
	}
	qty := formatQuantity(*i.Quantity)
	if i.Unit != nil && *i.Unit != "" {
		return qty + " " + *i.Unit
	}
	return qty
}

// formatQuantity formats a float quantity for display
func formatQuantity(q float64) string {
	if q == float64(int(q)) {
		return strconv.Itoa(int(q))
	}
	return strconv.FormatFloat(q, 'f', -1, 64)
}

// RecipeIngredientWithQuantity extends Ingredient with quantity info from a recipe
type RecipeIngredientWithQuantity struct {
	Ingredient
	Quantity *float64
	Unit     *string
}

// AggregatedIngredient represents an ingredient with combined quantities from multiple recipes
type AggregatedIngredient struct {
	IngredientID   uuid.UUID
	IngredientName string
	CategoryID     *uuid.UUID
	CategoryName   *string
	Quantities     []QuantityUnit
	RecipeSources  []RecipeSource
}

// QuantityUnit represents a quantity with its unit
type QuantityUnit struct {
	Quantity *float64
	Unit     *string
}

// RecipeSource tracks quantity contribution from a specific recipe
type RecipeSource struct {
	RecipeID   uuid.UUID
	RecipeName string
	Quantity   *float64
	Unit       *string
}
