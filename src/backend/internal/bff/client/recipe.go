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
func (c *RecipeClient) GetByID(ctx context.Context, userID, id string) (*recipepb.Recipe, error) {
	c.logger.Debug("getting recipe by id", "id", id, "userId", userID)

	resp, err := c.client.GetRecipe(ctx, &recipepb.GetRecipeRequest{
		RecipeId: id,
		UserId:   userID,
	})
	if err != nil {
		return nil, fmt.Errorf("get recipe by id: %w", err)
	}

	return resp, nil
}

// ListResponse contains the paginated recipes response.
type ListResponse struct {
	Recipes    []*recipepb.Recipe
	PageIndex  int32
	PageSize   int32
	TotalCount int32
	TotalPages int32
}

// ListRecipes retrieves recipes with pagination and optional filters.
func (c *RecipeClient) ListRecipes(ctx context.Context, req *recipepb.ListRecipesRequest) (*ListResponse, error) {
	c.logger.Debug("listing recipes", "pageIndex", req.GetPageIndex(), "pageSize", req.GetPageSize(), "userId", req.GetUserId())

	resp, err := c.client.ListRecipes(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("list recipes: %w", err)
	}

	return &ListResponse{
		Recipes:    resp.GetRecipes(),
		PageIndex:  resp.GetPageIndex(),
		PageSize:   resp.GetPageSize(),
		TotalCount: resp.GetTotalCount(),
		TotalPages: resp.GetTotalPages(),
	}, nil
}

// Create creates a new recipe
func (c *RecipeClient) Create(ctx context.Context, req *recipepb.CreateRecipeRequest) (*recipepb.Recipe, error) {
	c.logger.Debug("creating recipe", "userId", req.GetUserId())

	resp, err := c.client.CreateRecipe(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("create recipe: %w", err)
	}

	return resp, nil
}

// Update updates an existing recipe.
func (c *RecipeClient) Update(ctx context.Context, req *recipepb.UpdateRecipeRequest) (*recipepb.Recipe, error) {
	c.logger.Debug("updating recipe", "recipeId", req.GetRecipeId(), "userId", req.GetUserId())

	resp, err := c.client.UpdateRecipe(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("update recipe: %w", err)
	}

	return resp, nil
}

// Delete deletes an existing recipe.
func (c *RecipeClient) Delete(ctx context.Context, req *recipepb.DeleteRecipeRequest) error {
	c.logger.Debug("deleting recipe", "recipeId", req.GetRecipeId(), "userId", req.GetUserId())

	_, err := c.client.DeleteRecipe(ctx, req)
	if err != nil {
		return fmt.Errorf("delete recipe: %w", err)
	}

	return nil
}

// GetSimilar retrieves recipes similar to a given recipe
func (c *RecipeClient) GetSimilar(ctx context.Context, userID, recipeID string, amount int32) ([]*recipepb.Recipe, error) {
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

// GetCuisines retrieves available cuisines.
func (c *RecipeClient) GetCuisines(ctx context.Context, userID string) ([]*recipepb.Cuisine, error) {
	c.logger.Debug("getting cuisines", "userId", userID)

	resp, err := c.client.GetCuisines(ctx, &recipepb.GetCuisinesRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("get cuisines: %w", err)
	}

	return resp.GetCuisines(), nil
}

// CreateCuisine creates a new cuisine.
func (c *RecipeClient) CreateCuisine(ctx context.Context, userID, name string) (*recipepb.Cuisine, error) {
	c.logger.Debug("creating cuisine", "name", name, "userId", userID)

	resp, err := c.client.CreateCuisine(ctx, &recipepb.CreateCuisineRequest{
		UserId: userID,
		Name:   name,
	})
	if err != nil {
		return nil, fmt.Errorf("create cuisine: %w", err)
	}

	return resp, nil
}
