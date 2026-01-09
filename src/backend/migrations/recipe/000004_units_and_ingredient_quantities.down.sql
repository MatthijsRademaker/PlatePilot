ALTER TABLE recipe_ingredients
    DROP COLUMN IF EXISTS unit,
    DROP COLUMN IF EXISTS quantity;

DROP INDEX IF EXISTS ix_units_user_id;
DROP TABLE IF EXISTS units;
