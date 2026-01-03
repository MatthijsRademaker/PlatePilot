package agents

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/platepilot/backend/internal/llm"
)

//go:embed prompts/mealplanner.txt
var mealplannerPrompt string

// MealPlannerAgent creates weekly meal plans based on available recipes and preferences
type MealPlannerAgent struct {
	*BaseAgent
}

// MealPlanRequest contains the input for meal planning
type MealPlanRequest struct {
	// AvailableRecipes is the list of recipes to choose from
	AvailableRecipes []RecipeInfo `json:"available_recipes"`

	// Days is the number of days to plan (1-14)
	Days int `json:"days"`

	// MealTypes specifies which meals to include
	MealTypes []string `json:"meal_types"` // "breakfast", "lunch", "dinner", "snack"

	// Preferences defines user preferences
	Preferences MealPlanPreferences `json:"preferences"`

	// StartDay is the day to start planning from (default "Monday")
	StartDay string `json:"start_day,omitempty"`
}

// MealPlanPreferences defines user preferences for meal planning
type MealPlanPreferences struct {
	// DietaryRestrictions like "vegetarian", "vegan", "gluten-free"
	DietaryRestrictions []string `json:"dietary_restrictions,omitempty"`

	// Allergies to strictly avoid
	Allergies []string `json:"allergies,omitempty"`

	// PreferredCuisines like "Italian", "Mexican", "Asian"
	PreferredCuisines []string `json:"preferred_cuisines,omitempty"`

	// AvoidIngredients lists ingredients to avoid
	AvoidIngredients []string `json:"avoid_ingredients,omitempty"`

	// CookingSkillLevel: "beginner", "intermediate", "advanced"
	CookingSkillLevel string `json:"cooking_skill_level,omitempty"`

	// MaxPrepTimeWeekday in minutes for weekday meals
	MaxPrepTimeWeekday int `json:"max_prep_time_weekday,omitempty"`

	// MaxPrepTimeWeekend in minutes for weekend meals
	MaxPrepTimeWeekend int `json:"max_prep_time_weekend,omitempty"`

	// TargetCaloriesPerDay for the entire day
	TargetCaloriesPerDay int `json:"target_calories_per_day,omitempty"`

	// ServingsPerMeal default number of servings
	ServingsPerMeal int `json:"servings_per_meal,omitempty"`
}

// NewMealPlannerAgent creates a new meal planner agent
func NewMealPlannerAgent(client *llm.Client, cache Cache) *MealPlannerAgent {
	return &MealPlannerAgent{
		BaseAgent: NewBaseAgent("mealplanner", mealplannerPrompt, client, cache),
	}
}

// Execute runs the meal planner agent
func (a *MealPlannerAgent) Execute(ctx context.Context, input AgentInput) (*AgentOutput, error) {
	var result MealPlanOutput
	return a.ExecuteWithJSON(ctx, input, &result)
}

// Plan creates a meal plan based on the request
func (a *MealPlannerAgent) Plan(ctx context.Context, req MealPlanRequest) (*MealPlanOutput, error) {
	if len(req.AvailableRecipes) == 0 {
		return &MealPlanOutput{Days: []MealPlanDay{}}, nil
	}

	// Set defaults
	if req.Days <= 0 {
		req.Days = 7
	}
	if req.Days > 14 {
		req.Days = 14
	}

	if len(req.MealTypes) == 0 {
		req.MealTypes = []string{"breakfast", "lunch", "dinner"}
	}

	if req.StartDay == "" {
		req.StartDay = "Monday"
	}

	input := AgentInput{
		UserMessage: fmt.Sprintf(
			"Create a %d-day meal plan starting from %s. Include these meal types: %v.",
			req.Days, req.StartDay, req.MealTypes,
		),
		Context: map[string]any{
			"available_recipes": req.AvailableRecipes,
			"preferences":       req.Preferences,
			"meal_types":        req.MealTypes,
		},
	}

	var result MealPlanOutput
	output, err := a.ExecuteWithJSON(ctx, input, &result)
	if err != nil {
		return nil, err
	}

	// Validate returned recipe IDs exist in available recipes
	validRecipes := make(map[string]bool)
	for _, r := range req.AvailableRecipes {
		validRecipes[r.ID] = true
	}

	// Filter out invalid recipe IDs
	validatedResult := MealPlanOutput{
		Days: make([]MealPlanDay, 0, len(result.Days)),
	}

	for _, day := range result.Days {
		validMeals := make([]Meal, 0, len(day.Meals))
		for _, meal := range day.Meals {
			if validRecipes[meal.RecipeID] {
				validMeals = append(validMeals, meal)
			}
		}
		if len(validMeals) > 0 {
			validatedResult.Days = append(validatedResult.Days, MealPlanDay{
				Day:   day.Day,
				Meals: validMeals,
			})
		}
	}

	// Update output with validated result
	output.Structured = &validatedResult

	return &validatedResult, nil
}

// PlanWeek is a convenience method for creating a week-long meal plan
func (a *MealPlannerAgent) PlanWeek(ctx context.Context, recipes []RecipeInfo, preferences MealPlanPreferences) (*MealPlanOutput, error) {
	return a.Plan(ctx, MealPlanRequest{
		AvailableRecipes: recipes,
		Days:             7,
		MealTypes:        []string{"breakfast", "lunch", "dinner"},
		Preferences:      preferences,
		StartDay:         "Monday",
	})
}

// GetDayNames returns the day names for a meal plan starting from the given day
func GetDayNames(startDay string, numDays int) []string {
	days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

	// Find start index
	startIdx := 0
	for i, d := range days {
		if d == startDay {
			startIdx = i
			break
		}
	}

	result := make([]string, numDays)
	for i := 0; i < numDays; i++ {
		result[i] = days[(startIdx+i)%7]
	}
	return result
}
