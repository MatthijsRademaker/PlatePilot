package handler

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/platepilot/backend/internal/mealplanner/domain"
	pb "github.com/platepilot/backend/internal/mealplanner/pb"
)

// GRPCHandler implements the MealPlannerService gRPC interface
type GRPCHandler struct {
	pb.UnimplementedMealPlannerServiceServer
	planner *domain.Planner
	logger  *slog.Logger
}

// NewGRPCHandler creates a new gRPC handler
func NewGRPCHandler(planner *domain.Planner, logger *slog.Logger) *GRPCHandler {
	return &GRPCHandler{
		planner: planner,
		logger:  logger,
	}
}

// SuggestRecipes suggests recipes based on constraints
func (h *GRPCHandler) SuggestRecipes(ctx context.Context, req *pb.SuggestionsRequest) (*pb.SuggestionsResponse, error) {
	h.logger.Debug("suggest recipes request",
		"dailyConstraints", len(req.GetDailyConstraints()),
		"alreadySelected", len(req.GetAlreadySelectedRecipeIds()),
		"amount", req.GetAmount(),
	)

	// Convert protobuf request to domain request
	domainReq, err := h.toDomainRequest(req)
	if err != nil {
		h.logger.Error("failed to convert request", "error", err)
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	// Get suggestions from the planner
	recipeIDs, err := h.planner.SuggestMeals(ctx, domainReq)
	if err != nil {
		h.logger.Error("failed to suggest meals", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to suggest meals: %v", err)
	}

	// Convert UUIDs to strings
	recipeIDStrings := make([]string, len(recipeIDs))
	for i, id := range recipeIDs {
		recipeIDStrings[i] = id.String()
	}

	h.logger.Info("suggested recipes",
		"count", len(recipeIDStrings),
		"amount", req.GetAmount(),
	)

	return &pb.SuggestionsResponse{
		RecipeIds: recipeIDStrings,
	}, nil
}

func (h *GRPCHandler) toDomainRequest(req *pb.SuggestionsRequest) (domain.SuggestionRequest, error) {
	// Parse already selected recipe IDs
	alreadySelected := make([]uuid.UUID, 0, len(req.GetAlreadySelectedRecipeIds()))
	for _, idStr := range req.GetAlreadySelectedRecipeIds() {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return domain.SuggestionRequest{}, err
		}
		alreadySelected = append(alreadySelected, id)
	}

	// Parse daily constraints
	dailyConstraints := make([]domain.DailyConstraints, 0, len(req.GetDailyConstraints()))
	for _, dc := range req.GetDailyConstraints() {
		ingredientConstraints := make([]uuid.UUID, 0, len(dc.GetIngredientConstraints()))
		for _, ic := range dc.GetIngredientConstraints() {
			id, err := uuid.Parse(ic.GetEntityId())
			if err != nil {
				return domain.SuggestionRequest{}, err
			}
			ingredientConstraints = append(ingredientConstraints, id)
		}

		cuisineConstraints := make([]uuid.UUID, 0, len(dc.GetCuisineConstraints()))
		for _, cc := range dc.GetCuisineConstraints() {
			id, err := uuid.Parse(cc.GetEntityId())
			if err != nil {
				return domain.SuggestionRequest{}, err
			}
			cuisineConstraints = append(cuisineConstraints, id)
		}

		dailyConstraints = append(dailyConstraints, domain.DailyConstraints{
			IngredientConstraints: ingredientConstraints,
			CuisineConstraints:    cuisineConstraints,
		})
	}

	amount := int(req.GetAmount())
	if amount <= 0 {
		amount = 5
	}
	if amount > 50 {
		amount = 50
	}

	return domain.SuggestionRequest{
		DailyConstraints:       dailyConstraints,
		AlreadySelectedRecipes: alreadySelected,
		Amount:                 amount,
	}, nil
}
