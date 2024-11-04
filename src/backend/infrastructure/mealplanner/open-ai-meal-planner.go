package mealplanner

import (
	"PlatePilot/domain/mealplanner"
	"PlatePilot/domain/recipes"
	"context"

	"github.com/openai/openai-go" // imported as openai
	"github.com/openai/openai-go/option"
)

type OpenAiMealPlanner struct {
	client            *openai.Client
	recipesRepository recipes.RecipeRepository
}

func CreateMealPlanner(openAiKey string) *mealplanner.MealPlanner {

	client := openai.NewClient(
		option.WithAPIKey("My API Key"), // defaults to os.LookupEnv("OPENAI_API_KEY")
	)
	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage("Say this is a test"),
		}),
		Model: openai.F(openai.ChatModelGPT4o),
	})
	if err != nil {
		panic(err.Error())
	}
}
