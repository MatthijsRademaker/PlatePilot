package client

import (
	"context"
	"fmt"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	mealplannerpb "github.com/platepilot/backend/internal/mealplanner/pb"
)

// MealPlannerClient wraps the gRPC client for the MealPlanner API
type MealPlannerClient struct {
	conn   *grpc.ClientConn
	client mealplannerpb.MealPlannerServiceClient
	logger *slog.Logger
}

// NewMealPlannerClient creates a new MealPlanner API client
func NewMealPlannerClient(address string, logger *slog.Logger) (*MealPlannerClient, error) {
	conn, err := grpc.NewClient(address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("connect to mealplanner api: %w", err)
	}

	return &MealPlannerClient{
		conn:   conn,
		client: mealplannerpb.NewMealPlannerServiceClient(conn),
		logger: logger,
	}, nil
}

// Close closes the gRPC connection
func (c *MealPlannerClient) Close() error {
	return c.conn.Close()
}

// SuggestRecipes suggests recipes based on constraints
func (c *MealPlannerClient) SuggestRecipes(ctx context.Context, req *mealplannerpb.SuggestionsRequest) ([]string, error) {
	c.logger.Debug("suggesting recipes",
		"dailyConstraints", len(req.GetDailyConstraints()),
		"alreadySelected", len(req.GetAlreadySelectedRecipeIds()),
		"amount", req.GetAmount(),
		"userId", req.GetUserId(),
	)

	resp, err := c.client.SuggestRecipes(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("suggest recipes: %w", err)
	}

	return resp.GetRecipeIds(), nil
}

// GetWeekPlan retrieves a week plan for a user.
func (c *MealPlannerClient) GetWeekPlan(ctx context.Context, userID, startDate string) (*mealplannerpb.WeekPlan, error) {
	c.logger.Debug("getting week plan", "startDate", startDate, "userId", userID)

	resp, err := c.client.GetWeekPlan(ctx, &mealplannerpb.GetWeekPlanRequest{
		UserId:    userID,
		StartDate: startDate,
	})
	if err != nil {
		return nil, fmt.Errorf("get week plan: %w", err)
	}

	return resp.GetPlan(), nil
}

// UpsertWeekPlan creates or updates a week plan.
func (c *MealPlannerClient) UpsertWeekPlan(ctx context.Context, req *mealplannerpb.UpsertWeekPlanRequest) (*mealplannerpb.WeekPlan, error) {
	c.logger.Debug("upserting week plan", "userId", req.GetUserId())

	resp, err := c.client.UpsertWeekPlan(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("upsert week plan: %w", err)
	}

	return resp.GetPlan(), nil
}
