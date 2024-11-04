package mealplanner

import recipes "PlatePilot/domain/recipes"

type MealPlan struct {
	RecipeIds []uint
}

// Factory method with validation
func NewMealPlan(recipe1 recipes.Recipe) *MealPlan {
	recipes := []uint{recipe1.Id}
	return &MealPlan{
		RecipeIds: recipes,
	}
}

func (mealPlan *MealPlan) AddRecipe(recipe recipes.Recipe) *MealPlan {
	mealPlan.RecipeIds = append(mealPlan.RecipeIds, recipe.Id)

	return mealPlan
}

func (mealPlan *MealPlan) AddRecipes(recipes []recipes.Recipe) *MealPlan {
	for _, r := range recipes {
		mealPlan.RecipeIds = append(mealPlan.RecipeIds, r.Id)
	}

	return mealPlan
}
