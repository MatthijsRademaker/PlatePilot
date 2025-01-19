using Domain;
using Infrastructure.Entities;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Configuration;

namespace Infrastructure;

/// <remarks>
/// Add migrations using the following command inside the 'MealPlannerInfrastructure' project directory:
///
/// dotnet ef migrations add --startup-project ../Application --context RecipeContext [migration-name]
/// </remarks>
public class RecipeContext(DbContextOptions<RecipeContext> options) : DbContext(options)
{
    public DbSet<Recipe> Recipes { get; set; }

    // Add function mapping, note the not supported exception is thrown when the function is called client side
    [DbFunction("recipe_matches", Schema = "public")]
    public static bool RecipeMatches(Guid recipeId, string constraintType, Guid entityId) =>
        throw new NotSupportedException();

    protected override void OnModelCreating(ModelBuilder builder)
    {
        builder.HasPostgresExtension("vector");
        builder.ApplyConfiguration(new RecipeEntityTypeConfiguration());

        // Add the outbox table to this context
        // builder.UseIntegrationEventLogs();
    }

    public async Task MigrateAsync()
    {
        try
        {
            if (Database.IsNpgsql())
            {
                await Database.MigrateAsync();
            }
        }
        catch (Exception ex)
        {
            throw new Exception(
                $"An error occurred while migrating the database: {ex.Message}",
                ex
            );
        }
    }
}
