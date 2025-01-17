using Microsoft.EntityFrameworkCore.Migrations;

namespace MealPlannerInfrastructure.Migrations;

public partial class AddMatchesFunction : Migration
{
    protected override void Up(MigrationBuilder migrationBuilder)
    {
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
        migrationBuilder.Sql("DROP FUNCTION IF EXISTS recipe_matches;");
        migrationBuilder.Sql("DROP TYPE IF EXISTS constraint_type;");
    }
}
