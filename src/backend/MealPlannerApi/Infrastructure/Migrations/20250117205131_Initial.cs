using System.Collections.Generic;
using Microsoft.EntityFrameworkCore.Migrations;
using Npgsql.EntityFrameworkCore.PostgreSQL.Metadata;
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
            migrationBuilder.EnsureSchema(
                name: "materialized");

            migrationBuilder.AlterDatabase()
                .Annotation("Npgsql:PostgresExtension:vector", ",,");

            migrationBuilder.CreateTable(
                name: "recipe_view",
                schema: "materialized",
                columns: table => new
                {
                    Id = table.Column<int>(type: "integer", nullable: false)
                        .Annotation("Npgsql:ValueGenerationStrategy", NpgsqlValueGenerationStrategy.IdentityByDefaultColumn),
                    SearchVector = table.Column<Vector>(type: "Vector(128)", nullable: false),
                    CuisineId = table.Column<int>(type: "integer", nullable: false),
                    MainIngredientId = table.Column<int>(type: "integer", nullable: false),
                    IngredientIds = table.Column<List<int>>(type: "integer[]", nullable: false),
                    AllergyIds = table.Column<List<int>>(type: "integer[]", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_recipe_view", x => x.Id);
                });

            migrationBuilder.CreateIndex(
                name: "IX_recipe_view_Id",
                schema: "materialized",
                table: "recipe_view",
                column: "Id",
                unique: true);

            migrationBuilder.CreateIndex(
                name: "IX_recipe_view_SearchVector",
                schema: "materialized",
                table: "recipe_view",
                column: "SearchVector")
                .Annotation("Npgsql:IndexMethod", "ivfflat")
                .Annotation("Npgsql:IndexOperators", new[] { "vector_cosine_ops" });
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropTable(
                name: "recipe_view",
                schema: "materialized");
        }
    }
}
