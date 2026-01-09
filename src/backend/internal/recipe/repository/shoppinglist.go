package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/platepilot/backend/internal/common/domain"
)

var (
	ErrShoppingListNotFound     = errors.New("shopping list not found")
	ErrShoppingListItemNotFound = errors.New("shopping list item not found")
)

// ShoppingListRepository provides access to shopping list data
type ShoppingListRepository struct {
	pool *pgxpool.Pool
}

// NewShoppingListRepository creates a new shopping list repository
func NewShoppingListRepository(pool *pgxpool.Pool) *ShoppingListRepository {
	return &ShoppingListRepository{pool: pool}
}

// GetByID retrieves a shopping list by ID with all items
func (r *ShoppingListRepository) GetByID(ctx context.Context, userID, id uuid.UUID) (*domain.ShoppingList, error) {
	query := `
		SELECT id, user_id, name, week_start_date, created_at, updated_at, completed_at
		FROM shopping_lists
		WHERE id = $1 AND user_id = $2
	`

	var list domain.ShoppingList
	err := r.pool.QueryRow(ctx, query, id, userID).Scan(
		&list.ID, &list.UserID, &list.Name, &list.WeekStartDate,
		&list.CreatedAt, &list.UpdatedAt, &list.CompletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrShoppingListNotFound
		}
		return nil, fmt.Errorf("query shopping list: %w", err)
	}

	// Load items
	items, err := r.getListItems(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("load items: %w", err)
	}
	list.Items = items

	// Load source recipes
	recipes, err := r.getListRecipes(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("load recipes: %w", err)
	}
	list.Recipes = recipes

	return &list, nil
}

// GetAll retrieves all shopping lists for a user
func (r *ShoppingListRepository) GetAll(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.ShoppingList, error) {
	query := `
		SELECT id, user_id, name, week_start_date, created_at, updated_at, completed_at
		FROM shopping_lists
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query shopping lists: %w", err)
	}
	defer rows.Close()

	var lists []domain.ShoppingList
	for rows.Next() {
		var list domain.ShoppingList
		err := rows.Scan(
			&list.ID, &list.UserID, &list.Name, &list.WeekStartDate,
			&list.CreatedAt, &list.UpdatedAt, &list.CompletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan shopping list: %w", err)
		}
		lists = append(lists, list)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate shopping lists: %w", err)
	}

	// Load item counts for each list (for display purposes)
	for i := range lists {
		items, err := r.getListItems(ctx, lists[i].ID)
		if err != nil {
			return nil, fmt.Errorf("load items for list %s: %w", lists[i].ID, err)
		}
		lists[i].Items = items
	}

	return lists, nil
}

// Create creates a new shopping list
func (r *ShoppingListRepository) Create(ctx context.Context, list *domain.ShoppingList) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if list.ID == uuid.Nil {
		list.ID = uuid.New()
	}

	query := `
		INSERT INTO shopping_lists (id, user_id, name, week_start_date)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at, updated_at
	`
	err = tx.QueryRow(ctx, query, list.ID, list.UserID, list.Name, list.WeekStartDate).Scan(
		&list.CreatedAt, &list.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert shopping list: %w", err)
	}

	// Insert items
	for i := range list.Items {
		item := &list.Items[i]
		if item.ID == uuid.Nil {
			item.ID = uuid.New()
		}
		item.ShoppingListID = list.ID

		itemQuery := `
			INSERT INTO shopping_list_items (id, shopping_list_id, ingredient_id, custom_name, quantity, unit, checked, notes, is_custom)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING created_at, updated_at
		`
		err = tx.QueryRow(ctx, itemQuery,
			item.ID, item.ShoppingListID, item.IngredientID, item.CustomName,
			item.Quantity, item.Unit, item.Checked, item.Notes, item.IsCustom,
		).Scan(&item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			return fmt.Errorf("insert shopping list item: %w", err)
		}

		// Insert item sources
		for _, source := range item.Sources {
			sourceQuery := `
				INSERT INTO shopping_list_item_sources (shopping_list_item_id, recipe_id, quantity, unit)
				VALUES ($1, $2, $3, $4)
			`
			_, err = tx.Exec(ctx, sourceQuery, item.ID, source.RecipeID, source.Quantity, source.Unit)
			if err != nil {
				return fmt.Errorf("insert item source: %w", err)
			}
		}
	}

	// Insert recipe associations
	for _, recipe := range list.Recipes {
		_, err = tx.Exec(ctx,
			`INSERT INTO shopping_list_recipes (shopping_list_id, recipe_id) VALUES ($1, $2)`,
			list.ID, recipe.ID,
		)
		if err != nil {
			return fmt.Errorf("insert shopping list recipe: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// Update updates a shopping list's basic properties
func (r *ShoppingListRepository) Update(ctx context.Context, list *domain.ShoppingList) error {
	query := `
		UPDATE shopping_lists
		SET name = $2, completed_at = $3
		WHERE id = $1 AND user_id = $4
	`
	result, err := r.pool.Exec(ctx, query, list.ID, list.Name, list.CompletedAt, list.UserID)
	if err != nil {
		return fmt.Errorf("update shopping list: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrShoppingListNotFound
	}

	return nil
}

// Delete removes a shopping list
func (r *ShoppingListRepository) Delete(ctx context.Context, userID, id uuid.UUID) error {
	result, err := r.pool.Exec(ctx, `DELETE FROM shopping_lists WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("delete shopping list: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrShoppingListNotFound
	}

	return nil
}

// AddItem adds an item to a shopping list
func (r *ShoppingListRepository) AddItem(ctx context.Context, userID uuid.UUID, item *domain.ShoppingListItem) error {
	// Verify the shopping list belongs to the user
	var listUserID uuid.UUID
	err := r.pool.QueryRow(ctx, `SELECT user_id FROM shopping_lists WHERE id = $1`, item.ShoppingListID).Scan(&listUserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrShoppingListNotFound
		}
		return fmt.Errorf("verify shopping list ownership: %w", err)
	}
	if listUserID != userID {
		return ErrShoppingListNotFound
	}

	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}

	query := `
		INSERT INTO shopping_list_items (id, shopping_list_id, ingredient_id, custom_name, quantity, unit, checked, notes, is_custom)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING created_at, updated_at
	`
	err = r.pool.QueryRow(ctx, query,
		item.ID, item.ShoppingListID, item.IngredientID, item.CustomName,
		item.Quantity, item.Unit, item.Checked, item.Notes, item.IsCustom,
	).Scan(&item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert shopping list item: %w", err)
	}

	return nil
}

// UpdateItem updates a shopping list item
func (r *ShoppingListRepository) UpdateItem(ctx context.Context, userID uuid.UUID, item *domain.ShoppingListItem) error {
	query := `
		UPDATE shopping_list_items sli
		SET quantity = $2, unit = $3, checked = $4, notes = $5
		FROM shopping_lists sl
		WHERE sli.id = $1
		  AND sli.shopping_list_id = sl.id
		  AND sl.user_id = $6
	`
	result, err := r.pool.Exec(ctx, query, item.ID, item.Quantity, item.Unit, item.Checked, item.Notes, userID)
	if err != nil {
		return fmt.Errorf("update shopping list item: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrShoppingListItemNotFound
	}

	return nil
}

// ToggleItemChecked toggles the checked state of an item
func (r *ShoppingListRepository) ToggleItemChecked(ctx context.Context, userID, itemID uuid.UUID) (bool, error) {
	query := `
		UPDATE shopping_list_items sli
		SET checked = NOT sli.checked
		FROM shopping_lists sl
		WHERE sli.id = $1
		  AND sli.shopping_list_id = sl.id
		  AND sl.user_id = $2
		RETURNING sli.checked
	`

	var newCheckedState bool
	err := r.pool.QueryRow(ctx, query, itemID, userID).Scan(&newCheckedState)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, ErrShoppingListItemNotFound
		}
		return false, fmt.Errorf("toggle item checked: %w", err)
	}

	return newCheckedState, nil
}

// DeleteItem removes an item from a shopping list
func (r *ShoppingListRepository) DeleteItem(ctx context.Context, userID, itemID uuid.UUID) error {
	query := `
		DELETE FROM shopping_list_items sli
		USING shopping_lists sl
		WHERE sli.id = $1
		  AND sli.shopping_list_id = sl.id
		  AND sl.user_id = $2
	`
	result, err := r.pool.Exec(ctx, query, itemID, userID)
	if err != nil {
		return fmt.Errorf("delete shopping list item: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrShoppingListItemNotFound
	}

	return nil
}

// GetAggregatedIngredients retrieves and aggregates ingredients from multiple recipes
func (r *ShoppingListRepository) GetAggregatedIngredients(ctx context.Context, userID uuid.UUID, recipeIDs []uuid.UUID) ([]domain.AggregatedIngredient, error) {
	if len(recipeIDs) == 0 {
		return nil, nil
	}

	// Build query to get all ingredients from the specified recipes
	// with their quantities and units from recipe_ingredient_lines
	query := `
		SELECT
			i.id as ingredient_id,
			i.name as ingredient_name,
			ic.id as category_id,
			ic.name as category_name,
			ril.quantity_value,
			ril.unit,
			r.id as recipe_id,
			r.name as recipe_name
		FROM recipe_ingredient_lines ril
		JOIN ingredients i ON ril.ingredient_id = i.id
		JOIN recipes r ON ril.recipe_id = r.id
		LEFT JOIN ingredient_categories ic ON i.category_id = ic.id
		WHERE ril.recipe_id = ANY($1)
		  AND (r.user_id = $2 OR EXISTS (SELECT 1 FROM recipe_shares rs WHERE rs.recipe_id = r.id AND rs.shared_with_user_id = $2))
		ORDER BY ic.display_order NULLS LAST, i.name
	`

	rows, err := r.pool.Query(ctx, query, recipeIDs, userID)
	if err != nil {
		return nil, fmt.Errorf("query ingredients: %w", err)
	}
	defer rows.Close()

	// Aggregate ingredients by ID and unit
	type ingredientKey struct {
		ID   uuid.UUID
		Unit string
	}
	aggregated := make(map[ingredientKey]*domain.AggregatedIngredient)
	ingredientOrder := []ingredientKey{}

	for rows.Next() {
		var ingredientID uuid.UUID
		var ingredientName string
		var categoryID *uuid.UUID
		var categoryName *string
		var quantity *float64
		var unit *string
		var recipeID uuid.UUID
		var recipeName string

		err := rows.Scan(
			&ingredientID, &ingredientName, &categoryID, &categoryName,
			&quantity, &unit, &recipeID, &recipeName,
		)
		if err != nil {
			return nil, fmt.Errorf("scan ingredient: %w", err)
		}

		unitStr := ""
		if unit != nil {
			unitStr = *unit
		}
		key := ingredientKey{ID: ingredientID, Unit: unitStr}

		if agg, exists := aggregated[key]; exists {
			// Add to existing aggregation
			if quantity != nil {
				found := false
				for i, qu := range agg.Quantities {
					// Same unit, add quantities
					if (qu.Unit == nil && unit == nil) || (qu.Unit != nil && unit != nil && *qu.Unit == *unit) {
						if qu.Quantity != nil && quantity != nil {
							newQty := *qu.Quantity + *quantity
							agg.Quantities[i].Quantity = &newQty
						} else if quantity != nil {
							agg.Quantities[i].Quantity = quantity
						}
						found = true
						break
					}
				}
				if !found {
					agg.Quantities = append(agg.Quantities, domain.QuantityUnit{Quantity: quantity, Unit: unit})
				}
			}
			agg.RecipeSources = append(agg.RecipeSources, domain.RecipeSource{
				RecipeID:   recipeID,
				RecipeName: recipeName,
				Quantity:   quantity,
				Unit:       unit,
			})
		} else {
			// New aggregation
			agg := &domain.AggregatedIngredient{
				IngredientID:   ingredientID,
				IngredientName: ingredientName,
				CategoryID:     categoryID,
				CategoryName:   categoryName,
				Quantities:     []domain.QuantityUnit{{Quantity: quantity, Unit: unit}},
				RecipeSources: []domain.RecipeSource{{
					RecipeID:   recipeID,
					RecipeName: recipeName,
					Quantity:   quantity,
					Unit:       unit,
				}},
			}
			aggregated[key] = agg
			ingredientOrder = append(ingredientOrder, key)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate ingredients: %w", err)
	}

	// Convert map to slice maintaining order
	result := make([]domain.AggregatedIngredient, 0, len(aggregated))
	for _, key := range ingredientOrder {
		result = append(result, *aggregated[key])
	}

	return result, nil
}

// GetRecipesByIDs retrieves multiple recipes by their IDs
func (r *ShoppingListRepository) GetRecipesByIDs(ctx context.Context, userID uuid.UUID, recipeIDs []uuid.UUID) ([]domain.Recipe, error) {
	if len(recipeIDs) == 0 {
		return nil, nil
	}

	query := `
		SELECT r.id, r.user_id, r.name
		FROM recipes r
		WHERE r.id = ANY($1)
		  AND (r.user_id = $2 OR EXISTS (SELECT 1 FROM recipe_shares rs WHERE rs.recipe_id = r.id AND rs.shared_with_user_id = $2))
	`

	rows, err := r.pool.Query(ctx, query, recipeIDs, userID)
	if err != nil {
		return nil, fmt.Errorf("query recipes: %w", err)
	}
	defer rows.Close()

	var recipes []domain.Recipe
	for rows.Next() {
		var recipe domain.Recipe
		if err := rows.Scan(&recipe.ID, &recipe.UserID, &recipe.Name); err != nil {
			return nil, fmt.Errorf("scan recipe: %w", err)
		}
		recipes = append(recipes, recipe)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate recipes: %w", err)
	}

	return recipes, nil
}

// Count returns the total number of shopping lists for a user
func (r *ShoppingListRepository) Count(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM shopping_lists WHERE user_id = $1`, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count shopping lists: %w", err)
	}
	return count, nil
}

// Helper methods

func (r *ShoppingListRepository) getListItems(ctx context.Context, listID uuid.UUID) ([]domain.ShoppingListItem, error) {
	query := `
		SELECT
			sli.id, sli.shopping_list_id, sli.ingredient_id, sli.custom_name,
			sli.quantity, sli.unit, sli.checked, sli.notes, sli.is_custom,
			sli.created_at, sli.updated_at,
			i.id, i.name, i.description,
			ic.id, ic.name, ic.display_order
		FROM shopping_list_items sli
		LEFT JOIN ingredients i ON sli.ingredient_id = i.id
		LEFT JOIN ingredient_categories ic ON i.category_id = ic.id
		WHERE sli.shopping_list_id = $1
		ORDER BY ic.display_order NULLS LAST, COALESCE(i.name, sli.custom_name)
	`

	rows, err := r.pool.Query(ctx, query, listID)
	if err != nil {
		return nil, fmt.Errorf("query list items: %w", err)
	}
	defer rows.Close()

	var items []domain.ShoppingListItem
	for rows.Next() {
		var item domain.ShoppingListItem
		var ingredientID, catID *uuid.UUID
		var ingredientName, ingredientDescription, catName *string
		var catOrder *int

		err := rows.Scan(
			&item.ID, &item.ShoppingListID, &item.IngredientID, &item.CustomName,
			&item.Quantity, &item.Unit, &item.Checked, &item.Notes, &item.IsCustom,
			&item.CreatedAt, &item.UpdatedAt,
			&ingredientID, &ingredientName, &ingredientDescription,
			&catID, &catName, &catOrder,
		)
		if err != nil {
			return nil, fmt.Errorf("scan list item: %w", err)
		}

		if ingredientID != nil && ingredientName != nil {
			description := ""
			if ingredientDescription != nil {
				description = *ingredientDescription
			}
			item.Ingredient = &domain.Ingredient{
				ID:          *ingredientID,
				Name:        *ingredientName,
				Description: description,
			}
		}

		if catID != nil {
			order := 0
			if catOrder != nil {
				order = *catOrder
			}
			name := ""
			if catName != nil {
				name = *catName
			}
			item.Category = &domain.IngredientCategory{
				ID:           *catID,
				Name:         name,
				DisplayOrder: order,
			}
		}

		// Load sources for this item
		sources, err := r.getItemSources(ctx, item.ID)
		if err != nil {
			return nil, fmt.Errorf("load item sources: %w", err)
		}
		item.Sources = sources

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate list items: %w", err)
	}

	return items, nil
}

func (r *ShoppingListRepository) getItemSources(ctx context.Context, itemID uuid.UUID) ([]domain.ShoppingListItemSource, error) {
	query := `
		SELECT slis.recipe_id, r.name, slis.quantity, slis.unit
		FROM shopping_list_item_sources slis
		JOIN recipes r ON slis.recipe_id = r.id
		WHERE slis.shopping_list_item_id = $1
	`

	rows, err := r.pool.Query(ctx, query, itemID)
	if err != nil {
		return nil, fmt.Errorf("query item sources: %w", err)
	}
	defer rows.Close()

	var sources []domain.ShoppingListItemSource
	for rows.Next() {
		var source domain.ShoppingListItemSource
		err := rows.Scan(&source.RecipeID, &source.RecipeName, &source.Quantity, &source.Unit)
		if err != nil {
			return nil, fmt.Errorf("scan item source: %w", err)
		}
		sources = append(sources, source)
	}

	return sources, nil
}

func (r *ShoppingListRepository) getListRecipes(ctx context.Context, listID uuid.UUID) ([]domain.Recipe, error) {
	query := `
		SELECT r.id, r.user_id, r.name
		FROM recipes r
		JOIN shopping_list_recipes slr ON r.id = slr.recipe_id
		WHERE slr.shopping_list_id = $1
	`

	rows, err := r.pool.Query(ctx, query, listID)
	if err != nil {
		return nil, fmt.Errorf("query list recipes: %w", err)
	}
	defer rows.Close()

	var recipes []domain.Recipe
	for rows.Next() {
		var recipe domain.Recipe
		if err := rows.Scan(&recipe.ID, &recipe.UserID, &recipe.Name); err != nil {
			return nil, fmt.Errorf("scan recipe: %w", err)
		}
		recipes = append(recipes, recipe)
	}

	return recipes, nil
}

// GetIngredientCategories retrieves all ingredient categories
func (r *ShoppingListRepository) GetIngredientCategories(ctx context.Context) ([]domain.IngredientCategory, error) {
	query := `SELECT id, name, display_order, created_at FROM ingredient_categories ORDER BY display_order`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query ingredient categories: %w", err)
	}
	defer rows.Close()

	var categories []domain.IngredientCategory
	for rows.Next() {
		var cat domain.IngredientCategory
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.DisplayOrder, &cat.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan category: %w", err)
		}
		categories = append(categories, cat)
	}

	return categories, nil
}
