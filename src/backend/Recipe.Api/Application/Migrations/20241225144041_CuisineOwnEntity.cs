using Microsoft.EntityFrameworkCore.Migrations;
using Npgsql.EntityFrameworkCore.PostgreSQL.Metadata;

#nullable disable

namespace RecipeApi.Application.Migrations
{
    /// <inheritdoc />
    public partial class CuisineOwnEntity : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropColumn(
                name: "Cuisine",
                table: "Recipes");

            migrationBuilder.AddColumn<int>(
                name: "CuisineId",
                table: "Recipes",
                type: "integer",
                nullable: false,
                defaultValue: 0);

            migrationBuilder.CreateTable(
                name: "Cuisines",
                columns: table => new
                {
                    Id = table.Column<int>(type: "integer", nullable: false)
                        .Annotation("Npgsql:ValueGenerationStrategy", NpgsqlValueGenerationStrategy.IdentityByDefaultColumn),
                    Name = table.Column<string>(type: "text", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_Cuisines", x => x.Id);
                });

            migrationBuilder.CreateIndex(
                name: "IX_Recipes_CuisineId",
                table: "Recipes",
                column: "CuisineId");

            migrationBuilder.AddForeignKey(
                name: "FK_Recipes_Cuisines_CuisineId",
                table: "Recipes",
                column: "CuisineId",
                principalTable: "Cuisines",
                principalColumn: "Id",
                onDelete: ReferentialAction.Cascade);
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropForeignKey(
                name: "FK_Recipes_Cuisines_CuisineId",
                table: "Recipes");

            migrationBuilder.DropTable(
                name: "Cuisines");

            migrationBuilder.DropIndex(
                name: "IX_Recipes_CuisineId",
                table: "Recipes");

            migrationBuilder.DropColumn(
                name: "CuisineId",
                table: "Recipes");

            migrationBuilder.AddColumn<string>(
                name: "Cuisine",
                table: "Recipes",
                type: "text",
                nullable: false,
                defaultValue: "");
        }
    }
}
