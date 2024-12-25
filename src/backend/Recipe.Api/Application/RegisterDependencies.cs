using Domain;
using Infrastructure;
using Microsoft.EntityFrameworkCore;

namespace Application
{
    public static class RegisterDependencies
    {
        public static void AddApplicationServices(this IHostApplicationBuilder builder)
        {
            builder.AddNpgsqlDbContext<RecipeContext>(
                "recipedb",
                configureDbContextOptions: dbContextOptionsBuilder =>
                {
                    dbContextOptionsBuilder.UseNpgsql(builder =>
                    {
                        builder.EnableRetryOnFailure();
                        builder.MigrationsAssembly("Recipe.Application");
                        // TODO look into once AI
                        // builder.UseVector();
                    });
                }
            );

            builder.EnrichNpgsqlDbContext<RecipeContext>();

            builder.Services.AddScoped<IRecipeRepository, RecipeRepository>();
        }
    }
}
