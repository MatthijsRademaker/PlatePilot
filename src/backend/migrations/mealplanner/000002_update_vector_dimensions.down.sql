-- Revert vector dimensions from 1536 back to 128
--
-- NOTE: This migration resets all search vectors. Data will need to be repopulated.

-- Drop existing vector index
DROP INDEX IF EXISTS ix_recipes_search_vector;

-- Create temporary column with old dimensions
ALTER TABLE recipes ADD COLUMN search_vector_old vector(128);

-- Initialize with zero vectors
UPDATE recipes SET search_vector_old = (
    SELECT array_to_string(array_fill(0::real, ARRAY[128]), ',')::vector(128)
);

-- Drop new column and rename old one
ALTER TABLE recipes DROP COLUMN search_vector;
ALTER TABLE recipes RENAME COLUMN search_vector_old TO search_vector;

-- Make column NOT NULL
ALTER TABLE recipes ALTER COLUMN search_vector SET NOT NULL;

-- Recreate vector similarity search index with original dimensions
CREATE INDEX ix_recipes_search_vector ON recipes
    USING ivfflat (search_vector vector_cosine_ops)
    WITH (lists = 100);
