namespace Infrastructure.RecipeApi.Api;

using Microsoft.Extensions.DependencyInjection;
using Microsoft.Kiota.Http.HttpClientLibrary;
using MobileBFF.Infrastructure.Recipes.Api;

public static class KiotaServiceCollectionExtensions
{
    public static IServiceCollection RegisterRecipeApi(this IServiceCollection services)
    {
        services.AddKiotaHandlers();
        // TODO correlationId
        services
            .AddHttpClient<RecipeApi.Client.RecipeApi>()
            .ConfigureHttpClient(client =>
            {
                client.BaseAddress = new Uri("http+https://recipeapi");
            })
            .AttachKiotaHandlers();

        services.AddTransient(sp => sp.GetRequiredService<RecipeApiClientFactory>().GetClient());
        return services;
    }

    public static IServiceCollection AddKiotaHandlers(this IServiceCollection services)
    {
        // Dynamically load the Kiota handlers from the Client Factory
        var kiotaHandlers = KiotaClientFactory.GetDefaultHandlerTypes();
        // And register them in the DI container
        foreach (var handler in kiotaHandlers)
        {
            services.AddTransient(handler);
        }

        return services;
    }

    public static IHttpClientBuilder AttachKiotaHandlers(this IHttpClientBuilder builder)
    {
        // Dynamically load the Kiota handlers from the Client Factory
        var kiotaHandlers = KiotaClientFactory.GetDefaultHandlerTypes();
        // And attach them to the http client builder
        foreach (var handler in kiotaHandlers)
        {
            builder.AddHttpMessageHandler(
                (sp) => (DelegatingHandler)sp.GetRequiredService(handler)
            );
        }

        return builder;
    }
}
