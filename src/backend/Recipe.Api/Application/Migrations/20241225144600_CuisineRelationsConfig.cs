using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace RecipeApi.Application.Migrations
{
    /// <inheritdoc />
    public partial class CuisineRelationsConfig : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.AlterColumn<string>(
                name: "Name",
                table: "Cuisines",
                type: "character varying(50)",
                maxLength: 50,
                nullable: false,
                oldClrType: typeof(string),
                oldType: "text");

            migrationBuilder.CreateIndex(
                name: "IX_Cuisines_Name",
                table: "Cuisines",
                column: "Name");
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropIndex(
                name: "IX_Cuisines_Name",
                table: "Cuisines");

            migrationBuilder.AlterColumn<string>(
                name: "Name",
                table: "Cuisines",
                type: "text",
                nullable: false,
                oldClrType: typeof(string),
                oldType: "character varying(50)",
                oldMaxLength: 50);
        }
    }
}
