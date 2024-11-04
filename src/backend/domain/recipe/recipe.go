package recipe

import (
	"errors"
	"time"
)

type Recipe struct {
	Name         string
	Ingredients  []Ingredient
	Instructions []string
	CookingTime  time.Duration
	Cuisines     []string
	KCalories    uint
}

// Factory method with validation
func NewRecipe(name string, ingredients []Ingredient, instructions []string, cookingTime time.Duration, cuisines []string, kCalories uint) (*Recipe, error) {
	if name == "" {
		return nil, errors.New("recipe name cannot be empty")
	}

	if len(ingredients) == 0 {
		return nil, errors.New("recipe must have at least one ingredient")
	}

	return &Recipe{
		Name:         name,
		Ingredients:  ingredients,
		Instructions: instructions,
		CookingTime:  cookingTime,
		Cuisines:     cuisines,
		KCalories:    kCalories,
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
