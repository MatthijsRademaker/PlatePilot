-- MealPlanner API Initial Migration - Rollback

DROP TRIGGER IF EXISTS update_recipes_updated_at ON recipes;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP FUNCTION IF EXISTS recipe_matches(UUID, constraint_type, UUID);
DROP TYPE IF EXISTS constraint_type;
DROP TABLE IF EXISTS recipes;

-- Note: We don't drop the extensions as they may be used by other databases/schemas
