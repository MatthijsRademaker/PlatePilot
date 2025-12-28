package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"
)

// Recipe represents a recipe in the read model
type Recipe struct {
	ID                 uuid.UUID
	Name               string
	Description        string
	PrepTime           string
	CookTime           string
	SearchVector       pgvector.Vector
	CuisineID          uuid.UUID
	CuisineName        string
	MainIngredientID   uuid.UUID
	MainIngredientName string
	IngredientIDs      []uuid.UUID
	AllergyIDs         []uuid.UUID
	Directions         []string
	ImageURL           string
	Tags               []string
	Calories           int
}

// Repository provides access to the mealplanner read model
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a new repository
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// GetByID retrieves a recipe by its ID
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Recipe, error) {
	query := `
		SELECT
			id, name, description, prep_time, cook_time,
			search_vector, cuisine_id, cuisine_name,
			main_ingredient_id, main_ingredient_name,
			ingredient_ids, allergy_ids, directions,
			image_url, tags, calories
		FROM recipes
		WHERE id = $1
	`

	var recipe Recipe
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&recipe.ID, &recipe.Name, &recipe.Description,
		&recipe.PrepTime, &recipe.CookTime,
		&recipe.SearchVector, &recipe.CuisineID, &recipe.CuisineName,
		&recipe.MainIngredientID, &recipe.MainIngredientName,
		&recipe.IngredientIDs, &recipe.AllergyIDs, &recipe.Directions,
		&recipe.ImageURL, &recipe.Tags, &recipe.Calories,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("recipe not found: %s", id)
		}
		return nil, fmt.Errorf("query recipe: %w", err)
	}

	return &recipe, nil
}

// GetAll retrieves all recipes with pagination
func (r *Repository) GetAll(ctx context.Context, limit, offset int) ([]Recipe, error) {
	query := `
		SELECT
			id, name, description, prep_time, cook_time,
			search_vector, cuisine_id, cuisine_name,
			main_ingredient_id, main_ingredient_name,
			ingredient_ids, allergy_ids, directions,
			image_url, tags, calories
		FROM recipes
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query recipes: %w", err)
	}
	defer rows.Close()

	return scanRecipes(rows)
}

// GetByCuisine retrieves recipes by cuisine ID
func (r *Repository) GetByCuisine(ctx context.Context, cuisineID uuid.UUID, limit, offset int) ([]Recipe, error) {
	query := `
		SELECT
			id, name, description, prep_time, cook_time,
			search_vector, cuisine_id, cuisine_name,
			main_ingredient_id, main_ingredient_name,
			ingredient_ids, allergy_ids, directions,
			image_url, tags, calories
		FROM recipes
		WHERE cuisine_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, cuisineID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query recipes by cuisine: %w", err)
	}
	defer rows.Close()

	return scanRecipes(rows)
}

// GetByIngredient retrieves recipes containing a specific ingredient
func (r *Repository) GetByIngredient(ctx context.Context, ingredientID uuid.UUID, limit, offset int) ([]Recipe, error) {
	query := `
		SELECT
			id, name, description, prep_time, cook_time,
			search_vector, cuisine_id, cuisine_name,
			main_ingredient_id, main_ingredient_name,
			ingredient_ids, allergy_ids, directions,
			image_url, tags, calories
		FROM recipes
		WHERE main_ingredient_id = $1 OR $1 = ANY(ingredient_ids)
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, ingredientID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query recipes by ingredient: %w", err)
	}
	defer rows.Close()

	return scanRecipes(rows)
}

// GetExcludingAllergy retrieves recipes that don't contain a specific allergy
func (r *Repository) GetExcludingAllergy(ctx context.Context, allergyID uuid.UUID, limit, offset int) ([]Recipe, error) {
	query := `
		SELECT
			id, name, description, prep_time, cook_time,
			search_vector, cuisine_id, cuisine_name,
			main_ingredient_id, main_ingredient_name,
			ingredient_ids, allergy_ids, directions,
			image_url, tags, calories
		FROM recipes
		WHERE NOT ($1 = ANY(allergy_ids))
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, allergyID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query recipes excluding allergy: %w", err)
	}
	defer rows.Close()

	return scanRecipes(rows)
}

// GetSimilar retrieves recipes similar to a given recipe using vector similarity
func (r *Repository) GetSimilar(ctx context.Context, recipeID uuid.UUID, limit int) ([]Recipe, error) {
	// First get the vector for the target recipe
	var targetVector pgvector.Vector
	err := r.pool.QueryRow(ctx,
		`SELECT search_vector FROM recipes WHERE id = $1`,
		recipeID,
	).Scan(&targetVector)
	if err != nil {
		return nil, fmt.Errorf("get target vector: %w", err)
	}

	query := `
		SELECT
			id, name, description, prep_time, cook_time,
			search_vector, cuisine_id, cuisine_name,
			main_ingredient_id, main_ingredient_name,
			ingredient_ids, allergy_ids, directions,
			image_url, tags, calories
		FROM recipes
		WHERE id != $1
		ORDER BY search_vector <=> $2
		LIMIT $3
	`

	rows, err := r.pool.Query(ctx, query, recipeID, targetVector, limit)
	if err != nil {
		return nil, fmt.Errorf("query similar recipes: %w", err)
	}
	defer rows.Close()

	return scanRecipes(rows)
}

// GetVectorByID retrieves just the search vector for a recipe
func (r *Repository) GetVectorByID(ctx context.Context, id uuid.UUID) (pgvector.Vector, error) {
	var vector pgvector.Vector
	err := r.pool.QueryRow(ctx,
		`SELECT search_vector FROM recipes WHERE id = $1`,
		id,
	).Scan(&vector)
	if err != nil {
		return pgvector.Vector{}, fmt.Errorf("get vector: %w", err)
	}
	return vector, nil
}

// Upsert inserts or updates a recipe in the read model
func (r *Repository) Upsert(ctx context.Context, recipe *Recipe) error {
	query := `
		INSERT INTO recipes (
			id, name, description, prep_time, cook_time,
			search_vector, cuisine_id, cuisine_name,
			main_ingredient_id, main_ingredient_name,
			ingredient_ids, allergy_ids, directions,
			image_url, tags, calories
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
		)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			prep_time = EXCLUDED.prep_time,
			cook_time = EXCLUDED.cook_time,
			search_vector = EXCLUDED.search_vector,
			cuisine_id = EXCLUDED.cuisine_id,
			cuisine_name = EXCLUDED.cuisine_name,
			main_ingredient_id = EXCLUDED.main_ingredient_id,
			main_ingredient_name = EXCLUDED.main_ingredient_name,
			ingredient_ids = EXCLUDED.ingredient_ids,
			allergy_ids = EXCLUDED.allergy_ids,
			directions = EXCLUDED.directions,
			image_url = EXCLUDED.image_url,
			tags = EXCLUDED.tags,
			calories = EXCLUDED.calories,
			updated_at = NOW()
	`

	_, err := r.pool.Exec(ctx, query,
		recipe.ID, recipe.Name, recipe.Description,
		recipe.PrepTime, recipe.CookTime,
		recipe.SearchVector, recipe.CuisineID, recipe.CuisineName,
		recipe.MainIngredientID, recipe.MainIngredientName,
		recipe.IngredientIDs, recipe.AllergyIDs, recipe.Directions,
		recipe.ImageURL, recipe.Tags, recipe.Calories,
	)
	if err != nil {
		return fmt.Errorf("upsert recipe: %w", err)
	}

	return nil
}

// Delete removes a recipe from the read model
func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM recipes WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete recipe: %w", err)
	}
	return nil
}

func scanRecipes(rows pgx.Rows) ([]Recipe, error) {
	var recipes []Recipe
	for rows.Next() {
		var recipe Recipe
		err := rows.Scan(
			&recipe.ID, &recipe.Name, &recipe.Description,
			&recipe.PrepTime, &recipe.CookTime,
			&recipe.SearchVector, &recipe.CuisineID, &recipe.CuisineName,
			&recipe.MainIngredientID, &recipe.MainIngredientName,
			&recipe.IngredientIDs, &recipe.AllergyIDs, &recipe.Directions,
			&recipe.ImageURL, &recipe.Tags, &recipe.Calories,
		)
		if err != nil {
			return nil, fmt.Errorf("scan recipe: %w", err)
		}
		recipes = append(recipes, recipe)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate recipes: %w", err)
	}

	return recipes, nil
}
