-- Recipe Repository Queries

-- name: GetRecipeByID :one
SELECT
    r.id,
    r.name,
    r.description,
    r.prep_time,
    r.cook_time,
    r.main_ingredient_id,
    r.cuisine_id,
    r.directions,
    r.nutritional_info_calories,
    r.metadata_search_vector,
    r.metadata_image_url,
    r.metadata_tags,
    r.metadata_published_date,
    r.created_at,
    r.updated_at
FROM recipes r
WHERE r.id = $1;

-- name: GetAllRecipes :many
SELECT
    r.id,
    r.name,
    r.description,
    r.prep_time,
    r.cook_time,
    r.main_ingredient_id,
    r.cuisine_id,
    r.directions,
    r.nutritional_info_calories,
    r.metadata_search_vector,
    r.metadata_image_url,
    r.metadata_tags,
    r.metadata_published_date,
    r.created_at,
    r.updated_at
FROM recipes r
ORDER BY r.created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetRecipesByCuisine :many
SELECT
    r.id,
    r.name,
    r.description,
    r.prep_time,
    r.cook_time,
    r.main_ingredient_id,
    r.cuisine_id,
    r.directions,
    r.nutritional_info_calories,
    r.metadata_search_vector,
    r.metadata_image_url,
    r.metadata_tags,
    r.metadata_published_date,
    r.created_at,
    r.updated_at
FROM recipes r
WHERE r.cuisine_id = $1
ORDER BY r.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetRecipesByMainIngredient :many
SELECT
    r.id,
    r.name,
    r.description,
    r.prep_time,
    r.cook_time,
    r.main_ingredient_id,
    r.cuisine_id,
    r.directions,
    r.nutritional_info_calories,
    r.metadata_search_vector,
    r.metadata_image_url,
    r.metadata_tags,
    r.metadata_published_date,
    r.created_at,
    r.updated_at
FROM recipes r
WHERE r.main_ingredient_id = $1
ORDER BY r.created_at DESC
LIMIT $2 OFFSET $3;

-- name: CreateRecipe :one
INSERT INTO recipes (
    id,
    name,
    description,
    prep_time,
    cook_time,
    main_ingredient_id,
    cuisine_id,
    directions,
    nutritional_info_calories,
    metadata_search_vector,
    metadata_image_url,
    metadata_tags,
    metadata_published_date
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
)
RETURNING *;

-- name: UpdateRecipe :one
UPDATE recipes
SET
    name = $2,
    description = $3,
    prep_time = $4,
    cook_time = $5,
    main_ingredient_id = $6,
    cuisine_id = $7,
    directions = $8,
    nutritional_info_calories = $9,
    metadata_search_vector = $10,
    metadata_image_url = $11,
    metadata_tags = $12,
    metadata_published_date = $13
WHERE id = $1
RETURNING *;

-- name: DeleteRecipe :exec
DELETE FROM recipes WHERE id = $1;

-- name: GetSimilarRecipes :many
SELECT
    r.id,
    r.name,
    r.description,
    r.prep_time,
    r.cook_time,
    r.main_ingredient_id,
    r.cuisine_id,
    r.directions,
    r.nutritional_info_calories,
    r.metadata_search_vector,
    r.metadata_image_url,
    r.metadata_tags,
    r.metadata_published_date,
    r.created_at,
    r.updated_at,
    1 - (r.metadata_search_vector <=> $1) as similarity
FROM recipes r
WHERE r.id != $2
ORDER BY r.metadata_search_vector <=> $1
LIMIT $3;

-- Ingredient queries

-- name: GetIngredientByID :one
SELECT id, name, quantity, created_at
FROM ingredients
WHERE id = $1;

-- name: GetIngredientByName :one
SELECT id, name, quantity, created_at
FROM ingredients
WHERE name = $1;

-- name: CreateIngredient :one
INSERT INTO ingredients (id, name, quantity)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetRecipeIngredients :many
SELECT i.id, i.name, i.quantity, i.created_at
FROM ingredients i
JOIN recipe_ingredients ri ON i.id = ri.ingredient_id
WHERE ri.recipe_id = $1;

-- name: AddRecipeIngredient :exec
INSERT INTO recipe_ingredients (recipe_id, ingredient_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: RemoveRecipeIngredients :exec
DELETE FROM recipe_ingredients WHERE recipe_id = $1;

-- Cuisine queries

-- name: GetCuisineByID :one
SELECT id, name, created_at
FROM cuisines
WHERE id = $1;

-- name: GetCuisineByName :one
SELECT id, name, created_at
FROM cuisines
WHERE name = $1;

-- name: CreateCuisine :one
INSERT INTO cuisines (id, name)
VALUES ($1, $2)
RETURNING *;

-- name: GetAllCuisines :many
SELECT id, name, created_at
FROM cuisines
ORDER BY name;

-- Allergy queries

-- name: GetAllergyByID :one
SELECT id, name, created_at
FROM allergies
WHERE id = $1;

-- name: GetAllergyByName :one
SELECT id, name, created_at
FROM allergies
WHERE name = $1;

-- name: CreateAllergy :one
INSERT INTO allergies (id, name)
VALUES ($1, $2)
RETURNING *;

-- name: GetIngredientAllergies :many
SELECT a.id, a.name, a.created_at
FROM allergies a
JOIN ingredient_allergies ia ON a.id = ia.allergy_id
WHERE ia.ingredient_id = $1;

-- name: AddIngredientAllergy :exec
INSERT INTO ingredient_allergies (ingredient_id, allergy_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;
