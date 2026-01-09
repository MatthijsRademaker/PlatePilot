-- Drop Recipe API schema (squashed)
DROP TABLE IF EXISTS recipe_shares;
DROP TABLE IF EXISTS recipe_nutrition;
DROP TABLE IF EXISTS ingredient_nutrition;
DROP TABLE IF EXISTS recipe_steps;
DROP TABLE IF EXISTS recipe_ingredient_lines;
DROP TABLE IF EXISTS recipes;
DROP TABLE IF EXISTS ingredient_allergies;
DROP TABLE IF EXISTS allergies;
DROP TABLE IF EXISTS ingredients;
DROP TABLE IF EXISTS cuisines;
DROP TABLE IF EXISTS user_refresh_tokens;
DROP TABLE IF EXISTS user_oauth_accounts;
DROP TABLE IF EXISTS user_credentials;
DROP TABLE IF EXISTS users;

DROP FUNCTION IF EXISTS update_updated_at_column();
