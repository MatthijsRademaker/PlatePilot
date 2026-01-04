package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

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
func (h *RecipeHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "recipe id is required")
		return
	}

	recipe, err := h.client.GetByID(r.Context(), id)
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
func (h *RecipeHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	pageIndex := parseIntParam(r, "pageIndex", 1)
	pageSize := parseIntParam(r, "pageSize", 20)

	if pageSize > 100 {
		pageSize = 100
	}

	resp, err := h.client.GetAll(r.Context(), int32(pageIndex), int32(pageSize))
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
func (h *RecipeHandler) GetSimilar(w http.ResponseWriter, r *http.Request) {
	recipeID := r.URL.Query().Get("recipe")
	if recipeID == "" {
		writeError(w, http.StatusBadRequest, "recipe parameter is required")
		return
	}

	amount := parseIntParam(r, "amount", 5)
	if amount > 50 {
		amount = 50
	}

	recipes, err := h.client.GetSimilar(r.Context(), recipeID, int32(amount))
	if err != nil {
		h.logger.Error("failed to get similar recipes", "recipeId", recipeID, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to fetch similar recipes")
		return
	}

	writeJSON(w, http.StatusOK, toRecipesJSON(recipes))
}

// GetByCuisine handles GET /v1/recipe/cuisine/{id}
func (h *RecipeHandler) GetByCuisine(w http.ResponseWriter, r *http.Request) {
	cuisineID := chi.URLParam(r, "id")
	if cuisineID == "" {
		writeError(w, http.StatusBadRequest, "cuisine id is required")
		return
	}

	recipes, err := h.client.GetByCuisine(r.Context(), cuisineID)
	if err != nil {
		h.logger.Error("failed to get recipes by cuisine", "cuisineId", cuisineID, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to fetch recipes")
		return
	}

	writeJSON(w, http.StatusOK, toRecipesJSON(recipes))
}

// GetByIngredient handles GET /v1/recipe/ingredient/{id}
func (h *RecipeHandler) GetByIngredient(w http.ResponseWriter, r *http.Request) {
	ingredientID := chi.URLParam(r, "id")
	if ingredientID == "" {
		writeError(w, http.StatusBadRequest, "ingredient id is required")
		return
	}

	recipes, err := h.client.GetByIngredient(r.Context(), ingredientID)
	if err != nil {
		h.logger.Error("failed to get recipes by ingredient", "ingredientId", ingredientID, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to fetch recipes")
		return
	}

	writeJSON(w, http.StatusOK, toRecipesJSON(recipes))
}

// GetByAllergy handles GET /v1/recipe/allergy/{id}
func (h *RecipeHandler) GetByAllergy(w http.ResponseWriter, r *http.Request) {
	allergyID := chi.URLParam(r, "id")
	if allergyID == "" {
		writeError(w, http.StatusBadRequest, "allergy id is required")
		return
	}

	recipes, err := h.client.GetByAllergy(r.Context(), allergyID)
	if err != nil {
		h.logger.Error("failed to get recipes by allergy", "allergyId", allergyID, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to fetch recipes")
		return
	}

	writeJSON(w, http.StatusOK, toRecipesJSON(recipes))
}

// Create handles POST /v1/recipe/create
func (h *RecipeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRecipeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := req.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	recipe, err := h.client.Create(r.Context(), req.ToProto())
	if err != nil {
		h.logger.Error("failed to create recipe", "name", req.Name, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to create recipe")
		return
	}

	writeJSON(w, http.StatusCreated, toRecipeJSON(recipe))
}

// CreateRecipeRequest is the request body for creating a recipe
type CreateRecipeRequest struct {
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	PrepTime         string   `json:"prepTime"`
	CookTime         string   `json:"cookTime"`
	MainIngredientID string   `json:"mainIngredientId"`
	CuisineID        string   `json:"cuisineId"`
	IngredientIDs    []string `json:"ingredientIds"`
	Directions       []string `json:"directions"`
}

// Validate validates the create recipe request
func (r *CreateRecipeRequest) Validate() error {
	if r.Name == "" {
		return &ValidationError{Field: "name", Message: "name is required"}
	}
	if r.MainIngredientID == "" {
		return &ValidationError{Field: "mainIngredientId", Message: "mainIngredientId is required"}
	}
	if r.CuisineID == "" {
		return &ValidationError{Field: "cuisineId", Message: "cuisineId is required"}
	}
	return nil
}

// ToProto converts the request to a protobuf message
func (r *CreateRecipeRequest) ToProto() *recipepb.CreateRecipeRequest {
	return &recipepb.CreateRecipeRequest{
		Name:             r.Name,
		Description:      r.Description,
		PrepTime:         r.PrepTime,
		CookTime:         r.CookTime,
		MainIngredientId: r.MainIngredientID,
		CuisineId:        r.CuisineID,
		IngredientIds:    r.IngredientIDs,
		Directions:       r.Directions,
	}
}

// RecipeJSON is the JSON response for a recipe
type RecipeJSON struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	Description    string          `json:"description"`
	PrepTime       string          `json:"prepTime"`
	CookTime       string          `json:"cookTime"`
	MainIngredient *IngredientJSON `json:"mainIngredient,omitempty"`
	Cuisine        *CuisineJSON    `json:"cuisine,omitempty"`
	Ingredients    []IngredientJSON `json:"ingredients"`
	Directions     []string        `json:"directions"`
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
