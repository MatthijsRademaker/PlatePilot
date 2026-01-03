-- Recipe API Initial Migration
-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS vector;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Cuisines table
CREATE TABLE cuisines (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX ix_cuisines_name ON cuisines (name);

-- Ingredients table
CREATE TABLE ingredients (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) NOT NULL,
    quantity TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX ix_ingredients_name ON ingredients (name);

-- Allergies table
CREATE TABLE allergies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX ix_allergies_name ON allergies (name);

-- Recipes table
CREATE TABLE recipes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    prep_time TEXT NOT NULL DEFAULT '',
    cook_time TEXT NOT NULL DEFAULT '',
    main_ingredient_id UUID NOT NULL REFERENCES ingredients(id) ON DELETE CASCADE,
    cuisine_id UUID NOT NULL REFERENCES cuisines(id) ON DELETE CASCADE,
    directions TEXT[] NOT NULL DEFAULT '{}',
    -- Owned types flattened
    nutritional_info_calories INTEGER NOT NULL DEFAULT 0,
    metadata_search_vector vector(128) NOT NULL,
    metadata_image_url TEXT,
    metadata_tags TEXT[] NOT NULL DEFAULT '{}',
    metadata_published_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX ix_recipes_name ON recipes (name);
CREATE INDEX ix_recipes_cuisine_id ON recipes (cuisine_id);
CREATE INDEX ix_recipes_main_ingredient_id ON recipes (main_ingredient_id);

-- Vector similarity search index (IVFFlat for cosine distance)
CREATE INDEX ix_recipes_search_vector ON recipes
    USING ivfflat (metadata_search_vector vector_cosine_ops)
    WITH (lists = 100);

-- Recipe-Ingredient many-to-many relationship
CREATE TABLE recipe_ingredients (
    recipe_id UUID NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    ingredient_id UUID NOT NULL REFERENCES ingredients(id) ON DELETE CASCADE,
    PRIMARY KEY (recipe_id, ingredient_id)
);

CREATE INDEX ix_recipe_ingredients_recipe_id ON recipe_ingredients (recipe_id);
CREATE INDEX ix_recipe_ingredients_ingredient_id ON recipe_ingredients (ingredient_id);

-- Ingredient-Allergy many-to-many relationship
CREATE TABLE ingredient_allergies (
    ingredient_id UUID NOT NULL REFERENCES ingredients(id) ON DELETE CASCADE,
    allergy_id UUID NOT NULL REFERENCES allergies(id) ON DELETE CASCADE,
    PRIMARY KEY (ingredient_id, allergy_id)
);

CREATE INDEX ix_ingredient_allergies_allergy_id ON ingredient_allergies (allergy_id);

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
