-- Down migration for shopping list feature

-- Drop triggers
DROP TRIGGER IF EXISTS update_shopping_list_items_updated_at ON shopping_list_items;
DROP TRIGGER IF EXISTS update_shopping_lists_updated_at ON shopping_lists;

-- Drop indexes and tables in reverse order
DROP INDEX IF EXISTS ix_shopping_list_item_sources_recipe_id;
DROP TABLE IF EXISTS shopping_list_item_sources;

DROP INDEX IF EXISTS ix_shopping_list_recipes_recipe_id;
DROP TABLE IF EXISTS shopping_list_recipes;

DROP INDEX IF EXISTS ix_shopping_list_items_checked;
DROP INDEX IF EXISTS ix_shopping_list_items_ingredient_id;
DROP INDEX IF EXISTS ix_shopping_list_items_list_id;
DROP TABLE IF EXISTS shopping_list_items;

DROP INDEX IF EXISTS ix_shopping_lists_created_at;
DROP INDEX IF EXISTS ix_shopping_lists_user_id;
DROP TABLE IF EXISTS shopping_lists;

-- Remove category from ingredients
DROP INDEX IF EXISTS ix_ingredients_category_id;
ALTER TABLE ingredients DROP COLUMN IF EXISTS category_id;

-- Drop ingredient categories
DROP INDEX IF EXISTS ix_ingredient_categories_name;
DROP TABLE IF EXISTS ingredient_categories;
