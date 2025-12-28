-- Recipe API Initial Migration - Rollback

DROP TRIGGER IF EXISTS update_recipes_updated_at ON recipes;
DROP FUNCTION IF EXISTS update_updated_at_column();

DROP TABLE IF EXISTS ingredient_allergies;
DROP TABLE IF EXISTS recipe_ingredients;
DROP TABLE IF EXISTS recipes;
DROP TABLE IF EXISTS allergies;
DROP TABLE IF EXISTS ingredients;
DROP TABLE IF EXISTS cuisines;

-- Note: We don't drop the extensions as they may be used by other databases/schemas
