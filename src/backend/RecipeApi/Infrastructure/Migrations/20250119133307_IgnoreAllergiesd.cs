using System;
using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace RecipeInfrastructure.Migrations
{
    /// <inheritdoc />
    public partial class IgnoreAllergiesd : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropForeignKey(
                name: "FK_Allergies_Recipes_RecipeId",
                table: "Allergies");

            migrationBuilder.DropIndex(
                name: "IX_Allergies_RecipeId",
                table: "Allergies");

            migrationBuilder.DropColumn(
                name: "RecipeId",
                table: "Allergies");
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.AddColumn<Guid>(
                name: "RecipeId",
                table: "Allergies",
                type: "uuid",
                nullable: true);

            migrationBuilder.CreateIndex(
                name: "IX_Allergies_RecipeId",
                table: "Allergies",
                column: "RecipeId");

            migrationBuilder.AddForeignKey(
                name: "FK_Allergies_Recipes_RecipeId",
                table: "Allergies",
                column: "RecipeId",
                principalTable: "Recipes",
                principalColumn: "Id");
        }
    }
}
