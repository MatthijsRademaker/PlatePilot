-- Add units table and per-recipe ingredient quantities/units
CREATE TABLE units (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(20) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (user_id, name)
);

CREATE INDEX ix_units_user_id ON units (user_id);

ALTER TABLE recipe_ingredients
    ADD COLUMN quantity TEXT NOT NULL DEFAULT '',
    ADD COLUMN unit TEXT NOT NULL DEFAULT '';
