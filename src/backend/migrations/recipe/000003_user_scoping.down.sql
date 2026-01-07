DROP TRIGGER IF EXISTS update_user_credentials_updated_at ON user_credentials;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

DROP INDEX IF EXISTS ix_recipe_shares_user_id;
DROP TABLE IF EXISTS recipe_shares;

DROP INDEX IF EXISTS ix_recipes_user_id;
ALTER TABLE recipes DROP COLUMN IF EXISTS user_id;

DROP INDEX IF EXISTS ix_user_refresh_tokens_user_id;
DROP INDEX IF EXISTS ix_user_refresh_tokens_token_hash;
DROP TABLE IF EXISTS user_refresh_tokens;
DROP TABLE IF EXISTS user_oauth_accounts;
DROP TABLE IF EXISTS user_credentials;
DROP TABLE IF EXISTS users;

DROP EXTENSION IF EXISTS citext;
