using AzureServiceBusEventBus.Abstractions;
using Common.Events;
using Domain;
using Infrastructure;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;

namespace Application
{
    public static class DependencyInjection
    {
        public static void AddInfrastructure(this IHostApplicationBuilder builder)
        {
            builder.AddNpgsqlDbContext<RecipeContext>(
                "recipedb",
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

            builder.Services.AddMediatR(cfg =>
                cfg.RegisterServicesFromAssembly(typeof(DependencyInjection).Assembly)
            );
            
            builder.AddEventBus("recipe-api", false);
        }
    }
}
