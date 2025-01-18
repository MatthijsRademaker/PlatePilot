using Microsoft.EntityFrameworkCore.Migrations;
using Pgvector;

namespace Infrastructure.Migrations;

public partial class InitialCreate : Migration
{
    protected override void Up(MigrationBuilder migrationBuilder)
    {
        // Enable required extensions
        migrationBuilder.AlterDatabase().Annotation("Npgsql:PostgresExtension:vector", ",,");

        // Create base recipes table
        migrationBuilder.CreateTable(
            name: "recipes",
            columns: table => new
            {
                Id = table.Column<Guid>(type: "uuid", nullable: false),
                SearchVector = table.Column<Vector>(type: "vector(128)", nullable: false),
                CuisineId = table.Column<int>(type: "integer", nullable: false),
                MainIngredientId = table.Column<int>(type: "integer", nullable: false),
                IngredientIds = table.Column<List<int>>(type: "integer[]", nullable: false),
                AllergyIds = table.Column<List<int>>(type: "integer[]", nullable: false),
            }
        );

        // Create materialized view
        migrationBuilder.Sql(
            @"
            CREATE SCHEMA materialized;
            
            CREATE MATERIALIZED VIEW materialized.recipe_view AS 
            SELECT *
            FROM recipes;

            CREATE UNIQUE INDEX recipe_view_id_idx ON materialized.recipe_view(id);
            CREATE INDEX recipe_view_search_idx ON materialized.recipe_view 
            USING ivfflat (search_vector vector_cosine_ops);
        "
        );

        migrationBuilder.Sql(
            @"
            CREATE TYPE constraint_type AS ENUM ('AllergyConstraint', 'CuisineConstraint', 'IngredientConstraint');

            CREATE OR REPLACE FUNCTION recipe_matches(
                recipe_id uuid,
                constraint_type constraint_type,
                entity_id uuid
            ) RETURNS boolean
            LANGUAGE plpgsql
            AS $$
            BEGIN
                RETURN CASE constraint_type
                    WHEN 'AllergyConstraint' THEN
                        EXISTS (
                            SELECT 1 
                            FROM materialized.recipe_view r
                            WHERE r.id = $1 
                            AND $3 = ANY(r.allergy_ids)
                        )
                    WHEN 'CuisineConstraint' THEN
                        EXISTS (
                            SELECT 1 
                            FROM materialized.recipe_view r
                            WHERE r.id = $1 
                            AND r.cuisine_id = $3
                        )
                    WHEN 'IngredientConstraint' THEN
                        EXISTS (
                            SELECT 1 
                            FROM materialized.recipe_view r
                            WHERE r.id = $1 
                            AND ($3 = r.main_ingredient_id OR $3 = ANY(r.ingredient_ids))
                        )
                END;
            END;
            $$;"
        );
    }

    protected override void Down(MigrationBuilder migrationBuilder)
    {
        migrationBuilder.Sql("DROP MATERIALIZED VIEW IF EXISTS materialized.recipe_view;");
        migrationBuilder.Sql("DROP SCHEMA IF EXISTS materialized CASCADE;");
        migrationBuilder.Sql("DROP FUNCTION IF EXISTS recipe_matches;");
        migrationBuilder.Sql("DROP TYPE IF EXISTS constraint_type;");
        migrationBuilder.DropTable(name: "recipes");
    }
}
