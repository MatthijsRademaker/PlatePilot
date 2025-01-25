using System.Reflection;
using AzureServiceBusEventBus.Abstractions;
using Common.Events;
using Domain;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;

namespace Infrastructure;

public static class DependencyInjection
{
    
    public static void AddInfrastructure(this IHostApplicationBuilder builder)
    {
        builder.AddNpgsqlDbContext<RecipeContext>(
            "mealplannerdb",
            configureDbContextOptions: dbContextOptionsBuilder =>
            {
                dbContextOptionsBuilder.UseNpgsql(builder =>
                {
                    builder.EnableRetryOnFailure();
                    builder.UseVector();
                });
            }
        );

        builder.EnrichNpgsqlDbContext<RecipeContext>();

        builder.Services.AddScoped<IMealPlanner, MealPlanner>();
        builder.Services.AddScoped<IEventHandler<RecipeCreatedEvent>, RecipeEventHandler>();
        
        if (Assembly.GetEntryAssembly()?.GetName().Name != "GetDocument.Insider")
        {
            builder.AddEventBus("meal-planner-api", true);
        }
    }
}
