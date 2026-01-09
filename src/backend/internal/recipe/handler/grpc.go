package handler

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/platepilot/backend/internal/common/domain"
	"github.com/platepilot/backend/internal/common/vector"
	pb "github.com/platepilot/backend/internal/recipe/pb"
	"github.com/platepilot/backend/internal/recipe/repository"
)

// GRPCHandler implements the RecipeService gRPC interface
type GRPCHandler struct {
	pb.UnimplementedRecipeServiceServer
	repo      RecipeRepository
	vectorGen vector.Generator
	publisher EventPublisher
	logger    *slog.Logger
}

const (
	defaultCuisineName = "General"
	guidedModeTag      = "guided-mode"
)

// NewGRPCHandler creates a new gRPC handler
func NewGRPCHandler(repo RecipeRepository, vectorGen vector.Generator, publisher EventPublisher, logger *slog.Logger) *GRPCHandler {
	return &GRPCHandler{
		repo:      repo,
		vectorGen: vectorGen,
		publisher: publisher,
		logger:    logger,
	}
}

// GetRecipeById retrieves a recipe by ID
func (h *GRPCHandler) GetRecipeById(ctx context.Context, req *pb.GetRecipeByIdRequest) (*pb.RecipeResponse, error) {
	h.logger.Debug("get recipe by id", "recipeId", req.GetRecipeId(), "userId", req.GetUserId())

	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	id, err := uuid.Parse(req.GetRecipeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid recipe ID: %v", err)
	}

	recipe, err := h.repo.GetByID(ctx, userID, id)
	if err != nil {
		if errors.Is(err, repository.ErrRecipeNotFound) {
			return nil, status.Errorf(codes.NotFound, "recipe not found")
		}
		h.logger.Error("failed to get recipe", "error", err, "recipeId", id)
		return nil, status.Errorf(codes.Internal, "failed to get recipe")
	}

	return toRecipeResponse(recipe), nil
}

// GetAllRecipes retrieves all recipes with pagination
func (h *GRPCHandler) GetAllRecipes(ctx context.Context, req *pb.GetAllRecipesRequest) (*pb.GetAllRecipesResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	pageIndex := int(req.GetPageIndex())
	pageSize := int(req.GetPageSize())

	if pageIndex < 1 {
		pageIndex = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (pageIndex - 1) * pageSize

	h.logger.Debug("get all recipes", "pageIndex", pageIndex, "pageSize", pageSize, "offset", offset)

	recipes, err := h.repo.GetAll(ctx, userID, pageSize, offset)
	if err != nil {
		h.logger.Error("failed to get recipes", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get recipes")
	}

	totalCount, err := h.repo.Count(ctx, userID)
	if err != nil {
		h.logger.Error("failed to count recipes", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to count recipes")
	}

	totalPages := int32((totalCount + int64(pageSize) - 1) / int64(pageSize))

	return toRecipesResponseWithPagination(recipes, int32(pageIndex), int32(pageSize), int32(totalCount), totalPages), nil
}

// CreateRecipe creates a new recipe
func (h *GRPCHandler) CreateRecipe(ctx context.Context, req *pb.CreateRecipeRequest) (*pb.RecipeResponse, error) {
	h.logger.Debug("create recipe", "name", req.GetName(), "userId", req.GetUserId())

	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	// Validate required fields
	if req.GetName() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "name is required")
	}

	ingredients, err := h.resolveIngredients(ctx, req)
	if err != nil {
		return nil, err
	}

	mainIngredient, err := h.resolveMainIngredient(ctx, req, ingredients)
	if err != nil {
		return nil, err
	}

	cuisine, err := h.resolveCuisine(ctx, req)
	if err != nil {
		return nil, err
	}

	ingredients = ensureMainIngredientIncluded(ingredients, mainIngredient)

	// Build the recipe
	recipe := &domain.Recipe{
		ID:             uuid.New(),
		UserID:         userID,
		Name:           req.GetName(),
		Description:    req.GetDescription(),
		PrepTime:       req.GetPrepTime(),
		CookTime:       req.GetCookTime(),
		MainIngredient: mainIngredient,
		Cuisine:        cuisine,
		Ingredients:    ingredients,
		Directions:     req.GetDirections(),
		Metadata: domain.Metadata{
			PublishedDate: time.Now().UTC(),
			Tags:          normalizeTags(req.GetTags(), req.GetGuidedMode()),
		},
	}

	// Generate vector embedding
	recipe.Metadata.SearchVector = h.vectorGen.GenerateForRecipe(recipe)

	// Save recipe
	if err := h.repo.Create(ctx, recipe); err != nil {
		h.logger.Error("failed to create recipe", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to create recipe")
	}

	h.logger.Info("recipe created", "recipeId", recipe.ID, "name", recipe.Name)

	// Publish event (non-blocking - don't fail the request if publishing fails)
	if h.publisher != nil {
		if err := h.publisher.PublishRecipeCreated(ctx, recipe); err != nil {
			h.logger.Error("failed to publish recipe created event",
				"error", err,
				"recipeId", recipe.ID,
			)
			// Don't return error - the recipe was created successfully
		}
	}

	return toRecipeResponse(recipe), nil
}

// GetSimilarRecipes retrieves similar recipes using vector similarity
func (h *GRPCHandler) GetSimilarRecipes(ctx context.Context, req *pb.GetSimilarRecipesRequest) (*pb.GetAllRecipesResponse, error) {
	h.logger.Debug("get similar recipes", "recipeId", req.GetRecipeId(), "amount", req.GetAmount(), "userId", req.GetUserId())

	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	recipeID, err := uuid.Parse(req.GetRecipeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid recipe ID: %v", err)
	}

	amount := int(req.GetAmount())
	if amount < 1 {
		amount = 5
	}
	if amount > 50 {
		amount = 50
	}

	recipes, err := h.repo.GetSimilar(ctx, userID, recipeID, amount)
	if err != nil {
		if errors.Is(err, repository.ErrRecipeNotFound) {
			return nil, status.Errorf(codes.NotFound, "recipe not found")
		}
		h.logger.Error("failed to get similar recipes", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get similar recipes")
	}

	return toRecipesResponse(recipes), nil
}

// GetRecipesByCuisine retrieves recipes by cuisine
func (h *GRPCHandler) GetRecipesByCuisine(ctx context.Context, req *pb.GetRecipesByCuisineRequest) (*pb.GetAllRecipesResponse, error) {
	h.logger.Debug("get recipes by cuisine", "cuisineId", req.GetCuisineId(), "userId", req.GetUserId())

	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	cuisineID, err := uuid.Parse(req.GetCuisineId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid cuisine ID: %v", err)
	}

	// Using default pagination for now
	recipes, err := h.repo.GetByCuisine(ctx, userID, cuisineID, 100, 0)
	if err != nil {
		h.logger.Error("failed to get recipes by cuisine", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get recipes by cuisine")
	}

	return toRecipesResponse(recipes), nil
}

// GetRecipesByIngredient retrieves recipes containing a specific ingredient
func (h *GRPCHandler) GetRecipesByIngredient(ctx context.Context, req *pb.GetRecipesByIngredientRequest) (*pb.GetAllRecipesResponse, error) {
	h.logger.Debug("get recipes by ingredient", "ingredientId", req.GetIngredientId(), "userId", req.GetUserId())

	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	ingredientID, err := uuid.Parse(req.GetIngredientId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid ingredient ID: %v", err)
	}

	// Using default pagination for now
	recipes, err := h.repo.GetByIngredient(ctx, userID, ingredientID, 100, 0)
	if err != nil {
		h.logger.Error("failed to get recipes by ingredient", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get recipes by ingredient")
	}

	return toRecipesResponse(recipes), nil
}

// GetRecipesByAllergy retrieves recipes excluding a specific allergy
func (h *GRPCHandler) GetRecipesByAllergy(ctx context.Context, req *pb.GetRecipesByAllergyRequest) (*pb.GetAllRecipesResponse, error) {
	h.logger.Debug("get recipes excluding allergy", "allergyId", req.GetAllergyId(), "userId", req.GetUserId())

	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	allergyID, err := uuid.Parse(req.GetAllergyId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid allergy ID: %v", err)
	}

	// Using default pagination for now - returns recipes WITHOUT this allergy
	recipes, err := h.repo.GetExcludingAllergy(ctx, userID, allergyID, 100, 0)
	if err != nil {
		h.logger.Error("failed to get recipes excluding allergy", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get recipes excluding allergy")
	}

	return toRecipesResponse(recipes), nil
}

// GetUnits retrieves available ingredient units.
func (h *GRPCHandler) GetUnits(ctx context.Context, req *pb.GetUnitsRequest) (*pb.GetUnitsResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	units, err := h.repo.GetUnits(ctx, userID)
	if err != nil {
		h.logger.Error("failed to get units", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get units")
	}

	resp := &pb.GetUnitsResponse{Units: make([]*pb.Unit, len(units))}
	for i, unit := range units {
		resp.Units[i] = &pb.Unit{
			Id:   unit.ID.String(),
			Name: unit.Name,
		}
	}

	return resp, nil
}

// CreateUnit creates a new ingredient unit.
func (h *GRPCHandler) CreateUnit(ctx context.Context, req *pb.CreateUnitRequest) (*pb.Unit, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	name := strings.TrimSpace(req.GetName())
	if name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "name is required")
	}

	existing, err := h.repo.GetUnitByName(ctx, userID, name)
	if err == nil {
		return &pb.Unit{
			Id:   existing.ID.String(),
			Name: existing.Name,
		}, nil
	}
	if err != nil && !errors.Is(err, repository.ErrUnitNotFound) {
		h.logger.Error("failed to get unit", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get unit")
	}

	unit := &domain.Unit{
		ID:     uuid.New(),
		UserID: userID,
		Name:   name,
	}
	if err := h.repo.CreateUnit(ctx, unit); err != nil {
		h.logger.Error("failed to create unit", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to create unit")
	}

	return &pb.Unit{
		Id:   unit.ID.String(),
		Name: unit.Name,
	}, nil
}

func (h *GRPCHandler) resolveIngredients(ctx context.Context, req *pb.CreateRecipeRequest) ([]domain.Ingredient, error) {
	var ingredients []domain.Ingredient
	indexByID := make(map[uuid.UUID]int)

	appendIngredient := func(ingredient *domain.Ingredient, quantity string, unit string) {
		if ingredient == nil {
			return
		}
		if index, ok := indexByID[ingredient.ID]; ok {
			if quantity != "" {
				ingredients[index].Quantity = quantity
			}
			if unit != "" {
				ingredients[index].Unit = unit
			}
			return
		}
		ingredient.Quantity = quantity
		ingredient.Unit = unit
		indexByID[ingredient.ID] = len(ingredients)
		ingredients = append(ingredients, *ingredient)
	}

	for _, input := range req.GetIngredients() {
		idStr := strings.TrimSpace(input.GetId())
		name := strings.TrimSpace(input.GetName())
		quantity := strings.TrimSpace(input.GetQuantity())
		unit := strings.TrimSpace(input.GetUnit())

		if idStr != "" {
			ingredientID, err := uuid.Parse(idStr)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "invalid ingredient ID: %v", err)
			}
			ingredient, err := h.repo.GetIngredientByID(ctx, ingredientID)
			if err != nil {
				if errors.Is(err, repository.ErrIngredientNotFound) {
					return nil, status.Errorf(codes.NotFound, "ingredient not found: %s", idStr)
				}
				h.logger.Error("failed to get ingredient", "error", err, "ingredientId", idStr)
				return nil, status.Errorf(codes.Internal, "failed to get ingredient")
			}
			appendIngredient(ingredient, quantity, unit)
			continue
		}

		if name == "" {
			continue
		}
		ingredient, err := h.repo.GetOrCreateIngredient(ctx, name, quantity)
		if err != nil {
			h.logger.Error("failed to get or create ingredient", "error", err, "ingredientName", name)
			return nil, status.Errorf(codes.Internal, "failed to create ingredient")
		}
		appendIngredient(ingredient, quantity, unit)
	}

	for _, idStr := range req.GetIngredientIds() {
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			continue
		}
		ingredientID, err := uuid.Parse(idStr)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid ingredient ID: %v", err)
		}
		ingredient, err := h.repo.GetIngredientByID(ctx, ingredientID)
		if err != nil {
			if errors.Is(err, repository.ErrIngredientNotFound) {
				return nil, status.Errorf(codes.NotFound, "ingredient not found: %s", idStr)
			}
			h.logger.Error("failed to get ingredient", "error", err, "ingredientId", idStr)
			return nil, status.Errorf(codes.Internal, "failed to get ingredient")
		}
		appendIngredient(ingredient, "", "")
	}

	for _, name := range req.GetIngredientNames() {
		cleaned := strings.TrimSpace(name)
		if cleaned == "" {
			continue
		}
		ingredient, err := h.repo.GetOrCreateIngredient(ctx, cleaned, "")
		if err != nil {
			h.logger.Error("failed to get or create ingredient", "error", err, "ingredientName", cleaned)
			return nil, status.Errorf(codes.Internal, "failed to create ingredient")
		}
		appendIngredient(ingredient, "", "")
	}

	return ingredients, nil
}

func (h *GRPCHandler) resolveMainIngredient(
	ctx context.Context,
	req *pb.CreateRecipeRequest,
	ingredients []domain.Ingredient,
) (*domain.Ingredient, error) {
	mainIngredientID := strings.TrimSpace(req.GetMainIngredientId())
	if mainIngredientID != "" {
		parsedID, err := uuid.Parse(mainIngredientID)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid main ingredient ID: %v", err)
		}
		ingredient, err := h.repo.GetIngredientByID(ctx, parsedID)
		if err != nil {
			if errors.Is(err, repository.ErrIngredientNotFound) {
				return nil, status.Errorf(codes.NotFound, "main ingredient not found")
			}
			h.logger.Error("failed to get main ingredient", "error", err)
			return nil, status.Errorf(codes.Internal, "failed to get main ingredient")
		}
		return ingredient, nil
	}

	mainIngredientName := strings.TrimSpace(req.GetMainIngredientName())
	if mainIngredientName != "" {
		ingredient, err := h.repo.GetOrCreateIngredient(ctx, mainIngredientName, "")
		if err != nil {
			h.logger.Error("failed to get or create main ingredient", "error", err, "ingredientName", mainIngredientName)
			return nil, status.Errorf(codes.Internal, "failed to create main ingredient")
		}
		return ingredient, nil
	}

	if len(ingredients) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "at least one ingredient is required")
	}

	return &ingredients[0], nil
}

func (h *GRPCHandler) resolveCuisine(ctx context.Context, req *pb.CreateRecipeRequest) (*domain.Cuisine, error) {
	cuisineID := strings.TrimSpace(req.GetCuisineId())
	if cuisineID != "" {
		parsedID, err := uuid.Parse(cuisineID)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid cuisine ID: %v", err)
		}
		cuisine, err := h.repo.GetCuisineByID(ctx, parsedID)
		if err != nil {
			if errors.Is(err, repository.ErrCuisineNotFound) {
				return nil, status.Errorf(codes.NotFound, "cuisine not found")
			}
			h.logger.Error("failed to get cuisine", "error", err)
			return nil, status.Errorf(codes.Internal, "failed to get cuisine")
		}
		return cuisine, nil
	}

	cuisineName := strings.TrimSpace(req.GetCuisineName())
	if cuisineName == "" {
		cuisineName = defaultCuisineName
	}
	cuisine, err := h.repo.GetOrCreateCuisine(ctx, cuisineName)
	if err != nil {
		h.logger.Error("failed to get or create cuisine", "error", err, "cuisineName", cuisineName)
		return nil, status.Errorf(codes.Internal, "failed to create cuisine")
	}
	return cuisine, nil
}

func ensureMainIngredientIncluded(ingredients []domain.Ingredient, mainIngredient *domain.Ingredient) []domain.Ingredient {
	if mainIngredient == nil {
		return ingredients
	}

	for _, ingredient := range ingredients {
		if ingredient.ID == mainIngredient.ID {
			mainIngredient.Quantity = ingredient.Quantity
			mainIngredient.Unit = ingredient.Unit
			return ingredients
		}
	}

	return append([]domain.Ingredient{*mainIngredient}, ingredients...)
}

func normalizeTags(tags []string, guidedMode bool) []string {
	if len(tags) == 0 && !guidedMode {
		return nil
	}

	normalized := make([]string, 0, len(tags)+1)
	seen := make(map[string]struct{})

	for _, tag := range tags {
		cleaned := strings.TrimSpace(tag)
		if cleaned == "" {
			continue
		}
		cleaned = strings.ToLower(cleaned)
		if _, ok := seen[cleaned]; ok {
			continue
		}
		seen[cleaned] = struct{}{}
		normalized = append(normalized, cleaned)
	}

	if guidedMode {
		if _, ok := seen[guidedModeTag]; !ok {
			normalized = append(normalized, guidedModeTag)
		}
	}

	if len(normalized) == 0 {
		return nil
	}

	return normalized
}

// Conversion helpers

func toRecipeResponse(r *domain.Recipe) *pb.RecipeResponse {
	resp := &pb.RecipeResponse{
		Id:          r.ID.String(),
		Name:        r.Name,
		Description: r.Description,
		PrepTime:    r.PrepTime,
		CookTime:    r.CookTime,
		Directions:  r.Directions,
	}

	if r.MainIngredient != nil {
		resp.MainIngredient = &pb.Ingredient{
			Id:       r.MainIngredient.ID.String(),
			Name:     r.MainIngredient.Name,
			Quantity: r.MainIngredient.Quantity,
			Unit:     r.MainIngredient.Unit,
		}
	}

	if r.Cuisine != nil {
		resp.Cuisine = &pb.Cuisine{
			Id:   r.Cuisine.ID.String(),
			Name: r.Cuisine.Name,
		}
	}

	resp.Ingredients = make([]*pb.Ingredient, len(r.Ingredients))
	for i, ing := range r.Ingredients {
		resp.Ingredients[i] = &pb.Ingredient{
			Id:       ing.ID.String(),
			Name:     ing.Name,
			Quantity: ing.Quantity,
			Unit:     ing.Unit,
		}
	}

	return resp
}

func toRecipesResponse(recipes []domain.Recipe) *pb.GetAllRecipesResponse {
	resp := &pb.GetAllRecipesResponse{
		Recipes: make([]*pb.RecipeResponse, len(recipes)),
	}

	for i := range recipes {
		resp.Recipes[i] = toRecipeResponse(&recipes[i])
	}

	return resp
}

func toRecipesResponseWithPagination(recipes []domain.Recipe, pageIndex, pageSize, totalCount, totalPages int32) *pb.GetAllRecipesResponse {
	resp := toRecipesResponse(recipes)
	resp.PageIndex = pageIndex
	resp.PageSize = pageSize
	resp.TotalCount = totalCount
	resp.TotalPages = totalPages
	return resp
}
