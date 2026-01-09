package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/platepilot/backend/internal/common/domain"
	"github.com/platepilot/backend/internal/recipe/repository"
)

// ShoppingListRepository defines the interface for shopping list data access
type ShoppingListRepository interface {
	GetByID(ctx context.Context, userID, id uuid.UUID) (*domain.ShoppingList, error)
	GetAll(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.ShoppingList, error)
	Create(ctx context.Context, list *domain.ShoppingList) error
	Update(ctx context.Context, list *domain.ShoppingList) error
	Delete(ctx context.Context, userID, id uuid.UUID) error
	AddItem(ctx context.Context, userID uuid.UUID, item *domain.ShoppingListItem) error
	UpdateItem(ctx context.Context, userID uuid.UUID, item *domain.ShoppingListItem) error
	ToggleItemChecked(ctx context.Context, userID, itemID uuid.UUID) (bool, error)
	DeleteItem(ctx context.Context, userID, itemID uuid.UUID) error
	GetAggregatedIngredients(ctx context.Context, userID uuid.UUID, recipeIDs []uuid.UUID) ([]domain.AggregatedIngredient, error)
	GetRecipesByIDs(ctx context.Context, userID uuid.UUID, recipeIDs []uuid.UUID) ([]domain.Recipe, error)
	Count(ctx context.Context, userID uuid.UUID) (int64, error)
}

// ShoppingListHandler handles REST requests for shopping lists
type ShoppingListHandler struct {
	repo   ShoppingListRepository
	logger *slog.Logger
}

// NewShoppingListHandler creates a new shopping list handler
func NewShoppingListHandler(repo ShoppingListRepository, logger *slog.Logger) *ShoppingListHandler {
	return &ShoppingListHandler{
		repo:   repo,
		logger: logger,
	}
}

// ---- JSON Types ----

// ShoppingListJSON is the JSON response for a shopping list
type ShoppingListJSON struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	WeekStartDate *string                `json:"weekStartDate,omitempty"`
	Items         []ShoppingListItemJSON `json:"items"`
	Recipes       []RecipeRefJSON        `json:"recipes"`
	TotalItems    int                    `json:"totalItems"`
	CheckedItems  int                    `json:"checkedItems"`
	CreatedAt     string                 `json:"createdAt"`
	UpdatedAt     string                 `json:"updatedAt"`
	CompletedAt   *string                `json:"completedAt,omitempty"`
}

// ShoppingListItemJSON is the JSON response for a shopping list item
type ShoppingListItemJSON struct {
	ID           string                   `json:"id"`
	Ingredient   *IngredientJSON          `json:"ingredient,omitempty"`
	Category     *IngredientCategoryJSON  `json:"category,omitempty"`
	CustomName   *string                  `json:"customName,omitempty"`
	Quantity     *float64                 `json:"quantity,omitempty"`
	Unit         *string                  `json:"unit,omitempty"`
	DisplayQty   string                   `json:"displayQuantity"`
	Checked      bool                     `json:"checked"`
	Notes        *string                  `json:"notes,omitempty"`
	IsCustom     bool                     `json:"isCustom"`
	Sources      []ItemSourceJSON         `json:"sources,omitempty"`
	CreatedAt    string                   `json:"createdAt"`
}

// IngredientJSON is the JSON response for ingredient refs.
type IngredientJSON struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// IngredientCategoryJSON is the JSON response for an ingredient category
type IngredientCategoryJSON struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	DisplayOrder int    `json:"displayOrder"`
}

// ItemSourceJSON is the JSON response for an item source
type ItemSourceJSON struct {
	RecipeID   string   `json:"recipeId"`
	RecipeName string   `json:"recipeName"`
	Quantity   *float64 `json:"quantity,omitempty"`
	Unit       *string  `json:"unit,omitempty"`
}

// RecipeRefJSON is a minimal recipe reference
type RecipeRefJSON struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ShoppingListSummaryJSON is a summary for list views
type ShoppingListSummaryJSON struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	TotalItems   int     `json:"totalItems"`
	CheckedItems int     `json:"checkedItems"`
	CreatedAt    string  `json:"createdAt"`
	UpdatedAt    string  `json:"updatedAt"`
	CompletedAt  *string `json:"completedAt,omitempty"`
}

// PaginatedShoppingListsJSON is the paginated response for shopping lists
type PaginatedShoppingListsJSON struct {
	Items      []ShoppingListSummaryJSON `json:"items"`
	PageIndex  int32                     `json:"pageIndex"`
	PageSize   int32                     `json:"pageSize"`
	TotalCount int32                     `json:"totalCount"`
	TotalPages int32                     `json:"totalPages"`
}

// ---- Request Types ----

// CreateFromRecipesRequest is the request to create a shopping list from recipes
type CreateFromRecipesRequest struct {
	Name          string   `json:"name"`
	RecipeIDs     []string `json:"recipeIds"`
	WeekStartDate *string  `json:"weekStartDate,omitempty"`
}

// CreateShoppingListRequest is the request to create an empty shopping list
type CreateShoppingListRequest struct {
	Name string `json:"name"`
}

// UpdateShoppingListRequest is the request to update a shopping list
type UpdateShoppingListRequest struct {
	Name string `json:"name"`
}

// AddItemRequest is the request to add an item to a shopping list
type AddItemRequest struct {
	IngredientID *string  `json:"ingredientId,omitempty"`
	CustomName   *string  `json:"customName,omitempty"`
	Quantity     *float64 `json:"quantity,omitempty"`
	Unit         *string  `json:"unit,omitempty"`
	Notes        *string  `json:"notes,omitempty"`
}

// UpdateItemRequest is the request to update an item
type UpdateItemRequest struct {
	Quantity *float64 `json:"quantity,omitempty"`
	Unit     *string  `json:"unit,omitempty"`
	Notes    *string  `json:"notes,omitempty"`
	Checked  *bool    `json:"checked,omitempty"`
}

// ---- Handlers ----

// GetByID handles GET /v1/shoppinglist/{id}
// @Summary      Get shopping list by ID
// @Description  Retrieves a single shopping list by its unique identifier
// @Tags         shoppinglists
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Shopping List ID (UUID)"
// @Success      200  {object}  ShoppingListJSON
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /shoppinglist/{id} [get]
func (h *ShoppingListHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid shopping list id")
		return
	}

	list, err := h.repo.GetByID(r.Context(), userID, id)
	if err != nil {
		if errors.Is(err, repository.ErrShoppingListNotFound) {
			writeError(w, http.StatusNotFound, "shopping list not found")
			return
		}
		h.logger.Error("failed to get shopping list", "id", id, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to get shopping list")
		return
	}

	writeJSON(w, http.StatusOK, toShoppingListJSON(list))
}

// GetAll handles GET /v1/shoppinglist
// @Summary      Get all shopping lists (paginated)
// @Description  Retrieves a paginated list of all shopping lists for the user
// @Tags         shoppinglists
// @Accept       json
// @Produce      json
// @Param        pageIndex  query     int  false  "Page number (1-indexed)"  default(1)
// @Param        pageSize   query     int  false  "Items per page (max 100)" default(20)
// @Success      200  {object}  PaginatedShoppingListsJSON
// @Failure      500  {object}  ErrorResponse
// @Router       /shoppinglist [get]
func (h *ShoppingListHandler) GetAll(w http.ResponseWriter, r *http.Request) {
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

	offset := (pageIndex - 1) * pageSize

	lists, err := h.repo.GetAll(r.Context(), userID, pageSize, offset)
	if err != nil {
		h.logger.Error("failed to get shopping lists", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to get shopping lists")
		return
	}

	totalCount, err := h.repo.Count(r.Context(), userID)
	if err != nil {
		h.logger.Error("failed to count shopping lists", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to count shopping lists")
		return
	}

	totalPages := int32((totalCount + int64(pageSize) - 1) / int64(pageSize))

	summaries := make([]ShoppingListSummaryJSON, len(lists))
	for i, list := range lists {
		summaries[i] = toShoppingListSummaryJSON(&list)
	}

	writeJSON(w, http.StatusOK, PaginatedShoppingListsJSON{
		Items:      summaries,
		PageIndex:  int32(pageIndex),
		PageSize:   int32(pageSize),
		TotalCount: int32(totalCount),
		TotalPages: totalPages,
	})
}

// CreateFromRecipes handles POST /v1/shoppinglist/from-recipes
// @Summary      Create shopping list from recipes
// @Description  Creates a new shopping list with aggregated ingredients from specified recipes
// @Tags         shoppinglists
// @Accept       json
// @Produce      json
// @Param        request  body      CreateFromRecipesRequest  true  "Recipe IDs to create list from"
// @Success      201      {object}  ShoppingListJSON
// @Failure      400      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /shoppinglist/from-recipes [post]
func (h *ShoppingListHandler) CreateFromRecipes(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req CreateFromRecipesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if len(req.RecipeIDs) == 0 {
		writeError(w, http.StatusBadRequest, "at least one recipe ID is required")
		return
	}

	// Parse recipe IDs
	recipeIDs := make([]uuid.UUID, len(req.RecipeIDs))
	for i, idStr := range req.RecipeIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid recipe ID: "+idStr)
			return
		}
		recipeIDs[i] = id
	}

	// Get aggregated ingredients
	aggregated, err := h.repo.GetAggregatedIngredients(r.Context(), userID, recipeIDs)
	if err != nil {
		h.logger.Error("failed to get aggregated ingredients", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to aggregate ingredients")
		return
	}

	// Get recipes for reference
	recipes, err := h.repo.GetRecipesByIDs(r.Context(), userID, recipeIDs)
	if err != nil {
		h.logger.Error("failed to get recipes", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to get recipes")
		return
	}

	// Generate name if not provided
	name := req.Name
	if name == "" {
		name = "Shopping List - " + time.Now().Format("Jan 2, 2006")
	}

	// Parse week start date if provided
	var weekStartDate *time.Time
	if req.WeekStartDate != nil && *req.WeekStartDate != "" {
		t, err := time.Parse("2006-01-02", *req.WeekStartDate)
		if err == nil {
			weekStartDate = &t
		}
	}

	// Build shopping list
	list := &domain.ShoppingList{
		ID:            uuid.New(),
		UserID:        userID,
		Name:          name,
		WeekStartDate: weekStartDate,
		Items:         make([]domain.ShoppingListItem, 0, len(aggregated)),
		Recipes:       recipes,
	}

	// Convert aggregated ingredients to items
	for _, agg := range aggregated {
		// For each unique quantity/unit combination, create an item
		for _, qu := range agg.Quantities {
			item := domain.ShoppingListItem{
				ID:           uuid.New(),
				IngredientID: &agg.IngredientID,
				Ingredient: &domain.Ingredient{
					ID:   agg.IngredientID,
					Name: agg.IngredientName,
				},
				Quantity: qu.Quantity,
				Unit:     qu.Unit,
				Checked:  false,
				IsCustom: false,
			}

			// Add category if available
			if agg.CategoryID != nil {
				catName := ""
				if agg.CategoryName != nil {
					catName = *agg.CategoryName
				}
				item.Category = &domain.IngredientCategory{
					ID:   *agg.CategoryID,
					Name: catName,
				}
			}

			// Add sources
			for _, src := range agg.RecipeSources {
				item.Sources = append(item.Sources, domain.ShoppingListItemSource{
					RecipeID:   src.RecipeID,
					RecipeName: src.RecipeName,
					Quantity:   src.Quantity,
					Unit:       src.Unit,
				})
			}

			list.Items = append(list.Items, item)
		}
	}

	// Save to database
	if err := h.repo.Create(r.Context(), list); err != nil {
		h.logger.Error("failed to create shopping list", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to create shopping list")
		return
	}

	h.logger.Info("shopping list created", "id", list.ID, "name", list.Name, "items", len(list.Items))

	writeJSON(w, http.StatusCreated, toShoppingListJSON(list))
}

// Create handles POST /v1/shoppinglist
// @Summary      Create empty shopping list
// @Description  Creates a new empty shopping list
// @Tags         shoppinglists
// @Accept       json
// @Produce      json
// @Param        request  body      CreateShoppingListRequest  true  "Shopping list details"
// @Success      201      {object}  ShoppingListJSON
// @Failure      400      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /shoppinglist [post]
func (h *ShoppingListHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req CreateShoppingListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" {
		req.Name = "Shopping List - " + time.Now().Format("Jan 2, 2006")
	}

	list := &domain.ShoppingList{
		ID:     uuid.New(),
		UserID: userID,
		Name:   req.Name,
		Items:  []domain.ShoppingListItem{},
	}

	if err := h.repo.Create(r.Context(), list); err != nil {
		h.logger.Error("failed to create shopping list", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to create shopping list")
		return
	}

	writeJSON(w, http.StatusCreated, toShoppingListJSON(list))
}

// Update handles PATCH /v1/shoppinglist/{id}
// @Summary      Update shopping list
// @Description  Updates a shopping list's name
// @Tags         shoppinglists
// @Accept       json
// @Produce      json
// @Param        id       path      string                     true  "Shopping List ID (UUID)"
// @Param        request  body      UpdateShoppingListRequest  true  "Update details"
// @Success      200      {object}  ShoppingListJSON
// @Failure      400      {object}  ErrorResponse
// @Failure      404      {object}  ErrorResponse
// @Router       /shoppinglist/{id} [patch]
func (h *ShoppingListHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid shopping list id")
		return
	}

	var req UpdateShoppingListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Get existing list
	list, err := h.repo.GetByID(r.Context(), userID, id)
	if err != nil {
		if errors.Is(err, repository.ErrShoppingListNotFound) {
			writeError(w, http.StatusNotFound, "shopping list not found")
			return
		}
		h.logger.Error("failed to get shopping list", "id", id, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to get shopping list")
		return
	}

	// Update fields
	if req.Name != "" {
		list.Name = req.Name
	}

	if err := h.repo.Update(r.Context(), list); err != nil {
		h.logger.Error("failed to update shopping list", "id", id, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to update shopping list")
		return
	}

	// Reload to get updated timestamp
	list, _ = h.repo.GetByID(r.Context(), userID, id)

	writeJSON(w, http.StatusOK, toShoppingListJSON(list))
}

// Delete handles DELETE /v1/shoppinglist/{id}
// @Summary      Delete shopping list
// @Description  Deletes a shopping list and all its items
// @Tags         shoppinglists
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Shopping List ID (UUID)"
// @Success      204  "No content"
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /shoppinglist/{id} [delete]
func (h *ShoppingListHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid shopping list id")
		return
	}

	if err := h.repo.Delete(r.Context(), userID, id); err != nil {
		if errors.Is(err, repository.ErrShoppingListNotFound) {
			writeError(w, http.StatusNotFound, "shopping list not found")
			return
		}
		h.logger.Error("failed to delete shopping list", "id", id, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to delete shopping list")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AddItem handles POST /v1/shoppinglist/{id}/items
// @Summary      Add item to shopping list
// @Description  Adds a new item to a shopping list
// @Tags         shoppinglists
// @Accept       json
// @Produce      json
// @Param        id       path      string          true  "Shopping List ID (UUID)"
// @Param        request  body      AddItemRequest  true  "Item details"
// @Success      201      {object}  ShoppingListItemJSON
// @Failure      400      {object}  ErrorResponse
// @Failure      404      {object}  ErrorResponse
// @Router       /shoppinglist/{id}/items [post]
func (h *ShoppingListHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	listIDStr := chi.URLParam(r, "id")
	listID, err := uuid.Parse(listIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid shopping list id")
		return
	}

	var req AddItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate: need either ingredient ID or custom name
	if req.IngredientID == nil && req.CustomName == nil {
		writeError(w, http.StatusBadRequest, "either ingredientId or customName is required")
		return
	}

	item := &domain.ShoppingListItem{
		ID:             uuid.New(),
		ShoppingListID: listID,
		CustomName:     req.CustomName,
		Quantity:       req.Quantity,
		Unit:           req.Unit,
		Notes:          req.Notes,
		Checked:        false,
		IsCustom:       req.CustomName != nil,
	}

	if req.IngredientID != nil {
		ingredientID, err := uuid.Parse(*req.IngredientID)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid ingredient id")
			return
		}
		item.IngredientID = &ingredientID
		item.IsCustom = false
	}

	if err := h.repo.AddItem(r.Context(), userID, item); err != nil {
		if errors.Is(err, repository.ErrShoppingListNotFound) {
			writeError(w, http.StatusNotFound, "shopping list not found")
			return
		}
		h.logger.Error("failed to add item", "listId", listID, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to add item")
		return
	}

	writeJSON(w, http.StatusCreated, toShoppingListItemJSON(item))
}

// UpdateItem handles PATCH /v1/shoppinglist/{id}/items/{itemId}
// @Summary      Update shopping list item
// @Description  Updates an item's quantity, unit, notes, or checked state
// @Tags         shoppinglists
// @Accept       json
// @Produce      json
// @Param        id       path      string              true  "Shopping List ID (UUID)"
// @Param        itemId   path      string              true  "Item ID (UUID)"
// @Param        request  body      UpdateItemRequest   true  "Update details"
// @Success      200      {object}  ShoppingListItemJSON
// @Failure      400      {object}  ErrorResponse
// @Failure      404      {object}  ErrorResponse
// @Router       /shoppinglist/{id}/items/{itemId} [patch]
func (h *ShoppingListHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	itemIDStr := chi.URLParam(r, "itemId")
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	var req UpdateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	item := &domain.ShoppingListItem{
		ID:       itemID,
		Quantity: req.Quantity,
		Unit:     req.Unit,
		Notes:    req.Notes,
	}

	if req.Checked != nil {
		item.Checked = *req.Checked
	}

	if err := h.repo.UpdateItem(r.Context(), userID, item); err != nil {
		if errors.Is(err, repository.ErrShoppingListItemNotFound) {
			writeError(w, http.StatusNotFound, "item not found")
			return
		}
		h.logger.Error("failed to update item", "itemId", itemID, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to update item")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"id": itemID.String()})
}

// ToggleItem handles POST /v1/shoppinglist/{id}/items/{itemId}/toggle
// @Summary      Toggle item checked state
// @Description  Toggles the checked state of an item
// @Tags         shoppinglists
// @Accept       json
// @Produce      json
// @Param        id       path      string  true  "Shopping List ID (UUID)"
// @Param        itemId   path      string  true  "Item ID (UUID)"
// @Success      200      {object}  map[string]bool
// @Failure      400      {object}  ErrorResponse
// @Failure      404      {object}  ErrorResponse
// @Router       /shoppinglist/{id}/items/{itemId}/toggle [post]
func (h *ShoppingListHandler) ToggleItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	itemIDStr := chi.URLParam(r, "itemId")
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	checked, err := h.repo.ToggleItemChecked(r.Context(), userID, itemID)
	if err != nil {
		if errors.Is(err, repository.ErrShoppingListItemNotFound) {
			writeError(w, http.StatusNotFound, "item not found")
			return
		}
		h.logger.Error("failed to toggle item", "itemId", itemID, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to toggle item")
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"checked": checked})
}

// DeleteItem handles DELETE /v1/shoppinglist/{id}/items/{itemId}
// @Summary      Delete shopping list item
// @Description  Removes an item from a shopping list
// @Tags         shoppinglists
// @Accept       json
// @Produce      json
// @Param        id       path      string  true  "Shopping List ID (UUID)"
// @Param        itemId   path      string  true  "Item ID (UUID)"
// @Success      204  "No content"
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /shoppinglist/{id}/items/{itemId} [delete]
func (h *ShoppingListHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	itemIDStr := chi.URLParam(r, "itemId")
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	if err := h.repo.DeleteItem(r.Context(), userID, itemID); err != nil {
		if errors.Is(err, repository.ErrShoppingListItemNotFound) {
			writeError(w, http.StatusNotFound, "item not found")
			return
		}
		h.logger.Error("failed to delete item", "itemId", itemID, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to delete item")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ---- Conversion Helpers ----

func toShoppingListJSON(list *domain.ShoppingList) ShoppingListJSON {
	items := make([]ShoppingListItemJSON, len(list.Items))
	checkedCount := 0
	for i, item := range list.Items {
		items[i] = toShoppingListItemJSON(&item)
		if item.Checked {
			checkedCount++
		}
	}

	recipes := make([]RecipeRefJSON, len(list.Recipes))
	for i, r := range list.Recipes {
		recipes[i] = RecipeRefJSON{
			ID:   r.ID.String(),
			Name: r.Name,
		}
	}

	var weekStartDate *string
	if list.WeekStartDate != nil {
		s := list.WeekStartDate.Format("2006-01-02")
		weekStartDate = &s
	}

	var completedAt *string
	if list.CompletedAt != nil {
		s := list.CompletedAt.Format(time.RFC3339)
		completedAt = &s
	}

	return ShoppingListJSON{
		ID:            list.ID.String(),
		Name:          list.Name,
		WeekStartDate: weekStartDate,
		Items:         items,
		Recipes:       recipes,
		TotalItems:    len(list.Items),
		CheckedItems:  checkedCount,
		CreatedAt:     list.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     list.UpdatedAt.Format(time.RFC3339),
		CompletedAt:   completedAt,
	}
}

func toShoppingListItemJSON(item *domain.ShoppingListItem) ShoppingListItemJSON {
	var ingredient *IngredientJSON
	if item.Ingredient != nil {
		ingredient = &IngredientJSON{
			ID:   item.Ingredient.ID.String(),
			Name: item.Ingredient.Name,
		}
	}

	var category *IngredientCategoryJSON
	if item.Category != nil {
		category = &IngredientCategoryJSON{
			ID:           item.Category.ID.String(),
			Name:         item.Category.Name,
			DisplayOrder: item.Category.DisplayOrder,
		}
	}

	sources := make([]ItemSourceJSON, len(item.Sources))
	for i, src := range item.Sources {
		sources[i] = ItemSourceJSON{
			RecipeID:   src.RecipeID.String(),
			RecipeName: src.RecipeName,
			Quantity:   src.Quantity,
			Unit:       src.Unit,
		}
	}

	displayQty := item.DisplayQuantity()

	return ShoppingListItemJSON{
		ID:         item.ID.String(),
		Ingredient: ingredient,
		Category:   category,
		CustomName: item.CustomName,
		Quantity:   item.Quantity,
		Unit:       item.Unit,
		DisplayQty: displayQty,
		Checked:    item.Checked,
		Notes:      item.Notes,
		IsCustom:   item.IsCustom,
		Sources:    sources,
		CreatedAt:  item.CreatedAt.Format(time.RFC3339),
	}
}

func toShoppingListSummaryJSON(list *domain.ShoppingList) ShoppingListSummaryJSON {
	checkedCount := 0
	for _, item := range list.Items {
		if item.Checked {
			checkedCount++
		}
	}

	var completedAt *string
	if list.CompletedAt != nil {
		s := list.CompletedAt.Format(time.RFC3339)
		completedAt = &s
	}

	return ShoppingListSummaryJSON{
		ID:           list.ID.String(),
		Name:         list.Name,
		TotalItems:   len(list.Items),
		CheckedItems: checkedCount,
		CreatedAt:    list.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    list.UpdatedAt.Format(time.RFC3339),
		CompletedAt:  completedAt,
	}
}
