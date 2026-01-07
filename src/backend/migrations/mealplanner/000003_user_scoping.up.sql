-- Add user scoping to read model
ALTER TABLE recipes ADD COLUMN user_id UUID NOT NULL;
CREATE INDEX ix_recipes_user_id ON recipes (user_id);
