package handler

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/platepilot/backend/internal/common/domain"
	"github.com/platepilot/backend/internal/common/vector"
	pb "github.com/platepilot/backend/internal/recipe/pb"
	"github.com/platepilot/backend/internal/recipe/repository"
)

// GRPCHandler implements the RecipeService gRPC interface.
type GRPCHandler struct {
	pb.UnimplementedRecipeServiceServer
	repo      RecipeRepository
	vectorGen vector.Generator
	publisher EventPublisher
	logger    *slog.Logger
}

const (
	defaultCuisineName = "General"
)

// NewGRPCHandler creates a new gRPC handler.
func NewGRPCHandler(repo RecipeRepository, vectorGen vector.Generator, publisher EventPublisher, logger *slog.Logger) *GRPCHandler {
	return &GRPCHandler{
		repo:      repo,
		vectorGen: vectorGen,
		publisher: publisher,
		logger:    logger,
	}
}

// GetRecipe retrieves a recipe by ID.
func (h *GRPCHandler) GetRecipe(ctx context.Context, req *pb.GetRecipeRequest) (*pb.Recipe, error) {
	h.logger.Debug("get recipe", "recipeId", req.GetRecipeId(), "userId", req.GetUserId())

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

// ListRecipes retrieves recipes with pagination and optional filters.
func (h *GRPCHandler) ListRecipes(ctx context.Context, req *pb.ListRecipesRequest) (*pb.ListRecipesResponse, error) {
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

	filter, err := buildRecipeFilter(req)
	if err != nil {
		return nil, err
	}

	recipes, err := h.repo.List(ctx, userID, filter, pageSize, offset)
	if err != nil {
		h.logger.Error("failed to list recipes", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to list recipes")
	}

	totalCount, err := h.repo.Count(ctx, userID, filter)
	if err != nil {
		h.logger.Error("failed to count recipes", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to count recipes")
	}

	totalPages := int32((totalCount + int64(pageSize) - 1) / int64(pageSize))
	return toRecipesResponseWithPagination(recipes, int32(pageIndex), int32(pageSize), int32(totalCount), totalPages), nil
}

// CreateRecipe creates a new recipe.
func (h *GRPCHandler) CreateRecipe(ctx context.Context, req *pb.CreateRecipeRequest) (*pb.Recipe, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	recipe, err := h.buildRecipeFromInput(ctx, userID, req.GetRecipe())
	if err != nil {
		return nil, err
	}

	recipe.ID = uuid.New()
	recipe.SearchVector = h.vectorGen.GenerateForRecipe(recipe)

	if err := h.repo.Create(ctx, recipe); err != nil {
		h.logger.Error("failed to create recipe", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to create recipe")
	}

	h.logger.Info("recipe created", "recipeId", recipe.ID, "name", recipe.Name)

	if h.publisher != nil {
		if err := h.publisher.PublishRecipeUpserted(ctx, recipe); err != nil {
			h.logger.Error("failed to publish recipe upserted event",
				"error", err,
				"recipeId", recipe.ID,
			)
		}
	}

	return toRecipeResponse(recipe), nil
}

// UpdateRecipe updates an existing recipe.
func (h *GRPCHandler) UpdateRecipe(ctx context.Context, req *pb.UpdateRecipeRequest) (*pb.Recipe, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	recipeID, err := uuid.Parse(req.GetRecipeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid recipe ID: %v", err)
	}

	recipe, err := h.buildRecipeFromInput(ctx, userID, req.GetRecipe())
	if err != nil {
		return nil, err
	}
	recipe.ID = recipeID
	recipe.UserID = userID
	recipe.SearchVector = h.vectorGen.GenerateForRecipe(recipe)

	if err := h.repo.Update(ctx, recipe); err != nil {
		if errors.Is(err, repository.ErrRecipeNotFound) {
			return nil, status.Errorf(codes.NotFound, "recipe not found")
		}
		h.logger.Error("failed to update recipe", "error", err, "recipeId", recipeID)
		return nil, status.Errorf(codes.Internal, "failed to update recipe")
	}

	if h.publisher != nil {
		if err := h.publisher.PublishRecipeUpserted(ctx, recipe); err != nil {
			h.logger.Error("failed to publish recipe upserted event",
				"error", err,
				"recipeId", recipe.ID,
			)
		}
	}

	return toRecipeResponse(recipe), nil
}

// DeleteRecipe deletes a recipe.
func (h *GRPCHandler) DeleteRecipe(ctx context.Context, req *pb.DeleteRecipeRequest) (*emptypb.Empty, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	recipeID, err := uuid.Parse(req.GetRecipeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid recipe ID: %v", err)
	}

	if err := h.repo.Delete(ctx, userID, recipeID); err != nil {
		if errors.Is(err, repository.ErrRecipeNotFound) {
			return nil, status.Errorf(codes.NotFound, "recipe not found")
		}
		h.logger.Error("failed to delete recipe", "error", err, "recipeId", recipeID)
		return nil, status.Errorf(codes.Internal, "failed to delete recipe")
	}

	if h.publisher != nil {
		if err := h.publisher.PublishRecipeDeleted(ctx, recipeID, userID); err != nil {
			h.logger.Error("failed to publish recipe deleted event",
				"error", err,
				"recipeId", recipeID,
			)
		}
	}

	return &emptypb.Empty{}, nil
}

// GetSimilarRecipes retrieves similar recipes using vector similarity.
func (h *GRPCHandler) GetSimilarRecipes(ctx context.Context, req *pb.GetSimilarRecipesRequest) (*pb.ListRecipesResponse, error) {
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

// GetCuisines retrieves available cuisines.
func (h *GRPCHandler) GetCuisines(ctx context.Context, req *pb.GetCuisinesRequest) (*pb.GetCuisinesResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	cuisines, err := h.repo.GetCuisines(ctx, userID)
	if err != nil {
		h.logger.Error("failed to get cuisines", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get cuisines")
	}

	resp := &pb.GetCuisinesResponse{Cuisines: make([]*pb.Cuisine, len(cuisines))}
	for i, cuisine := range cuisines {
		resp.Cuisines[i] = &pb.Cuisine{
			Id:   cuisine.ID.String(),
			Name: cuisine.Name,
		}
	}

	return resp, nil
}

// CreateCuisine creates a new cuisine.
func (h *GRPCHandler) CreateCuisine(ctx context.Context, req *pb.CreateCuisineRequest) (*pb.Cuisine, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	name := strings.TrimSpace(req.GetName())
	if name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "name is required")
	}

	cuisine, err := h.repo.GetOrCreateCuisine(ctx, userID, name)
	if err != nil {
		h.logger.Error("failed to create cuisine", "error", err, "cuisineName", name)
		return nil, status.Errorf(codes.Internal, "failed to create cuisine")
	}

	return &pb.Cuisine{
		Id:   cuisine.ID.String(),
		Name: cuisine.Name,
	}, nil
}

func (h *GRPCHandler) buildRecipeFromInput(ctx context.Context, userID uuid.UUID, input *pb.RecipeInput) (*domain.Recipe, error) {
	if input == nil {
		return nil, status.Errorf(codes.InvalidArgument, "recipe is required")
	}

	name := strings.TrimSpace(input.GetName())
	if name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "name is required")
	}

	ingredientLines, err := h.resolveIngredientLines(ctx, userID, input.GetIngredientLines())
	if err != nil {
		return nil, err
	}

	mainIngredient, err := h.resolveMainIngredient(ctx, userID, input, ingredientLines)
	if err != nil {
		return nil, err
	}

	cuisine, err := h.resolveCuisine(ctx, userID, input)
	if err != nil {
		return nil, err
	}

	steps, err := h.resolveSteps(input.GetSteps())
	if err != nil {
		return nil, err
	}

	servings := int(input.GetServings())
	if servings < 1 {
		servings = 1
	}

	prepMinutes := int(input.GetPrepTimeMinutes())
	if prepMinutes < 0 {
		prepMinutes = 0
	}

	cookMinutes := int(input.GetCookTimeMinutes())
	if cookMinutes < 0 {
		cookMinutes = 0
	}

	var yieldQuantity *float64
	if input.GetYieldQuantity() != nil {
		value := input.GetYieldQuantity().GetValue()
		yieldQuantity = &value
	}

	recipe := &domain.Recipe{
		UserID:           userID,
		Name:             name,
		Description:      strings.TrimSpace(input.GetDescription()),
		PrepTimeMinutes:  prepMinutes,
		CookTimeMinutes:  cookMinutes,
		TotalTimeMinutes: prepMinutes + cookMinutes,
		Servings:         servings,
		YieldQuantity:    yieldQuantity,
		YieldUnit:        strings.TrimSpace(input.GetYieldUnit()),
		MainIngredient:   mainIngredient,
		Cuisine:          cuisine,
		IngredientLines:  ingredientLines,
		Steps:            steps,
		Tags:             normalizeTags(input.GetTags()),
		ImageURL:         strings.TrimSpace(input.GetImageUrl()),
		Nutrition:        nutritionFromProto(input.GetNutrition()),
	}

	return recipe, nil
}

func (h *GRPCHandler) resolveIngredientLines(ctx context.Context, userID uuid.UUID, inputs []*pb.IngredientLineInput) ([]domain.RecipeIngredientLine, error) {
	lines := make([]domain.RecipeIngredientLine, 0, len(inputs))

	for index, input := range inputs {
		if input == nil {
			continue
		}

		var ingredient *domain.Ingredient
		if idStr := strings.TrimSpace(input.GetIngredientId()); idStr != "" {
			ingredientID, err := uuid.Parse(idStr)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "invalid ingredient ID: %v", err)
			}
			ingredient, err = h.repo.GetIngredientByID(ctx, userID, ingredientID)
			if err != nil {
				if errors.Is(err, repository.ErrIngredientNotFound) {
					return nil, status.Errorf(codes.NotFound, "ingredient not found: %s", idStr)
				}
				h.logger.Error("failed to get ingredient", "error", err, "ingredientId", idStr)
				return nil, status.Errorf(codes.Internal, "failed to get ingredient")
			}
		} else if name := strings.TrimSpace(input.GetIngredientName()); name != "" {
			var err error
			ingredient, err = h.repo.GetOrCreateIngredient(ctx, userID, name)
			if err != nil {
				h.logger.Error("failed to get or create ingredient", "error", err, "ingredientName", name)
				return nil, status.Errorf(codes.Internal, "failed to create ingredient")
			}
		} else {
			return nil, status.Errorf(codes.InvalidArgument, "ingredient line is missing id or name")
		}

		var quantityValue *float64
		if input.GetQuantityValue() != nil {
			value := input.GetQuantityValue().GetValue()
			quantityValue = &value
		}

		sortOrder := int(input.GetSortOrder())
		if sortOrder == 0 {
			sortOrder = index + 1
		}

		lines = append(lines, domain.RecipeIngredientLine{
			Ingredient:    *ingredient,
			QuantityValue: quantityValue,
			QuantityText:  strings.TrimSpace(input.GetQuantityText()),
			Unit:          strings.TrimSpace(input.GetUnit()),
			IsOptional:    input.GetIsOptional(),
			Note:          strings.TrimSpace(input.GetNote()),
			SortOrder:     sortOrder,
		})
	}

	return lines, nil
}

func (h *GRPCHandler) resolveMainIngredient(
	ctx context.Context,
	userID uuid.UUID,
	input *pb.RecipeInput,
	lines []domain.RecipeIngredientLine,
) (*domain.Ingredient, error) {
	if idStr := strings.TrimSpace(input.GetMainIngredientId()); idStr != "" {
		ingredientID, err := uuid.Parse(idStr)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid main ingredient ID: %v", err)
		}
		ingredient, err := h.repo.GetIngredientByID(ctx, userID, ingredientID)
		if err != nil {
			if errors.Is(err, repository.ErrIngredientNotFound) {
				return nil, status.Errorf(codes.NotFound, "main ingredient not found")
			}
			h.logger.Error("failed to get main ingredient", "error", err)
			return nil, status.Errorf(codes.Internal, "failed to get main ingredient")
		}
		return ingredient, nil
	}

	if name := strings.TrimSpace(input.GetMainIngredientName()); name != "" {
		ingredient, err := h.repo.GetOrCreateIngredient(ctx, userID, name)
		if err != nil {
			h.logger.Error("failed to get or create main ingredient", "error", err, "ingredientName", name)
			return nil, status.Errorf(codes.Internal, "failed to create main ingredient")
		}
		return ingredient, nil
	}

	if len(lines) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "main ingredient is required")
	}

	return &lines[0].Ingredient, nil
}

func (h *GRPCHandler) resolveCuisine(ctx context.Context, userID uuid.UUID, input *pb.RecipeInput) (*domain.Cuisine, error) {
	cuisineID := strings.TrimSpace(input.GetCuisineId())
	if cuisineID != "" {
		parsedID, err := uuid.Parse(cuisineID)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid cuisine ID: %v", err)
		}
		cuisine, err := h.repo.GetCuisineByID(ctx, userID, parsedID)
		if err != nil {
			if errors.Is(err, repository.ErrCuisineNotFound) {
				return nil, status.Errorf(codes.NotFound, "cuisine not found")
			}
			h.logger.Error("failed to get cuisine", "error", err)
			return nil, status.Errorf(codes.Internal, "failed to get cuisine")
		}
		return cuisine, nil
	}

	cuisineName := strings.TrimSpace(input.GetCuisineName())
	if cuisineName == "" {
		cuisineName = defaultCuisineName
	}
	cuisine, err := h.repo.GetOrCreateCuisine(ctx, userID, cuisineName)
	if err != nil {
		h.logger.Error("failed to get or create cuisine", "error", err, "cuisineName", cuisineName)
		return nil, status.Errorf(codes.Internal, "failed to create cuisine")
	}
	return cuisine, nil
}

func (h *GRPCHandler) resolveSteps(inputs []*pb.RecipeStepInput) ([]domain.RecipeStep, error) {
	steps := make([]domain.RecipeStep, 0, len(inputs))

	for index, input := range inputs {
		if input == nil {
			continue
		}

		stepIndex := int(input.GetStepIndex())
		if stepIndex == 0 {
			stepIndex = index + 1
		}

		var duration *int
		if input.GetDurationSeconds() != nil {
			value := int(input.GetDurationSeconds().GetValue())
			duration = &value
		}

		var temperature *float64
		if input.GetTemperatureValue() != nil {
			value := input.GetTemperatureValue().GetValue()
			temperature = &value
		}

		steps = append(steps, domain.RecipeStep{
			StepIndex:        stepIndex,
			Instruction:      strings.TrimSpace(input.GetInstruction()),
			DurationSeconds:  duration,
			TemperatureValue: temperature,
			TemperatureUnit:  strings.TrimSpace(input.GetTemperatureUnit()),
			MediaURL:         strings.TrimSpace(input.GetMediaUrl()),
		})
	}

	return steps, nil
}

func buildRecipeFilter(req *pb.ListRecipesRequest) (domain.RecipeFilter, error) {
	filter := domain.RecipeFilter{
		Tags: normalizeTags(req.GetTags()),
	}

	if value := strings.TrimSpace(req.GetCuisineId()); value != "" {
		id, err := uuid.Parse(value)
		if err != nil {
			return filter, status.Errorf(codes.InvalidArgument, "invalid cuisine ID: %v", err)
		}
		filter.CuisineID = &id
	}

	if value := strings.TrimSpace(req.GetIngredientId()); value != "" {
		id, err := uuid.Parse(value)
		if err != nil {
			return filter, status.Errorf(codes.InvalidArgument, "invalid ingredient ID: %v", err)
		}
		filter.IngredientID = &id
	}

	if value := strings.TrimSpace(req.GetAllergyId()); value != "" {
		id, err := uuid.Parse(value)
		if err != nil {
			return filter, status.Errorf(codes.InvalidArgument, "invalid allergy ID: %v", err)
		}
		filter.AllergyID = &id
	}

	return filter, nil
}

func nutritionFromProto(input *pb.RecipeNutrition) domain.RecipeNutrition {
	if input == nil {
		return domain.RecipeNutrition{}
	}
	return domain.RecipeNutrition{
		CaloriesTotal:      int(input.GetCaloriesTotal()),
		CaloriesPerServing: int(input.GetCaloriesPerServing()),
		ProteinG:           input.GetProteinG(),
		CarbsG:             input.GetCarbsG(),
		FatG:               input.GetFatG(),
		FiberG:             input.GetFiberG(),
		SugarG:             input.GetSugarG(),
		SodiumMg:           input.GetSodiumMg(),
	}
}

func normalizeTags(tags []string) []string {
	if len(tags) == 0 {
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

	if len(normalized) == 0 {
		return nil
	}

	return normalized
}

// Conversion helpers

func toRecipeResponse(r *domain.Recipe) *pb.Recipe {
	if r == nil {
		return nil
	}

	resp := &pb.Recipe{
		Id:               r.ID.String(),
		UserId:           r.UserID.String(),
		Name:             r.Name,
		Description:      r.Description,
		PrepTimeMinutes:  int32(r.PrepTimeMinutes),
		CookTimeMinutes:  int32(r.CookTimeMinutes),
		TotalTimeMinutes: int32(r.TotalTimeMinutes),
		Servings:         int32(r.Servings),
		YieldUnit:        r.YieldUnit,
		Tags:             r.Tags,
		ImageUrl:         r.ImageURL,
		Nutrition: &pb.RecipeNutrition{
			CaloriesTotal:      int32(r.Nutrition.CaloriesTotal),
			CaloriesPerServing: int32(r.Nutrition.CaloriesPerServing),
			ProteinG:           r.Nutrition.ProteinG,
			CarbsG:             r.Nutrition.CarbsG,
			FatG:               r.Nutrition.FatG,
			FiberG:             r.Nutrition.FiberG,
			SugarG:             r.Nutrition.SugarG,
			SodiumMg:           r.Nutrition.SodiumMg,
		},
	}

	if r.YieldQuantity != nil {
		resp.YieldQuantity = wrapperspb.Double(*r.YieldQuantity)
	}

	if r.MainIngredient != nil {
		resp.MainIngredient = &pb.IngredientRef{
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

	lines := make([]*pb.IngredientLine, len(r.IngredientLines))
	for i, line := range r.IngredientLines {
		lineResp := &pb.IngredientLine{
			Ingredient: &pb.IngredientRef{
				Id:   line.Ingredient.ID.String(),
				Name: line.Ingredient.Name,
			},
			QuantityText: line.QuantityText,
			Unit:         line.Unit,
			IsOptional:   line.IsOptional,
			Note:         line.Note,
			SortOrder:    int32(line.SortOrder),
		}
		if line.QuantityValue != nil {
			lineResp.QuantityValue = wrapperspb.Double(*line.QuantityValue)
		}
		lines[i] = lineResp
	}
	resp.IngredientLines = lines

	steps := make([]*pb.RecipeStep, len(r.Steps))
	for i, step := range r.Steps {
		stepResp := &pb.RecipeStep{
			StepIndex:       int32(step.StepIndex),
			Instruction:     step.Instruction,
			TemperatureUnit: step.TemperatureUnit,
			MediaUrl:        step.MediaURL,
		}
		if step.DurationSeconds != nil {
			stepResp.DurationSeconds = wrapperspb.Int32(int32(*step.DurationSeconds))
		}
		if step.TemperatureValue != nil {
			stepResp.TemperatureValue = wrapperspb.Double(*step.TemperatureValue)
		}
		steps[i] = stepResp
	}
	resp.Steps = steps

	return resp
}

func toRecipesResponse(recipes []domain.Recipe) *pb.ListRecipesResponse {
	resp := &pb.ListRecipesResponse{
		Recipes: make([]*pb.Recipe, len(recipes)),
	}

	for i := range recipes {
		resp.Recipes[i] = toRecipeResponse(&recipes[i])
	}

	return resp
}

func toRecipesResponseWithPagination(recipes []domain.Recipe, pageIndex, pageSize, totalCount, totalPages int32) *pb.ListRecipesResponse {
	resp := toRecipesResponse(recipes)
	resp.PageIndex = pageIndex
	resp.PageSize = pageSize
	resp.TotalCount = totalCount
	resp.TotalPages = totalPages
	return resp
}
