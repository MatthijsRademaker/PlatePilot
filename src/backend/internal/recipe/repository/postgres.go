package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pgvector/pgvector-go"
	"github.com/platepilot/backend/internal/common/domain"
)

var (
	ErrRecipeNotFound     = errors.New("recipe not found")
	ErrIngredientNotFound = errors.New("ingredient not found")
	ErrCuisineNotFound    = errors.New("cuisine not found")
	ErrAllergyNotFound    = errors.New("allergy not found")
	ErrUserNotFound       = errors.New("user not found")
)

// Repository provides access to the recipe write model
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a new repository
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func accessClause(alias string, userParam int) string {
	return fmt.Sprintf("(%s.user_id = $%d OR EXISTS (SELECT 1 FROM recipe_shares rs WHERE rs.recipe_id = %s.id AND rs.shared_with_user_id = $%d))", alias, userParam, alias, userParam)
}

func activeClause(alias string) string {
	return fmt.Sprintf("%s.deleted_at IS NULL", alias)
}

// GetByID retrieves a recipe by ID with all related entities
func (r *Repository) GetByID(ctx context.Context, userID, id uuid.UUID) (*domain.Recipe, error) {
	query := `
		SELECT
			r.id, r.user_id, r.name, r.description,
			r.prep_time_minutes, r.cook_time_minutes, r.total_time_minutes,
			r.servings, r.yield_quantity, r.yield_unit,
			r.image_url, r.tags, r.search_vector,
			r.created_at, r.updated_at, r.deleted_at,
			c.id, c.user_id, c.name, c.created_at,
			mi.id, mi.user_id, mi.name, mi.description, mi.created_at, mi.updated_at,
			rn.calories_total, rn.calories_per_serving,
			rn.protein_g, rn.carbs_g, rn.fat_g, rn.fiber_g, rn.sugar_g, rn.sodium_mg
		FROM recipes r
		JOIN cuisines c ON r.cuisine_id = c.id
		JOIN ingredients mi ON r.main_ingredient_id = mi.id
		LEFT JOIN recipe_nutrition rn ON rn.recipe_id = r.id
		WHERE r.id = $1
		  AND ` + activeClause("r") + `
		  AND ` + accessClause("r", 2) + `
	`

	var recipe domain.Recipe
	var cuisine domain.Cuisine
	var mainIngredient domain.Ingredient
	var searchVector pgvector.Vector
	var imageURL *string
	var tags []string
	var yieldUnit *string
	var yieldQuantity pgtype.Numeric
	var deletedAt *time.Time
	var protein pgtype.Numeric
	var carbs pgtype.Numeric
	var fat pgtype.Numeric
	var fiber pgtype.Numeric
	var sugar pgtype.Numeric
	var sodium pgtype.Numeric
	var caloriesTotal int
	var caloriesPerServing int

	err := r.pool.QueryRow(ctx, query, id, userID).Scan(
		&recipe.ID, &recipe.UserID, &recipe.Name, &recipe.Description,
		&recipe.PrepTimeMinutes, &recipe.CookTimeMinutes, &recipe.TotalTimeMinutes,
		&recipe.Servings, &yieldQuantity, &yieldUnit,
		&imageURL, &tags, &searchVector,
		&recipe.CreatedAt, &recipe.UpdatedAt, &deletedAt,
		&cuisine.ID, &cuisine.UserID, &cuisine.Name, &cuisine.CreatedAt,
		&mainIngredient.ID, &mainIngredient.UserID, &mainIngredient.Name, &mainIngredient.Description, &mainIngredient.CreatedAt, &mainIngredient.UpdatedAt,
		&caloriesTotal, &caloriesPerServing,
		&protein, &carbs, &fat, &fiber, &sugar, &sodium,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRecipeNotFound
		}
		return nil, fmt.Errorf("query recipe: %w", err)
	}

	recipe.Cuisine = &cuisine
	recipe.MainIngredient = &mainIngredient
	recipe.SearchVector = searchVector
	if imageURL != nil {
		recipe.ImageURL = *imageURL
	}
	if yieldUnit != nil {
		recipe.YieldUnit = *yieldUnit
	}
	yieldPtr, err := numericToFloatPtr(yieldQuantity)
	if err != nil {
		return nil, fmt.Errorf("parse yield quantity: %w", err)
	}
	recipe.YieldQuantity = yieldPtr
	recipe.Tags = tags
	recipe.DeletedAt = deletedAt
	recipe.Nutrition = domain.RecipeNutrition{
		CaloriesTotal:      caloriesTotal,
		CaloriesPerServing: caloriesPerServing,
		ProteinG:           numericToFloat(protein),
		CarbsG:             numericToFloat(carbs),
		FatG:               numericToFloat(fat),
		FiberG:             numericToFloat(fiber),
		SugarG:             numericToFloat(sugar),
		SodiumMg:           numericToFloat(sodium),
	}

	lines, err := r.getRecipeIngredientLines(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("load ingredient lines: %w", err)
	}
	recipe.IngredientLines = lines

	steps, err := r.getRecipeSteps(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("load steps: %w", err)
	}
	recipe.Steps = steps

	return &recipe, nil
}

// GetAll retrieves all recipes with pagination
func (r *Repository) List(ctx context.Context, userID uuid.UUID, filter domain.RecipeFilter, limit, offset int) ([]domain.Recipe, error) {
	var sb strings.Builder
	args := []any{userID}
	argPos := 2

	sb.WriteString(`
		SELECT
			r.id, r.user_id, r.name, r.description,
			r.prep_time_minutes, r.cook_time_minutes, r.total_time_minutes,
			r.servings, r.yield_quantity, r.yield_unit,
			r.image_url, r.tags, r.search_vector,
			r.created_at, r.updated_at, r.deleted_at,
			c.id, c.user_id, c.name, c.created_at,
			mi.id, mi.user_id, mi.name, mi.description, mi.created_at, mi.updated_at,
			rn.calories_total, rn.calories_per_serving,
			rn.protein_g, rn.carbs_g, rn.fat_g, rn.fiber_g, rn.sugar_g, rn.sodium_mg
		FROM recipes r
		JOIN cuisines c ON r.cuisine_id = c.id
		JOIN ingredients mi ON r.main_ingredient_id = mi.id
		LEFT JOIN recipe_nutrition rn ON rn.recipe_id = r.id
		WHERE ` + activeClause("r") + `
		  AND ` + accessClause("r", 1) + `
	`)

	if filter.CuisineID != nil {
		sb.WriteString(fmt.Sprintf(" AND r.cuisine_id = $%d", argPos))
		args = append(args, *filter.CuisineID)
		argPos++
	}

	if filter.IngredientID != nil {
		sb.WriteString(fmt.Sprintf(`
			AND (
				r.main_ingredient_id = $%d OR EXISTS (
					SELECT 1 FROM recipe_ingredient_lines ril
					WHERE ril.recipe_id = r.id AND ril.ingredient_id = $%d
				)
			)
		`, argPos, argPos))
		args = append(args, *filter.IngredientID)
		argPos++
	}

	if filter.AllergyID != nil {
		sb.WriteString(fmt.Sprintf(`
			AND NOT EXISTS (
				SELECT 1 FROM ingredient_allergies ia
				WHERE ia.ingredient_id = r.main_ingredient_id AND ia.allergy_id = $%d
			)
			AND NOT EXISTS (
				SELECT 1
				FROM recipe_ingredient_lines ril
				JOIN ingredient_allergies ia ON ril.ingredient_id = ia.ingredient_id
				WHERE ril.recipe_id = r.id AND ia.allergy_id = $%d
			)
		`, argPos, argPos))
		args = append(args, *filter.AllergyID)
		argPos++
	}

	if len(filter.Tags) > 0 {
		sb.WriteString(fmt.Sprintf(" AND r.tags @> $%d", argPos))
		args = append(args, filter.Tags)
		argPos++
	}

	sb.WriteString(fmt.Sprintf(" ORDER BY r.created_at DESC LIMIT $%d OFFSET $%d", argPos, argPos+1))
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, sb.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("query recipes: %w", err)
	}
	defer rows.Close()

	return r.scanRecipes(ctx, rows)
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

	if recipe.TotalTimeMinutes == 0 {
		recipe.TotalTimeMinutes = recipe.PrepTimeMinutes + recipe.CookTimeMinutes
	}

	query := `
		INSERT INTO recipes (
			id, user_id, name, description,
			prep_time_minutes, cook_time_minutes, total_time_minutes,
			servings, yield_quantity, yield_unit,
			main_ingredient_id, cuisine_id,
			image_url, tags, search_vector
		) VALUES (
			$1, $2, $3, $4,
			$5, $6, $7,
			$8, $9, $10,
			$11, $12,
			$13, $14, $15
		)
	`

	var imageURL *string
	if recipe.ImageURL != "" {
		imageURL = &recipe.ImageURL
	}

	var yieldUnit *string
	if recipe.YieldUnit != "" {
		yieldUnit = &recipe.YieldUnit
	}

	tags := recipe.Tags
	if tags == nil {
		tags = []string{}
	}

	_, err = tx.Exec(ctx, query,
		recipe.ID, recipe.UserID, recipe.Name, recipe.Description,
		recipe.PrepTimeMinutes, recipe.CookTimeMinutes, recipe.TotalTimeMinutes,
		recipe.Servings, recipe.YieldQuantity, yieldUnit,
		recipe.MainIngredient.ID, recipe.Cuisine.ID,
		imageURL, tags, recipe.SearchVector,
	)
	if err != nil {
		return fmt.Errorf("insert recipe: %w", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO recipe_nutrition (
			recipe_id, calories_total, calories_per_serving,
			protein_g, carbs_g, fat_g, fiber_g, sugar_g, sodium_mg
		) VALUES (
			$1, $2, $3,
			$4, $5, $6, $7, $8, $9
		)
	`, recipe.ID, recipe.Nutrition.CaloriesTotal, recipe.Nutrition.CaloriesPerServing,
		recipe.Nutrition.ProteinG, recipe.Nutrition.CarbsG, recipe.Nutrition.FatG,
		recipe.Nutrition.FiberG, recipe.Nutrition.SugarG, recipe.Nutrition.SodiumMg,
	)
	if err != nil {
		return fmt.Errorf("insert recipe nutrition: %w", err)
	}

	for _, line := range recipe.IngredientLines {
		lineID := line.ID
		if lineID == uuid.Nil {
			lineID = uuid.New()
		}
		_, err = tx.Exec(ctx, `
			INSERT INTO recipe_ingredient_lines (
				id, recipe_id, ingredient_id,
				quantity_value, quantity_text, unit,
				is_optional, note, sort_order
			) VALUES (
				$1, $2, $3,
				$4, $5, $6,
				$7, $8, $9
			)
		`, lineID, recipe.ID, line.Ingredient.ID,
			line.QuantityValue, line.QuantityText, line.Unit,
			line.IsOptional, line.Note, line.SortOrder,
		)
		if err != nil {
			return fmt.Errorf("insert recipe ingredient line: %w", err)
		}
	}

	for _, step := range recipe.Steps {
		stepID := step.ID
		if stepID == uuid.Nil {
			stepID = uuid.New()
		}
		_, err = tx.Exec(ctx, `
			INSERT INTO recipe_steps (
				id, recipe_id, step_index, instruction,
				duration_seconds, temperature_value, temperature_unit, media_url
			) VALUES (
				$1, $2, $3, $4,
				$5, $6, $7, $8
			)
		`, stepID, recipe.ID, step.StepIndex, step.Instruction,
			step.DurationSeconds, step.TemperatureValue, step.TemperatureUnit, step.MediaURL,
		)
		if err != nil {
			return fmt.Errorf("insert recipe step: %w", err)
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
			name = $2, description = $3,
			prep_time_minutes = $4, cook_time_minutes = $5, total_time_minutes = $6,
			servings = $7, yield_quantity = $8, yield_unit = $9,
			main_ingredient_id = $10, cuisine_id = $11,
			image_url = $12, tags = $13, search_vector = $14
		WHERE id = $1 AND user_id = $15 AND deleted_at IS NULL
	`

	var imageURL *string
	if recipe.ImageURL != "" {
		imageURL = &recipe.ImageURL
	}

	if recipe.TotalTimeMinutes == 0 {
		recipe.TotalTimeMinutes = recipe.PrepTimeMinutes + recipe.CookTimeMinutes
	}

	var yieldUnit *string
	if recipe.YieldUnit != "" {
		yieldUnit = &recipe.YieldUnit
	}

	tags := recipe.Tags
	if tags == nil {
		tags = []string{}
	}

	result, err := tx.Exec(ctx, query,
		recipe.ID, recipe.Name, recipe.Description,
		recipe.PrepTimeMinutes, recipe.CookTimeMinutes, recipe.TotalTimeMinutes,
		recipe.Servings, recipe.YieldQuantity, yieldUnit,
		recipe.MainIngredient.ID, recipe.Cuisine.ID,
		imageURL, tags, recipe.SearchVector,
		recipe.UserID,
	)
	if err != nil {
		return fmt.Errorf("update recipe: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrRecipeNotFound
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO recipe_nutrition (
			recipe_id, calories_total, calories_per_serving,
			protein_g, carbs_g, fat_g, fiber_g, sugar_g, sodium_mg
		) VALUES (
			$1, $2, $3,
			$4, $5, $6, $7, $8, $9
		)
		ON CONFLICT (recipe_id) DO UPDATE SET
			calories_total = EXCLUDED.calories_total,
			calories_per_serving = EXCLUDED.calories_per_serving,
			protein_g = EXCLUDED.protein_g,
			carbs_g = EXCLUDED.carbs_g,
			fat_g = EXCLUDED.fat_g,
			fiber_g = EXCLUDED.fiber_g,
			sugar_g = EXCLUDED.sugar_g,
			sodium_mg = EXCLUDED.sodium_mg,
			updated_at = NOW()
	`, recipe.ID, recipe.Nutrition.CaloriesTotal, recipe.Nutrition.CaloriesPerServing,
		recipe.Nutrition.ProteinG, recipe.Nutrition.CarbsG, recipe.Nutrition.FatG,
		recipe.Nutrition.FiberG, recipe.Nutrition.SugarG, recipe.Nutrition.SodiumMg,
	)
	if err != nil {
		return fmt.Errorf("upsert recipe nutrition: %w", err)
	}

	// Replace recipe ingredient lines
	_, err = tx.Exec(ctx, `DELETE FROM recipe_ingredient_lines WHERE recipe_id = $1`, recipe.ID)
	if err != nil {
		return fmt.Errorf("delete recipe ingredient lines: %w", err)
	}

	for _, line := range recipe.IngredientLines {
		lineID := line.ID
		if lineID == uuid.Nil {
			lineID = uuid.New()
		}
		_, err = tx.Exec(ctx, `
			INSERT INTO recipe_ingredient_lines (
				id, recipe_id, ingredient_id,
				quantity_value, quantity_text, unit,
				is_optional, note, sort_order
			) VALUES (
				$1, $2, $3,
				$4, $5, $6,
				$7, $8, $9
			)
		`, lineID, recipe.ID, line.Ingredient.ID,
			line.QuantityValue, line.QuantityText, line.Unit,
			line.IsOptional, line.Note, line.SortOrder,
		)
		if err != nil {
			return fmt.Errorf("insert recipe ingredient line: %w", err)
		}
	}

	_, err = tx.Exec(ctx, `DELETE FROM recipe_steps WHERE recipe_id = $1`, recipe.ID)
	if err != nil {
		return fmt.Errorf("delete recipe steps: %w", err)
	}

	for _, step := range recipe.Steps {
		stepID := step.ID
		if stepID == uuid.Nil {
			stepID = uuid.New()
		}
		_, err = tx.Exec(ctx, `
			INSERT INTO recipe_steps (
				id, recipe_id, step_index, instruction,
				duration_seconds, temperature_value, temperature_unit, media_url
			) VALUES (
				$1, $2, $3, $4,
				$5, $6, $7, $8
			)
		`, stepID, recipe.ID, step.StepIndex, step.Instruction,
			step.DurationSeconds, step.TemperatureValue, step.TemperatureUnit, step.MediaURL,
		)
		if err != nil {
			return fmt.Errorf("insert recipe step: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// Delete removes a recipe
func (r *Repository) Delete(ctx context.Context, userID, id uuid.UUID) error {
	result, err := r.pool.Exec(ctx, `
		UPDATE recipes
		SET deleted_at = NOW()
		WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
	`, id, userID)
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
func (r *Repository) GetSimilar(ctx context.Context, userID, recipeID uuid.UUID, limit int) ([]domain.Recipe, error) {
	// First get the vector for the target recipe
	var targetVector pgvector.Vector
	err := r.pool.QueryRow(ctx,
		`SELECT search_vector FROM recipes r WHERE r.id = $1 AND `+activeClause("r")+` AND `+accessClause("r", 2),
		recipeID, userID,
	).Scan(&targetVector)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRecipeNotFound
		}
		return nil, fmt.Errorf("get target vector: %w", err)
	}

	query := `
		SELECT
			r.id, r.user_id, r.name, r.description,
			r.prep_time_minutes, r.cook_time_minutes, r.total_time_minutes,
			r.servings, r.yield_quantity, r.yield_unit,
			r.image_url, r.tags, r.search_vector,
			r.created_at, r.updated_at, r.deleted_at,
			c.id, c.user_id, c.name, c.created_at,
			mi.id, mi.user_id, mi.name, mi.description, mi.created_at, mi.updated_at,
			rn.calories_total, rn.calories_per_serving,
			rn.protein_g, rn.carbs_g, rn.fat_g, rn.fiber_g, rn.sugar_g, rn.sodium_mg
		FROM recipes r
		JOIN cuisines c ON r.cuisine_id = c.id
		JOIN ingredients mi ON r.main_ingredient_id = mi.id
		LEFT JOIN recipe_nutrition rn ON rn.recipe_id = r.id
		WHERE r.id != $1
		  AND ` + activeClause("r") + `
		  AND ` + accessClause("r", 2) + `
		ORDER BY r.search_vector <=> $3
		LIMIT $4
	`

	rows, err := r.pool.Query(ctx, query, recipeID, userID, targetVector, limit)
	if err != nil {
		return nil, fmt.Errorf("query similar recipes: %w", err)
	}
	defer rows.Close()

	return r.scanRecipes(ctx, rows)
}

// GetByCuisine retrieves recipes by cuisine ID
// Count returns the total number of recipes
func (r *Repository) Count(ctx context.Context, userID uuid.UUID, filter domain.RecipeFilter) (int64, error) {
	var sb strings.Builder
	args := []any{userID}
	argPos := 2

	sb.WriteString(`SELECT COUNT(*) FROM recipes r WHERE ` + activeClause("r") + ` AND ` + accessClause("r", 1))

	if filter.CuisineID != nil {
		sb.WriteString(fmt.Sprintf(" AND r.cuisine_id = $%d", argPos))
		args = append(args, *filter.CuisineID)
		argPos++
	}
	if filter.IngredientID != nil {
		sb.WriteString(fmt.Sprintf(`
			AND (
				r.main_ingredient_id = $%d OR EXISTS (
					SELECT 1 FROM recipe_ingredient_lines ril
					WHERE ril.recipe_id = r.id AND ril.ingredient_id = $%d
				)
			)
		`, argPos, argPos))
		args = append(args, *filter.IngredientID)
		argPos++
	}
	if filter.AllergyID != nil {
		sb.WriteString(fmt.Sprintf(`
			AND NOT EXISTS (
				SELECT 1 FROM ingredient_allergies ia
				WHERE ia.ingredient_id = r.main_ingredient_id AND ia.allergy_id = $%d
			)
			AND NOT EXISTS (
				SELECT 1
				FROM recipe_ingredient_lines ril
				JOIN ingredient_allergies ia ON ril.ingredient_id = ia.ingredient_id
				WHERE ril.recipe_id = r.id AND ia.allergy_id = $%d
			)
		`, argPos, argPos))
		args = append(args, *filter.AllergyID)
		argPos++
	}
	if len(filter.Tags) > 0 {
		sb.WriteString(fmt.Sprintf(" AND r.tags @> $%d", argPos))
		args = append(args, filter.Tags)
		argPos++
	}

	var count int64
	err := r.pool.QueryRow(ctx, sb.String(), args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count recipes: %w", err)
	}
	return count, nil
}

// User operations

// GetUserByEmail retrieves a user by email
func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, email, display_name, created_at, updated_at FROM users WHERE email = $1`

	var user domain.User
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.DisplayName, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("query user: %w", err)
	}

	return &user, nil
}

// CreateUser creates a new user record
func (r *Repository) CreateUser(ctx context.Context, user *domain.User) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	query := `INSERT INTO users (id, email, display_name) VALUES ($1, $2, $3)`
	_, err := r.pool.Exec(ctx, query, user.ID, user.Email, user.DisplayName)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}

	return nil
}

// GetUserPasswordHash retrieves the password hash for a user
func (r *Repository) GetUserPasswordHash(ctx context.Context, userID uuid.UUID) (string, error) {
	var hash string
	err := r.pool.QueryRow(ctx, `SELECT password_hash FROM user_credentials WHERE user_id = $1`, userID).Scan(&hash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrUserNotFound
		}
		return "", fmt.Errorf("query user credentials: %w", err)
	}
	return hash, nil
}

// CreateUserCredentials creates password credentials for a user
func (r *Repository) CreateUserCredentials(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO user_credentials (user_id, password_hash) VALUES ($1, $2)`,
		userID, passwordHash,
	)
	if err != nil {
		return fmt.Errorf("insert user credentials: %w", err)
	}
	return nil
}

// Ingredient operations

// GetIngredientByID retrieves an ingredient by ID
func (r *Repository) GetIngredientByID(ctx context.Context, userID, id uuid.UUID) (*domain.Ingredient, error) {
	query := `SELECT id, user_id, name, description, created_at, updated_at FROM ingredients WHERE id = $1 AND user_id = $2`

	var ingredient domain.Ingredient
	err := r.pool.QueryRow(ctx, query, id, userID).Scan(
		&ingredient.ID, &ingredient.UserID, &ingredient.Name, &ingredient.Description, &ingredient.CreatedAt, &ingredient.UpdatedAt,
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
func (r *Repository) GetIngredientByName(ctx context.Context, userID uuid.UUID, name string) (*domain.Ingredient, error) {
	query := `SELECT id, user_id, name, description, created_at, updated_at FROM ingredients WHERE user_id = $1 AND name = $2`

	var ingredient domain.Ingredient
	err := r.pool.QueryRow(ctx, query, userID, name).Scan(
		&ingredient.ID, &ingredient.UserID, &ingredient.Name, &ingredient.Description, &ingredient.CreatedAt, &ingredient.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrIngredientNotFound
		}
		return nil, fmt.Errorf("query ingredient: %w", err)
	}

	allergies, err := r.getIngredientAllergies(ctx, ingredient.ID)
	if err != nil {
		return nil, fmt.Errorf("load allergies: %w", err)
	}
	ingredient.Allergies = allergies

	return &ingredient, nil
}

// CreateIngredient creates a new ingredient
func (r *Repository) CreateIngredient(ctx context.Context, ingredient *domain.Ingredient) error {
	if ingredient.ID == uuid.Nil {
		ingredient.ID = uuid.New()
	}

	query := `INSERT INTO ingredients (id, user_id, name, description) VALUES ($1, $2, $3, $4)`
	_, err := r.pool.Exec(ctx, query, ingredient.ID, ingredient.UserID, ingredient.Name, ingredient.Description)
	if err != nil {
		return fmt.Errorf("insert ingredient: %w", err)
	}

	return nil
}

// GetOrCreateIngredient gets an existing ingredient by name or creates it
func (r *Repository) GetOrCreateIngredient(ctx context.Context, userID uuid.UUID, name string) (*domain.Ingredient, error) {
	ingredient, err := r.GetIngredientByName(ctx, userID, name)
	if err == nil {
		return ingredient, nil
	}
	if !errors.Is(err, ErrIngredientNotFound) {
		return nil, err
	}

	// Create new ingredient
	ingredient = &domain.Ingredient{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        name,
		Description: "",
	}
	if err := r.CreateIngredient(ctx, ingredient); err != nil {
		return nil, err
	}

	return ingredient, nil
}

// Cuisine operations

// GetCuisineByID retrieves a cuisine by ID for a user.
func (r *Repository) GetCuisineByID(ctx context.Context, userID, id uuid.UUID) (*domain.Cuisine, error) {
	query := `SELECT id, user_id, name, created_at FROM cuisines WHERE id = $1 AND user_id = $2`

	var cuisine domain.Cuisine
	err := r.pool.QueryRow(ctx, query, id, userID).Scan(&cuisine.ID, &cuisine.UserID, &cuisine.Name, &cuisine.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCuisineNotFound
		}
		return nil, fmt.Errorf("query cuisine: %w", err)
	}

	return &cuisine, nil
}

// GetCuisineByName retrieves a cuisine by name for a user.
func (r *Repository) GetCuisineByName(ctx context.Context, userID uuid.UUID, name string) (*domain.Cuisine, error) {
	query := `SELECT id, user_id, name, created_at FROM cuisines WHERE user_id = $1 AND name = $2`

	var cuisine domain.Cuisine
	err := r.pool.QueryRow(ctx, query, userID, name).Scan(&cuisine.ID, &cuisine.UserID, &cuisine.Name, &cuisine.CreatedAt)
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

	query := `INSERT INTO cuisines (id, user_id, name) VALUES ($1, $2, $3)`
	_, err := r.pool.Exec(ctx, query, cuisine.ID, cuisine.UserID, cuisine.Name)
	if err != nil {
		return fmt.Errorf("insert cuisine: %w", err)
	}

	return nil
}

// GetCuisines retrieves all cuisines for a user.
func (r *Repository) GetCuisines(ctx context.Context, userID uuid.UUID) ([]domain.Cuisine, error) {
	query := `SELECT id, user_id, name, created_at FROM cuisines WHERE user_id = $1 ORDER BY name ASC`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query cuisines: %w", err)
	}
	defer rows.Close()

	var cuisines []domain.Cuisine
	for rows.Next() {
		var cuisine domain.Cuisine
		if err := rows.Scan(&cuisine.ID, &cuisine.UserID, &cuisine.Name, &cuisine.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan cuisine: %w", err)
		}
		cuisines = append(cuisines, cuisine)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate cuisines: %w", err)
	}

	return cuisines, nil
}

// GetOrCreateCuisine gets an existing cuisine by name or creates it for a user.
func (r *Repository) GetOrCreateCuisine(ctx context.Context, userID uuid.UUID, name string) (*domain.Cuisine, error) {
	cuisine, err := r.GetCuisineByName(ctx, userID, name)
	if err == nil {
		return cuisine, nil
	}
	if !errors.Is(err, ErrCuisineNotFound) {
		return nil, err
	}

	// Create new cuisine
	cuisine = &domain.Cuisine{
		ID:     uuid.New(),
		UserID: userID,
		Name:   name,
	}
	if err := r.CreateCuisine(ctx, cuisine); err != nil {
		return nil, err
	}

	return cuisine, nil
}

// GetAllCuisines retrieves all cuisines for a user.
func (r *Repository) GetAllCuisines(ctx context.Context, userID uuid.UUID) ([]domain.Cuisine, error) {
	query := `SELECT id, user_id, name, created_at FROM cuisines WHERE user_id = $1 ORDER BY name`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query cuisines: %w", err)
	}
	defer rows.Close()

	var cuisines []domain.Cuisine
	for rows.Next() {
		var cuisine domain.Cuisine
		if err := rows.Scan(&cuisine.ID, &cuisine.UserID, &cuisine.Name, &cuisine.CreatedAt); err != nil {
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

func (r *Repository) getRecipeIngredientLines(ctx context.Context, recipeID uuid.UUID) ([]domain.RecipeIngredientLine, error) {
	query := `
		SELECT
			ril.id, ril.ingredient_id,
			ril.quantity_value, ril.quantity_text, ril.unit,
			ril.is_optional, ril.note, ril.sort_order,
			i.user_id, i.name, i.description, i.created_at, i.updated_at
		FROM recipe_ingredient_lines ril
		JOIN ingredients i ON i.id = ril.ingredient_id
		WHERE ril.recipe_id = $1
		ORDER BY ril.sort_order
	`

	rows, err := r.pool.Query(ctx, query, recipeID)
	if err != nil {
		return nil, fmt.Errorf("query recipe ingredient lines: %w", err)
	}
	defer rows.Close()

	var lines []domain.RecipeIngredientLine
	for rows.Next() {
		var line domain.RecipeIngredientLine
		var ingredient domain.Ingredient
		var quantityValue pgtype.Numeric

		if err := rows.Scan(
			&line.ID, &ingredient.ID,
			&quantityValue, &line.QuantityText, &line.Unit,
			&line.IsOptional, &line.Note, &line.SortOrder,
			&ingredient.UserID, &ingredient.Name, &ingredient.Description, &ingredient.CreatedAt, &ingredient.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan ingredient line: %w", err)
		}

		valuePtr, err := numericToFloatPtr(quantityValue)
		if err != nil {
			return nil, fmt.Errorf("parse ingredient quantity: %w", err)
		}
		line.QuantityValue = valuePtr

		allergies, err := r.getIngredientAllergies(ctx, ingredient.ID)
		if err != nil {
			return nil, fmt.Errorf("load ingredient allergies: %w", err)
		}
		ingredient.Allergies = allergies

		line.Ingredient = ingredient
		lines = append(lines, line)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate ingredient lines: %w", err)
	}

	return lines, nil
}

func (r *Repository) getRecipeSteps(ctx context.Context, recipeID uuid.UUID) ([]domain.RecipeStep, error) {
	query := `
		SELECT id, step_index, instruction, duration_seconds, temperature_value, temperature_unit, media_url
		FROM recipe_steps
		WHERE recipe_id = $1
		ORDER BY step_index
	`

	rows, err := r.pool.Query(ctx, query, recipeID)
	if err != nil {
		return nil, fmt.Errorf("query recipe steps: %w", err)
	}
	defer rows.Close()

	var steps []domain.RecipeStep
	for rows.Next() {
		var step domain.RecipeStep
		var duration pgtype.Int4
		var temperature pgtype.Numeric

		if err := rows.Scan(
			&step.ID, &step.StepIndex, &step.Instruction,
			&duration, &temperature, &step.TemperatureUnit, &step.MediaURL,
		); err != nil {
			return nil, fmt.Errorf("scan recipe step: %w", err)
		}

		if duration.Valid {
			value := int(duration.Int32)
			step.DurationSeconds = &value
		}

		tempValue, err := numericToFloatPtr(temperature)
		if err != nil {
			return nil, fmt.Errorf("parse temperature value: %w", err)
		}
		step.TemperatureValue = tempValue

		steps = append(steps, step)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate recipe steps: %w", err)
	}

	return steps, nil
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
		var yieldUnit *string
		var yieldQuantity pgtype.Numeric
		var deletedAt *time.Time
		var protein pgtype.Numeric
		var carbs pgtype.Numeric
		var fat pgtype.Numeric
		var fiber pgtype.Numeric
		var sugar pgtype.Numeric
		var sodium pgtype.Numeric
		var caloriesTotal int
		var caloriesPerServing int

		err := rows.Scan(
			&recipe.ID, &recipe.UserID, &recipe.Name, &recipe.Description,
			&recipe.PrepTimeMinutes, &recipe.CookTimeMinutes, &recipe.TotalTimeMinutes,
			&recipe.Servings, &yieldQuantity, &yieldUnit,
			&imageURL, &tags, &searchVector,
			&recipe.CreatedAt, &recipe.UpdatedAt, &deletedAt,
			&cuisine.ID, &cuisine.UserID, &cuisine.Name, &cuisine.CreatedAt,
			&mainIngredient.ID, &mainIngredient.UserID, &mainIngredient.Name, &mainIngredient.Description, &mainIngredient.CreatedAt, &mainIngredient.UpdatedAt,
			&caloriesTotal, &caloriesPerServing,
			&protein, &carbs, &fat, &fiber, &sugar, &sodium,
		)
		if err != nil {
			return nil, fmt.Errorf("scan recipe: %w", err)
		}

		recipe.Cuisine = &cuisine
		recipe.MainIngredient = &mainIngredient
		recipe.SearchVector = searchVector
		if imageURL != nil {
			recipe.ImageURL = *imageURL
		}
		if yieldUnit != nil {
			recipe.YieldUnit = *yieldUnit
		}
		yieldPtr, err := numericToFloatPtr(yieldQuantity)
		if err != nil {
			return nil, fmt.Errorf("parse yield quantity: %w", err)
		}
		recipe.YieldQuantity = yieldPtr
		recipe.Tags = tags
		recipe.DeletedAt = deletedAt
		recipe.Nutrition = domain.RecipeNutrition{
			CaloriesTotal:      caloriesTotal,
			CaloriesPerServing: caloriesPerServing,
			ProteinG:           numericToFloat(protein),
			CarbsG:             numericToFloat(carbs),
			FatG:               numericToFloat(fat),
			FiberG:             numericToFloat(fiber),
			SugarG:             numericToFloat(sugar),
			SodiumMg:           numericToFloat(sodium),
		}

		recipes = append(recipes, recipe)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate recipes: %w", err)
	}

	for i := range recipes {
		lines, err := r.getRecipeIngredientLines(ctx, recipes[i].ID)
		if err != nil {
			return nil, err
		}
		recipes[i].IngredientLines = lines

		steps, err := r.getRecipeSteps(ctx, recipes[i].ID)
		if err != nil {
			return nil, err
		}
		recipes[i].Steps = steps
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
