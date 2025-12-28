-- Create databases for each service
CREATE DATABASE recipedb;
CREATE DATABASE mealplannerdb;

-- Enable pgvector extension on each database
\c recipedb;
CREATE EXTENSION IF NOT EXISTS vector;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

\c mealplannerdb;
CREATE EXTENSION IF NOT EXISTS vector;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Grant permissions
\c recipedb;
GRANT ALL PRIVILEGES ON DATABASE recipedb TO platepilot;
GRANT ALL ON SCHEMA public TO platepilot;

\c mealplannerdb;
GRANT ALL PRIVILEGES ON DATABASE mealplannerdb TO platepilot;
GRANT ALL ON SCHEMA public TO platepilot;
