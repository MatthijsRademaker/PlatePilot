package agents

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/google/uuid"

	"github.com/platepilot/backend/internal/llm"
)

//go:embed prompts/recipe_suggestions.txt
var recipeSuggestionsPrompt string

// RecipeSuggestionsAgent suggests recipes based on user preferences and constraints
type RecipeSuggestionsAgent struct {
	*BaseAgent
}

// RecipeSuggestionRequest contains the input for recipe suggestions
type RecipeSuggestionRequest struct {
	// AvailableRecipes is the list of recipes to choose from
	AvailableRecipes []RecipeInfo `json:"available_recipes"`

	// Constraints defines filtering criteria
	Constraints SuggestionConstraints `json:"constraints"`

	// Amount is the number of suggestions to return
	Amount int `json:"amount"`

	// AlreadySelected contains IDs of recipes already chosen (to avoid duplicates)
	AlreadySelected []string `json:"already_selected,omitempty"`
}

// RecipeInfo contains recipe information for the agent
type RecipeInfo struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Cuisine     string   `json:"cuisine,omitempty"`
	Ingredients []string `json:"ingredients,omitempty"`
	PrepTime    string   `json:"prep_time,omitempty"`
	CookTime    string   `json:"cook_time,omitempty"`
	Calories    int      `json:"calories,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

// SuggestionConstraints defines filtering criteria for suggestions
type SuggestionConstraints struct {
	// DietaryRestrictions like "vegetarian", "vegan", "gluten-free"
	DietaryRestrictions []string `json:"dietary_restrictions,omitempty"`

	// PreferredCuisines like "Italian", "Mexican", "Asian"
	PreferredCuisines []string `json:"preferred_cuisines,omitempty"`

	// AvoidIngredients lists ingredients to avoid
	AvoidIngredients []string `json:"avoid_ingredients,omitempty"`

	// MaxPrepTime in minutes
	MaxPrepTime int `json:"max_prep_time,omitempty"`

	// MaxCalories per serving
	MaxCalories int `json:"max_calories,omitempty"`

	// PreferQuick suggests quick recipes first
	PreferQuick bool `json:"prefer_quick,omitempty"`
}

// NewRecipeSuggestionsAgent creates a new recipe suggestions agent
func NewRecipeSuggestionsAgent(client *llm.Client, cache Cache) *RecipeSuggestionsAgent {
	return &RecipeSuggestionsAgent{
		BaseAgent: NewBaseAgent("recipe_suggestions", recipeSuggestionsPrompt, client, cache),
	}
}

// Execute runs the recipe suggestions agent
func (a *RecipeSuggestionsAgent) Execute(ctx context.Context, input AgentInput) (*AgentOutput, error) {
	var result RecipeSuggestionOutput
	return a.ExecuteWithJSON(ctx, input, &result)
}

// Suggest returns recipe suggestions based on the request
func (a *RecipeSuggestionsAgent) Suggest(ctx context.Context, req RecipeSuggestionRequest) (*RecipeSuggestionOutput, error) {
	if len(req.AvailableRecipes) == 0 {
		return &RecipeSuggestionOutput{Recipes: []SuggestedRecipe{}}, nil
	}

	if req.Amount <= 0 {
		req.Amount = 5
	}

	input := AgentInput{
		UserMessage: fmt.Sprintf("Please suggest %d recipes from the available options that best match the given constraints.", req.Amount),
		Context: map[string]any{
			"available_recipes": req.AvailableRecipes,
			"constraints":       req.Constraints,
			"already_selected":  req.AlreadySelected,
		},
	}

	var result RecipeSuggestionOutput
	output, err := a.ExecuteWithJSON(ctx, input, &result)
	if err != nil {
		return nil, err
	}

	// Validate returned recipe IDs exist in available recipes
	validRecipes := make(map[string]bool)
	for _, r := range req.AvailableRecipes {
		validRecipes[r.ID] = true
	}

	validatedResult := RecipeSuggestionOutput{
		Recipes: make([]SuggestedRecipe, 0, len(result.Recipes)),
	}

	for _, suggestion := range result.Recipes {
		if validRecipes[suggestion.RecipeID] {
			validatedResult.Recipes = append(validatedResult.Recipes, suggestion)
		}
	}

	// Update output with validated result
	output.Structured = &validatedResult

	return &validatedResult, nil
}

// SuggestFromIDs is a convenience method that accepts UUIDs directly
func (a *RecipeSuggestionsAgent) SuggestFromIDs(ctx context.Context, recipes []RecipeInfo, constraints SuggestionConstraints, alreadySelected []uuid.UUID, amount int) ([]uuid.UUID, error) {
	selectedStrings := make([]string, len(alreadySelected))
	for i, id := range alreadySelected {
		selectedStrings[i] = id.String()
	}

	result, err := a.Suggest(ctx, RecipeSuggestionRequest{
		AvailableRecipes: recipes,
		Constraints:      constraints,
		Amount:           amount,
		AlreadySelected:  selectedStrings,
	})
	if err != nil {
		return nil, err
	}

	ids := make([]uuid.UUID, 0, len(result.Recipes))
	for _, r := range result.Recipes {
		if id, err := uuid.Parse(r.RecipeID); err == nil {
			ids = append(ids, id)
		}
	}

	return ids, nil
}
