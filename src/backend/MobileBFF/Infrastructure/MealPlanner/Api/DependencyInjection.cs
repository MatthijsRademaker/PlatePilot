namespace Infrastructure.RecipeApi.Api;

using Application.Endpoints.V1.Protos;
using Microsoft.Extensions.DependencyInjection;

public static class MealPlannerApiServiceCollectionExtensions
{
    public static IServiceCollection RegisterMealPlannerApi(this IServiceCollection services)
    {
        services.AddGrpcClient<MealPlannerService.MealPlannerServiceClient>(client =>
        {
            client.Address = new Uri("http+https://meal-planner-api");
        });

        return services;
    }
}
