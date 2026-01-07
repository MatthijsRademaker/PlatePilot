package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/platepilot/backend/internal/bff/client"
	mealplannerpb "github.com/platepilot/backend/internal/mealplanner/pb"
)

// MealPlanHandler handles REST requests for meal planning
type MealPlanHandler struct {
	client *client.MealPlannerClient
	logger *slog.Logger
}

// NewMealPlanHandler creates a new meal plan handler
func NewMealPlanHandler(client *client.MealPlannerClient, logger *slog.Logger) *MealPlanHandler {
	return &MealPlanHandler{
		client: client,
		logger: logger,
	}
}

// Suggest handles POST /v1/mealplan/suggest
// @Summary      Suggest recipes for meal planning
// @Description  Suggests recipes based on constraints and already selected recipes
// @Tags         mealplan
// @Accept       json
// @Produce      json
// @Param        request  body      SuggestRequest  true  "Suggestion request with constraints"
// @Success      200      {object}  SuggestResponse
// @Failure      400      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /mealplan/suggest [post]
func (h *MealPlanHandler) Suggest(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req SuggestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Amount <= 0 {
		req.Amount = 5
	}
	if req.Amount > 20 {
		req.Amount = 20
	}

	recipeIDs, err := h.client.SuggestRecipes(r.Context(), req.ToProto(userID.String()))
	if err != nil {
		h.logger.Error("failed to suggest recipes", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to suggest recipes")
		return
	}

	writeJSON(w, http.StatusOK, SuggestResponse{RecipeIDs: recipeIDs})
}

// SuggestRequest is the request body for suggesting recipes
type SuggestRequest struct {
	DailyConstraints         []DailyConstraint `json:"dailyConstraints"`
	AlreadySelectedRecipeIDs []string          `json:"alreadySelectedRecipeIds"`
	Amount                   int32             `json:"amount"`
}

// DailyConstraint represents constraints for a single day
type DailyConstraint struct {
	IngredientConstraints []EntityConstraint `json:"ingredientConstraints"`
	CuisineConstraints    []EntityConstraint `json:"cuisineConstraints"`
}

// EntityConstraint represents a constraint on an entity
type EntityConstraint struct {
	EntityID string `json:"entityId"`
}

// ToProto converts the request to a protobuf message
func (r *SuggestRequest) ToProto(userID string) *mealplannerpb.SuggestionsRequest {
	dailyConstraints := make([]*mealplannerpb.DailyConstraints, len(r.DailyConstraints))
	for i, dc := range r.DailyConstraints {
		ingredientConstraints := make([]*mealplannerpb.IngredientConstraint, len(dc.IngredientConstraints))
		for j, ic := range dc.IngredientConstraints {
			ingredientConstraints[j] = &mealplannerpb.IngredientConstraint{
				EntityId: ic.EntityID,
			}
		}

		cuisineConstraints := make([]*mealplannerpb.CuisineConstraint, len(dc.CuisineConstraints))
		for j, cc := range dc.CuisineConstraints {
			cuisineConstraints[j] = &mealplannerpb.CuisineConstraint{
				EntityId: cc.EntityID,
			}
		}

		dailyConstraints[i] = &mealplannerpb.DailyConstraints{
			IngredientConstraints: ingredientConstraints,
			CuisineConstraints:    cuisineConstraints,
		}
	}

	return &mealplannerpb.SuggestionsRequest{
		UserId:                   userID,
		DailyConstraints:         dailyConstraints,
		AlreadySelectedRecipeIds: r.AlreadySelectedRecipeIDs,
		Amount:                   r.Amount,
	}
}

// SuggestResponse is the response for suggesting recipes
type SuggestResponse struct {
	RecipeIDs []string `json:"recipeIds"`
}
