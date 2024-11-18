package mealplanner

import (
	"PlatePilot/domain/mealplanner"
	"PlatePilot/domain/recipes"

	"github.com/openai/openai-go" // imported as openai
	"github.com/openai/openai-go/option"
)

type OpenAiMealPlanner struct {
	client            *openai.Client
	recipesRepository *recipes.RecipeRepository
}

func CreateMealPlanner(openAiKey string, recipeRepository *recipes.RecipeRepository) *mealplanner.MealPlanner {

	client := openai.NewClient(
		option.WithAPIKey("My API Key"), // defaults to os.LookupEnv("OPENAI_API_KEY")
	)

	return &OpenAiMealPlanner{
		client:            client,
		recipesRepository: recipeRepository,
	}
}

func (planner *mealplanner.MealPlanner) Suggest(mealPlan *mealplanner.MealPlan, amountToSuggest int, constraints ...mealplanner.MealPlanConstraints) (*mealplanner.MealPlan, error) {

}
