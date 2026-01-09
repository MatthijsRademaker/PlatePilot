package seed

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/google/uuid"

	"github.com/platepilot/backend/internal/common/auth"
	"github.com/platepilot/backend/internal/common/domain"
	"github.com/platepilot/backend/internal/common/vector"
	"github.com/platepilot/backend/internal/recipe/events"
	"github.com/platepilot/backend/internal/recipe/repository"
)

// Seeder handles database seeding from JSON files
type Seeder struct {
	repo      *repository.Repository
	vectorGen vector.Generator
	publisher *events.Publisher
	logger    *slog.Logger
}

// NewSeeder creates a new database seeder
func NewSeeder(repo *repository.Repository, vectorGen vector.Generator, publisher *events.Publisher, logger *slog.Logger) *Seeder {
	return &Seeder{
		repo:      repo,
		vectorGen: vectorGen,
		publisher: publisher,
		logger:    logger,
	}
}

// SeedFromFile seeds the database from a JSON file
func (s *Seeder) SeedFromFile(ctx context.Context, filePath string) error {
	seedUser, err := s.ensureSeedUser(ctx)
	if err != nil {
		return fmt.Errorf("ensure seed user: %w", err)
	}

	// Check if already seeded
	count, err := s.repo.Count(ctx, seedUser.ID, domain.RecipeFilter{})
	if err != nil {
		return fmt.Errorf("check recipe count: %w", err)
	}

	if count > 0 {
		s.logger.Info("database already seeded", "recipeCount", count)
		return nil
	}

	// Read JSON file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read seed file: %w", err)
	}

	var seedData SeedData
	if err := json.Unmarshal(data, &seedData); err != nil {
		return fmt.Errorf("unmarshal seed data: %w", err)
	}

	s.logger.Info("seeding database", "recipeCount", len(seedData.Recipes))

	// Track already created entities to avoid duplicates
	createdIngredients := make(map[string]*domain.Ingredient)
	createdCuisines := make(map[string]*domain.Cuisine)

	for _, recipeData := range seedData.Recipes {
		if err := s.seedRecipe(ctx, seedUser.ID, recipeData, createdIngredients, createdCuisines); err != nil {
			s.logger.Error("failed to seed recipe",
				"error", err,
				"recipeName", recipeData.Name,
			)
			// Continue with other recipes
			continue
		}
	}

	s.logger.Info("database seeding complete",
		"recipesCreated", len(seedData.Recipes),
		"ingredientsCreated", len(createdIngredients),
		"cuisinesCreated", len(createdCuisines),
	)

	return nil
}

func (s *Seeder) seedRecipe(
	ctx context.Context,
	userID uuid.UUID,
	data RecipeData,
	createdIngredients map[string]*domain.Ingredient,
	createdCuisines map[string]*domain.Cuisine,
) error {
	// Get or create cuisine
	cuisine, err := s.getOrCreateCuisine(ctx, userID, data.Cuisine, createdCuisines)
	if err != nil {
		return fmt.Errorf("get or create cuisine: %w", err)
	}

	// Get or create main ingredient
	mainIngredient, err := s.getOrCreateIngredient(ctx, userID, data.MainIngredient, createdIngredients)
	if err != nil {
		return fmt.Errorf("get or create main ingredient: %w", err)
	}

	// Get or create recipe ingredient lines
	var ingredientLines []domain.RecipeIngredientLine
	sortOrder := 1
	lineIngredientIDs := make(map[uuid.UUID]struct{})
	for _, ingData := range data.Ingredients {
		ingredient, err := s.getOrCreateIngredient(ctx, userID, ingData, createdIngredients)
		if err != nil {
			return fmt.Errorf("get or create ingredient %s: %w", ingData.Name, err)
		}
		ingredientLines = append(ingredientLines, domain.RecipeIngredientLine{
			Ingredient:   *ingredient,
			QuantityText: ingData.Quantity,
			SortOrder:    sortOrder,
		})
		lineIngredientIDs[ingredient.ID] = struct{}{}
		sortOrder++
	}

	if _, exists := lineIngredientIDs[mainIngredient.ID]; !exists && mainIngredient != nil {
		ingredientLines = append(ingredientLines, domain.RecipeIngredientLine{
			Ingredient:   *mainIngredient,
			QuantityText: data.MainIngredient.Quantity,
			SortOrder:    sortOrder,
		})
	}

	// Build recipe
	recipeID := data.ID
	if recipeID == uuid.Nil {
		recipeID = uuid.New()
	}

	recipe := &domain.Recipe{
		ID:             recipeID,
		UserID:         userID,
		Name:           data.Name,
		Description:    data.Description,
		PrepTimeMinutes:  parseMinutes(data.PrepTime),
		CookTimeMinutes:  parseMinutes(data.CookTime),
		TotalTimeMinutes: parseMinutes(data.PrepTime) + parseMinutes(data.CookTime),
		Servings:         1,
		MainIngredient:   mainIngredient,
		Cuisine:        cuisine,
		IngredientLines: ingredientLines,
		Steps:           buildSteps(data.Directions),
		Tags:            data.Metadata.Tags,
		ImageURL:        data.Metadata.ImageURL,
		Nutrition: domain.RecipeNutrition{
			CaloriesTotal: data.NutritionalInfo.Calories,
		},
	}

	// Generate vector embedding
	recipe.SearchVector = s.vectorGen.GenerateForRecipe(recipe)

	// Create recipe
	if err := s.repo.Create(ctx, recipe); err != nil {
		return fmt.Errorf("create recipe: %w", err)
	}

	s.logger.Debug("seeded recipe", "id", recipe.ID, "name", recipe.Name)

	// Optionally publish event
	if s.publisher != nil {
		if err := s.publisher.PublishRecipeUpserted(ctx, recipe); err != nil {
			s.logger.Warn("failed to publish recipe created event during seeding",
				"error", err,
				"recipeId", recipe.ID,
			)
			// Don't fail the seed operation
		}
	}

	return nil
}

func (s *Seeder) ensureSeedUser(ctx context.Context) (*domain.User, error) {
	email := os.Getenv("PLATEPILOT_SEED_USER_EMAIL")
	if email == "" {
		email = "seed@platepilot.local"
	}

	password := os.Getenv("PLATEPILOT_SEED_USER_PASSWORD")
	if password == "" {
		password = "platepilot"
	}

	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if !errors.Is(err, repository.ErrUserNotFound) {
			return nil, err
		}

		user = &domain.User{
			ID:          uuid.New(),
			Email:       email,
			DisplayName: "Seed User",
		}

		if err := s.repo.CreateUser(ctx, user); err != nil {
			return nil, err
		}
	}

	if _, err := s.repo.GetUserPasswordHash(ctx, user.ID); err != nil {
		if !errors.Is(err, repository.ErrUserNotFound) {
			return nil, err
		}
		hash, err := auth.HashPassword(password)
		if err != nil {
			return nil, err
		}
		if err := s.repo.CreateUserCredentials(ctx, user.ID, hash); err != nil {
			return nil, err
		}
		s.logger.Info("seed user credentials created", "email", email)
	}

	return user, nil
}

func buildSteps(directions []string) []domain.RecipeStep {
	steps := make([]domain.RecipeStep, 0, len(directions))
	for i, instruction := range directions {
		cleaned := strings.TrimSpace(instruction)
		if cleaned == "" {
			continue
		}
		steps = append(steps, domain.RecipeStep{
			StepIndex:   i + 1,
			Instruction: cleaned,
		})
	}
	return steps
}

func parseMinutes(input string) int {
	var digits []rune
	for _, r := range input {
		if unicode.IsDigit(r) {
			digits = append(digits, r)
		} else if len(digits) > 0 {
			break
		}
	}
	if len(digits) == 0 {
		return 0
	}
	value, err := strconv.Atoi(string(digits))
	if err != nil {
		return 0
	}
	return value
}

func (s *Seeder) getOrCreateCuisine(
	ctx context.Context,
	userID uuid.UUID,
	data CuisineData,
	cache map[string]*domain.Cuisine,
) (*domain.Cuisine, error) {
	// Check cache first
	if cuisine, ok := cache[data.Name]; ok {
		return cuisine, nil
	}

	// Try to get from database
	cuisine, err := s.repo.GetCuisineByName(ctx, userID, data.Name)
	if err == nil {
		cache[data.Name] = cuisine
		return cuisine, nil
	}

	// Create new cuisine
	cuisineID := data.ID
	if cuisineID == uuid.Nil {
		cuisineID = uuid.New()
	}

	cuisine = &domain.Cuisine{
		ID:     cuisineID,
		UserID: userID,
		Name:   data.Name,
	}

	if err := s.repo.CreateCuisine(ctx, cuisine); err != nil {
		return nil, err
	}

	cache[data.Name] = cuisine
	s.logger.Debug("created cuisine", "id", cuisine.ID, "name", cuisine.Name)

	return cuisine, nil
}

func (s *Seeder) getOrCreateIngredient(
	ctx context.Context,
	userID uuid.UUID,
	data IngredientData,
	cache map[string]*domain.Ingredient,
) (*domain.Ingredient, error) {
	// Check cache first
	if ingredient, ok := cache[data.Name]; ok {
		return ingredient, nil
	}

	// Try to get from database
	ingredient, err := s.repo.GetIngredientByName(ctx, userID, data.Name)
	if err == nil {
		cache[data.Name] = ingredient
		return ingredient, nil
	}

	// Create new ingredient
	ingredientID := data.ID
	if ingredientID == uuid.Nil {
		ingredientID = uuid.New()
	}

	ingredient = &domain.Ingredient{
		ID:          ingredientID,
		UserID:      userID,
		Name:        data.Name,
		Description: "",
	}

	if err := s.repo.CreateIngredient(ctx, ingredient); err != nil {
		return nil, err
	}

	cache[data.Name] = ingredient
	s.logger.Debug("created ingredient", "id", ingredient.ID, "name", ingredient.Name)

	return ingredient, nil
}

// JSON data structures

// SeedData is the root structure of the seed JSON
type SeedData struct {
	Recipes []RecipeData `json:"recipes"`
}

// RecipeData represents a recipe in the seed JSON
type RecipeData struct {
	ID              uuid.UUID        `json:"id"`
	Name            string           `json:"name"`
	Description     string           `json:"description"`
	PrepTime        string           `json:"prepTime"`
	CookTime        string           `json:"cookTime"`
	MainIngredient  IngredientData   `json:"mainIngredient"`
	Cuisine         CuisineData      `json:"cuisine"`
	Ingredients     []IngredientData `json:"ingredients"`
	Directions      []string         `json:"directions"`
	Metadata        MetadataData     `json:"metadata"`
	NutritionalInfo NutritionalData  `json:"nutritionalInfo"`
}

// IngredientData represents an ingredient in the seed JSON
type IngredientData struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Quantity string    `json:"quantity"`
}

// CuisineData represents a cuisine in the seed JSON
type CuisineData struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// MetadataData represents metadata in the seed JSON
type MetadataData struct {
	ImageURL string   `json:"imageUrl"`
	Tags     []string `json:"tags"`
}

// NutritionalData represents nutritional info in the seed JSON
type NutritionalData struct {
	Calories int `json:"Calories"`
}
