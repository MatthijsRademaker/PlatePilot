package handler

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/platepilot/backend/internal/mealplanner/domain"
	pb "github.com/platepilot/backend/internal/mealplanner/pb"
	"github.com/platepilot/backend/internal/mealplanner/repository"
)

// GRPCHandler implements the MealPlannerService gRPC interface
type GRPCHandler struct {
	pb.UnimplementedMealPlannerServiceServer
	planner   MealPlanner
	planStore MealPlanStore
	logger    *slog.Logger
}

// NewGRPCHandler creates a new gRPC handler
func NewGRPCHandler(planner MealPlanner, planStore MealPlanStore, logger *slog.Logger) *GRPCHandler {
	return &GRPCHandler{
		planner:   planner,
		planStore: planStore,
		logger:    logger,
	}
}

// SuggestRecipes suggests recipes based on constraints
func (h *GRPCHandler) SuggestRecipes(ctx context.Context, req *pb.SuggestionsRequest) (*pb.SuggestionsResponse, error) {
	h.logger.Debug("suggest recipes request",
		"dailyConstraints", len(req.GetDailyConstraints()),
		"alreadySelected", len(req.GetAlreadySelectedRecipeIds()),
		"amount", req.GetAmount(),
		"userId", req.GetUserId(),
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

// GetWeekPlan retrieves a saved meal plan for the given week.
func (h *GRPCHandler) GetWeekPlan(ctx context.Context, req *pb.GetWeekPlanRequest) (*pb.GetWeekPlanResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	startDate, err := parseDate(req.GetStartDate())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid start date: %v", err)
	}

	plan, err := h.planStore.GetWeekPlan(ctx, userID, startDate)
	if err != nil {
		if errors.Is(err, repository.ErrMealPlanNotFound) {
			emptyPlan := &domain.WeekPlan{
				UserID:    userID,
				StartDate: startDate,
				EndDate:   startDate.AddDate(0, 0, 6),
				Slots:     []domain.MealSlot{},
			}
			return &pb.GetWeekPlanResponse{Plan: toWeekPlanProto(emptyPlan)}, nil
		}
		h.logger.Error("failed to get week plan", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get week plan")
	}

	return &pb.GetWeekPlanResponse{Plan: toWeekPlanProto(plan)}, nil
}

// UpsertWeekPlan creates or updates a week plan.
func (h *GRPCHandler) UpsertWeekPlan(ctx context.Context, req *pb.UpsertWeekPlanRequest) (*pb.UpsertWeekPlanResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	planInput := req.GetPlan()
	if planInput == nil {
		return nil, status.Errorf(codes.InvalidArgument, "plan is required")
	}

	startDate, err := parseDate(planInput.GetStartDate())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid start date: %v", err)
	}

	endDate := startDate.AddDate(0, 0, 6)
	if planInput.GetEndDate() != "" {
		parsedEnd, err := parseDate(planInput.GetEndDate())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid end date: %v", err)
		}
		endDate = parsedEnd
	}

	slots := make([]domain.MealSlot, 0, len(planInput.GetSlots()))
	for _, slot := range planInput.GetSlots() {
		if slot.GetRecipeId() == "" {
			continue
		}
		slotDate, err := parseDate(slot.GetDate())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid slot date: %v", err)
		}
		recipeID, err := uuid.Parse(slot.GetRecipeId())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid recipe ID: %v", err)
		}
		slots = append(slots, domain.MealSlot{
			Date:     slotDate,
			MealType: slot.GetMealType(),
			RecipeID: recipeID,
		})
	}

	plan := domain.WeekPlan{
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
		Slots:     slots,
	}

	updated, err := h.planStore.UpsertWeekPlan(ctx, plan)
	if err != nil {
		h.logger.Error("failed to upsert week plan", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to upsert week plan")
	}

	return &pb.UpsertWeekPlanResponse{Plan: toWeekPlanProto(updated)}, nil
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

	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return domain.SuggestionRequest{}, err
	}

	return domain.SuggestionRequest{
		UserID:                 userID,
		DailyConstraints:       dailyConstraints,
		AlreadySelectedRecipes: alreadySelected,
		Amount:                 amount,
	}, nil
}

func parseDate(value string) (time.Time, error) {
	const layout = "2006-01-02"
	return time.Parse(layout, value)
}

func toWeekPlanProto(plan *domain.WeekPlan) *pb.WeekPlan {
	slots := make([]*pb.MealSlot, 0, len(plan.Slots))
	for _, slot := range plan.Slots {
		recipe := &pb.MealPlanRecipe{
			Id:          slot.RecipeID.String(),
			Name:        slot.RecipeName,
			Description: slot.RecipeDescription,
		}
		slots = append(slots, &pb.MealSlot{
			Date:     slot.Date.Format("2006-01-02"),
			MealType: slot.MealType,
			Recipe:   recipe,
		})
	}

	return &pb.WeekPlan{
		StartDate: plan.StartDate.Format("2006-01-02"),
		EndDate:   plan.EndDate.Format("2006-01-02"),
		Slots:     slots,
	}
}
