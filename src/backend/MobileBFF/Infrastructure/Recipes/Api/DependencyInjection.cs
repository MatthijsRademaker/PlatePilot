namespace Infrastructure.RecipeApi.Api;

using Application.Endpoints.V1.Protos;
using Microsoft.Extensions.DependencyInjection;

public static class KiotaServiceCollectionExtensions
{
    public static IServiceCollection RegisterRecipeApi(this IServiceCollection services)
    {
        services.AddGrpcClient<RecipeService.RecipeServiceClient>(client =>
        {
            client.Address = new Uri("https://recipe-api");
        })
        .ConfigurePrimaryHttpMessageHandler(() =>
        {
            var handler = new HttpClientHandler();
            handler.ServerCertificateCustomValidationCallback =
                HttpClientHandler.DangerousAcceptAnyServerCertificateValidator;

            return handler;
        });;

        return services;
    }
}
