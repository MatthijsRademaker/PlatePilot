package client

import (
	"context"
	"fmt"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	recipepb "github.com/platepilot/backend/internal/recipe/pb"
)

// RecipeClient wraps the gRPC client for the Recipe API
type RecipeClient struct {
	conn   *grpc.ClientConn
	client recipepb.RecipeServiceClient
	logger *slog.Logger
}

// NewRecipeClient creates a new Recipe API client
func NewRecipeClient(address string, logger *slog.Logger) (*RecipeClient, error) {
	conn, err := grpc.NewClient(address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("connect to recipe api: %w", err)
	}

	return &RecipeClient{
		conn:   conn,
		client: recipepb.NewRecipeServiceClient(conn),
		logger: logger,
	}, nil
}

// Close closes the gRPC connection
func (c *RecipeClient) Close() error {
	return c.conn.Close()
}

// GetByID retrieves a recipe by its ID
func (c *RecipeClient) GetByID(ctx context.Context, userID, id string) (*recipepb.RecipeResponse, error) {
	c.logger.Debug("getting recipe by id", "id", id, "userId", userID)

	resp, err := c.client.GetRecipeById(ctx, &recipepb.GetRecipeByIdRequest{
		RecipeId: id,
		UserId:   userID,
	})
	if err != nil {
		return nil, fmt.Errorf("get recipe by id: %w", err)
	}

	return resp, nil
}

// GetAllResponse contains the paginated recipes response
type GetAllResponse struct {
	Recipes    []*recipepb.RecipeResponse
	PageIndex  int32
	PageSize   int32
	TotalCount int32
	TotalPages int32
}

// GetAll retrieves all recipes with pagination
func (c *RecipeClient) GetAll(ctx context.Context, userID string, pageIndex, pageSize int32) (*GetAllResponse, error) {
	c.logger.Debug("getting all recipes", "pageIndex", pageIndex, "pageSize", pageSize, "userId", userID)

	resp, err := c.client.GetAllRecipes(ctx, &recipepb.GetAllRecipesRequest{
		UserId:    userID,
		PageIndex: pageIndex,
		PageSize:  pageSize,
	})
	if err != nil {
		return nil, fmt.Errorf("get all recipes: %w", err)
	}

	return &GetAllResponse{
		Recipes:    resp.GetRecipes(),
		PageIndex:  resp.GetPageIndex(),
		PageSize:   resp.GetPageSize(),
		TotalCount: resp.GetTotalCount(),
		TotalPages: resp.GetTotalPages(),
	}, nil
}

// Create creates a new recipe
func (c *RecipeClient) Create(ctx context.Context, req *recipepb.CreateRecipeRequest) (*recipepb.RecipeResponse, error) {
	c.logger.Debug("creating recipe", "name", req.GetName())

	resp, err := c.client.CreateRecipe(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("create recipe: %w", err)
	}

	return resp, nil
}

// GetSimilar retrieves recipes similar to a given recipe
func (c *RecipeClient) GetSimilar(ctx context.Context, userID, recipeID string, amount int32) ([]*recipepb.RecipeResponse, error) {
	c.logger.Debug("getting similar recipes", "recipeId", recipeID, "amount", amount, "userId", userID)

	resp, err := c.client.GetSimilarRecipes(ctx, &recipepb.GetSimilarRecipesRequest{
		RecipeId: recipeID,
		Amount:   amount,
		UserId:   userID,
	})
	if err != nil {
		return nil, fmt.Errorf("get similar recipes: %w", err)
	}

	return resp.GetRecipes(), nil
}

// GetByCuisine retrieves recipes by cuisine
func (c *RecipeClient) GetByCuisine(ctx context.Context, userID, cuisineID string) ([]*recipepb.RecipeResponse, error) {
	c.logger.Debug("getting recipes by cuisine", "cuisineId", cuisineID, "userId", userID)

	resp, err := c.client.GetRecipesByCuisine(ctx, &recipepb.GetRecipesByCuisineRequest{
		CuisineId: cuisineID,
		UserId:    userID,
	})
	if err != nil {
		return nil, fmt.Errorf("get recipes by cuisine: %w", err)
	}

	return resp.GetRecipes(), nil
}

// GetByIngredient retrieves recipes by ingredient
func (c *RecipeClient) GetByIngredient(ctx context.Context, userID, ingredientID string) ([]*recipepb.RecipeResponse, error) {
	c.logger.Debug("getting recipes by ingredient", "ingredientId", ingredientID, "userId", userID)

	resp, err := c.client.GetRecipesByIngredient(ctx, &recipepb.GetRecipesByIngredientRequest{
		IngredientId: ingredientID,
		UserId:       userID,
	})
	if err != nil {
		return nil, fmt.Errorf("get recipes by ingredient: %w", err)
	}

	return resp.GetRecipes(), nil
}

// GetByAllergy retrieves recipes excluding an allergy
func (c *RecipeClient) GetByAllergy(ctx context.Context, userID, allergyID string) ([]*recipepb.RecipeResponse, error) {
	c.logger.Debug("getting recipes by allergy", "allergyId", allergyID, "userId", userID)

	resp, err := c.client.GetRecipesByAllergy(ctx, &recipepb.GetRecipesByAllergyRequest{
		AllergyId: allergyID,
		UserId:    userID,
	})
	if err != nil {
		return nil, fmt.Errorf("get recipes by allergy: %w", err)
	}

	return resp.GetRecipes(), nil
}

// GetUnits retrieves available ingredient units.
func (c *RecipeClient) GetUnits(ctx context.Context, userID string) ([]*recipepb.Unit, error) {
	c.logger.Debug("getting units", "userId", userID)

	resp, err := c.client.GetUnits(ctx, &recipepb.GetUnitsRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("get units: %w", err)
	}

	return resp.GetUnits(), nil
}

// CreateUnit creates a new ingredient unit.
func (c *RecipeClient) CreateUnit(ctx context.Context, userID, name string) (*recipepb.Unit, error) {
	c.logger.Debug("creating unit", "name", name, "userId", userID)

	resp, err := c.client.CreateUnit(ctx, &recipepb.CreateUnitRequest{
		UserId: userID,
		Name:   name,
	})
	if err != nil {
		return nil, fmt.Errorf("create unit: %w", err)
	}

	return resp, nil
}
