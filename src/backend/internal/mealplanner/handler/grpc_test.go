package handler_test

import (
	"testing"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/platepilot/backend/internal/mealplanner/pb"
	"github.com/platepilot/backend/internal/mealplanner/testutil"
)

// =============================================================================
// SuggestRecipes Tests
// =============================================================================

func TestSuggestRecipes_ValidRequest_ReturnsSuggestions(t *testing.T) {
	// Given
	tc := givenMealPlannerAPI()
	id1, id2 := uuid.New(), uuid.New()
	givenPlannerWillSuggest(tc, id1, id2)

	// When
	resp, err := whenRequestingSuggestions(tc, &pb.SuggestionsRequest{
		Amount: 5,
	})

	// Then
	thenNoError(t, err)
	thenResponseContainsSuggestions(t, resp, 2)
	thenResponseContainsRecipeID(t, resp, id1)
	thenResponseContainsRecipeID(t, resp, id2)
}

func TestSuggestRecipes_WithDailyConstraints_PassesConstraintsToPlanner(t *testing.T) {
	// Given
	tc := givenMealPlannerAPI()
	cuisineID := uuid.New()
	ingredientID := uuid.New()

	// When
	_, err := whenRequestingSuggestions(tc, &pb.SuggestionsRequest{
		DailyConstraints: []*pb.DailyConstraints{
			{
				CuisineConstraints: []*pb.CuisineConstraint{
					{EntityId: cuisineID.String()},
				},
				IngredientConstraints: []*pb.IngredientConstraint{
					{EntityId: ingredientID.String()},
				},
			},
		},
		Amount: 5,
	})

	// Then
	thenNoError(t, err)
	thenPlannerWasCalled(t, tc)
	thenPlannerReceivedConstraints(t, tc, 1)
}

func TestSuggestRecipes_WithAlreadySelected_PassesToPlanner(t *testing.T) {
	// Given
	tc := givenMealPlannerAPI()
	selectedID := uuid.New()

	// When
	_, err := whenRequestingSuggestions(tc, &pb.SuggestionsRequest{
		AlreadySelectedRecipeIds: []string{selectedID.String()},
		Amount:                   5,
	})

	// Then
	thenNoError(t, err)
	thenPlannerReceivedAlreadySelected(t, tc, 1)
}

func TestSuggestRecipes_InvalidSelectedRecipeID_ReturnsInvalidArgument(t *testing.T) {
	// Given
	tc := givenMealPlannerAPI()

	// When
	_, err := whenRequestingSuggestions(tc, &pb.SuggestionsRequest{
		AlreadySelectedRecipeIds: []string{"not-a-valid-uuid"},
		Amount:                   5,
	})

	// Then
	thenErrorHasCode(t, err, codes.InvalidArgument)
	thenPlannerWasNotCalled(t, tc)
}

func TestSuggestRecipes_InvalidCuisineConstraintID_ReturnsInvalidArgument(t *testing.T) {
	// Given
	tc := givenMealPlannerAPI()

	// When
	_, err := whenRequestingSuggestions(tc, &pb.SuggestionsRequest{
		DailyConstraints: []*pb.DailyConstraints{
			{
				CuisineConstraints: []*pb.CuisineConstraint{
					{EntityId: "invalid-uuid"},
				},
			},
		},
		Amount: 5,
	})

	// Then
	thenErrorHasCode(t, err, codes.InvalidArgument)
	thenPlannerWasNotCalled(t, tc)
}

func TestSuggestRecipes_InvalidIngredientConstraintID_ReturnsInvalidArgument(t *testing.T) {
	// Given
	tc := givenMealPlannerAPI()

	// When
	_, err := whenRequestingSuggestions(tc, &pb.SuggestionsRequest{
		DailyConstraints: []*pb.DailyConstraints{
			{
				IngredientConstraints: []*pb.IngredientConstraint{
					{EntityId: "invalid-uuid"},
				},
			},
		},
		Amount: 5,
	})

	// Then
	thenErrorHasCode(t, err, codes.InvalidArgument)
	thenPlannerWasNotCalled(t, tc)
}

func TestSuggestRecipes_ZeroAmount_DefaultsToFive(t *testing.T) {
	// Given
	tc := givenMealPlannerAPI()

	// When
	_, err := whenRequestingSuggestions(tc, &pb.SuggestionsRequest{
		Amount: 0,
	})

	// Then
	thenNoError(t, err)
	thenPlannerReceivedAmount(t, tc, 5)
}

func TestSuggestRecipes_NegativeAmount_DefaultsToFive(t *testing.T) {
	// Given
	tc := givenMealPlannerAPI()

	// When
	_, err := whenRequestingSuggestions(tc, &pb.SuggestionsRequest{
		Amount: -10,
	})

	// Then
	thenNoError(t, err)
	thenPlannerReceivedAmount(t, tc, 5)
}

func TestSuggestRecipes_ExcessiveAmount_CapsAtFifty(t *testing.T) {
	// Given
	tc := givenMealPlannerAPI()

	// When
	_, err := whenRequestingSuggestions(tc, &pb.SuggestionsRequest{
		Amount: 100,
	})

	// Then
	thenNoError(t, err)
	thenPlannerReceivedAmount(t, tc, 50)
}

func TestSuggestRecipes_PlannerError_ReturnsInternal(t *testing.T) {
	// Given
	tc := givenMealPlannerAPI()
	givenPlannerFails(tc)

	// When
	_, err := whenRequestingSuggestions(tc, &pb.SuggestionsRequest{
		Amount: 5,
	})

	// Then
	thenErrorHasCode(t, err, codes.Internal)
}

func TestSuggestRecipes_EmptyResult_ReturnsEmptyList(t *testing.T) {
	// Given
	tc := givenMealPlannerAPI()
	// Planner returns empty by default

	// When
	resp, err := whenRequestingSuggestions(tc, &pb.SuggestionsRequest{
		Amount: 5,
	})

	// Then
	thenNoError(t, err)
	thenResponseContainsSuggestions(t, resp, 0)
}

func TestSuggestRecipes_MultipleDailyConstraints_AllParsedCorrectly(t *testing.T) {
	// Given
	tc := givenMealPlannerAPI()
	cuisine1, cuisine2 := uuid.New(), uuid.New()
	ingredient1 := uuid.New()

	// When
	_, err := whenRequestingSuggestions(tc, &pb.SuggestionsRequest{
		DailyConstraints: []*pb.DailyConstraints{
			{
				CuisineConstraints: []*pb.CuisineConstraint{
					{EntityId: cuisine1.String()},
				},
			},
			{
				CuisineConstraints: []*pb.CuisineConstraint{
					{EntityId: cuisine2.String()},
				},
				IngredientConstraints: []*pb.IngredientConstraint{
					{EntityId: ingredient1.String()},
				},
			},
		},
		Amount: 5,
	})

	// Then
	thenNoError(t, err)
	thenPlannerReceivedConstraints(t, tc, 2)
}

// =============================================================================
// Given Helpers (Setup)
// =============================================================================

func givenMealPlannerAPI() *testutil.HandlerTestContext {
	return testutil.NewHandlerTestContext()
}

func givenPlannerWillSuggest(tc *testutil.HandlerTestContext, ids ...uuid.UUID) {
	tc.Planner.SetSuggestedRecipes(ids...)
}

func givenPlannerFails(tc *testutil.HandlerTestContext) {
	tc.Planner.FailOnSuggestMeals = true
}

// =============================================================================
// When Helpers (Action)
// =============================================================================

func whenRequestingSuggestions(tc *testutil.HandlerTestContext, req *pb.SuggestionsRequest) (*pb.SuggestionsResponse, error) {
	return tc.Handler.SuggestRecipes(tc.Ctx, req)
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

func thenErrorHasCode(t *testing.T, err error, expectedCode codes.Code) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error with code %v, got nil", expectedCode)
	}
	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("expected gRPC status error, got %v", err)
	}
	if st.Code() != expectedCode {
		t.Fatalf("expected error code %v, got %v", expectedCode, st.Code())
	}
}

func thenResponseContainsSuggestions(t *testing.T, resp *pb.SuggestionsResponse, count int) {
	t.Helper()
	if resp == nil {
		t.Fatal("expected response, got nil")
	}
	if len(resp.RecipeIds) != count {
		t.Fatalf("expected %d suggestions, got %d", count, len(resp.RecipeIds))
	}
}

func thenResponseContainsRecipeID(t *testing.T, resp *pb.SuggestionsResponse, id uuid.UUID) {
	t.Helper()
	for _, recipeID := range resp.RecipeIds {
		if recipeID == id.String() {
			return
		}
	}
	t.Fatalf("expected response to contain recipe ID %s", id)
}

func thenPlannerWasCalled(t *testing.T, tc *testutil.HandlerTestContext) {
	t.Helper()
	if tc.Planner.SuggestMealsCallCount() == 0 {
		t.Fatal("expected planner to be called, but it was not")
	}
}

func thenPlannerWasNotCalled(t *testing.T, tc *testutil.HandlerTestContext) {
	t.Helper()
	if tc.Planner.SuggestMealsCallCount() > 0 {
		t.Fatal("expected planner not to be called, but it was")
	}
}

func thenPlannerReceivedConstraints(t *testing.T, tc *testutil.HandlerTestContext, count int) {
	t.Helper()
	if len(tc.Planner.SuggestMealsCalls) == 0 {
		t.Fatal("expected planner to be called")
	}
	lastCall := tc.Planner.SuggestMealsCalls[len(tc.Planner.SuggestMealsCalls)-1]
	if len(lastCall.DailyConstraints) != count {
		t.Fatalf("expected %d daily constraints, got %d", count, len(lastCall.DailyConstraints))
	}
}

func thenPlannerReceivedAlreadySelected(t *testing.T, tc *testutil.HandlerTestContext, count int) {
	t.Helper()
	if len(tc.Planner.SuggestMealsCalls) == 0 {
		t.Fatal("expected planner to be called")
	}
	lastCall := tc.Planner.SuggestMealsCalls[len(tc.Planner.SuggestMealsCalls)-1]
	if len(lastCall.AlreadySelectedRecipes) != count {
		t.Fatalf("expected %d already selected recipes, got %d", count, len(lastCall.AlreadySelectedRecipes))
	}
}

func thenPlannerReceivedAmount(t *testing.T, tc *testutil.HandlerTestContext, expectedAmount int) {
	t.Helper()
	if len(tc.Planner.SuggestMealsCalls) == 0 {
		t.Fatal("expected planner to be called")
	}
	lastCall := tc.Planner.SuggestMealsCalls[len(tc.Planner.SuggestMealsCalls)-1]
	if lastCall.Amount != expectedAmount {
		t.Fatalf("expected amount %d, got %d", expectedAmount, lastCall.Amount)
	}
}
