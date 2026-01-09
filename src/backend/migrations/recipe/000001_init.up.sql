-- Recipe API Initial Schema (squashed)
-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS vector;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS citext;

-- Updated at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Users
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email CITEXT NOT NULL UNIQUE,
    display_name TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE user_credentials (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TRIGGER update_user_credentials_updated_at
    BEFORE UPDATE ON user_credentials
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE user_oauth_accounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider TEXT NOT NULL,
    provider_user_id TEXT NOT NULL,
    email CITEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (provider, provider_user_id),
    UNIQUE (user_id, provider)
);

CREATE TABLE user_refresh_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    last_used_at TIMESTAMPTZ,
    user_agent TEXT,
    ip_address TEXT
);

CREATE UNIQUE INDEX ix_user_refresh_tokens_token_hash ON user_refresh_tokens (token_hash);
CREATE INDEX ix_user_refresh_tokens_user_id ON user_refresh_tokens (user_id);

-- Cuisines (user-scoped)
CREATE TABLE cuisines (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE UNIQUE INDEX ux_cuisines_user_id_name ON cuisines (user_id, name);
CREATE INDEX ix_cuisines_user_id ON cuisines (user_id);
CREATE INDEX ix_cuisines_name ON cuisines (name);

-- Ingredients (user-scoped)
CREATE TABLE ingredients (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (user_id, name)
);

CREATE INDEX ix_ingredients_user_id ON ingredients (user_id);
CREATE INDEX ix_ingredients_name ON ingredients (name);

CREATE TRIGGER update_ingredients_updated_at
    BEFORE UPDATE ON ingredients
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Allergies (global)
CREATE TABLE allergies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX ix_allergies_name ON allergies (name);

-- Ingredient-Allergy many-to-many
CREATE TABLE ingredient_allergies (
    ingredient_id UUID NOT NULL REFERENCES ingredients(id) ON DELETE CASCADE,
    allergy_id UUID NOT NULL REFERENCES allergies(id) ON DELETE CASCADE,
    PRIMARY KEY (ingredient_id, allergy_id)
);

CREATE INDEX ix_ingredient_allergies_allergy_id ON ingredient_allergies (allergy_id);

-- Recipes
CREATE TABLE recipes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    prep_time_minutes INTEGER NOT NULL DEFAULT 0,
    cook_time_minutes INTEGER NOT NULL DEFAULT 0,
    total_time_minutes INTEGER NOT NULL DEFAULT 0,
    servings INTEGER NOT NULL DEFAULT 1,
    yield_quantity NUMERIC(12,3),
    yield_unit TEXT,
    main_ingredient_id UUID NOT NULL REFERENCES ingredients(id) ON DELETE RESTRICT,
    cuisine_id UUID NOT NULL REFERENCES cuisines(id) ON DELETE RESTRICT,
    image_url TEXT,
    tags TEXT[] NOT NULL DEFAULT '{}',
    search_vector vector(1536) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX ix_recipes_user_id ON recipes (user_id);
CREATE INDEX ix_recipes_name ON recipes (name);
CREATE INDEX ix_recipes_cuisine_id ON recipes (cuisine_id);
CREATE INDEX ix_recipes_main_ingredient_id ON recipes (main_ingredient_id);
CREATE INDEX ix_recipes_tags ON recipes USING GIN (tags);

CREATE INDEX ix_recipes_search_vector ON recipes
    USING ivfflat (search_vector vector_cosine_ops)
    WITH (lists = 100);

CREATE TRIGGER update_recipes_updated_at
    BEFORE UPDATE ON recipes
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Recipe ingredient lines
CREATE TABLE recipe_ingredient_lines (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    recipe_id UUID NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    ingredient_id UUID NOT NULL REFERENCES ingredients(id) ON DELETE RESTRICT,
    quantity_value NUMERIC(12,3),
    quantity_text TEXT,
    unit TEXT,
    is_optional BOOLEAN NOT NULL DEFAULT FALSE,
    note TEXT,
    sort_order INTEGER NOT NULL
);

CREATE UNIQUE INDEX ux_recipe_ingredient_lines_recipe_sort ON recipe_ingredient_lines (recipe_id, sort_order);
CREATE INDEX ix_recipe_ingredient_lines_recipe_id ON recipe_ingredient_lines (recipe_id);
CREATE INDEX ix_recipe_ingredient_lines_ingredient_id ON recipe_ingredient_lines (ingredient_id);

-- Structured recipe steps
CREATE TABLE recipe_steps (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    recipe_id UUID NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    step_index INTEGER NOT NULL,
    instruction TEXT NOT NULL,
    duration_seconds INTEGER,
    temperature_value NUMERIC(8,2),
    temperature_unit TEXT,
    media_url TEXT
);

CREATE UNIQUE INDEX ux_recipe_steps_recipe_step ON recipe_steps (recipe_id, step_index);
CREATE INDEX ix_recipe_steps_recipe_id ON recipe_steps (recipe_id);

-- Ingredient nutrition (optional enrichment)
CREATE TABLE ingredient_nutrition (
    ingredient_id UUID PRIMARY KEY REFERENCES ingredients(id) ON DELETE CASCADE,
    serving_size_value NUMERIC(12,3) NOT NULL,
    serving_unit TEXT NOT NULL,
    calories INTEGER NOT NULL DEFAULT 0,
    protein_g NUMERIC(8,2) NOT NULL DEFAULT 0,
    carbs_g NUMERIC(8,2) NOT NULL DEFAULT 0,
    fat_g NUMERIC(8,2) NOT NULL DEFAULT 0,
    fiber_g NUMERIC(8,2) NOT NULL DEFAULT 0,
    sugar_g NUMERIC(8,2) NOT NULL DEFAULT 0,
    sodium_mg NUMERIC(10,2) NOT NULL DEFAULT 0
);

-- Recipe nutrition (stored aggregate)
CREATE TABLE recipe_nutrition (
    recipe_id UUID PRIMARY KEY REFERENCES recipes(id) ON DELETE CASCADE,
    calories_total INTEGER NOT NULL DEFAULT 0,
    calories_per_serving INTEGER NOT NULL DEFAULT 0,
    protein_g NUMERIC(8,2) NOT NULL DEFAULT 0,
    carbs_g NUMERIC(8,2) NOT NULL DEFAULT 0,
    fat_g NUMERIC(8,2) NOT NULL DEFAULT 0,
    fiber_g NUMERIC(8,2) NOT NULL DEFAULT 0,
    sugar_g NUMERIC(8,2) NOT NULL DEFAULT 0,
    sodium_mg NUMERIC(10,2) NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Sharing support (future)
CREATE TABLE recipe_shares (
    recipe_id UUID NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    shared_with_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (recipe_id, shared_with_user_id)
);

CREATE INDEX ix_recipe_shares_user_id ON recipe_shares (shared_with_user_id);
