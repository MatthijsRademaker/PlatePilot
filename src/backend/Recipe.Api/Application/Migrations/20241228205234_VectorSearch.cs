using Microsoft.EntityFrameworkCore.Migrations;
using Pgvector;

#nullable disable

namespace RecipeApplication.Migrations
{
    /// <inheritdoc />
    public partial class VectorSearch : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.AddColumn<Vector>(
                name: "SearchVector",
                table: "Recipes",
                type: "Vector(384)",
                nullable: false);
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropColumn(
                name: "SearchVector",
                table: "Recipes");
        }
    }
}
