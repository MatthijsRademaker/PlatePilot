package recipe

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Recipe struct {
	Id           uuid.UUID
	Name         string
	Ingredients  []Ingredient
	Instructions []string
	CookingTime  time.Duration
}

// Factory method with validation
func NewRecipe(name string, ingredients []Ingredient, instructions []string, cookingTime time.Duration) (*Recipe, error) {
	if name == "" {
		return nil, errors.New("recipe name cannot be empty")
	}

	if len(ingredients) == 0 {
		return nil, errors.New("recipe must have at least one ingredient")
	}

	return &Recipe{
		Id:           uuid.New(),
		Name:         name,
		Ingredients:  ingredients,
		Instructions: instructions,
		CookingTime:  cookingTime,
	}, nil
}

// Rich domain methods
func (r *Recipe) ChangeName(newName string) error {
	if newName == "" {
		return errors.New("recipe name cannot be empty")
	}
	r.Name = newName
	return nil
}

func (r *Recipe) AddIngredient(ingredient Ingredient) {
	r.Ingredients = append(r.Ingredients, ingredient)
}

func (r *Recipe) AddInstruction(instruction string) {
	r.Instructions = append(r.Instructions, instruction)
}
