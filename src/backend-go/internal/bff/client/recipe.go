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
func (c *RecipeClient) GetByID(ctx context.Context, id string) (*recipepb.RecipeResponse, error) {
	c.logger.Debug("getting recipe by id", "id", id)

	resp, err := c.client.GetRecipeById(ctx, &recipepb.GetRecipeByIdRequest{
		RecipeId: id,
	})
	if err != nil {
		return nil, fmt.Errorf("get recipe by id: %w", err)
	}

	return resp, nil
}

// GetAll retrieves all recipes with pagination
func (c *RecipeClient) GetAll(ctx context.Context, pageIndex, pageSize int32) ([]*recipepb.RecipeResponse, error) {
	c.logger.Debug("getting all recipes", "pageIndex", pageIndex, "pageSize", pageSize)

	resp, err := c.client.GetAllRecipes(ctx, &recipepb.GetAllRecipesRequest{
		PageIndex: pageIndex,
		PageSize:  pageSize,
	})
	if err != nil {
		return nil, fmt.Errorf("get all recipes: %w", err)
	}

	return resp.GetRecipes(), nil
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
func (c *RecipeClient) GetSimilar(ctx context.Context, recipeID string, amount int32) ([]*recipepb.RecipeResponse, error) {
	c.logger.Debug("getting similar recipes", "recipeId", recipeID, "amount", amount)

	resp, err := c.client.GetSimilarRecipes(ctx, &recipepb.GetSimilarRecipesRequest{
		RecipeId: recipeID,
		Amount:   amount,
	})
	if err != nil {
		return nil, fmt.Errorf("get similar recipes: %w", err)
	}

	return resp.GetRecipes(), nil
}

// GetByCuisine retrieves recipes by cuisine
func (c *RecipeClient) GetByCuisine(ctx context.Context, cuisineID string) ([]*recipepb.RecipeResponse, error) {
	c.logger.Debug("getting recipes by cuisine", "cuisineId", cuisineID)

	resp, err := c.client.GetRecipesByCuisine(ctx, &recipepb.GetRecipesByCuisineRequest{
		CuisineId: cuisineID,
	})
	if err != nil {
		return nil, fmt.Errorf("get recipes by cuisine: %w", err)
	}

	return resp.GetRecipes(), nil
}

// GetByIngredient retrieves recipes by ingredient
func (c *RecipeClient) GetByIngredient(ctx context.Context, ingredientID string) ([]*recipepb.RecipeResponse, error) {
	c.logger.Debug("getting recipes by ingredient", "ingredientId", ingredientID)

	resp, err := c.client.GetRecipesByIngredient(ctx, &recipepb.GetRecipesByIngredientRequest{
		IngredientId: ingredientID,
	})
	if err != nil {
		return nil, fmt.Errorf("get recipes by ingredient: %w", err)
	}

	return resp.GetRecipes(), nil
}

// GetByAllergy retrieves recipes excluding an allergy
func (c *RecipeClient) GetByAllergy(ctx context.Context, allergyID string) ([]*recipepb.RecipeResponse, error) {
	c.logger.Debug("getting recipes by allergy", "allergyId", allergyID)

	resp, err := c.client.GetRecipesByAllergy(ctx, &recipepb.GetRecipesByAllergyRequest{
		AllergyId: allergyID,
	})
	if err != nil {
		return nil, fmt.Errorf("get recipes by allergy: %w", err)
	}

	return resp.GetRecipes(), nil
}
