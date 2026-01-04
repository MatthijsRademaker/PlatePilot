package handler

import (
	"context"
	"errors"
	"log/slog"
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
	h.logger.Debug("get recipe by id", "recipeId", req.GetRecipeId())

	id, err := uuid.Parse(req.GetRecipeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid recipe ID: %v", err)
	}

	recipe, err := h.repo.GetByID(ctx, id)
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

	recipes, err := h.repo.GetAll(ctx, pageSize, offset)
	if err != nil {
		h.logger.Error("failed to get recipes", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get recipes")
	}

	totalCount, err := h.repo.Count(ctx)
	if err != nil {
		h.logger.Error("failed to count recipes", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to count recipes")
	}

	totalPages := int32((totalCount + int64(pageSize) - 1) / int64(pageSize))

	return toRecipesResponseWithPagination(recipes, int32(pageIndex), int32(pageSize), int32(totalCount), totalPages), nil
}

// CreateRecipe creates a new recipe
func (h *GRPCHandler) CreateRecipe(ctx context.Context, req *pb.CreateRecipeRequest) (*pb.RecipeResponse, error) {
	h.logger.Debug("create recipe", "name", req.GetName())

	// Validate required fields
	if req.GetName() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "name is required")
	}

	// Parse main ingredient ID
	mainIngredientID, err := uuid.Parse(req.GetMainIngredientId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid main ingredient ID: %v", err)
	}

	// Parse cuisine ID
	cuisineID, err := uuid.Parse(req.GetCuisineId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid cuisine ID: %v", err)
	}

	// Get main ingredient
	mainIngredient, err := h.repo.GetIngredientByID(ctx, mainIngredientID)
	if err != nil {
		if errors.Is(err, repository.ErrIngredientNotFound) {
			return nil, status.Errorf(codes.NotFound, "main ingredient not found")
		}
		h.logger.Error("failed to get main ingredient", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get main ingredient")
	}

	// Get cuisine
	cuisine, err := h.repo.GetCuisineByID(ctx, cuisineID)
	if err != nil {
		if errors.Is(err, repository.ErrCuisineNotFound) {
			return nil, status.Errorf(codes.NotFound, "cuisine not found")
		}
		h.logger.Error("failed to get cuisine", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get cuisine")
	}

	// Parse and get ingredients
	var ingredients []domain.Ingredient
	for _, idStr := range req.GetIngredientIds() {
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
		ingredients = append(ingredients, *ingredient)
	}

	// Build the recipe
	recipe := &domain.Recipe{
		ID:             uuid.New(),
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
	h.logger.Debug("get similar recipes", "recipeId", req.GetRecipeId(), "amount", req.GetAmount())

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

	recipes, err := h.repo.GetSimilar(ctx, recipeID, amount)
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
	h.logger.Debug("get recipes by cuisine", "cuisineId", req.GetCuisineId())

	cuisineID, err := uuid.Parse(req.GetCuisineId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid cuisine ID: %v", err)
	}

	// Using default pagination for now
	recipes, err := h.repo.GetByCuisine(ctx, cuisineID, 100, 0)
	if err != nil {
		h.logger.Error("failed to get recipes by cuisine", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get recipes by cuisine")
	}

	return toRecipesResponse(recipes), nil
}

// GetRecipesByIngredient retrieves recipes containing a specific ingredient
func (h *GRPCHandler) GetRecipesByIngredient(ctx context.Context, req *pb.GetRecipesByIngredientRequest) (*pb.GetAllRecipesResponse, error) {
	h.logger.Debug("get recipes by ingredient", "ingredientId", req.GetIngredientId())

	ingredientID, err := uuid.Parse(req.GetIngredientId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid ingredient ID: %v", err)
	}

	// Using default pagination for now
	recipes, err := h.repo.GetByIngredient(ctx, ingredientID, 100, 0)
	if err != nil {
		h.logger.Error("failed to get recipes by ingredient", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get recipes by ingredient")
	}

	return toRecipesResponse(recipes), nil
}

// GetRecipesByAllergy retrieves recipes excluding a specific allergy
func (h *GRPCHandler) GetRecipesByAllergy(ctx context.Context, req *pb.GetRecipesByAllergyRequest) (*pb.GetAllRecipesResponse, error) {
	h.logger.Debug("get recipes excluding allergy", "allergyId", req.GetAllergyId())

	allergyID, err := uuid.Parse(req.GetAllergyId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid allergy ID: %v", err)
	}

	// Using default pagination for now - returns recipes WITHOUT this allergy
	recipes, err := h.repo.GetExcludingAllergy(ctx, allergyID, 100, 0)
	if err != nil {
		h.logger.Error("failed to get recipes excluding allergy", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get recipes excluding allergy")
	}

	return toRecipesResponse(recipes), nil
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
			Id:   r.MainIngredient.ID.String(),
			Name: r.MainIngredient.Name,
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
			Id:   ing.ID.String(),
			Name: ing.Name,
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
