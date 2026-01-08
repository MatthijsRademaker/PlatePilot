package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/platepilot/backend/internal/bff/client"
	recipepb "github.com/platepilot/backend/internal/recipe/pb"
)

// RecipeHandler handles REST requests for recipes
type RecipeHandler struct {
	client *client.RecipeClient
	logger *slog.Logger
}

// NewRecipeHandler creates a new recipe handler
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

// PaginatedRecipesJSON is the paginated response for recipes
type PaginatedRecipesJSON struct {
	Items      []RecipeJSON `json:"items"`
	PageIndex  int32        `json:"pageIndex"`
	PageSize   int32        `json:"pageSize"`
	TotalCount int32        `json:"totalCount"`
	TotalPages int32        `json:"totalPages"`
}

// GetAll handles GET /v1/recipe/all
// @Summary      Get all recipes (paginated)
// @Description  Retrieves a paginated list of all recipes
// @Tags         recipes
// @Accept       json
// @Produce      json
// @Param        pageIndex  query     int     false  "Page number (1-indexed)"  default(1)
// @Param        pageSize   query     int     false  "Items per page (max 100)" default(20)
// @Success      200  {object}  PaginatedRecipesJSON
// @Failure      500  {object}  ErrorResponse
// @Router       /recipe/all [get]
func (h *RecipeHandler) GetAll(w http.ResponseWriter, r *http.Request) {
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

	resp, err := h.client.GetAll(r.Context(), userID.String(), int32(pageIndex), int32(pageSize))
	if err != nil {
		h.logger.Error("failed to get recipes", "error", err)
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

// GetByCuisine handles GET /v1/recipe/cuisine/{id}
// @Summary      Get recipes by cuisine
// @Description  Retrieves all recipes belonging to a specific cuisine
// @Tags         recipes
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Cuisine ID (UUID)"
// @Success      200  {array}   RecipeJSON
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /recipe/cuisine/{id} [get]
func (h *RecipeHandler) GetByCuisine(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	cuisineID := chi.URLParam(r, "id")
	if cuisineID == "" {
		writeError(w, http.StatusBadRequest, "cuisine id is required")
		return
	}

	recipes, err := h.client.GetByCuisine(r.Context(), userID.String(), cuisineID)
	if err != nil {
		h.logger.Error("failed to get recipes by cuisine", "cuisineId", cuisineID, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to fetch recipes")
		return
	}

	writeJSON(w, http.StatusOK, toRecipesJSON(recipes))
}

// GetByIngredient handles GET /v1/recipe/ingredient/{id}
// @Summary      Get recipes by ingredient
// @Description  Retrieves all recipes containing a specific ingredient
// @Tags         recipes
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Ingredient ID (UUID)"
// @Success      200  {array}   RecipeJSON
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /recipe/ingredient/{id} [get]
func (h *RecipeHandler) GetByIngredient(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	ingredientID := chi.URLParam(r, "id")
	if ingredientID == "" {
		writeError(w, http.StatusBadRequest, "ingredient id is required")
		return
	}

	recipes, err := h.client.GetByIngredient(r.Context(), userID.String(), ingredientID)
	if err != nil {
		h.logger.Error("failed to get recipes by ingredient", "ingredientId", ingredientID, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to fetch recipes")
		return
	}

	writeJSON(w, http.StatusOK, toRecipesJSON(recipes))
}

// GetByAllergy handles GET /v1/recipe/allergy/{id}
// @Summary      Get recipes avoiding allergen
// @Description  Retrieves all recipes that do not contain a specific allergen
// @Tags         recipes
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Allergy ID (UUID)"
// @Success      200  {array}   RecipeJSON
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /recipe/allergy/{id} [get]
func (h *RecipeHandler) GetByAllergy(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	allergyID := chi.URLParam(r, "id")
	if allergyID == "" {
		writeError(w, http.StatusBadRequest, "allergy id is required")
		return
	}

	recipes, err := h.client.GetByAllergy(r.Context(), userID.String(), allergyID)
	if err != nil {
		h.logger.Error("failed to get recipes by allergy", "allergyId", allergyID, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to fetch recipes")
		return
	}

	writeJSON(w, http.StatusOK, toRecipesJSON(recipes))
}

// Create handles POST /v1/recipe/create
// @Summary      Create a new recipe
// @Description  Creates a new recipe with the provided details
// @Tags         recipes
// @Accept       json
// @Produce      json
// @Param        recipe  body      CreateRecipeRequest  true  "Recipe to create"
// @Success      201     {object}  RecipeJSON
// @Failure      400     {object}  ErrorResponse
// @Failure      500     {object}  ErrorResponse
// @Router       /recipe/create [post]
func (h *RecipeHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req CreateRecipeRequest
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

// CreateRecipeRequest is the request body for creating a recipe
type CreateRecipeRequest struct {
	Name               string   `json:"name"`
	Description        string   `json:"description"`
	PrepTime           string   `json:"prepTime"`
	CookTime           string   `json:"cookTime"`
	MainIngredientID   string   `json:"mainIngredientId"`
	MainIngredientName string   `json:"mainIngredientName"`
	CuisineID          string   `json:"cuisineId"`
	CuisineName        string   `json:"cuisineName"`
	IngredientIDs      []string `json:"ingredientIds"`
	IngredientNames    []string `json:"ingredientNames"`
	Directions         []string `json:"directions"`
	Tags               []string `json:"tags"`
	GuidedMode         bool     `json:"guidedMode"`
}

// Validate validates the create recipe request
func (r *CreateRecipeRequest) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return &ValidationError{Field: "name", Message: "name is required"}
	}
	hasMainIngredient := strings.TrimSpace(r.MainIngredientID) != "" || strings.TrimSpace(r.MainIngredientName) != ""
	hasIngredients := len(r.IngredientIDs) > 0 || hasNonEmpty(r.IngredientNames)
	if !hasMainIngredient && !hasIngredients {
		return &ValidationError{Field: "ingredients", Message: "at least one ingredient is required"}
	}
	if !hasNonEmpty(r.Directions) {
		return &ValidationError{Field: "directions", Message: "at least one direction is required"}
	}
	return nil
}

// ToProto converts the request to a protobuf message
func (r *CreateRecipeRequest) ToProto(userID string) *recipepb.CreateRecipeRequest {
	return &recipepb.CreateRecipeRequest{
		UserId:             userID,
		Name:               r.Name,
		Description:        r.Description,
		PrepTime:           r.PrepTime,
		CookTime:           r.CookTime,
		MainIngredientId:   r.MainIngredientID,
		MainIngredientName: strings.TrimSpace(r.MainIngredientName),
		CuisineId:          r.CuisineID,
		CuisineName:        strings.TrimSpace(r.CuisineName),
		IngredientIds:      r.IngredientIDs,
		IngredientNames:    sanitizeStrings(r.IngredientNames),
		Directions:         sanitizeStrings(r.Directions),
		Tags:               sanitizeStrings(r.Tags),
		GuidedMode:         r.GuidedMode,
	}
}

func hasNonEmpty(values []string) bool {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return true
		}
	}
	return false
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

// RecipeJSON is the JSON response for a recipe
type RecipeJSON struct {
	ID             string           `json:"id"`
	Name           string           `json:"name"`
	Description    string           `json:"description"`
	PrepTime       string           `json:"prepTime"`
	CookTime       string           `json:"cookTime"`
	MainIngredient *IngredientJSON  `json:"mainIngredient,omitempty"`
	Cuisine        *CuisineJSON     `json:"cuisine,omitempty"`
	Ingredients    []IngredientJSON `json:"ingredients"`
	Directions     []string         `json:"directions"`
}

// IngredientJSON is the JSON response for an ingredient
type IngredientJSON struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CuisineJSON is the JSON response for a cuisine
type CuisineJSON struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func toRecipeJSON(r *recipepb.RecipeResponse) RecipeJSON {
	ingredients := make([]IngredientJSON, len(r.GetIngredients()))
	for i, ing := range r.GetIngredients() {
		ingredients[i] = IngredientJSON{
			ID:   ing.GetId(),
			Name: ing.GetName(),
		}
	}

	var mainIngredient *IngredientJSON
	if r.GetMainIngredient() != nil {
		mainIngredient = &IngredientJSON{
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

	return RecipeJSON{
		ID:             r.GetId(),
		Name:           r.GetName(),
		Description:    r.GetDescription(),
		PrepTime:       r.GetPrepTime(),
		CookTime:       r.GetCookTime(),
		MainIngredient: mainIngredient,
		Cuisine:        cuisine,
		Ingredients:    ingredients,
		Directions:     r.GetDirections(),
	}
}

func toRecipesJSON(recipes []*recipepb.RecipeResponse) []RecipeJSON {
	result := make([]RecipeJSON, len(recipes))
	for i, r := range recipes {
		result[i] = toRecipeJSON(r)
	}
	return result
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
