using System;
using System.Collections.Generic;
using Microsoft.EntityFrameworkCore.Migrations;
using Pgvector;

#nullable disable

namespace MealPlannerInfrastructure.Migrations
{
    /// <inheritdoc />
    public partial class Initial : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.AlterDatabase()
                .Annotation("Npgsql:PostgresExtension:vector", ",,");

            migrationBuilder.CreateTable(
                name: "recipes",
                columns: table => new
                {
                    Id = table.Column<Guid>(type: "uuid", nullable: false),
                    SearchVector = table.Column<Vector>(type: "Vector(128)", nullable: false),
                    CuisineId = table.Column<Guid>(type: "uuid", nullable: false),
                    MainIngredientId = table.Column<Guid>(type: "uuid", nullable: false),
                    IngredientIds = table.Column<List<Guid>>(type: "uuid[]", nullable: false),
                    AllergyIds = table.Column<List<Guid>>(type: "uuid[]", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_recipes", x => x.Id);
                });

            migrationBuilder.CreateIndex(
                name: "IX_recipes_Id",
                table: "recipes",
                column: "Id",
                unique: true);

            migrationBuilder.CreateIndex(
                name: "IX_recipes_SearchVector",
                table: "recipes",
                column: "SearchVector")
                .Annotation("Npgsql:IndexMethod", "ivfflat")
                .Annotation("Npgsql:IndexOperators", new[] { "vector_cosine_ops" });
            
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

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.Sql("DROP FUNCTION IF EXISTS recipe_matches;");
            migrationBuilder.Sql("DROP TYPE IF EXISTS constraint_type;");
            migrationBuilder.DropTable(
                name: "recipes");
        }
    }
}
