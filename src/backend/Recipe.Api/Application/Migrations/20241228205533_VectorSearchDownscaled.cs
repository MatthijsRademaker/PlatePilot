using Microsoft.EntityFrameworkCore.Migrations;
using Pgvector;

#nullable disable

namespace RecipeApplication.Migrations
{
    /// <inheritdoc />
    public partial class VectorSearchDownscaled : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.AlterColumn<Vector>(
                name: "SearchVector",
                table: "Recipes",
                type: "Vector(128)",
                nullable: false,
                oldClrType: typeof(Vector),
                oldType: "Vector(384)");
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.AlterColumn<Vector>(
                name: "SearchVector",
                table: "Recipes",
                type: "Vector(384)",
                nullable: false,
                oldClrType: typeof(Vector),
                oldType: "Vector(128)");
        }
    }
}
