package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"
	"github.com/platepilot/backend/internal/common/domain"
)

var (
	ErrRecipeNotFound     = errors.New("recipe not found")
	ErrIngredientNotFound = errors.New("ingredient not found")
	ErrCuisineNotFound    = errors.New("cuisine not found")
	ErrAllergyNotFound    = errors.New("allergy not found")
)

// Repository provides access to the recipe write model
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a new repository
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// GetByID retrieves a recipe by ID with all related entities
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Recipe, error) {
	query := `
		SELECT
			r.id, r.name, r.description, r.prep_time, r.cook_time,
			r.directions, r.nutritional_info_calories,
			r.metadata_search_vector, r.metadata_image_url,
			r.metadata_tags, r.metadata_published_date,
			r.created_at, r.updated_at,
			c.id, c.name, c.created_at,
			mi.id, mi.name, mi.quantity, mi.created_at
		FROM recipes r
		JOIN cuisines c ON r.cuisine_id = c.id
		JOIN ingredients mi ON r.main_ingredient_id = mi.id
		WHERE r.id = $1
	`

	var recipe domain.Recipe
	var cuisine domain.Cuisine
	var mainIngredient domain.Ingredient
	var searchVector pgvector.Vector
	var imageURL *string
	var tags []string
	var publishedDate time.Time

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&recipe.ID, &recipe.Name, &recipe.Description,
		&recipe.PrepTime, &recipe.CookTime,
		&recipe.Directions, &recipe.NutritionalInfo.Calories,
		&searchVector, &imageURL, &tags, &publishedDate,
		&recipe.CreatedAt, &recipe.UpdatedAt,
		&cuisine.ID, &cuisine.Name, &cuisine.CreatedAt,
		&mainIngredient.ID, &mainIngredient.Name, &mainIngredient.Quantity, &mainIngredient.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRecipeNotFound
		}
		return nil, fmt.Errorf("query recipe: %w", err)
	}

	recipe.Cuisine = &cuisine
	recipe.MainIngredient = &mainIngredient
	recipe.Metadata.SearchVector = searchVector
	if imageURL != nil {
		recipe.Metadata.ImageURL = *imageURL
	}
	recipe.Metadata.Tags = tags
	recipe.Metadata.PublishedDate = publishedDate

	// Load recipe ingredients
	ingredients, err := r.getRecipeIngredients(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("load ingredients: %w", err)
	}
	recipe.Ingredients = ingredients

	return &recipe, nil
}

// GetAll retrieves all recipes with pagination
func (r *Repository) GetAll(ctx context.Context, limit, offset int) ([]domain.Recipe, error) {
	query := `
		SELECT
			r.id, r.name, r.description, r.prep_time, r.cook_time,
			r.directions, r.nutritional_info_calories,
			r.metadata_search_vector, r.metadata_image_url,
			r.metadata_tags, r.metadata_published_date,
			r.created_at, r.updated_at,
			c.id, c.name, c.created_at,
			mi.id, mi.name, mi.quantity, mi.created_at
		FROM recipes r
		JOIN cuisines c ON r.cuisine_id = c.id
		JOIN ingredients mi ON r.main_ingredient_id = mi.id
		ORDER BY r.created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query recipes: %w", err)
	}
	defer rows.Close()

	recipes, err := r.scanRecipes(ctx, rows)
	if err != nil {
		return nil, err
	}

	return recipes, nil
}

// Create creates a new recipe with all relationships
func (r *Repository) Create(ctx context.Context, recipe *domain.Recipe) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Ensure recipe has an ID
	if recipe.ID == uuid.Nil {
		recipe.ID = uuid.New()
	}

	// Insert recipe
	query := `
		INSERT INTO recipes (
			id, name, description, prep_time, cook_time,
			main_ingredient_id, cuisine_id, directions,
			nutritional_info_calories, metadata_search_vector,
			metadata_image_url, metadata_tags, metadata_published_date
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		)
	`

	var imageURL *string
	if recipe.Metadata.ImageURL != "" {
		imageURL = &recipe.Metadata.ImageURL
	}

	// Ensure tags is not nil (PostgreSQL requires non-null for TEXT[] with NOT NULL)
	tags := recipe.Metadata.Tags
	if tags == nil {
		tags = []string{}
	}

	_, err = tx.Exec(ctx, query,
		recipe.ID, recipe.Name, recipe.Description,
		recipe.PrepTime, recipe.CookTime,
		recipe.MainIngredient.ID, recipe.Cuisine.ID, recipe.Directions,
		recipe.NutritionalInfo.Calories, recipe.Metadata.SearchVector,
		imageURL, tags, recipe.Metadata.PublishedDate,
	)
	if err != nil {
		return fmt.Errorf("insert recipe: %w", err)
	}

	// Insert recipe-ingredient relationships
	for _, ingredient := range recipe.Ingredients {
		_, err = tx.Exec(ctx,
			`INSERT INTO recipe_ingredients (recipe_id, ingredient_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
			recipe.ID, ingredient.ID,
		)
		if err != nil {
			return fmt.Errorf("insert recipe ingredient: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// Update updates an existing recipe
func (r *Repository) Update(ctx context.Context, recipe *domain.Recipe) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		UPDATE recipes SET
			name = $2, description = $3, prep_time = $4, cook_time = $5,
			main_ingredient_id = $6, cuisine_id = $7, directions = $8,
			nutritional_info_calories = $9, metadata_search_vector = $10,
			metadata_image_url = $11, metadata_tags = $12, metadata_published_date = $13
		WHERE id = $1
	`

	var imageURL *string
	if recipe.Metadata.ImageURL != "" {
		imageURL = &recipe.Metadata.ImageURL
	}

	result, err := tx.Exec(ctx, query,
		recipe.ID, recipe.Name, recipe.Description,
		recipe.PrepTime, recipe.CookTime,
		recipe.MainIngredient.ID, recipe.Cuisine.ID, recipe.Directions,
		recipe.NutritionalInfo.Calories, recipe.Metadata.SearchVector,
		imageURL, recipe.Metadata.Tags, recipe.Metadata.PublishedDate,
	)
	if err != nil {
		return fmt.Errorf("update recipe: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrRecipeNotFound
	}

	// Replace recipe ingredients
	_, err = tx.Exec(ctx, `DELETE FROM recipe_ingredients WHERE recipe_id = $1`, recipe.ID)
	if err != nil {
		return fmt.Errorf("delete recipe ingredients: %w", err)
	}

	for _, ingredient := range recipe.Ingredients {
		_, err = tx.Exec(ctx,
			`INSERT INTO recipe_ingredients (recipe_id, ingredient_id) VALUES ($1, $2)`,
			recipe.ID, ingredient.ID,
		)
		if err != nil {
			return fmt.Errorf("insert recipe ingredient: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// Delete removes a recipe
func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.pool.Exec(ctx, `DELETE FROM recipes WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete recipe: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrRecipeNotFound
	}

	return nil
}

// Query operations

// GetSimilar retrieves recipes similar to a given recipe using vector similarity
func (r *Repository) GetSimilar(ctx context.Context, recipeID uuid.UUID, limit int) ([]domain.Recipe, error) {
	// First get the vector for the target recipe
	var targetVector pgvector.Vector
	err := r.pool.QueryRow(ctx,
		`SELECT metadata_search_vector FROM recipes WHERE id = $1`,
		recipeID,
	).Scan(&targetVector)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRecipeNotFound
		}
		return nil, fmt.Errorf("get target vector: %w", err)
	}

	query := `
		SELECT
			r.id, r.name, r.description, r.prep_time, r.cook_time,
			r.directions, r.nutritional_info_calories,
			r.metadata_search_vector, r.metadata_image_url,
			r.metadata_tags, r.metadata_published_date,
			r.created_at, r.updated_at,
			c.id, c.name, c.created_at,
			mi.id, mi.name, mi.quantity, mi.created_at
		FROM recipes r
		JOIN cuisines c ON r.cuisine_id = c.id
		JOIN ingredients mi ON r.main_ingredient_id = mi.id
		WHERE r.id != $1
		ORDER BY r.metadata_search_vector <=> $2
		LIMIT $3
	`

	rows, err := r.pool.Query(ctx, query, recipeID, targetVector, limit)
	if err != nil {
		return nil, fmt.Errorf("query similar recipes: %w", err)
	}
	defer rows.Close()

	return r.scanRecipes(ctx, rows)
}

// GetByCuisine retrieves recipes by cuisine ID
func (r *Repository) GetByCuisine(ctx context.Context, cuisineID uuid.UUID, limit, offset int) ([]domain.Recipe, error) {
	query := `
		SELECT
			r.id, r.name, r.description, r.prep_time, r.cook_time,
			r.directions, r.nutritional_info_calories,
			r.metadata_search_vector, r.metadata_image_url,
			r.metadata_tags, r.metadata_published_date,
			r.created_at, r.updated_at,
			c.id, c.name, c.created_at,
			mi.id, mi.name, mi.quantity, mi.created_at
		FROM recipes r
		JOIN cuisines c ON r.cuisine_id = c.id
		JOIN ingredients mi ON r.main_ingredient_id = mi.id
		WHERE r.cuisine_id = $1
		ORDER BY r.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, cuisineID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query recipes by cuisine: %w", err)
	}
	defer rows.Close()

	return r.scanRecipes(ctx, rows)
}

// GetByIngredient retrieves recipes containing a specific ingredient (main or in list)
func (r *Repository) GetByIngredient(ctx context.Context, ingredientID uuid.UUID, limit, offset int) ([]domain.Recipe, error) {
	query := `
		SELECT DISTINCT
			r.id, r.name, r.description, r.prep_time, r.cook_time,
			r.directions, r.nutritional_info_calories,
			r.metadata_search_vector, r.metadata_image_url,
			r.metadata_tags, r.metadata_published_date,
			r.created_at, r.updated_at,
			c.id, c.name, c.created_at,
			mi.id, mi.name, mi.quantity, mi.created_at
		FROM recipes r
		JOIN cuisines c ON r.cuisine_id = c.id
		JOIN ingredients mi ON r.main_ingredient_id = mi.id
		LEFT JOIN recipe_ingredients ri ON r.id = ri.recipe_id
		WHERE r.main_ingredient_id = $1 OR ri.ingredient_id = $1
		ORDER BY r.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, ingredientID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query recipes by ingredient: %w", err)
	}
	defer rows.Close()

	return r.scanRecipes(ctx, rows)
}

// GetExcludingAllergy retrieves recipes that don't contain ingredients with a specific allergy
func (r *Repository) GetExcludingAllergy(ctx context.Context, allergyID uuid.UUID, limit, offset int) ([]domain.Recipe, error) {
	// Find recipes that have NO ingredients with the given allergy
	query := `
		SELECT
			r.id, r.name, r.description, r.prep_time, r.cook_time,
			r.directions, r.nutritional_info_calories,
			r.metadata_search_vector, r.metadata_image_url,
			r.metadata_tags, r.metadata_published_date,
			r.created_at, r.updated_at,
			c.id, c.name, c.created_at,
			mi.id, mi.name, mi.quantity, mi.created_at
		FROM recipes r
		JOIN cuisines c ON r.cuisine_id = c.id
		JOIN ingredients mi ON r.main_ingredient_id = mi.id
		WHERE NOT EXISTS (
			-- Check main ingredient for allergy
			SELECT 1 FROM ingredient_allergies ia
			WHERE ia.ingredient_id = r.main_ingredient_id AND ia.allergy_id = $1
		)
		AND NOT EXISTS (
			-- Check recipe ingredients for allergy
			SELECT 1 FROM recipe_ingredients ri
			JOIN ingredient_allergies ia ON ri.ingredient_id = ia.ingredient_id
			WHERE ri.recipe_id = r.id AND ia.allergy_id = $1
		)
		ORDER BY r.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, allergyID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query recipes excluding allergy: %w", err)
	}
	defer rows.Close()

	return r.scanRecipes(ctx, rows)
}

// GetByAllergy retrieves recipes that contain ingredients with a specific allergy
// (useful for finding what TO avoid)
func (r *Repository) GetByAllergy(ctx context.Context, allergyID uuid.UUID, limit, offset int) ([]domain.Recipe, error) {
	query := `
		SELECT DISTINCT
			r.id, r.name, r.description, r.prep_time, r.cook_time,
			r.directions, r.nutritional_info_calories,
			r.metadata_search_vector, r.metadata_image_url,
			r.metadata_tags, r.metadata_published_date,
			r.created_at, r.updated_at,
			c.id, c.name, c.created_at,
			mi.id, mi.name, mi.quantity, mi.created_at
		FROM recipes r
		JOIN cuisines c ON r.cuisine_id = c.id
		JOIN ingredients mi ON r.main_ingredient_id = mi.id
		LEFT JOIN recipe_ingredients ri ON r.id = ri.recipe_id
		LEFT JOIN ingredient_allergies ia ON (ri.ingredient_id = ia.ingredient_id OR r.main_ingredient_id = ia.ingredient_id)
		WHERE ia.allergy_id = $1
		ORDER BY r.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, allergyID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query recipes by allergy: %w", err)
	}
	defer rows.Close()

	return r.scanRecipes(ctx, rows)
}

// Count returns the total number of recipes
func (r *Repository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM recipes`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count recipes: %w", err)
	}
	return count, nil
}

// Ingredient operations

// GetIngredientByID retrieves an ingredient by ID
func (r *Repository) GetIngredientByID(ctx context.Context, id uuid.UUID) (*domain.Ingredient, error) {
	query := `SELECT id, name, quantity, created_at FROM ingredients WHERE id = $1`

	var ingredient domain.Ingredient
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&ingredient.ID, &ingredient.Name, &ingredient.Quantity, &ingredient.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrIngredientNotFound
		}
		return nil, fmt.Errorf("query ingredient: %w", err)
	}

	// Load allergies for this ingredient
	allergies, err := r.getIngredientAllergies(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("load allergies: %w", err)
	}
	ingredient.Allergies = allergies

	return &ingredient, nil
}

// GetIngredientByName retrieves an ingredient by name
func (r *Repository) GetIngredientByName(ctx context.Context, name string) (*domain.Ingredient, error) {
	query := `SELECT id, name, quantity, created_at FROM ingredients WHERE name = $1`

	var ingredient domain.Ingredient
	err := r.pool.QueryRow(ctx, query, name).Scan(
		&ingredient.ID, &ingredient.Name, &ingredient.Quantity, &ingredient.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrIngredientNotFound
		}
		return nil, fmt.Errorf("query ingredient: %w", err)
	}

	return &ingredient, nil
}

// CreateIngredient creates a new ingredient
func (r *Repository) CreateIngredient(ctx context.Context, ingredient *domain.Ingredient) error {
	if ingredient.ID == uuid.Nil {
		ingredient.ID = uuid.New()
	}

	query := `INSERT INTO ingredients (id, name, quantity) VALUES ($1, $2, $3)`
	_, err := r.pool.Exec(ctx, query, ingredient.ID, ingredient.Name, ingredient.Quantity)
	if err != nil {
		return fmt.Errorf("insert ingredient: %w", err)
	}

	return nil
}

// GetOrCreateIngredient gets an existing ingredient by name or creates it
func (r *Repository) GetOrCreateIngredient(ctx context.Context, name string, quantity string) (*domain.Ingredient, error) {
	ingredient, err := r.GetIngredientByName(ctx, name)
	if err == nil {
		return ingredient, nil
	}
	if !errors.Is(err, ErrIngredientNotFound) {
		return nil, err
	}

	// Create new ingredient
	ingredient = &domain.Ingredient{
		ID:       uuid.New(),
		Name:     name,
		Quantity: quantity,
	}
	if err := r.CreateIngredient(ctx, ingredient); err != nil {
		return nil, err
	}

	return ingredient, nil
}

// Cuisine operations

// GetCuisineByID retrieves a cuisine by ID
func (r *Repository) GetCuisineByID(ctx context.Context, id uuid.UUID) (*domain.Cuisine, error) {
	query := `SELECT id, name, created_at FROM cuisines WHERE id = $1`

	var cuisine domain.Cuisine
	err := r.pool.QueryRow(ctx, query, id).Scan(&cuisine.ID, &cuisine.Name, &cuisine.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCuisineNotFound
		}
		return nil, fmt.Errorf("query cuisine: %w", err)
	}

	return &cuisine, nil
}

// GetCuisineByName retrieves a cuisine by name
func (r *Repository) GetCuisineByName(ctx context.Context, name string) (*domain.Cuisine, error) {
	query := `SELECT id, name, created_at FROM cuisines WHERE name = $1`

	var cuisine domain.Cuisine
	err := r.pool.QueryRow(ctx, query, name).Scan(&cuisine.ID, &cuisine.Name, &cuisine.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCuisineNotFound
		}
		return nil, fmt.Errorf("query cuisine: %w", err)
	}

	return &cuisine, nil
}

// CreateCuisine creates a new cuisine
func (r *Repository) CreateCuisine(ctx context.Context, cuisine *domain.Cuisine) error {
	if cuisine.ID == uuid.Nil {
		cuisine.ID = uuid.New()
	}

	query := `INSERT INTO cuisines (id, name) VALUES ($1, $2)`
	_, err := r.pool.Exec(ctx, query, cuisine.ID, cuisine.Name)
	if err != nil {
		return fmt.Errorf("insert cuisine: %w", err)
	}

	return nil
}

// GetOrCreateCuisine gets an existing cuisine by name or creates it
func (r *Repository) GetOrCreateCuisine(ctx context.Context, name string) (*domain.Cuisine, error) {
	cuisine, err := r.GetCuisineByName(ctx, name)
	if err == nil {
		return cuisine, nil
	}
	if !errors.Is(err, ErrCuisineNotFound) {
		return nil, err
	}

	// Create new cuisine
	cuisine = &domain.Cuisine{
		ID:   uuid.New(),
		Name: name,
	}
	if err := r.CreateCuisine(ctx, cuisine); err != nil {
		return nil, err
	}

	return cuisine, nil
}

// GetAllCuisines retrieves all cuisines
func (r *Repository) GetAllCuisines(ctx context.Context) ([]domain.Cuisine, error) {
	query := `SELECT id, name, created_at FROM cuisines ORDER BY name`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query cuisines: %w", err)
	}
	defer rows.Close()

	var cuisines []domain.Cuisine
	for rows.Next() {
		var cuisine domain.Cuisine
		if err := rows.Scan(&cuisine.ID, &cuisine.Name, &cuisine.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan cuisine: %w", err)
		}
		cuisines = append(cuisines, cuisine)
	}

	return cuisines, nil
}

// Allergy operations

// GetAllergyByID retrieves an allergy by ID
func (r *Repository) GetAllergyByID(ctx context.Context, id uuid.UUID) (*domain.Allergy, error) {
	query := `SELECT id, name, created_at FROM allergies WHERE id = $1`

	var allergy domain.Allergy
	err := r.pool.QueryRow(ctx, query, id).Scan(&allergy.ID, &allergy.Name, &allergy.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAllergyNotFound
		}
		return nil, fmt.Errorf("query allergy: %w", err)
	}

	return &allergy, nil
}

// GetAllergyByName retrieves an allergy by name
func (r *Repository) GetAllergyByName(ctx context.Context, name string) (*domain.Allergy, error) {
	query := `SELECT id, name, created_at FROM allergies WHERE name = $1`

	var allergy domain.Allergy
	err := r.pool.QueryRow(ctx, query, name).Scan(&allergy.ID, &allergy.Name, &allergy.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAllergyNotFound
		}
		return nil, fmt.Errorf("query allergy: %w", err)
	}

	return &allergy, nil
}

// CreateAllergy creates a new allergy
func (r *Repository) CreateAllergy(ctx context.Context, allergy *domain.Allergy) error {
	if allergy.ID == uuid.Nil {
		allergy.ID = uuid.New()
	}

	query := `INSERT INTO allergies (id, name) VALUES ($1, $2)`
	_, err := r.pool.Exec(ctx, query, allergy.ID, allergy.Name)
	if err != nil {
		return fmt.Errorf("insert allergy: %w", err)
	}

	return nil
}

// GetOrCreateAllergy gets an existing allergy by name or creates it
func (r *Repository) GetOrCreateAllergy(ctx context.Context, name string) (*domain.Allergy, error) {
	allergy, err := r.GetAllergyByName(ctx, name)
	if err == nil {
		return allergy, nil
	}
	if !errors.Is(err, ErrAllergyNotFound) {
		return nil, err
	}

	// Create new allergy
	allergy = &domain.Allergy{
		ID:   uuid.New(),
		Name: name,
	}
	if err := r.CreateAllergy(ctx, allergy); err != nil {
		return nil, err
	}

	return allergy, nil
}

// AddIngredientAllergy links an allergy to an ingredient
func (r *Repository) AddIngredientAllergy(ctx context.Context, ingredientID, allergyID uuid.UUID) error {
	query := `INSERT INTO ingredient_allergies (ingredient_id, allergy_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.pool.Exec(ctx, query, ingredientID, allergyID)
	if err != nil {
		return fmt.Errorf("add ingredient allergy: %w", err)
	}
	return nil
}

// Helper methods

func (r *Repository) getRecipeIngredients(ctx context.Context, recipeID uuid.UUID) ([]domain.Ingredient, error) {
	query := `
		SELECT i.id, i.name, i.quantity, i.created_at
		FROM ingredients i
		JOIN recipe_ingredients ri ON i.id = ri.ingredient_id
		WHERE ri.recipe_id = $1
	`

	rows, err := r.pool.Query(ctx, query, recipeID)
	if err != nil {
		return nil, fmt.Errorf("query recipe ingredients: %w", err)
	}
	defer rows.Close()

	var ingredients []domain.Ingredient
	for rows.Next() {
		var ingredient domain.Ingredient
		if err := rows.Scan(&ingredient.ID, &ingredient.Name, &ingredient.Quantity, &ingredient.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan ingredient: %w", err)
		}
		ingredients = append(ingredients, ingredient)
	}

	return ingredients, nil
}

func (r *Repository) getIngredientAllergies(ctx context.Context, ingredientID uuid.UUID) ([]domain.Allergy, error) {
	query := `
		SELECT a.id, a.name, a.created_at
		FROM allergies a
		JOIN ingredient_allergies ia ON a.id = ia.allergy_id
		WHERE ia.ingredient_id = $1
	`

	rows, err := r.pool.Query(ctx, query, ingredientID)
	if err != nil {
		return nil, fmt.Errorf("query ingredient allergies: %w", err)
	}
	defer rows.Close()

	var allergies []domain.Allergy
	for rows.Next() {
		var allergy domain.Allergy
		if err := rows.Scan(&allergy.ID, &allergy.Name, &allergy.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan allergy: %w", err)
		}
		allergies = append(allergies, allergy)
	}

	return allergies, nil
}

func (r *Repository) scanRecipes(ctx context.Context, rows pgx.Rows) ([]domain.Recipe, error) {
	var recipes []domain.Recipe

	for rows.Next() {
		var recipe domain.Recipe
		var cuisine domain.Cuisine
		var mainIngredient domain.Ingredient
		var searchVector pgvector.Vector
		var imageURL *string
		var tags []string
		var publishedDate time.Time

		err := rows.Scan(
			&recipe.ID, &recipe.Name, &recipe.Description,
			&recipe.PrepTime, &recipe.CookTime,
			&recipe.Directions, &recipe.NutritionalInfo.Calories,
			&searchVector, &imageURL, &tags, &publishedDate,
			&recipe.CreatedAt, &recipe.UpdatedAt,
			&cuisine.ID, &cuisine.Name, &cuisine.CreatedAt,
			&mainIngredient.ID, &mainIngredient.Name, &mainIngredient.Quantity, &mainIngredient.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan recipe: %w", err)
		}

		recipe.Cuisine = &cuisine
		recipe.MainIngredient = &mainIngredient
		recipe.Metadata.SearchVector = searchVector
		if imageURL != nil {
			recipe.Metadata.ImageURL = *imageURL
		}
		recipe.Metadata.Tags = tags
		recipe.Metadata.PublishedDate = publishedDate

		recipes = append(recipes, recipe)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate recipes: %w", err)
	}

	// Load ingredients for each recipe (N+1 queries for now - can optimize later)
	for i := range recipes {
		ingredients, err := r.getRecipeIngredients(ctx, recipes[i].ID)
		if err != nil {
			return nil, err
		}
		recipes[i].Ingredients = ingredients
	}

	return recipes, nil
}
