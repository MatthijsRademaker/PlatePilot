-- Shopping List Feature Migration
-- Adds shopping list tables and enhances recipe_ingredients with quantity/unit

-- Add ingredient categories for shopping list grouping
CREATE TABLE ingredient_categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) NOT NULL UNIQUE,
    display_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX ix_ingredient_categories_name ON ingredient_categories (name);

-- Seed default categories
INSERT INTO ingredient_categories (name, display_order) VALUES
    ('Produce', 1),
    ('Dairy', 2),
    ('Meat & Seafood', 3),
    ('Bakery', 4),
    ('Frozen', 5),
    ('Pantry', 6),
    ('Beverages', 7),
    ('Snacks', 8),
    ('Condiments & Sauces', 9),
    ('Spices & Seasonings', 10),
    ('Other', 100);

-- Add category to ingredients table
ALTER TABLE ingredients ADD COLUMN category_id UUID REFERENCES ingredient_categories(id);
CREATE INDEX ix_ingredients_category_id ON ingredients (category_id);

-- Add quantity and unit to recipe_ingredients junction table
ALTER TABLE recipe_ingredients ADD COLUMN quantity DECIMAL(10, 2);
ALTER TABLE recipe_ingredients ADD COLUMN unit VARCHAR(20);

-- Shopping Lists table
CREATE TABLE shopping_lists (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    week_start_date DATE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    completed_at TIMESTAMPTZ
);

CREATE INDEX ix_shopping_lists_user_id ON shopping_lists (user_id);
CREATE INDEX ix_shopping_lists_created_at ON shopping_lists (created_at DESC);

-- Shopping List Items table
CREATE TABLE shopping_list_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    shopping_list_id UUID NOT NULL REFERENCES shopping_lists(id) ON DELETE CASCADE,
    ingredient_id UUID REFERENCES ingredients(id) ON DELETE SET NULL,
    custom_name VARCHAR(100),
    quantity DECIMAL(10, 2),
    unit VARCHAR(20),
    checked BOOLEAN NOT NULL DEFAULT FALSE,
    notes TEXT,
    is_custom BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    -- Either ingredient_id or custom_name must be set
    CONSTRAINT chk_item_identity CHECK (
        (ingredient_id IS NOT NULL AND is_custom = FALSE) OR
        (custom_name IS NOT NULL AND is_custom = TRUE)
    )
);

CREATE INDEX ix_shopping_list_items_list_id ON shopping_list_items (shopping_list_id);
CREATE INDEX ix_shopping_list_items_ingredient_id ON shopping_list_items (ingredient_id);
CREATE INDEX ix_shopping_list_items_checked ON shopping_list_items (shopping_list_id, checked);

-- Junction table: which recipes contributed to a shopping list
CREATE TABLE shopping_list_recipes (
    shopping_list_id UUID NOT NULL REFERENCES shopping_lists(id) ON DELETE CASCADE,
    recipe_id UUID NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    PRIMARY KEY (shopping_list_id, recipe_id)
);

CREATE INDEX ix_shopping_list_recipes_recipe_id ON shopping_list_recipes (recipe_id);

-- Track which items came from which recipes (for breakdown view)
CREATE TABLE shopping_list_item_sources (
    shopping_list_item_id UUID NOT NULL REFERENCES shopping_list_items(id) ON DELETE CASCADE,
    recipe_id UUID NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    quantity DECIMAL(10, 2),
    unit VARCHAR(20),
    PRIMARY KEY (shopping_list_item_id, recipe_id)
);

CREATE INDEX ix_shopping_list_item_sources_recipe_id ON shopping_list_item_sources (recipe_id);

-- Updated at triggers
CREATE TRIGGER update_shopping_lists_updated_at
    BEFORE UPDATE ON shopping_lists
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_shopping_list_items_updated_at
    BEFORE UPDATE ON shopping_list_items
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
