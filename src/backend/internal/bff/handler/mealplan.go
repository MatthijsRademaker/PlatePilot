package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

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

// GetWeek handles GET /v1/mealplan/week
// @Summary      Get meal plan week
// @Description  Retrieves a weekly meal plan for the given start date
// @Tags         mealplan
// @Accept       json
// @Produce      json
// @Param        startDate  query     string  true  "Week start date (YYYY-MM-DD)"
// @Success      200        {object}  WeekPlanJSON
// @Failure      400        {object}  ErrorResponse
// @Failure      500        {object}  ErrorResponse
// @Router       /mealplan/week [get]
func (h *MealPlanHandler) GetWeek(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	startDate := r.URL.Query().Get("startDate")
	if startDate == "" {
		writeError(w, http.StatusBadRequest, "startDate is required")
		return
	}
	if _, err := time.Parse("2006-01-02", startDate); err != nil {
		writeError(w, http.StatusBadRequest, "startDate must be YYYY-MM-DD")
		return
	}

	plan, err := h.client.GetWeekPlan(r.Context(), userID.String(), startDate)
	if err != nil {
		h.logger.Error("failed to get week plan", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to fetch meal plan")
		return
	}

	writeJSON(w, http.StatusOK, toWeekPlanJSON(plan))
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

// UpsertWeek handles PUT /v1/mealplan/week
// @Summary      Save meal plan week
// @Description  Creates or updates a weekly meal plan
// @Tags         mealplan
// @Accept       json
// @Produce      json
// @Param        request  body      UpsertWeekPlanRequest  true  "Week plan payload"
// @Success      200      {object}  WeekPlanJSON
// @Failure      400      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /mealplan/week [put]
func (h *MealPlanHandler) UpsertWeek(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req UpsertWeekPlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.StartDate == "" {
		writeError(w, http.StatusBadRequest, "startDate is required")
		return
	}
	if _, err := time.Parse("2006-01-02", req.StartDate); err != nil {
		writeError(w, http.StatusBadRequest, "startDate must be YYYY-MM-DD")
		return
	}
	if req.EndDate != "" {
		if _, err := time.Parse("2006-01-02", req.EndDate); err != nil {
			writeError(w, http.StatusBadRequest, "endDate must be YYYY-MM-DD")
			return
		}
	}

	planInput := &mealplannerpb.WeekPlanInput{
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Slots:     req.ToSlots(),
	}

	plan, err := h.client.UpsertWeekPlan(r.Context(), &mealplannerpb.UpsertWeekPlanRequest{
		UserId: userID.String(),
		Plan:   planInput,
	})
	if err != nil {
		h.logger.Error("failed to upsert week plan", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to save meal plan")
		return
	}

	writeJSON(w, http.StatusOK, toWeekPlanJSON(plan))
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

// UpsertWeekPlanRequest is the request body for saving a week plan.
type UpsertWeekPlanRequest struct {
	StartDate string         `json:"startDate"`
	EndDate   string         `json:"endDate"`
	Days      []DayPlanInput `json:"days"`
}

// DayPlanInput represents a day in the week plan payload.
type DayPlanInput struct {
	Date  string          `json:"date"`
	Meals []MealSlotInput `json:"meals"`
}

// MealSlotInput represents a meal slot in the week plan payload.
type MealSlotInput struct {
	MealType string `json:"mealType"`
	RecipeID string `json:"recipeId,omitempty"`
}

func (r *UpsertWeekPlanRequest) ToSlots() []*mealplannerpb.MealSlotInput {
	slots := make([]*mealplannerpb.MealSlotInput, 0)
	for _, day := range r.Days {
		for _, meal := range day.Meals {
			if meal.RecipeID == "" {
				continue
			}
			slots = append(slots, &mealplannerpb.MealSlotInput{
				Date:     day.Date,
				MealType: meal.MealType,
				RecipeId: meal.RecipeID,
			})
		}
	}
	return slots
}

// WeekPlanJSON is the JSON response for a week plan.
type WeekPlanJSON struct {
	StartDate string        `json:"startDate"`
	EndDate   string        `json:"endDate"`
	Days      []DayPlanJSON `json:"days"`
}

// DayPlanJSON is the JSON response for a day plan.
type DayPlanJSON struct {
	Date  string         `json:"date"`
	Meals []MealSlotJSON `json:"meals"`
}

// MealSlotJSON is the JSON response for a meal slot.
type MealSlotJSON struct {
	ID       string             `json:"id"`
	Date     string             `json:"date"`
	MealType string             `json:"mealType"`
	Recipe   *RecipeSummaryJSON `json:"recipe,omitempty"`
}

// RecipeSummaryJSON is minimal recipe info for meal plans.
type RecipeSummaryJSON struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

func toWeekPlanJSON(plan *mealplannerpb.WeekPlan) WeekPlanJSON {
	startDate := plan.GetStartDate()
	endDate := plan.GetEndDate()

	dateKeys := make([]string, 0, 7)
	start, err := time.Parse("2006-01-02", startDate)
	if err == nil {
		end, endErr := time.Parse("2006-01-02", endDate)
		if endErr == nil {
			for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
				dateKeys = append(dateKeys, d.Format("2006-01-02"))
			}
		}
	}

	if len(dateKeys) == 0 {
		dateKeys = []string{startDate}
	}

	slotMap := make(map[string]*RecipeSummaryJSON)
	for _, slot := range plan.GetSlots() {
		if slot.GetRecipe() == nil {
			continue
		}
		key := slot.GetDate() + "|" + slot.GetMealType()
		slotMap[key] = &RecipeSummaryJSON{
			ID:          slot.GetRecipe().GetId(),
			Name:        slot.GetRecipe().GetName(),
			Description: slot.GetRecipe().GetDescription(),
		}
	}

	mealTypes := []string{"breakfast", "lunch", "dinner"}
	days := make([]DayPlanJSON, 0, len(dateKeys))
	for _, date := range dateKeys {
		meals := make([]MealSlotJSON, 0, len(mealTypes))
		for _, mealType := range mealTypes {
			key := date + "|" + mealType
			meals = append(meals, MealSlotJSON{
				ID:       date + "-" + mealType,
				Date:     date,
				MealType: mealType,
				Recipe:   slotMap[key],
			})
		}
		days = append(days, DayPlanJSON{
			Date:  date,
			Meals: meals,
		})
	}

	return WeekPlanJSON{
		StartDate: startDate,
		EndDate:   endDate,
		Days:      days,
	}
}
