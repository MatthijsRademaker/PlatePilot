package domain_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"

	"github.com/platepilot/backend/internal/mealplanner/domain"
	"github.com/platepilot/backend/internal/mealplanner/repository"
	"github.com/platepilot/backend/internal/mealplanner/testutil"
)

// =============================================================================
// SuggestMeals Tests - Basic Functionality
// =============================================================================

func TestSuggestMeals_NoRecipes_ReturnsEmptyList(t *testing.T) {
	// Given
	tc := givenPlanner()

	// When
	result, err := whenSuggestingMeals(tc, domain.SuggestionRequest{
		Amount: 5,
	})

	// Then
	thenNoError(t, err)
	thenResultIsEmpty(t, result)
}

func TestSuggestMeals_UserScoped_ReturnsOnlyOwnedRecipes(t *testing.T) {
	// Given
	tc := givenPlanner()
	ownedRecipe := givenRecipeExists(tc, "Owned Recipe")
	givenRecipeExistsForUser(tc, "Other User Recipe", uuid.New())

	// When
	result, err := whenSuggestingMeals(tc, domain.SuggestionRequest{
		UserID: tc.UserID,
		Amount: 5,
	})

	// Then
	thenNoError(t, err)
	thenResultHasCount(t, result, 1)
	thenResultContains(t, result, ownedRecipe.ID)
}

func TestSuggestMeals_WithRecipes_ReturnsSuggestions(t *testing.T) {
	// Given
	tc := givenPlanner()
	recipe1 := givenRecipeExists(tc, "Recipe 1")
	recipe2 := givenRecipeExists(tc, "Recipe 2")

	// When
	result, err := whenSuggestingMeals(tc, domain.SuggestionRequest{
		Amount: 5,
	})

	// Then
	thenNoError(t, err)
	thenResultContains(t, result, recipe1.ID)
	thenResultContains(t, result, recipe2.ID)
}

func TestSuggestMeals_AmountLimitsResults(t *testing.T) {
	// Given
	tc := givenPlanner()
	givenRecipeExists(tc, "Recipe 1")
	givenRecipeExists(tc, "Recipe 2")
	givenRecipeExists(tc, "Recipe 3")

	// When
	result, err := whenSuggestingMeals(tc, domain.SuggestionRequest{
		Amount: 2,
	})

	// Then
	thenNoError(t, err)
	thenResultHasCount(t, result, 2)
}

func TestSuggestMeals_FewerRecipesThanAmount_ReturnsAll(t *testing.T) {
	// Given
	tc := givenPlanner()
	givenRecipeExists(tc, "Recipe 1")
	givenRecipeExists(tc, "Recipe 2")

	// When
	result, err := whenSuggestingMeals(tc, domain.SuggestionRequest{
		Amount: 10,
	})

	// Then
	thenNoError(t, err)
	thenResultHasCount(t, result, 2)
}

func TestSuggestMeals_RepositoryError_ReturnsError(t *testing.T) {
	// Given
	tc := givenPlanner()
	givenRepositoryFails(tc)

	// When
	_, err := whenSuggestingMeals(tc, domain.SuggestionRequest{
		Amount: 5,
	})

	// Then
	thenErrorOccurred(t, err)
}

// =============================================================================
// SuggestMeals Tests - Already Selected Filtering
// =============================================================================

func TestSuggestMeals_ExcludesAlreadySelected(t *testing.T) {
	// Given
	tc := givenPlanner()
	recipe1 := givenRecipeExists(tc, "Recipe 1")
	recipe2 := givenRecipeExists(tc, "Recipe 2")

	// When
	result, err := whenSuggestingMeals(tc, domain.SuggestionRequest{
		AlreadySelectedRecipes: []uuid.UUID{recipe1.ID},
		Amount:                 5,
	})

	// Then
	thenNoError(t, err)
	thenResultDoesNotContain(t, result, recipe1.ID)
	thenResultContains(t, result, recipe2.ID)
}

func TestSuggestMeals_AllRecipesSelected_ReturnsEmpty(t *testing.T) {
	// Given
	tc := givenPlanner()
	recipe1 := givenRecipeExists(tc, "Recipe 1")
	recipe2 := givenRecipeExists(tc, "Recipe 2")

	// When
	result, err := whenSuggestingMeals(tc, domain.SuggestionRequest{
		AlreadySelectedRecipes: []uuid.UUID{recipe1.ID, recipe2.ID},
		Amount:                 5,
	})

	// Then
	thenNoError(t, err)
	thenResultIsEmpty(t, result)
}

// =============================================================================
// SuggestMeals Tests - Cuisine Constraints
// =============================================================================

func TestSuggestMeals_CuisineConstraint_FiltersRecipes(t *testing.T) {
	// Given
	tc := givenPlanner()
	italianCuisineID := uuid.New()
	mexicanCuisineID := uuid.New()

	italianRecipe := givenRecipeExistsWithCuisine(tc, "Pasta", italianCuisineID)
	givenRecipeExistsWithCuisine(tc, "Tacos", mexicanCuisineID)

	// When
	result, err := whenSuggestingMeals(tc, domain.SuggestionRequest{
		DailyConstraints: []domain.DailyConstraints{
			{CuisineConstraints: []uuid.UUID{italianCuisineID}},
		},
		Amount: 5,
	})

	// Then
	thenNoError(t, err)
	thenResultHasCount(t, result, 1)
	thenResultContains(t, result, italianRecipe.ID)
}

func TestSuggestMeals_MultipleCuisineConstraints_MatchesAny(t *testing.T) {
	// Given
	tc := givenPlanner()
	italianCuisineID := uuid.New()
	mexicanCuisineID := uuid.New()
	japaneseCuisineID := uuid.New()

	italianRecipe := givenRecipeExistsWithCuisine(tc, "Pasta", italianCuisineID)
	mexicanRecipe := givenRecipeExistsWithCuisine(tc, "Tacos", mexicanCuisineID)
	givenRecipeExistsWithCuisine(tc, "Sushi", japaneseCuisineID)

	// When
	result, err := whenSuggestingMeals(tc, domain.SuggestionRequest{
		DailyConstraints: []domain.DailyConstraints{
			{CuisineConstraints: []uuid.UUID{italianCuisineID, mexicanCuisineID}},
		},
		Amount: 5,
	})

	// Then
	thenNoError(t, err)
	thenResultHasCount(t, result, 2)
	thenResultContains(t, result, italianRecipe.ID)
	thenResultContains(t, result, mexicanRecipe.ID)
}

// =============================================================================
// SuggestMeals Tests - Ingredient Constraints
// =============================================================================

func TestSuggestMeals_MainIngredientConstraint_FiltersRecipes(t *testing.T) {
	// Given
	tc := givenPlanner()
	chickenID := uuid.New()
	beefID := uuid.New()

	chickenRecipe := givenRecipeExistsWithMainIngredient(tc, "Chicken Curry", chickenID)
	givenRecipeExistsWithMainIngredient(tc, "Beef Stew", beefID)

	// When
	result, err := whenSuggestingMeals(tc, domain.SuggestionRequest{
		DailyConstraints: []domain.DailyConstraints{
			{IngredientConstraints: []uuid.UUID{chickenID}},
		},
		Amount: 5,
	})

	// Then
	thenNoError(t, err)
	thenResultHasCount(t, result, 1)
	thenResultContains(t, result, chickenRecipe.ID)
}

func TestSuggestMeals_SecondaryIngredientConstraint_FiltersRecipes(t *testing.T) {
	// Given
	tc := givenPlanner()
	garlicID := uuid.New()
	recipeWithGarlic := givenRecipeExistsWithIngredients(tc, "Garlic Bread", []uuid.UUID{garlicID})
	givenRecipeExists(tc, "Plain Bread")

	// When
	result, err := whenSuggestingMeals(tc, domain.SuggestionRequest{
		DailyConstraints: []domain.DailyConstraints{
			{IngredientConstraints: []uuid.UUID{garlicID}},
		},
		Amount: 5,
	})

	// Then
	thenNoError(t, err)
	thenResultHasCount(t, result, 1)
	thenResultContains(t, result, recipeWithGarlic.ID)
}

// =============================================================================
// SuggestMeals Tests - Multiple Daily Constraints
// =============================================================================

func TestSuggestMeals_MultipleDailyConstraints_MatchesAnyDay(t *testing.T) {
	// Given
	tc := givenPlanner()
	italianCuisineID := uuid.New()
	mexicanCuisineID := uuid.New()
	japaneseCuisineID := uuid.New()

	italianRecipe := givenRecipeExistsWithCuisine(tc, "Pasta", italianCuisineID)
	mexicanRecipe := givenRecipeExistsWithCuisine(tc, "Tacos", mexicanCuisineID)
	givenRecipeExistsWithCuisine(tc, "Sushi", japaneseCuisineID)

	// When - requesting Italian for day 1, Mexican for day 2
	result, err := whenSuggestingMeals(tc, domain.SuggestionRequest{
		DailyConstraints: []domain.DailyConstraints{
			{CuisineConstraints: []uuid.UUID{italianCuisineID}},
			{CuisineConstraints: []uuid.UUID{mexicanCuisineID}},
		},
		Amount: 5,
	})

	// Then - both Italian and Mexican recipes should match
	thenNoError(t, err)
	thenResultHasCount(t, result, 2)
	thenResultContains(t, result, italianRecipe.ID)
	thenResultContains(t, result, mexicanRecipe.ID)
}

func TestSuggestMeals_NoConstraints_ReturnsAllRecipes(t *testing.T) {
	// Given
	tc := givenPlanner()
	recipe1 := givenRecipeExists(tc, "Recipe 1")
	recipe2 := givenRecipeExists(tc, "Recipe 2")

	// When - no constraints
	result, err := whenSuggestingMeals(tc, domain.SuggestionRequest{
		DailyConstraints: []domain.DailyConstraints{},
		Amount:           5,
	})

	// Then
	thenNoError(t, err)
	thenResultContains(t, result, recipe1.ID)
	thenResultContains(t, result, recipe2.ID)
}

// =============================================================================
// SuggestMeals Tests - Combined Constraints
// =============================================================================

func TestSuggestMeals_CuisineAndIngredientConstraints_BothMustMatch(t *testing.T) {
	// Given
	tc := givenPlanner()
	italianCuisineID := uuid.New()
	mexicanCuisineID := uuid.New()
	chickenID := uuid.New()
	beefID := uuid.New()

	// Italian chicken - matches both
	italianChicken := givenRecipeExistsWithCuisineAndIngredient(tc, "Italian Chicken", italianCuisineID, chickenID)
	// Italian beef - only matches cuisine
	givenRecipeExistsWithCuisineAndIngredient(tc, "Italian Beef", italianCuisineID, beefID)
	// Mexican chicken - only matches ingredient
	givenRecipeExistsWithCuisineAndIngredient(tc, "Mexican Chicken", mexicanCuisineID, chickenID)

	// When - require Italian AND chicken
	result, err := whenSuggestingMeals(tc, domain.SuggestionRequest{
		DailyConstraints: []domain.DailyConstraints{
			{
				CuisineConstraints:    []uuid.UUID{italianCuisineID},
				IngredientConstraints: []uuid.UUID{chickenID},
			},
		},
		Amount: 5,
	})

	// Then - only Italian Chicken matches both constraints
	thenNoError(t, err)
	thenResultHasCount(t, result, 1)
	thenResultContains(t, result, italianChicken.ID)
}

// =============================================================================
// SuggestMeals Tests - Diversity Scoring
// =============================================================================

func TestSuggestMeals_DiversityScoring_PrefersDifferentRecipes(t *testing.T) {
	// Given
	tc := givenPlanner()

	// Create recipes with specific vectors for diversity testing
	baseVector := testutil.CreateTestVector(0)
	similarVector := testutil.CreateSimilarVector(baseVector)
	differentVector := testutil.CreateDifferentVector(baseVector)

	selectedRecipe := givenRecipeExistsWithVector(tc, "Selected", baseVector)
	similarRecipe := givenRecipeExistsWithVector(tc, "Similar", similarVector)
	differentRecipe := givenRecipeExistsWithVector(tc, "Different", differentVector)

	// When - with the selected recipe already chosen
	result, err := whenSuggestingMeals(tc, domain.SuggestionRequest{
		AlreadySelectedRecipes: []uuid.UUID{selectedRecipe.ID},
		Amount:                 2,
	})

	// Then - the different recipe should rank higher due to diversity
	thenNoError(t, err)
	thenResultHasCount(t, result, 2)
	thenResultHasFirst(t, result, differentRecipe.ID)
	thenResultContains(t, result, similarRecipe.ID)
}

func TestSuggestMeals_NoPreviousSelection_AllHaveMaxDiversity(t *testing.T) {
	// Given
	tc := givenPlanner()
	recipe1 := givenRecipeExists(tc, "Recipe 1")
	recipe2 := givenRecipeExists(tc, "Recipe 2")

	// When - no recipes previously selected
	result, err := whenSuggestingMeals(tc, domain.SuggestionRequest{
		AlreadySelectedRecipes: []uuid.UUID{},
		Amount:                 5,
	})

	// Then - all recipes have equal diversity score (1.0)
	thenNoError(t, err)
	thenResultContains(t, result, recipe1.ID)
	thenResultContains(t, result, recipe2.ID)
}

// =============================================================================
// Given Helpers (Setup)
// =============================================================================

func givenPlanner() *testutil.PlannerTestContext {
	return testutil.NewPlannerTestContext()
}

func givenRecipeExists(tc *testutil.PlannerTestContext, name string) repository.Recipe {
	recipe := testutil.NewRecipeBuilder().
		WithName(name).
		WithUserID(tc.UserID).
		Build()
	tc.Repo.AddRecipe(recipe)
	return recipe
}

func givenRecipeExistsForUser(tc *testutil.PlannerTestContext, name string, userID uuid.UUID) repository.Recipe {
	recipe := testutil.NewRecipeBuilder().
		WithName(name).
		WithUserID(userID).
		Build()
	tc.Repo.AddRecipe(recipe)
	return recipe
}

func givenRecipeExistsWithCuisine(tc *testutil.PlannerTestContext, name string, cuisineID uuid.UUID) repository.Recipe {
	recipe := testutil.NewRecipeBuilder().
		WithName(name).
		WithCuisineID(cuisineID).
		WithUserID(tc.UserID).
		Build()
	tc.Repo.AddRecipe(recipe)
	return recipe
}

func givenRecipeExistsWithMainIngredient(tc *testutil.PlannerTestContext, name string, ingredientID uuid.UUID) repository.Recipe {
	recipe := testutil.NewRecipeBuilder().
		WithName(name).
		WithMainIngredientID(ingredientID).
		WithUserID(tc.UserID).
		Build()
	tc.Repo.AddRecipe(recipe)
	return recipe
}

func givenRecipeExistsWithIngredients(tc *testutil.PlannerTestContext, name string, ingredientIDs []uuid.UUID) repository.Recipe {
	recipe := testutil.NewRecipeBuilder().
		WithName(name).
		WithIngredientIDs(ingredientIDs).
		WithUserID(tc.UserID).
		Build()
	tc.Repo.AddRecipe(recipe)
	return recipe
}

func givenRecipeExistsWithCuisineAndIngredient(tc *testutil.PlannerTestContext, name string, cuisineID, ingredientID uuid.UUID) repository.Recipe {
	recipe := testutil.NewRecipeBuilder().
		WithName(name).
		WithCuisineID(cuisineID).
		WithMainIngredientID(ingredientID).
		WithUserID(tc.UserID).
		Build()
	tc.Repo.AddRecipe(recipe)
	return recipe
}

func givenRecipeExistsWithVector(tc *testutil.PlannerTestContext, name string, vector pgvector.Vector) repository.Recipe {
	recipe := testutil.NewRecipeBuilder().
		WithName(name).
		WithSearchVector(vector).
		WithUserID(tc.UserID).
		Build()
	tc.Repo.AddRecipe(recipe)
	return recipe
}

func givenRepositoryFails(tc *testutil.PlannerTestContext) {
	tc.Repo.FailOnGetAll = true
}

// =============================================================================
// When Helpers (Action)
// =============================================================================

func whenSuggestingMeals(tc *testutil.PlannerTestContext, req domain.SuggestionRequest) ([]uuid.UUID, error) {
	if req.UserID == uuid.Nil {
		req.UserID = tc.UserID
	}
	return tc.Planner.SuggestMeals(tc.Ctx, req)
}

// =============================================================================
// Then Helpers (Assertions)
// =============================================================================

func thenNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func thenErrorOccurred(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func thenResultIsEmpty(t *testing.T, result []uuid.UUID) {
	t.Helper()
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %d items", len(result))
	}
}

func thenResultHasCount(t *testing.T, result []uuid.UUID, count int) {
	t.Helper()
	if len(result) != count {
		t.Fatalf("expected %d results, got %d", count, len(result))
	}
}

func thenResultContains(t *testing.T, result []uuid.UUID, id uuid.UUID) {
	t.Helper()
	for _, r := range result {
		if r == id {
			return
		}
	}
	t.Fatalf("expected result to contain %s", id)
}

func thenResultDoesNotContain(t *testing.T, result []uuid.UUID, id uuid.UUID) {
	t.Helper()
	for _, r := range result {
		if r == id {
			t.Fatalf("expected result NOT to contain %s", id)
		}
	}
}

func thenResultHasFirst(t *testing.T, result []uuid.UUID, id uuid.UUID) {
	t.Helper()
	if len(result) == 0 {
		t.Fatal("expected non-empty result")
	}
	if result[0] != id {
		t.Fatalf("expected first result to be %s, got %s", id, result[0])
	}
}
