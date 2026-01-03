-- Update vector dimensions from 128 to 1536 for LLM embeddings
-- Supports both Ollama (nomic-embed-text: 768) and OpenAI (text-embedding-3-small: 1536)
-- Smaller embeddings are padded to 1536 by the application
--
-- NOTE: This migration resets all search vectors. Data will be repopulated from recipe-events.

-- Drop existing vector index
DROP INDEX IF EXISTS ix_recipes_search_vector;

-- Create temporary column with new dimensions
ALTER TABLE recipes ADD COLUMN search_vector_new vector(1536);

-- Initialize with zero vectors (will be repopulated by events)
UPDATE recipes SET search_vector_new = (
    SELECT array_to_string(array_fill(0::real, ARRAY[1536]), ',')::vector(1536)
);

-- Drop old column and rename new one
ALTER TABLE recipes DROP COLUMN search_vector;
ALTER TABLE recipes RENAME COLUMN search_vector_new TO search_vector;

-- Make column NOT NULL
ALTER TABLE recipes ALTER COLUMN search_vector SET NOT NULL;

-- Recreate vector similarity search index with new dimensions
CREATE INDEX ix_recipes_search_vector ON recipes
    USING ivfflat (search_vector vector_cosine_ops)
    WITH (lists = 100);

COMMENT ON COLUMN recipes.search_vector IS 'LLM-generated embeddings (1536-dim for OpenAI, padded for smaller models)';
