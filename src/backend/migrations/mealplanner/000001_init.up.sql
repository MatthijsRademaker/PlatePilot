-- MealPlanner API Initial Migration (Read Model)
-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS vector;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Read model table for recipes (denormalized for query performance)
CREATE TABLE recipes (
    id UUID PRIMARY KEY,
    name VARCHAR(255),
    description TEXT,
    prep_time VARCHAR(50),
    cook_time VARCHAR(50),
    search_vector vector(128) NOT NULL,
    cuisine_id UUID NOT NULL,
    cuisine_name VARCHAR(100),
    main_ingredient_id UUID NOT NULL,
    main_ingredient_name VARCHAR(100),
    ingredient_ids UUID[] NOT NULL DEFAULT '{}',
    allergy_ids UUID[] NOT NULL DEFAULT '{}',
    directions TEXT[] DEFAULT '{}',
    image_url TEXT,
    tags TEXT[] DEFAULT '{}',
    published_date TIMESTAMPTZ,
    calories INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Primary key index
CREATE UNIQUE INDEX ix_recipes_id ON recipes (id);

-- Vector similarity search index (IVFFlat for cosine distance)
CREATE INDEX ix_recipes_search_vector ON recipes
    USING ivfflat (search_vector vector_cosine_ops)
    WITH (lists = 100);

-- Filtering indexes
CREATE INDEX ix_recipes_cuisine_id ON recipes (cuisine_id);
CREATE INDEX ix_recipes_main_ingredient_id ON recipes (main_ingredient_id);

-- Array indexes for ingredient and allergy filtering
CREATE INDEX ix_recipes_ingredient_ids ON recipes USING GIN (ingredient_ids);
CREATE INDEX ix_recipes_allergy_ids ON recipes USING GIN (allergy_ids);

-- Constraint type enum for filtering functions
CREATE TYPE constraint_type AS ENUM (
    'AllergyConstraint',
    'CuisineConstraint',
    'IngredientConstraint'
);

-- Function to check if a recipe matches a constraint
CREATE OR REPLACE FUNCTION recipe_matches(
    p_recipe_id UUID,
    p_constraint_type constraint_type,
    p_entity_id UUID
) RETURNS BOOLEAN
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN CASE p_constraint_type
        WHEN 'AllergyConstraint' THEN
            EXISTS (
                SELECT 1
                FROM recipes r
                WHERE r.id = p_recipe_id
                AND p_entity_id = ANY(r.allergy_ids)
            )
        WHEN 'CuisineConstraint' THEN
            EXISTS (
                SELECT 1
                FROM recipes r
                WHERE r.id = p_recipe_id
                AND r.cuisine_id = p_entity_id
            )
        WHEN 'IngredientConstraint' THEN
            EXISTS (
                SELECT 1
                FROM recipes r
                WHERE r.id = p_recipe_id
                AND (p_entity_id = r.main_ingredient_id OR p_entity_id = ANY(r.ingredient_ids))
            )
    END;
END;
$$;

-- Updated at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply trigger to recipes table
CREATE TRIGGER update_recipes_updated_at
    BEFORE UPDATE ON recipes
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
