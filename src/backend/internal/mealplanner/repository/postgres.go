package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"
)

// IngredientLine represents a recipe ingredient line in the read model.
type IngredientLine struct {
	RecipeID       uuid.UUID
	IngredientID   uuid.UUID
	IngredientName string
	QuantityValue  *float64
	QuantityText   string
	Unit           string
	IsOptional     bool
	Note           string
	SortOrder      int
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
func (r *Repository) GetByID(ctx context.Context, userID, id uuid.UUID) (*Recipe, error) {
	query := `
		SELECT
			id, user_id, name, description,
			prep_time_minutes, cook_time_minutes, total_time_minutes,
			servings, yield_quantity, yield_unit,
			search_vector, cuisine_id, cuisine_name,
			main_ingredient_id, main_ingredient_name,
			ingredient_ids, allergy_ids, tags, image_url,
			calories_total, calories_per_serving,
			protein_g, carbs_g, fat_g, fiber_g, sugar_g, sodium_mg
		FROM recipes
		WHERE id = $1 AND user_id = $2
	`

	var recipe Recipe
	var yieldQuantity pgtype.Numeric
	var yieldUnit *string
	var protein pgtype.Numeric
	var carbs pgtype.Numeric
	var fat pgtype.Numeric
	var fiber pgtype.Numeric
	var sugar pgtype.Numeric
	var sodium pgtype.Numeric
	err := r.pool.QueryRow(ctx, query, id, userID).Scan(
		&recipe.ID, &recipe.UserID, &recipe.Name, &recipe.Description,
		&recipe.PrepTimeMinutes, &recipe.CookTimeMinutes, &recipe.TotalTimeMinutes,
		&recipe.Servings, &yieldQuantity, &yieldUnit,
		&recipe.SearchVector, &recipe.CuisineID, &recipe.CuisineName,
		&recipe.MainIngredientID, &recipe.MainIngredientName,
		&recipe.IngredientIDs, &recipe.AllergyIDs, &recipe.Tags, &recipe.ImageURL,
		&recipe.CaloriesTotal, &recipe.CaloriesPerServing,
		&protein, &carbs, &fat, &fiber, &sugar, &sodium,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("recipe not found: %s", id)
		}
		return nil, fmt.Errorf("query recipe: %w", err)
	}

	if yieldUnit != nil {
		recipe.YieldUnit = *yieldUnit
	}
	yieldPtr, err := numericToFloatPtr(yieldQuantity)
	if err != nil {
		return nil, fmt.Errorf("parse yield quantity: %w", err)
	}
	recipe.YieldQuantity = yieldPtr
	recipe.ProteinG = numericToFloat(protein)
	recipe.CarbsG = numericToFloat(carbs)
	recipe.FatG = numericToFloat(fat)
	recipe.FiberG = numericToFloat(fiber)
	recipe.SugarG = numericToFloat(sugar)
	recipe.SodiumMg = numericToFloat(sodium)

	return &recipe, nil
}

// GetAll retrieves all recipes with pagination
func (r *Repository) GetAll(ctx context.Context, userID uuid.UUID, limit, offset int) ([]Recipe, error) {
	query := `
		SELECT
			id, user_id, name, description,
			prep_time_minutes, cook_time_minutes, total_time_minutes,
			servings, yield_quantity, yield_unit,
			search_vector, cuisine_id, cuisine_name,
			main_ingredient_id, main_ingredient_name,
			ingredient_ids, allergy_ids, tags, image_url,
			calories_total, calories_per_serving,
			protein_g, carbs_g, fat_g, fiber_g, sugar_g, sodium_mg
		FROM recipes
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query recipes: %w", err)
	}
	defer rows.Close()

	return scanRecipes(rows)
}

// GetByCuisine retrieves recipes by cuisine ID
func (r *Repository) GetByCuisine(ctx context.Context, userID, cuisineID uuid.UUID, limit, offset int) ([]Recipe, error) {
	query := `
		SELECT
			id, user_id, name, description,
			prep_time_minutes, cook_time_minutes, total_time_minutes,
			servings, yield_quantity, yield_unit,
			search_vector, cuisine_id, cuisine_name,
			main_ingredient_id, main_ingredient_name,
			ingredient_ids, allergy_ids, tags, image_url,
			calories_total, calories_per_serving,
			protein_g, carbs_g, fat_g, fiber_g, sugar_g, sodium_mg
		FROM recipes
		WHERE cuisine_id = $1 AND user_id = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.pool.Query(ctx, query, cuisineID, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query recipes by cuisine: %w", err)
	}
	defer rows.Close()

	return scanRecipes(rows)
}

// GetByIngredient retrieves recipes containing a specific ingredient
func (r *Repository) GetByIngredient(ctx context.Context, userID, ingredientID uuid.UUID, limit, offset int) ([]Recipe, error) {
	query := `
		SELECT
			id, user_id, name, description,
			prep_time_minutes, cook_time_minutes, total_time_minutes,
			servings, yield_quantity, yield_unit,
			search_vector, cuisine_id, cuisine_name,
			main_ingredient_id, main_ingredient_name,
			ingredient_ids, allergy_ids, tags, image_url,
			calories_total, calories_per_serving,
			protein_g, carbs_g, fat_g, fiber_g, sugar_g, sodium_mg
		FROM recipes
		WHERE (main_ingredient_id = $1 OR $1 = ANY(ingredient_ids))
		  AND user_id = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.pool.Query(ctx, query, ingredientID, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query recipes by ingredient: %w", err)
	}
	defer rows.Close()

	return scanRecipes(rows)
}

// GetExcludingAllergy retrieves recipes that don't contain a specific allergy
func (r *Repository) GetExcludingAllergy(ctx context.Context, userID, allergyID uuid.UUID, limit, offset int) ([]Recipe, error) {
	query := `
		SELECT
			id, user_id, name, description,
			prep_time_minutes, cook_time_minutes, total_time_minutes,
			servings, yield_quantity, yield_unit,
			search_vector, cuisine_id, cuisine_name,
			main_ingredient_id, main_ingredient_name,
			ingredient_ids, allergy_ids, tags, image_url,
			calories_total, calories_per_serving,
			protein_g, carbs_g, fat_g, fiber_g, sugar_g, sodium_mg
		FROM recipes
		WHERE NOT ($1 = ANY(allergy_ids))
		  AND user_id = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.pool.Query(ctx, query, allergyID, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query recipes excluding allergy: %w", err)
	}
	defer rows.Close()

	return scanRecipes(rows)
}

// GetSimilar retrieves recipes similar to a given recipe using vector similarity
func (r *Repository) GetSimilar(ctx context.Context, userID, recipeID uuid.UUID, limit int) ([]Recipe, error) {
	// First get the vector for the target recipe
	var targetVector pgvector.Vector
	err := r.pool.QueryRow(ctx,
		`SELECT search_vector FROM recipes WHERE id = $1 AND user_id = $2`,
		recipeID, userID,
	).Scan(&targetVector)
	if err != nil {
		return nil, fmt.Errorf("get target vector: %w", err)
	}

	query := `
		SELECT
			id, user_id, name, description,
			prep_time_minutes, cook_time_minutes, total_time_minutes,
			servings, yield_quantity, yield_unit,
			search_vector, cuisine_id, cuisine_name,
			main_ingredient_id, main_ingredient_name,
			ingredient_ids, allergy_ids, tags, image_url,
			calories_total, calories_per_serving,
			protein_g, carbs_g, fat_g, fiber_g, sugar_g, sodium_mg
		FROM recipes
		WHERE id != $1 AND user_id = $2
		ORDER BY search_vector <=> $3
		LIMIT $4
	`

	rows, err := r.pool.Query(ctx, query, recipeID, userID, targetVector, limit)
	if err != nil {
		return nil, fmt.Errorf("query similar recipes: %w", err)
	}
	defer rows.Close()

	return scanRecipes(rows)
}

// GetVectorByID retrieves just the search vector for a recipe
func (r *Repository) GetVectorByID(ctx context.Context, userID, id uuid.UUID) (pgvector.Vector, error) {
	var vector pgvector.Vector
	err := r.pool.QueryRow(ctx,
		`SELECT search_vector FROM recipes WHERE id = $1 AND user_id = $2`,
		id, userID,
	).Scan(&vector)
	if err != nil {
		return pgvector.Vector{}, fmt.Errorf("get vector: %w", err)
	}
	return vector, nil
}

// Upsert inserts or updates a recipe in the read model along with ingredient lines.
func (r *Repository) Upsert(ctx context.Context, recipe *Recipe, lines []IngredientLine) error {
	tags := recipe.Tags
	if tags == nil {
		tags = []string{}
	}

	query := `
		INSERT INTO recipes (
			id, user_id, name, description,
			prep_time_minutes, cook_time_minutes, total_time_minutes,
			servings, yield_quantity, yield_unit,
			search_vector, cuisine_id, cuisine_name,
			main_ingredient_id, main_ingredient_name,
			ingredient_ids, allergy_ids, tags, image_url,
			calories_total, calories_per_serving,
			protein_g, carbs_g, fat_g, fiber_g, sugar_g, sodium_mg
		) VALUES (
			$1, $2, $3, $4,
			$5, $6, $7,
			$8, $9, $10,
			$11, $12, $13,
			$14, $15,
			$16, $17,
			$18, $19,
			$20, $21, $22, $23, $24, $25, $26, $27
		)
		ON CONFLICT (id) DO UPDATE SET
			user_id = EXCLUDED.user_id,
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			prep_time_minutes = EXCLUDED.prep_time_minutes,
			cook_time_minutes = EXCLUDED.cook_time_minutes,
			total_time_minutes = EXCLUDED.total_time_minutes,
			servings = EXCLUDED.servings,
			yield_quantity = EXCLUDED.yield_quantity,
			yield_unit = EXCLUDED.yield_unit,
			search_vector = EXCLUDED.search_vector,
			cuisine_id = EXCLUDED.cuisine_id,
			cuisine_name = EXCLUDED.cuisine_name,
			main_ingredient_id = EXCLUDED.main_ingredient_id,
			main_ingredient_name = EXCLUDED.main_ingredient_name,
			ingredient_ids = EXCLUDED.ingredient_ids,
			allergy_ids = EXCLUDED.allergy_ids,
			tags = EXCLUDED.tags,
			image_url = EXCLUDED.image_url,
			calories_total = EXCLUDED.calories_total,
			calories_per_serving = EXCLUDED.calories_per_serving,
			protein_g = EXCLUDED.protein_g,
			carbs_g = EXCLUDED.carbs_g,
			fat_g = EXCLUDED.fat_g,
			fiber_g = EXCLUDED.fiber_g,
			sugar_g = EXCLUDED.sugar_g,
			sodium_mg = EXCLUDED.sodium_mg,
			updated_at = NOW()
	`

	_, err := r.pool.Exec(ctx, query,
		recipe.ID, recipe.UserID, recipe.Name, recipe.Description,
		recipe.PrepTimeMinutes, recipe.CookTimeMinutes, recipe.TotalTimeMinutes,
		recipe.Servings, recipe.YieldQuantity, recipe.YieldUnit,
		recipe.SearchVector, recipe.CuisineID, recipe.CuisineName,
		recipe.MainIngredientID, recipe.MainIngredientName,
		recipe.IngredientIDs, recipe.AllergyIDs, tags, recipe.ImageURL,
		recipe.CaloriesTotal, recipe.CaloriesPerServing,
		recipe.ProteinG, recipe.CarbsG, recipe.FatG, recipe.FiberG, recipe.SugarG, recipe.SodiumMg,
	)
	if err != nil {
		return fmt.Errorf("upsert recipe: %w", err)
	}

	_, err = r.pool.Exec(ctx, `DELETE FROM recipe_ingredient_lines WHERE recipe_id = $1`, recipe.ID)
	if err != nil {
		return fmt.Errorf("delete recipe ingredient lines: %w", err)
	}

	for _, line := range lines {
		_, err = r.pool.Exec(ctx, `
			INSERT INTO recipe_ingredient_lines (
				recipe_id, ingredient_id, ingredient_name,
				quantity_value, quantity_text, unit,
				is_optional, note, sort_order
			) VALUES (
				$1, $2, $3,
				$4, $5, $6,
				$7, $8, $9
			)
		`, line.RecipeID, line.IngredientID, line.IngredientName,
			line.QuantityValue, line.QuantityText, line.Unit,
			line.IsOptional, line.Note, line.SortOrder,
		)
		if err != nil {
			return fmt.Errorf("insert recipe ingredient line: %w", err)
		}
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
		var yieldQuantity pgtype.Numeric
		var yieldUnit *string
		var protein pgtype.Numeric
		var carbs pgtype.Numeric
		var fat pgtype.Numeric
		var fiber pgtype.Numeric
		var sugar pgtype.Numeric
		var sodium pgtype.Numeric
		err := rows.Scan(
			&recipe.ID, &recipe.UserID, &recipe.Name, &recipe.Description,
			&recipe.PrepTimeMinutes, &recipe.CookTimeMinutes, &recipe.TotalTimeMinutes,
			&recipe.Servings, &yieldQuantity, &yieldUnit,
			&recipe.SearchVector, &recipe.CuisineID, &recipe.CuisineName,
			&recipe.MainIngredientID, &recipe.MainIngredientName,
			&recipe.IngredientIDs, &recipe.AllergyIDs, &recipe.Tags, &recipe.ImageURL,
			&recipe.CaloriesTotal, &recipe.CaloriesPerServing,
			&protein, &carbs, &fat, &fiber, &sugar, &sodium,
		)
		if err != nil {
			return nil, fmt.Errorf("scan recipe: %w", err)
		}
		if yieldUnit != nil {
			recipe.YieldUnit = *yieldUnit
		}
		yieldPtr, err := numericToFloatPtr(yieldQuantity)
		if err != nil {
			return nil, fmt.Errorf("parse yield quantity: %w", err)
		}
		recipe.YieldQuantity = yieldPtr
		recipe.ProteinG = numericToFloat(protein)
		recipe.CarbsG = numericToFloat(carbs)
		recipe.FatG = numericToFloat(fat)
		recipe.FiberG = numericToFloat(fiber)
		recipe.SugarG = numericToFloat(sugar)
		recipe.SodiumMg = numericToFloat(sodium)
		recipes = append(recipes, recipe)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate recipes: %w", err)
	}

	return recipes, nil
}

func numericToFloat(n pgtype.Numeric) float64 {
	if !n.Valid {
		return 0
	}
	value, err := n.Float64Value()
	if err != nil {
		return 0
	}
	if !value.Valid {
		return 0
	}
	return value.Float64
}

func numericToFloatPtr(n pgtype.Numeric) (*float64, error) {
	if !n.Valid {
		return nil, nil
	}
	value, err := n.Float64Value()
	if err != nil {
		return nil, err
	}
	if !value.Valid {
		return nil, nil
	}
	converted := value.Float64
	return &converted, nil
}
