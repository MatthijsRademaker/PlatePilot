DROP INDEX IF EXISTS ix_recipes_user_id;
ALTER TABLE recipes DROP COLUMN IF EXISTS user_id;
