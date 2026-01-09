package handler_test

import (
	"testing"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/platepilot/backend/internal/common/domain"
	pb "github.com/platepilot/backend/internal/recipe/pb"
	"github.com/platepilot/backend/internal/recipe/testutil"
)

func TestGetRecipe_ValidID_ReturnsRecipe(t *testing.T) {
	tc := givenRecipeAPI()
	recipe := givenRecipeExists(tc)

	resp, err := tc.Handler.GetRecipe(tc.Ctx, &pb.GetRecipeRequest{
		UserId:   tc.UserID.String(),
		RecipeId: recipe.ID.String(),
	})

	thenNoError(t, err)
	thenRecipeMatches(t, resp, recipe)
}

func TestListRecipes_Pagination_ReturnsTotals(t *testing.T) {
	tc := givenRecipeAPI()
	givenRecipeExists(tc)
	givenRecipeExistsWithName(tc, "Second Recipe")

	resp, err := tc.Handler.ListRecipes(tc.Ctx, &pb.ListRecipesRequest{
		UserId:    tc.UserID.String(),
		PageIndex: 1,
		PageSize:  1,
	})

	thenNoError(t, err)
	if len(resp.GetRecipes()) != 1 {
		t.Fatalf("expected 1 recipe, got %d", len(resp.GetRecipes()))
	}
	if resp.GetTotalCount() != 2 {
		t.Fatalf("expected total count 2, got %d", resp.GetTotalCount())
	}
}

func TestCreateRecipe_ValidInput_PersistsAndPublishes(t *testing.T) {
	tc := givenRecipeAPI()

	resp, err := tc.Handler.CreateRecipe(tc.Ctx, &pb.CreateRecipeRequest{
		UserId: tc.UserID.String(),
		Recipe: &pb.RecipeInput{
			Name:             "Margherita",
			Description:      "Simple and fresh",
			PrepTimeMinutes:  10,
			CookTimeMinutes:  20,
			Servings:         2,
			MainIngredientName: "Tomato",
			CuisineName:        "Italian",
			IngredientLines: []*pb.IngredientLineInput{
				{
					IngredientName: "Tomato",
					QuantityValue:  wrapperspb.Double(2),
					Unit:           "pcs",
					SortOrder:      1,
				},
				{
					IngredientName: "Mozzarella",
					QuantityText:   "to taste",
					SortOrder:      2,
				},
			},
			Steps: []*pb.RecipeStepInput{
				{
					StepIndex:   1,
					Instruction: "Slice tomatoes",
				},
			},
			Tags: []string{"vegetarian"},
			Nutrition: &pb.RecipeNutrition{
				CaloriesTotal:      600,
				CaloriesPerServing: 300,
			},
		},
	})

	thenNoError(t, err)
	if resp.GetId() == "" {
		t.Fatalf("expected recipe id to be set")
	}
	if tc.Publisher.UpsertedEventCount() != 1 {
		t.Fatalf("expected 1 RecipeUpsertedEvent, got %d", tc.Publisher.UpsertedEventCount())
	}
}

func TestDeleteRecipe_NotFound_ReturnsNotFound(t *testing.T) {
	tc := givenRecipeAPI()
	nonExistentID := uuid.New()

	_, err := tc.Handler.DeleteRecipe(tc.Ctx, &pb.DeleteRecipeRequest{
		UserId:   tc.UserID.String(),
		RecipeId: nonExistentID.String(),
	})

	thenErrorHasCode(t, err, codes.NotFound)
}

// Helpers

func givenRecipeAPI() *testutil.TestContext {
	return testutil.NewTestContext()
}

func givenRecipeExists(tc *testutil.TestContext) *domain.Recipe {
	builder := testutil.NewRecipeBuilder().WithUserID(tc.UserID)
	cuisine := testutil.NewCuisineBuilder().WithName("Italian").Build()
	ingredient := testutil.NewIngredientBuilder().WithName("Tomato").Build()
	builder.WithCuisine(cuisine).WithMainIngredient(ingredient)
	recipe := builder.Build()
	tc.Repo.AddCuisine(cuisine)
	tc.Repo.AddIngredient(ingredient)
	tc.Repo.AddRecipe(recipe)
	return recipe
}

func givenRecipeExistsWithName(tc *testutil.TestContext, name string) *domain.Recipe {
	builder := testutil.NewRecipeBuilder().WithUserID(tc.UserID).WithName(name)
	recipe := builder.Build()
	tc.Repo.AddRecipe(recipe)
	return recipe
}

func thenNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func thenErrorHasCode(t *testing.T, err error, code codes.Code) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error with code %v, got nil", code)
	}
	if status.Code(err) != code {
		t.Fatalf("expected code %v, got %v", code, status.Code(err))
	}
}

func thenRecipeMatches(t *testing.T, resp *pb.Recipe, recipe *domain.Recipe) {
	t.Helper()
	if resp.GetId() != recipe.ID.String() {
		t.Fatalf("expected id %s, got %s", recipe.ID.String(), resp.GetId())
	}
	if resp.GetName() != recipe.Name {
		t.Fatalf("expected name %s, got %s", recipe.Name, resp.GetName())
	}
}
