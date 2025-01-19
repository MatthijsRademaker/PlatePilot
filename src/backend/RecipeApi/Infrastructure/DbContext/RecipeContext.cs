using Common.Domain;
using Domain;
using Infrastructure.Entities;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Configuration;

namespace Infrastructure;

/// <remarks>
/// Add migrations using the following command inside the 'Infrastructure' project directory:
///
/// dotnet ef migrations add --startup-project ../Application --context RecipeContext [migration-name]
/// </remarks>
public class RecipeContext(DbContextOptions<RecipeContext> options) : DbContext(options)
{
    public DbSet<Recipe> Recipes { get; set; }

    public DbSet<Ingredient> Ingredients { get; set; }

    public DbSet<Cuisine> Cuisines { get; set; }

    public DbSet<Allergy> Allergies { get; set; }

    protected override void OnModelCreating(ModelBuilder builder)
    {
        builder.HasPostgresExtension("vector");
        builder.ApplyConfiguration(new RecipeEntityTypeConfiguration());
        builder.ApplyConfiguration(new IngredientEntityTypeConfiguration());
        builder.ApplyConfiguration(new CuisineEntityTypeConfiguration());
        builder.ApplyConfiguration(new AllergyEntityTypeConfiguration());

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

    public IQueryable<Recipe> GetRecipesWithIncludes()
    {
        return Recipes
            .Include(r => r.Ingredients)
            .ThenInclude(i => i.Allergies)
            .Include(r => r.MainIngredient)
            .ThenInclude(i => i.Allergies)
            .Include(r => r.Cuisine);
    }
}
