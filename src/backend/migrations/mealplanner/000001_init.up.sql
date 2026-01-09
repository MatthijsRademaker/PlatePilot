-- MealPlanner API Initial Schema (squashed)
-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS vector;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Updated at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Read model table for recipes (denormalized)
CREATE TABLE recipes (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    name TEXT,
    description TEXT,
    prep_time_minutes INTEGER DEFAULT 0,
    cook_time_minutes INTEGER DEFAULT 0,
    total_time_minutes INTEGER DEFAULT 0,
    servings INTEGER DEFAULT 1,
    yield_quantity NUMERIC(12,3),
    yield_unit TEXT,
    search_vector vector(1536) NOT NULL,
    cuisine_id UUID NOT NULL,
    cuisine_name TEXT,
    main_ingredient_id UUID NOT NULL,
    main_ingredient_name TEXT,
    ingredient_ids UUID[] NOT NULL DEFAULT '{}',
    allergy_ids UUID[] NOT NULL DEFAULT '{}',
    tags TEXT[] NOT NULL DEFAULT '{}',
    image_url TEXT,
    calories_total INTEGER DEFAULT 0,
    calories_per_serving INTEGER DEFAULT 0,
    protein_g NUMERIC(8,2) DEFAULT 0,
    carbs_g NUMERIC(8,2) DEFAULT 0,
    fat_g NUMERIC(8,2) DEFAULT 0,
    fiber_g NUMERIC(8,2) DEFAULT 0,
    sugar_g NUMERIC(8,2) DEFAULT 0,
    sodium_mg NUMERIC(10,2) DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE UNIQUE INDEX ix_recipes_id ON recipes (id);
CREATE INDEX ix_recipes_user_id ON recipes (user_id);
CREATE INDEX ix_recipes_cuisine_id ON recipes (cuisine_id);
CREATE INDEX ix_recipes_main_ingredient_id ON recipes (main_ingredient_id);
CREATE INDEX ix_recipes_ingredient_ids ON recipes USING GIN (ingredient_ids);
CREATE INDEX ix_recipes_allergy_ids ON recipes USING GIN (allergy_ids);
CREATE INDEX ix_recipes_tags ON recipes USING GIN (tags);

CREATE INDEX ix_recipes_search_vector ON recipes
    USING ivfflat (search_vector vector_cosine_ops)
    WITH (lists = 100);

CREATE TRIGGER update_recipes_updated_at
    BEFORE UPDATE ON recipes
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Ingredient lines for shopping list aggregation
CREATE TABLE recipe_ingredient_lines (
    recipe_id UUID NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    ingredient_id UUID NOT NULL,
    ingredient_name TEXT NOT NULL,
    quantity_value NUMERIC(12,3),
    quantity_text TEXT,
    unit TEXT,
    is_optional BOOLEAN NOT NULL DEFAULT FALSE,
    note TEXT,
    sort_order INTEGER NOT NULL,
    PRIMARY KEY (recipe_id, sort_order)
);

CREATE INDEX ix_recipe_ingredient_lines_recipe_id ON recipe_ingredient_lines (recipe_id);
CREATE INDEX ix_recipe_ingredient_lines_ingredient_id ON recipe_ingredient_lines (ingredient_id);

-- Meal plan tables
CREATE TABLE meal_plans (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (user_id, start_date)
);

CREATE INDEX ix_meal_plans_user_id ON meal_plans (user_id);
CREATE INDEX ix_meal_plans_start_date ON meal_plans (start_date);

CREATE TRIGGER update_meal_plans_updated_at
    BEFORE UPDATE ON meal_plans
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE meal_plan_slots (
    plan_id UUID NOT NULL REFERENCES meal_plans(id) ON DELETE CASCADE,
    slot_date DATE NOT NULL,
    meal_type TEXT NOT NULL,
    recipe_id UUID NOT NULL,
    PRIMARY KEY (plan_id, slot_date, meal_type)
);

ALTER TABLE meal_plan_slots
    ADD CONSTRAINT meal_plan_slots_meal_type_check
    CHECK (meal_type IN ('breakfast', 'lunch', 'dinner', 'snack'));

CREATE INDEX ix_meal_plan_slots_plan_id ON meal_plan_slots (plan_id);
CREATE INDEX ix_meal_plan_slots_recipe_id ON meal_plan_slots (recipe_id);
