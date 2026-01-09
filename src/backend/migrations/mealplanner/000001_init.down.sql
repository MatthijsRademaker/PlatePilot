-- Drop MealPlanner schema (squashed)
DROP TABLE IF EXISTS meal_plan_slots;
DROP TABLE IF EXISTS meal_plans;
DROP TABLE IF EXISTS recipe_ingredient_lines;
DROP TABLE IF EXISTS recipes;

DROP FUNCTION IF EXISTS update_updated_at_column();
