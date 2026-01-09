package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/platepilot/backend/internal/bff/client"
	recipepb "github.com/platepilot/backend/internal/recipe/pb"
)

// RecipeHandler handles REST requests for recipes.
type RecipeHandler struct {
	client *client.RecipeClient
	logger *slog.Logger
}

// NewRecipeHandler creates a new recipe handler.
func NewRecipeHandler(client *client.RecipeClient, logger *slog.Logger) *RecipeHandler {
	return &RecipeHandler{
		client: client,
		logger: logger,
	}
}

// GetByID handles GET /v1/recipe/{id}
// @Summary      Get recipe by ID
// @Description  Retrieves a single recipe by its unique identifier
// @Tags         recipes
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Recipe ID (UUID)"
// @Success      200  {object}  RecipeJSON
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /recipe/{id} [get]
func (h *RecipeHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "recipe id is required")
		return
	}

	recipe, err := h.client.GetByID(r.Context(), userID.String(), id)
	if err != nil {
		h.logger.Error("failed to get recipe", "id", id, "error", err)
		writeError(w, http.StatusNotFound, "recipe not found")
		return
	}

	writeJSON(w, http.StatusOK, toRecipeJSON(recipe))
}

// PaginatedRecipesJSON is the paginated response for recipes.
type PaginatedRecipesJSON struct {
	Items      []RecipeJSON `json:"items"`
	PageIndex  int32        `json:"pageIndex"`
	PageSize   int32        `json:"pageSize"`
	TotalCount int32        `json:"totalCount"`
	TotalPages int32        `json:"totalPages"`
}

// List handles GET /v1/recipe
// @Summary      List recipes (paginated)
// @Description  Retrieves a paginated list of recipes with optional filters
// @Tags         recipes
// @Accept       json
// @Produce      json
// @Param        pageIndex   query     int     false  "Page number (1-indexed)"  default(1)
// @Param        pageSize    query     int     false  "Items per page (max 100)" default(20)
// @Param        cuisineId   query     string  false  "Cuisine ID filter"
// @Param        ingredientId query    string  false  "Ingredient ID filter"
// @Param        allergyId   query     string  false  "Allergy ID filter (exclude)"
// @Param        tags        query     string  false  "Comma-separated tags filter"
// @Success      200  {object}  PaginatedRecipesJSON
// @Failure      500  {object}  ErrorResponse
// @Router       /recipe [get]
func (h *RecipeHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	pageIndex := parseIntParam(r, "pageIndex", 1)
	pageSize := parseIntParam(r, "pageSize", 20)
	if pageSize > 100 {
		pageSize = 100
	}

	req := &recipepb.ListRecipesRequest{
		UserId:    userID.String(),
		PageIndex: int32(pageIndex),
		PageSize:  int32(pageSize),
		CuisineId: strings.TrimSpace(r.URL.Query().Get("cuisineId")),
		IngredientId: strings.TrimSpace(r.URL.Query().Get("ingredientId")),
		AllergyId: strings.TrimSpace(r.URL.Query().Get("allergyId")),
		Tags:      splitCommaList(r.URL.Query().Get("tags")),
	}

	resp, err := h.client.ListRecipes(r.Context(), req)
	if err != nil {
		h.logger.Error("failed to list recipes", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to fetch recipes")
		return
	}

	writeJSON(w, http.StatusOK, PaginatedRecipesJSON{
		Items:      toRecipesJSON(resp.Recipes),
		PageIndex:  resp.PageIndex,
		PageSize:   resp.PageSize,
		TotalCount: resp.TotalCount,
		TotalPages: resp.TotalPages,
	})
}

// GetSimilar handles GET /v1/recipe/similar
// @Summary      Get similar recipes
// @Description  Finds recipes similar to the specified recipe using vector search
// @Tags         recipes
// @Accept       json
// @Produce      json
// @Param        recipe  query     string  true   "Recipe ID to find similar recipes for"
// @Param        amount  query     int     false  "Number of similar recipes to return (max 50)" default(5)
// @Success      200  {array}   RecipeJSON
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /recipe/similar [get]
func (h *RecipeHandler) GetSimilar(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	recipeID := r.URL.Query().Get("recipe")
	if recipeID == "" {
		writeError(w, http.StatusBadRequest, "recipe parameter is required")
		return
	}

	amount := parseIntParam(r, "amount", 5)
	if amount > 50 {
		amount = 50
	}

	recipes, err := h.client.GetSimilar(r.Context(), userID.String(), recipeID, int32(amount))
	if err != nil {
		h.logger.Error("failed to get similar recipes", "recipeId", recipeID, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to fetch similar recipes")
		return
	}

	writeJSON(w, http.StatusOK, toRecipesJSON(recipes))
}

// Create handles POST /v1/recipe
// @Summary      Create a new recipe
// @Description  Creates a new recipe with the provided details
// @Tags         recipes
// @Accept       json
// @Produce      json
// @Param        recipe  body      RecipeInputJSON  true  "Recipe to create"
// @Success      201  {object}  RecipeJSON
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /recipe [post]
func (h *RecipeHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req RecipeInputJSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := req.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	recipe, err := h.client.Create(r.Context(), req.ToProto(userID.String()))
	if err != nil {
		h.logger.Error("failed to create recipe", "name", req.Name, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to create recipe")
		return
	}

	writeJSON(w, http.StatusCreated, toRecipeJSON(recipe))
}

// Update handles PUT /v1/recipe/{id}
// @Summary      Update a recipe
// @Description  Updates an existing recipe with the provided details
// @Tags         recipes
// @Accept       json
// @Produce      json
// @Param        id      path      string            true  "Recipe ID (UUID)"
// @Param        recipe  body      RecipeInputJSON   true  "Recipe to update"
// @Success      200  {object}  RecipeJSON
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /recipe/{id} [put]
func (h *RecipeHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "recipe id is required")
		return
	}

	var req RecipeInputJSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := req.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	recipe, err := h.client.Update(r.Context(), req.ToUpdateProto(userID.String(), id))
	if err != nil {
		h.logger.Error("failed to update recipe", "id", id, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to update recipe")
		return
	}

	writeJSON(w, http.StatusOK, toRecipeJSON(recipe))
}

// Delete handles DELETE /v1/recipe/{id}
// @Summary      Delete a recipe
// @Description  Deletes a recipe by ID
// @Tags         recipes
// @Param        id   path      string  true  "Recipe ID (UUID)"
// @Success      204
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /recipe/{id} [delete]
func (h *RecipeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "recipe id is required")
		return
	}

	if err := h.client.Delete(r.Context(), &recipepb.DeleteRecipeRequest{
		RecipeId: id,
		UserId:   userID.String(),
	}); err != nil {
		h.logger.Error("failed to delete recipe", "id", id, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to delete recipe")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetCuisines handles GET /v1/recipe/cuisines
// @Summary      Get cuisines
// @Description  Retrieves all available cuisines
// @Tags         recipes
// @Accept       json
// @Produce      json
// @Success      200  {object}  CuisinesResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /recipe/cuisines [get]
func (h *RecipeHandler) GetCuisines(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	cuisines, err := h.client.GetCuisines(r.Context(), userID.String())
	if err != nil {
		h.logger.Error("failed to get cuisines", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to fetch cuisines")
		return
	}

	writeJSON(w, http.StatusOK, CuisinesResponse{Items: toCuisinesJSON(cuisines)})
}

// CreateCuisine handles POST /v1/recipe/cuisines
// @Summary      Create cuisine
// @Description  Creates a new cuisine
// @Tags         recipes
// @Accept       json
// @Produce      json
// @Param        cuisine  body      CreateCuisineRequest  true  "Cuisine to create"
// @Success      201  {object}  CuisineJSON
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /recipe/cuisines [post]
func (h *RecipeHandler) CreateCuisine(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req CreateCuisineRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	cuisine, err := h.client.CreateCuisine(r.Context(), userID.String(), name)
	if err != nil {
		h.logger.Error("failed to create cuisine", "name", name, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to create cuisine")
		return
	}

	writeJSON(w, http.StatusCreated, CuisineJSON{ID: cuisine.GetId(), Name: cuisine.GetName()})
}

// CreateCuisineRequest is the request body for creating a cuisine.
type CreateCuisineRequest struct {
	Name string `json:"name"`
}

// RecipeInputJSON is the JSON request body for creating/updating a recipe.
type RecipeInputJSON struct {
	Name               string                   `json:"name"`
	Description        string                   `json:"description,omitempty"`
	PrepTimeMinutes    int                      `json:"prepTimeMinutes,omitempty"`
	CookTimeMinutes    int                      `json:"cookTimeMinutes,omitempty"`
	Servings           int                      `json:"servings,omitempty"`
	YieldQuantity      *float64                 `json:"yieldQuantity,omitempty"`
	YieldUnit          string                   `json:"yieldUnit,omitempty"`
	MainIngredientID   string                   `json:"mainIngredientId,omitempty"`
	MainIngredientName string                   `json:"mainIngredientName,omitempty"`
	CuisineID          string                   `json:"cuisineId,omitempty"`
	CuisineName        string                   `json:"cuisineName,omitempty"`
	IngredientLines    []IngredientLineInputJSON `json:"ingredientLines,omitempty"`
	Steps              []RecipeStepJSON          `json:"steps,omitempty"`
	Tags               []string                 `json:"tags,omitempty"`
	ImageURL           string                   `json:"imageUrl,omitempty"`
	Nutrition          *RecipeNutritionJSON     `json:"nutrition,omitempty"`
}

// Validate validates the recipe input.
func (r *RecipeInputJSON) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return errString("name is required")
	}
	return nil
}

// ToProto converts the request to proto.
func (r *RecipeInputJSON) ToProto(userID string) *recipepb.CreateRecipeRequest {
	return &recipepb.CreateRecipeRequest{
		UserId: userID,
		Recipe: r.toRecipeInputProto(),
	}
}

// ToUpdateProto converts the request to a proto update request.
func (r *RecipeInputJSON) ToUpdateProto(userID, recipeID string) *recipepb.UpdateRecipeRequest {
	return &recipepb.UpdateRecipeRequest{
		UserId:   userID,
		RecipeId: recipeID,
		Recipe:   r.toRecipeInputProto(),
	}
}

func (r *RecipeInputJSON) toRecipeInputProto() *recipepb.RecipeInput {
	lines := make([]*recipepb.IngredientLineInput, 0, len(r.IngredientLines))
	for _, line := range r.IngredientLines {
		lines = append(lines, line.ToProto())
	}

	steps := make([]*recipepb.RecipeStepInput, 0, len(r.Steps))
	for _, step := range r.Steps {
		steps = append(steps, step.ToProto())
	}

	var yieldQuantity *wrapperspb.DoubleValue
	if r.YieldQuantity != nil {
		yieldQuantity = wrapperspb.Double(*r.YieldQuantity)
	}

	return &recipepb.RecipeInput{
		Name:               strings.TrimSpace(r.Name),
		Description:        strings.TrimSpace(r.Description),
		PrepTimeMinutes:    int32(r.PrepTimeMinutes),
		CookTimeMinutes:    int32(r.CookTimeMinutes),
		Servings:           int32(r.Servings),
		YieldQuantity:      yieldQuantity,
		YieldUnit:          strings.TrimSpace(r.YieldUnit),
		MainIngredientId:   strings.TrimSpace(r.MainIngredientID),
		MainIngredientName: strings.TrimSpace(r.MainIngredientName),
		CuisineId:          strings.TrimSpace(r.CuisineID),
		CuisineName:        strings.TrimSpace(r.CuisineName),
		IngredientLines:    lines,
		Steps:              steps,
		Tags:               sanitizeStrings(r.Tags),
		ImageUrl:           strings.TrimSpace(r.ImageURL),
		Nutrition:          nutritionToProto(r.Nutrition),
	}
}

// IngredientLineInputJSON represents a recipe ingredient line input.
type IngredientLineInputJSON struct {
	IngredientID   string   `json:"ingredientId,omitempty"`
	IngredientName string   `json:"ingredientName,omitempty"`
	QuantityValue  *float64 `json:"quantityValue,omitempty"`
	QuantityText   string   `json:"quantityText,omitempty"`
	Unit           string   `json:"unit,omitempty"`
	IsOptional     bool     `json:"isOptional,omitempty"`
	Note           string   `json:"note,omitempty"`
	SortOrder      int      `json:"sortOrder,omitempty"`
}

// ToProto converts the ingredient line input to proto.
func (l *IngredientLineInputJSON) ToProto() *recipepb.IngredientLineInput {
	var quantityValue *wrapperspb.DoubleValue
	if l.QuantityValue != nil {
		quantityValue = wrapperspb.Double(*l.QuantityValue)
	}

	return &recipepb.IngredientLineInput{
		IngredientId:   strings.TrimSpace(l.IngredientID),
		IngredientName: strings.TrimSpace(l.IngredientName),
		QuantityValue:  quantityValue,
		QuantityText:   strings.TrimSpace(l.QuantityText),
		Unit:           strings.TrimSpace(l.Unit),
		IsOptional:     l.IsOptional,
		Note:           strings.TrimSpace(l.Note),
		SortOrder:      int32(l.SortOrder),
	}
}

// RecipeStepJSON represents a recipe step.
type RecipeStepJSON struct {
	StepIndex        int      `json:"stepIndex,omitempty"`
	Instruction      string   `json:"instruction"`
	DurationSeconds  *int     `json:"durationSeconds,omitempty"`
	TemperatureValue *float64 `json:"temperatureValue,omitempty"`
	TemperatureUnit  string   `json:"temperatureUnit,omitempty"`
	MediaURL         string   `json:"mediaUrl,omitempty"`
}

// ToProto converts the step to proto input.
func (s *RecipeStepJSON) ToProto() *recipepb.RecipeStepInput {
	var duration *wrapperspb.Int32Value
	if s.DurationSeconds != nil {
		duration = wrapperspb.Int32(int32(*s.DurationSeconds))
	}

	var temperature *wrapperspb.DoubleValue
	if s.TemperatureValue != nil {
		temperature = wrapperspb.Double(*s.TemperatureValue)
	}

	return &recipepb.RecipeStepInput{
		StepIndex:       int32(s.StepIndex),
		Instruction:     strings.TrimSpace(s.Instruction),
		DurationSeconds: duration,
		TemperatureValue: temperature,
		TemperatureUnit: strings.TrimSpace(s.TemperatureUnit),
		MediaUrl:        strings.TrimSpace(s.MediaURL),
	}
}

// RecipeNutritionJSON represents recipe nutrition info.
type RecipeNutritionJSON struct {
	CaloriesTotal      int     `json:"caloriesTotal,omitempty"`
	CaloriesPerServing int     `json:"caloriesPerServing,omitempty"`
	ProteinG           float64 `json:"proteinG,omitempty"`
	CarbsG             float64 `json:"carbsG,omitempty"`
	FatG               float64 `json:"fatG,omitempty"`
	FiberG             float64 `json:"fiberG,omitempty"`
	SugarG             float64 `json:"sugarG,omitempty"`
	SodiumMg           float64 `json:"sodiumMg,omitempty"`
}

func nutritionToProto(n *RecipeNutritionJSON) *recipepb.RecipeNutrition {
	if n == nil {
		return nil
	}
	return &recipepb.RecipeNutrition{
		CaloriesTotal:      int32(n.CaloriesTotal),
		CaloriesPerServing: int32(n.CaloriesPerServing),
		ProteinG:           n.ProteinG,
		CarbsG:             n.CarbsG,
		FatG:               n.FatG,
		FiberG:             n.FiberG,
		SugarG:             n.SugarG,
		SodiumMg:           n.SodiumMg,
	}
}

// RecipeJSON is the JSON response for a recipe.
type RecipeJSON struct {
	ID               string               `json:"id"`
	Name             string               `json:"name"`
	Description      string               `json:"description,omitempty"`
	PrepTimeMinutes  int32                `json:"prepTimeMinutes"`
	CookTimeMinutes  int32                `json:"cookTimeMinutes"`
	TotalTimeMinutes int32                `json:"totalTimeMinutes"`
	Servings         int32                `json:"servings"`
	YieldQuantity    *float64             `json:"yieldQuantity,omitempty"`
	YieldUnit        string               `json:"yieldUnit,omitempty"`
	MainIngredient   *IngredientRefJSON   `json:"mainIngredient,omitempty"`
	Cuisine          *CuisineJSON         `json:"cuisine,omitempty"`
	IngredientLines  []IngredientLineJSON `json:"ingredientLines"`
	Steps            []RecipeStepJSON     `json:"steps"`
	Tags             []string             `json:"tags,omitempty"`
	ImageURL         string               `json:"imageUrl,omitempty"`
	Nutrition        RecipeNutritionJSON  `json:"nutrition"`
}

// IngredientRefJSON is the JSON response for ingredient refs.
type IngredientRefJSON struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// IngredientLineJSON is the JSON response for ingredient lines.
type IngredientLineJSON struct {
	Ingredient    IngredientRefJSON `json:"ingredient"`
	QuantityValue *float64          `json:"quantityValue,omitempty"`
	QuantityText  string            `json:"quantityText,omitempty"`
	Unit          string            `json:"unit,omitempty"`
	IsOptional    bool              `json:"isOptional,omitempty"`
	Note          string            `json:"note,omitempty"`
	SortOrder     int32             `json:"sortOrder,omitempty"`
}

// CuisinesResponse is the response for cuisine listing.
type CuisinesResponse struct {
	Items []CuisineJSON `json:"items"`
}

// CuisineJSON is the JSON response for a cuisine.
type CuisineJSON struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func toRecipeJSON(r *recipepb.Recipe) RecipeJSON {
	ingredientLines := make([]IngredientLineJSON, len(r.GetIngredientLines()))
	for i, line := range r.GetIngredientLines() {
		var quantityValue *float64
		if line.GetQuantityValue() != nil {
			value := line.GetQuantityValue().GetValue()
			quantityValue = &value
		}
		ingredientLines[i] = IngredientLineJSON{
			Ingredient: IngredientRefJSON{
				ID:   line.GetIngredient().GetId(),
				Name: line.GetIngredient().GetName(),
			},
			QuantityValue: quantityValue,
			QuantityText:  line.GetQuantityText(),
			Unit:          line.GetUnit(),
			IsOptional:    line.GetIsOptional(),
			Note:          line.GetNote(),
			SortOrder:     line.GetSortOrder(),
		}
	}

	steps := make([]RecipeStepJSON, len(r.GetSteps()))
	for i, step := range r.GetSteps() {
		var duration *int
		if step.GetDurationSeconds() != nil {
			value := int(step.GetDurationSeconds().GetValue())
			duration = &value
		}
		var temperature *float64
		if step.GetTemperatureValue() != nil {
			value := step.GetTemperatureValue().GetValue()
			temperature = &value
		}
		steps[i] = RecipeStepJSON{
			StepIndex:        int(step.GetStepIndex()),
			Instruction:      step.GetInstruction(),
			DurationSeconds:  duration,
			TemperatureValue: temperature,
			TemperatureUnit:  step.GetTemperatureUnit(),
			MediaURL:         step.GetMediaUrl(),
		}
	}

	var yieldQuantity *float64
	if r.GetYieldQuantity() != nil {
		value := r.GetYieldQuantity().GetValue()
		yieldQuantity = &value
	}

	var mainIngredient *IngredientRefJSON
	if r.GetMainIngredient() != nil {
		mainIngredient = &IngredientRefJSON{
			ID:   r.GetMainIngredient().GetId(),
			Name: r.GetMainIngredient().GetName(),
		}
	}

	var cuisine *CuisineJSON
	if r.GetCuisine() != nil {
		cuisine = &CuisineJSON{
			ID:   r.GetCuisine().GetId(),
			Name: r.GetCuisine().GetName(),
		}
	}

	nutrition := RecipeNutritionJSON{}
	if r.GetNutrition() != nil {
		nutrition = RecipeNutritionJSON{
			CaloriesTotal:      int(r.GetNutrition().GetCaloriesTotal()),
			CaloriesPerServing: int(r.GetNutrition().GetCaloriesPerServing()),
			ProteinG:           r.GetNutrition().GetProteinG(),
			CarbsG:             r.GetNutrition().GetCarbsG(),
			FatG:               r.GetNutrition().GetFatG(),
			FiberG:             r.GetNutrition().GetFiberG(),
			SugarG:             r.GetNutrition().GetSugarG(),
			SodiumMg:           r.GetNutrition().GetSodiumMg(),
		}
	}

	return RecipeJSON{
		ID:               r.GetId(),
		Name:             r.GetName(),
		Description:      r.GetDescription(),
		PrepTimeMinutes:  r.GetPrepTimeMinutes(),
		CookTimeMinutes:  r.GetCookTimeMinutes(),
		TotalTimeMinutes: r.GetTotalTimeMinutes(),
		Servings:         r.GetServings(),
		YieldQuantity:    yieldQuantity,
		YieldUnit:        r.GetYieldUnit(),
		MainIngredient:   mainIngredient,
		Cuisine:          cuisine,
		IngredientLines:  ingredientLines,
		Steps:            steps,
		Tags:             r.GetTags(),
		ImageURL:         r.GetImageUrl(),
		Nutrition:        nutrition,
	}
}

func toRecipesJSON(recipes []*recipepb.Recipe) []RecipeJSON {
	result := make([]RecipeJSON, len(recipes))
	for i, r := range recipes {
		result[i] = toRecipeJSON(r)
	}
	return result
}

func toCuisinesJSON(cuisines []*recipepb.Cuisine) []CuisineJSON {
	items := make([]CuisineJSON, len(cuisines))
	for i, cuisine := range cuisines {
		items[i] = CuisineJSON{
			ID:   cuisine.GetId(),
			Name: cuisine.GetName(),
		}
	}
	return items
}

func parseIntParam(r *http.Request, name string, defaultValue int) int {
	val := r.URL.Query().Get(name)
	if val == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}
	return i
}

func sanitizeStrings(values []string) []string {
	cleaned := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		cleaned = append(cleaned, trimmed)
	}
	return cleaned
}

func splitCommaList(input string) []string {
	if strings.TrimSpace(input) == "" {
		return nil
	}
	parts := strings.Split(input, ",")
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}
	return sanitizeStrings(parts)
}

type errString string

func (e errString) Error() string { return string(e) }
