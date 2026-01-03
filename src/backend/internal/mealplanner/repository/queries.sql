-- MealPlanner Repository Queries (Read Model)

-- name: GetRecipeByID :one
SELECT
    id,
    name,
    description,
    prep_time,
    cook_time,
    search_vector,
    cuisine_id,
    cuisine_name,
    main_ingredient_id,
    main_ingredient_name,
    ingredient_ids,
    allergy_ids,
    directions,
    image_url,
    tags,
    published_date,
    calories,
    created_at,
    updated_at
FROM recipes
WHERE id = $1;

-- name: GetAllRecipes :many
SELECT
    id,
    name,
    description,
    prep_time,
    cook_time,
    search_vector,
    cuisine_id,
    cuisine_name,
    main_ingredient_id,
    main_ingredient_name,
    ingredient_ids,
    allergy_ids,
    directions,
    image_url,
    tags,
    published_date,
    calories,
    created_at,
    updated_at
FROM recipes
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetRecipesByCuisine :many
SELECT
    id,
    name,
    description,
    prep_time,
    cook_time,
    search_vector,
    cuisine_id,
    cuisine_name,
    main_ingredient_id,
    main_ingredient_name,
    ingredient_ids,
    allergy_ids,
    directions,
    image_url,
    tags,
    published_date,
    calories,
    created_at,
    updated_at
FROM recipes
WHERE cuisine_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetRecipesByIngredient :many
SELECT
    id,
    name,
    description,
    prep_time,
    cook_time,
    search_vector,
    cuisine_id,
    cuisine_name,
    main_ingredient_id,
    main_ingredient_name,
    ingredient_ids,
    allergy_ids,
    directions,
    image_url,
    tags,
    published_date,
    calories,
    created_at,
    updated_at
FROM recipes
WHERE main_ingredient_id = $1 OR $1 = ANY(ingredient_ids)
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetRecipesExcludingAllergy :many
SELECT
    id,
    name,
    description,
    prep_time,
    cook_time,
    search_vector,
    cuisine_id,
    cuisine_name,
    main_ingredient_id,
    main_ingredient_name,
    ingredient_ids,
    allergy_ids,
    directions,
    image_url,
    tags,
    published_date,
    calories,
    created_at,
    updated_at
FROM recipes
WHERE NOT ($1 = ANY(allergy_ids))
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetSimilarRecipes :many
SELECT
    id,
    name,
    description,
    prep_time,
    cook_time,
    search_vector,
    cuisine_id,
    cuisine_name,
    main_ingredient_id,
    main_ingredient_name,
    ingredient_ids,
    allergy_ids,
    directions,
    image_url,
    tags,
    published_date,
    calories,
    created_at,
    updated_at,
    1 - (search_vector <=> $1) as similarity
FROM recipes
WHERE id != $2
ORDER BY search_vector <=> $1
LIMIT $3;

-- name: UpsertRecipe :exec
INSERT INTO recipes (
    id,
    name,
    description,
    prep_time,
    cook_time,
    search_vector,
    cuisine_id,
    cuisine_name,
    main_ingredient_id,
    main_ingredient_name,
    ingredient_ids,
    allergy_ids,
    directions,
    image_url,
    tags,
    published_date,
    calories
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
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
    published_date = EXCLUDED.published_date,
    calories = EXCLUDED.calories,
    updated_at = NOW();

-- name: DeleteRecipe :exec
DELETE FROM recipes WHERE id = $1;

-- name: CountRecipes :one
SELECT COUNT(*) FROM recipes;

-- name: GetRecipeVectorByID :one
SELECT id, search_vector
FROM recipes
WHERE id = $1;
