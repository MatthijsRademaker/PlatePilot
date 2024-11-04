package mealplanner

import recipes "PlatePilot/domain/recipes"

type MealPlan struct {
	Recipes []recipes.Recipe
}

// Factory method with validation
func NewMealPlan(recipe1 recipes.Recipe) *MealPlan {
	recipes := []recipes.Recipe{recipe1}
	return &MealPlan{
		Recipes: recipes,
	}
}

func (mealPlan *MealPlan) AddRecipe(recipe recipes.Recipe) *MealPlan {
	mealPlan.Recipes = append(mealPlan.Recipes, recipe)

	return mealPlan
}

func (mealPlan *MealPlan) AddRecipes(recipes []recipes.Recipe) *MealPlan {
	mealPlan.Recipes = append(mealPlan.Recipes, recipes...)

	return mealPlan
}
