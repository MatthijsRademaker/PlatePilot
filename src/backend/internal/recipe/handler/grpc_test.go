package handler_test

import (
	"testing"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/platepilot/backend/internal/common/domain"
	pb "github.com/platepilot/backend/internal/recipe/pb"
	"github.com/platepilot/backend/internal/recipe/testutil"
)

// =============================================================================
// GetRecipeById Tests
// =============================================================================

func TestGetRecipeById_ValidID_ReturnsRecipe(t *testing.T) {
	// Given
	tc := givenRecipeAPI()
	recipe := givenRecipeExists(tc)

	// When
	resp, err := whenGettingRecipeById(tc, recipe.ID.String())

	// Then
	thenNoError(t, err)
	thenRecipeResponseMatches(t, resp, recipe)
}

func TestGetRecipeById_InvalidUUID_ReturnsInvalidArgument(t *testing.T) {
	// Given
	tc := givenRecipeAPI()

	// When
	_, err := whenGettingRecipeById(tc, "not-a-valid-uuid")

	// Then
	thenErrorHasCode(t, err, codes.InvalidArgument)
}

func TestGetRecipeById_NotFound_ReturnsNotFound(t *testing.T) {
	// Given
	tc := givenRecipeAPI()
	nonExistentID := uuid.New()

	// When
	_, err := whenGettingRecipeById(tc, nonExistentID.String())

	// Then
	thenErrorHasCode(t, err, codes.NotFound)
}

func TestGetRecipeById_RepositoryError_ReturnsInternal(t *testing.T) {
	// Given
	tc := givenRecipeAPI()
	recipe := givenRecipeExists(tc)
	givenRepositoryFailsOnGetByID(tc)

	// When
	_, err := whenGettingRecipeById(tc, recipe.ID.String())

	// Then
	thenErrorHasCode(t, err, codes.Internal)
}

// =============================================================================
// GetAllRecipes Tests
// =============================================================================

func TestGetAllRecipes_WithRecipes_ReturnsAll(t *testing.T) {
	// Given
	tc := givenRecipeAPI()
	recipe1 := givenRecipeExists(tc)
	recipe2 := givenRecipeExistsWithName(tc, "Second Recipe")

	// When
	resp, err := whenGettingAllRecipes(tc, 1, 20)

	// Then
	thenNoError(t, err)
	thenResponseContainsRecipes(t, resp, 2)
	thenResponseContainsRecipeWithID(t, resp, recipe1.ID)
	thenResponseContainsRecipeWithID(t, resp, recipe2.ID)
}

func TestGetAllRecipes_EmptyRepository_ReturnsEmpty(t *testing.T) {
	// Given
	tc := givenRecipeAPI()

	// When
	resp, err := whenGettingAllRecipes(tc, 1, 20)

	// Then
	thenNoError(t, err)
	thenResponseContainsRecipes(t, resp, 0)
}

func TestGetAllRecipes_InvalidPageIndex_DefaultsToFirst(t *testing.T) {
	// Given
	tc := givenRecipeAPI()
	givenRecipeExists(tc)

	// When (pageIndex 0 or negative should default to 1)
	resp, err := whenGettingAllRecipes(tc, 0, 20)

	// Then
	thenNoError(t, err)
	thenResponseContainsRecipes(t, resp, 1)
}

func TestGetAllRecipes_InvalidPageSize_DefaultsToTwenty(t *testing.T) {
	// Given
	tc := givenRecipeAPI()
	givenRecipeExists(tc)

	// When (pageSize 0 or > 100 should default to 20)
	resp, err := whenGettingAllRecipes(tc, 1, 0)

	// Then
	thenNoError(t, err)
	thenResponseContainsRecipes(t, resp, 1)
}

func TestGetAllRecipes_RepositoryError_ReturnsInternal(t *testing.T) {
	// Given
	tc := givenRecipeAPI()
	givenRepositoryFailsOnGetAll(tc)

	// When
	_, err := whenGettingAllRecipes(tc, 1, 20)

	// Then
	thenErrorHasCode(t, err, codes.Internal)
}

// =============================================================================
// CreateRecipe Tests
// =============================================================================

func TestCreateRecipe_ValidInput_PersistsAndPublishesEvent(t *testing.T) {
	// Given
	tc := givenRecipeAPI()
	mainIngredient := givenIngredientExists(tc, "Chicken")
	cuisine := givenCuisineExists(tc, "Italian")
	ingredient1 := givenIngredientExists(tc, "Salt")
	ingredient2 := givenIngredientExists(tc, "Pepper")

	req := &pb.CreateRecipeRequest{
		Name:             "New Recipe",
		Description:      "A delicious recipe",
		PrepTime:         "10 mins",
		CookTime:         "20 mins",
		MainIngredientId: mainIngredient.ID.String(),
		CuisineId:        cuisine.ID.String(),
		IngredientIds:    []string{ingredient1.ID.String(), ingredient2.ID.String()},
		Directions:       []string{"Step 1", "Step 2"},
	}

	// When
	resp, err := whenCreatingRecipe(tc, req)

	// Then
	thenNoError(t, err)
	thenRecipeIsPersisted(t, tc)
	thenRecipeCreatedEventIsPublished(t, tc)
	thenResponseHasRecipeName(t, resp, "New Recipe")
}

func TestCreateRecipe_MissingName_ReturnsValidationError(t *testing.T) {
	// Given
	tc := givenRecipeAPI()
	mainIngredient := givenIngredientExists(tc, "Chicken")
	cuisine := givenCuisineExists(tc, "Italian")

	req := &pb.CreateRecipeRequest{
		Name:             "", // Missing name
		MainIngredientId: mainIngredient.ID.String(),
		CuisineId:        cuisine.ID.String(),
	}

	// When
	_, err := whenCreatingRecipe(tc, req)

	// Then
	thenErrorHasCode(t, err, codes.InvalidArgument)
	thenNoRecipeIsPersisted(t, tc)
	thenNoEventIsPublished(t, tc)
}

func TestCreateRecipe_InvalidMainIngredientID_ReturnsInvalidArgument(t *testing.T) {
	// Given
	tc := givenRecipeAPI()

	req := &pb.CreateRecipeRequest{
		Name:             "Test Recipe",
		MainIngredientId: "not-a-valid-uuid",
		CuisineId:        uuid.New().String(),
	}

	// When
	_, err := whenCreatingRecipe(tc, req)

	// Then
	thenErrorHasCode(t, err, codes.InvalidArgument)
	thenNoRecipeIsPersisted(t, tc)
}

func TestCreateRecipe_InvalidCuisineID_ReturnsInvalidArgument(t *testing.T) {
	// Given
	tc := givenRecipeAPI()

	req := &pb.CreateRecipeRequest{
		Name:             "Test Recipe",
		MainIngredientId: uuid.New().String(),
		CuisineId:        "not-a-valid-uuid",
	}

	// When
	_, err := whenCreatingRecipe(tc, req)

	// Then
	thenErrorHasCode(t, err, codes.InvalidArgument)
	thenNoRecipeIsPersisted(t, tc)
}

func TestCreateRecipe_MainIngredientNotFound_ReturnsNotFound(t *testing.T) {
	// Given
	tc := givenRecipeAPI()
	cuisine := givenCuisineExists(tc, "Italian")

	req := &pb.CreateRecipeRequest{
		Name:             "Test Recipe",
		MainIngredientId: uuid.New().String(), // Non-existent
		CuisineId:        cuisine.ID.String(),
	}

	// When
	_, err := whenCreatingRecipe(tc, req)

	// Then
	thenErrorHasCode(t, err, codes.NotFound)
	thenNoRecipeIsPersisted(t, tc)
}

func TestCreateRecipe_CuisineNotFound_ReturnsNotFound(t *testing.T) {
	// Given
	tc := givenRecipeAPI()
	mainIngredient := givenIngredientExists(tc, "Chicken")

	req := &pb.CreateRecipeRequest{
		Name:             "Test Recipe",
		MainIngredientId: mainIngredient.ID.String(),
		CuisineId:        uuid.New().String(), // Non-existent
	}

	// When
	_, err := whenCreatingRecipe(tc, req)

	// Then
	thenErrorHasCode(t, err, codes.NotFound)
	thenNoRecipeIsPersisted(t, tc)
}

func TestCreateRecipe_IngredientNotFound_ReturnsNotFound(t *testing.T) {
	// Given
	tc := givenRecipeAPI()
	mainIngredient := givenIngredientExists(tc, "Chicken")
	cuisine := givenCuisineExists(tc, "Italian")

	req := &pb.CreateRecipeRequest{
		Name:             "Test Recipe",
		MainIngredientId: mainIngredient.ID.String(),
		CuisineId:        cuisine.ID.String(),
		IngredientIds:    []string{uuid.New().String()}, // Non-existent
	}

	// When
	_, err := whenCreatingRecipe(tc, req)

	// Then
	thenErrorHasCode(t, err, codes.NotFound)
	thenNoRecipeIsPersisted(t, tc)
}

func TestCreateRecipe_RepositoryError_ReturnsInternal(t *testing.T) {
	// Given
	tc := givenRecipeAPI()
	mainIngredient := givenIngredientExists(tc, "Chicken")
	cuisine := givenCuisineExists(tc, "Italian")
	givenRepositoryFailsOnCreate(tc)

	req := &pb.CreateRecipeRequest{
		Name:             "Test Recipe",
		MainIngredientId: mainIngredient.ID.String(),
		CuisineId:        cuisine.ID.String(),
	}

	// When
	_, err := whenCreatingRecipe(tc, req)

	// Then
	thenErrorHasCode(t, err, codes.Internal)
	thenNoEventIsPublished(t, tc)
}

func TestCreateRecipe_PublisherError_StillSucceeds(t *testing.T) {
	// Given
	tc := givenRecipeAPI()
	mainIngredient := givenIngredientExists(tc, "Chicken")
	cuisine := givenCuisineExists(tc, "Italian")
	givenPublisherFails(tc)

	req := &pb.CreateRecipeRequest{
		Name:             "Test Recipe",
		MainIngredientId: mainIngredient.ID.String(),
		CuisineId:        cuisine.ID.String(),
	}

	// When
	resp, err := whenCreatingRecipe(tc, req)

	// Then (event publishing failure should not fail the request)
	thenNoError(t, err)
	thenRecipeIsPersisted(t, tc)
	thenResponseHasRecipeName(t, resp, "Test Recipe")
}

func TestCreateRecipe_NoPublisher_StillSucceeds(t *testing.T) {
	// Given
	tc := givenRecipeAPIWithoutPublisher()
	mainIngredient := givenIngredientExists(tc, "Chicken")
	cuisine := givenCuisineExists(tc, "Italian")

	req := &pb.CreateRecipeRequest{
		Name:             "Test Recipe",
		MainIngredientId: mainIngredient.ID.String(),
		CuisineId:        cuisine.ID.String(),
	}

	// When
	resp, err := whenCreatingRecipe(tc, req)

	// Then
	thenNoError(t, err)
	thenRecipeIsPersisted(t, tc)
	thenResponseHasRecipeName(t, resp, "Test Recipe")
}

// =============================================================================
// GetSimilarRecipes Tests
// =============================================================================

func TestGetSimilarRecipes_ValidRecipe_ReturnsSimilar(t *testing.T) {
	// Given
	tc := givenRecipeAPI()
	recipe := givenRecipeExists(tc)
	givenRecipeExistsWithName(tc, "Similar Recipe 1")
	givenRecipeExistsWithName(tc, "Similar Recipe 2")

	// When
	resp, err := whenGettingSimilarRecipes(tc, recipe.ID.String(), 5)

	// Then
	thenNoError(t, err)
	thenResponseContainsRecipes(t, resp, 2)
}

func TestGetSimilarRecipes_InvalidID_ReturnsInvalidArgument(t *testing.T) {
	// Given
	tc := givenRecipeAPI()

	// When
	_, err := whenGettingSimilarRecipes(tc, "not-a-valid-uuid", 5)

	// Then
	thenErrorHasCode(t, err, codes.InvalidArgument)
}

func TestGetSimilarRecipes_RecipeNotFound_ReturnsNotFound(t *testing.T) {
	// Given
	tc := givenRecipeAPI()

	// When
	_, err := whenGettingSimilarRecipes(tc, uuid.New().String(), 5)

	// Then
	thenErrorHasCode(t, err, codes.NotFound)
}

func TestGetSimilarRecipes_DefaultAmount_ReturnsFive(t *testing.T) {
	// Given
	tc := givenRecipeAPI()
	recipe := givenRecipeExists(tc)

	// When (amount 0 should default to 5)
	_, err := whenGettingSimilarRecipes(tc, recipe.ID.String(), 0)

	// Then
	thenNoError(t, err)
}

func TestGetSimilarRecipes_ExcessiveAmount_CapsAtFifty(t *testing.T) {
	// Given
	tc := givenRecipeAPI()
	recipe := givenRecipeExists(tc)

	// When (amount > 50 should cap at 50)
	_, err := whenGettingSimilarRecipes(tc, recipe.ID.String(), 100)

	// Then
	thenNoError(t, err)
}

// =============================================================================
// GetRecipesByCuisine Tests
// =============================================================================

func TestGetRecipesByCuisine_ValidCuisine_ReturnsRecipes(t *testing.T) {
	// Given
	tc := givenRecipeAPI()
	cuisine := givenCuisineExists(tc, "Italian")
	givenRecipeExistsWithCuisine(tc, "Pasta", cuisine)
	givenRecipeExistsWithCuisine(tc, "Pizza", cuisine)

	// When
	resp, err := whenGettingRecipesByCuisine(tc, cuisine.ID.String())

	// Then
	thenNoError(t, err)
	thenResponseContainsRecipes(t, resp, 2)
}

func TestGetRecipesByCuisine_InvalidID_ReturnsInvalidArgument(t *testing.T) {
	// Given
	tc := givenRecipeAPI()

	// When
	_, err := whenGettingRecipesByCuisine(tc, "not-a-valid-uuid")

	// Then
	thenErrorHasCode(t, err, codes.InvalidArgument)
}

// =============================================================================
// GetRecipesByIngredient Tests
// =============================================================================

func TestGetRecipesByIngredient_ValidIngredient_ReturnsRecipes(t *testing.T) {
	// Given
	tc := givenRecipeAPI()
	ingredient := givenIngredientExists(tc, "Tomato")
	givenRecipeExistsWithMainIngredient(tc, "Tomato Soup", ingredient)

	// When
	resp, err := whenGettingRecipesByIngredient(tc, ingredient.ID.String())

	// Then
	thenNoError(t, err)
	thenResponseContainsRecipes(t, resp, 1)
}

func TestGetRecipesByIngredient_InvalidID_ReturnsInvalidArgument(t *testing.T) {
	// Given
	tc := givenRecipeAPI()

	// When
	_, err := whenGettingRecipesByIngredient(tc, "not-a-valid-uuid")

	// Then
	thenErrorHasCode(t, err, codes.InvalidArgument)
}

// =============================================================================
// GetRecipesByAllergy Tests
// =============================================================================

func TestGetRecipesByAllergy_ValidAllergy_ReturnsRecipes(t *testing.T) {
	// Given
	tc := givenRecipeAPI()
	givenRecipeExists(tc)

	// When
	resp, err := whenGettingRecipesByAllergy(tc, uuid.New().String())

	// Then
	thenNoError(t, err)
	// Returns recipes that DON'T contain the allergy
	thenResponseContainsRecipes(t, resp, 1)
}

func TestGetRecipesByAllergy_InvalidID_ReturnsInvalidArgument(t *testing.T) {
	// Given
	tc := givenRecipeAPI()

	// When
	_, err := whenGettingRecipesByAllergy(tc, "not-a-valid-uuid")

	// Then
	thenErrorHasCode(t, err, codes.InvalidArgument)
}

// =============================================================================
// Given Helpers (Setup)
// =============================================================================

func givenRecipeAPI() *testutil.TestContext {
	return testutil.NewTestContext()
}

func givenRecipeAPIWithoutPublisher() *testutil.TestContext {
	return testutil.NewTestContextWithoutPublisher()
}

func givenRecipeExists(tc *testutil.TestContext) *domain.Recipe {
	mainIngredient := testutil.NewIngredientBuilder().
		WithName("Main Ingredient").
		Build()
	tc.Repo.AddIngredient(mainIngredient)

	cuisine := testutil.NewCuisineBuilder().
		WithName("Test Cuisine").
		Build()
	tc.Repo.AddCuisine(cuisine)

	recipe := testutil.NewRecipeBuilder().
		WithName("Test Recipe").
		WithMainIngredient(mainIngredient).
		WithCuisine(cuisine).
		Build()
	tc.Repo.AddRecipe(recipe)

	return recipe
}

func givenRecipeExistsWithName(tc *testutil.TestContext, name string) *domain.Recipe {
	mainIngredient := testutil.NewIngredientBuilder().Build()
	tc.Repo.AddIngredient(mainIngredient)

	cuisine := testutil.NewCuisineBuilder().Build()
	tc.Repo.AddCuisine(cuisine)

	recipe := testutil.NewRecipeBuilder().
		WithName(name).
		WithMainIngredient(mainIngredient).
		WithCuisine(cuisine).
		Build()
	tc.Repo.AddRecipe(recipe)

	return recipe
}

func givenRecipeExistsWithCuisine(tc *testutil.TestContext, name string, cuisine *domain.Cuisine) *domain.Recipe {
	mainIngredient := testutil.NewIngredientBuilder().Build()
	tc.Repo.AddIngredient(mainIngredient)

	recipe := testutil.NewRecipeBuilder().
		WithName(name).
		WithMainIngredient(mainIngredient).
		WithCuisine(cuisine).
		Build()
	tc.Repo.AddRecipe(recipe)

	return recipe
}

func givenRecipeExistsWithMainIngredient(tc *testutil.TestContext, name string, mainIngredient *domain.Ingredient) *domain.Recipe {
	cuisine := testutil.NewCuisineBuilder().Build()
	tc.Repo.AddCuisine(cuisine)

	recipe := testutil.NewRecipeBuilder().
		WithName(name).
		WithMainIngredient(mainIngredient).
		WithCuisine(cuisine).
		Build()
	tc.Repo.AddRecipe(recipe)

	return recipe
}

func givenIngredientExists(tc *testutil.TestContext, name string) *domain.Ingredient {
	ingredient := testutil.NewIngredientBuilder().
		WithName(name).
		Build()
	tc.Repo.AddIngredient(ingredient)
	return ingredient
}

func givenCuisineExists(tc *testutil.TestContext, name string) *domain.Cuisine {
	cuisine := testutil.NewCuisineBuilder().
		WithName(name).
		Build()
	tc.Repo.AddCuisine(cuisine)
	return cuisine
}

func givenRepositoryFailsOnGetByID(tc *testutil.TestContext) {
	tc.Repo.FailOnGetByID = true
}

func givenRepositoryFailsOnGetAll(tc *testutil.TestContext) {
	tc.Repo.FailOnGetAll = true
}

func givenRepositoryFailsOnCreate(tc *testutil.TestContext) {
	tc.Repo.FailOnCreate = true
}

func givenPublisherFails(tc *testutil.TestContext) {
	tc.Publisher.FailOnPublishCreated = true
}

// =============================================================================
// When Helpers (Action)
// =============================================================================

func whenGettingRecipeById(tc *testutil.TestContext, id string) (*pb.RecipeResponse, error) {
	return tc.Handler.GetRecipeById(tc.Ctx, &pb.GetRecipeByIdRequest{
		RecipeId: id,
	})
}

func whenGettingAllRecipes(tc *testutil.TestContext, pageIndex, pageSize int32) (*pb.GetAllRecipesResponse, error) {
	return tc.Handler.GetAllRecipes(tc.Ctx, &pb.GetAllRecipesRequest{
		PageIndex: pageIndex,
		PageSize:  pageSize,
	})
}

func whenCreatingRecipe(tc *testutil.TestContext, req *pb.CreateRecipeRequest) (*pb.RecipeResponse, error) {
	return tc.Handler.CreateRecipe(tc.Ctx, req)
}

func whenGettingSimilarRecipes(tc *testutil.TestContext, recipeID string, amount int32) (*pb.GetAllRecipesResponse, error) {
	return tc.Handler.GetSimilarRecipes(tc.Ctx, &pb.GetSimilarRecipesRequest{
		RecipeId: recipeID,
		Amount:   amount,
	})
}

func whenGettingRecipesByCuisine(tc *testutil.TestContext, cuisineID string) (*pb.GetAllRecipesResponse, error) {
	return tc.Handler.GetRecipesByCuisine(tc.Ctx, &pb.GetRecipesByCuisineRequest{
		CuisineId: cuisineID,
	})
}

func whenGettingRecipesByIngredient(tc *testutil.TestContext, ingredientID string) (*pb.GetAllRecipesResponse, error) {
	return tc.Handler.GetRecipesByIngredient(tc.Ctx, &pb.GetRecipesByIngredientRequest{
		IngredientId: ingredientID,
	})
}

func whenGettingRecipesByAllergy(tc *testutil.TestContext, allergyID string) (*pb.GetAllRecipesResponse, error) {
	return tc.Handler.GetRecipesByAllergy(tc.Ctx, &pb.GetRecipesByAllergyRequest{
		AllergyId: allergyID,
	})
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

func thenRecipeResponseMatches(t *testing.T, resp *pb.RecipeResponse, recipe *domain.Recipe) {
	t.Helper()
	if resp == nil {
		t.Fatal("expected response, got nil")
	}
	if resp.Id != recipe.ID.String() {
		t.Fatalf("expected ID %s, got %s", recipe.ID.String(), resp.Id)
	}
	if resp.Name != recipe.Name {
		t.Fatalf("expected name %s, got %s", recipe.Name, resp.Name)
	}
}

func thenResponseContainsRecipes(t *testing.T, resp *pb.GetAllRecipesResponse, count int) {
	t.Helper()
	if resp == nil {
		t.Fatal("expected response, got nil")
	}
	if len(resp.Recipes) != count {
		t.Fatalf("expected %d recipes, got %d", count, len(resp.Recipes))
	}
}

func thenResponseContainsRecipeWithID(t *testing.T, resp *pb.GetAllRecipesResponse, id uuid.UUID) {
	t.Helper()
	for _, r := range resp.Recipes {
		if r.Id == id.String() {
			return
		}
	}
	t.Fatalf("expected response to contain recipe with ID %s", id)
}

func thenRecipeIsPersisted(t *testing.T, tc *testutil.TestContext) {
	t.Helper()
	if len(tc.Repo.CreateCalls) == 0 {
		t.Fatal("expected recipe to be persisted, but no Create calls were made")
	}
}

func thenNoRecipeIsPersisted(t *testing.T, tc *testutil.TestContext) {
	t.Helper()
	if len(tc.Repo.CreateCalls) > 0 {
		t.Fatal("expected no recipe to be persisted, but Create was called")
	}
}

func thenRecipeCreatedEventIsPublished(t *testing.T, tc *testutil.TestContext) {
	t.Helper()
	if tc.Publisher == nil {
		return // No publisher configured
	}
	if tc.Publisher.CreatedEventCount() != 1 {
		t.Fatalf("expected 1 RecipeCreatedEvent, got %d", tc.Publisher.CreatedEventCount())
	}
}

func thenNoEventIsPublished(t *testing.T, tc *testutil.TestContext) {
	t.Helper()
	if tc.Publisher == nil {
		return // No publisher configured
	}
	if tc.Publisher.CreatedEventCount() > 0 {
		t.Fatal("expected no events to be published")
	}
}

func thenResponseHasRecipeName(t *testing.T, resp *pb.RecipeResponse, name string) {
	t.Helper()
	if resp == nil {
		t.Fatal("expected response, got nil")
	}
	if resp.Name != name {
		t.Fatalf("expected recipe name %s, got %s", name, resp.Name)
	}
}
