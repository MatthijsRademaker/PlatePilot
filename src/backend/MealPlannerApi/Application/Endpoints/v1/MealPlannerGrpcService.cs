using Application.Endpoints.V1.Protos;
using Domain;
using Grpc.Core;

namespace Application.Endpoints.V1;

public class MealPlannerGrpcService(IMealPlanner mealPlanner)
    : MealPlannerService.MealPlannerServiceBase
{
    public override async Task<SuggestionsResponse> SuggestRecipes(
        SuggestionsRequest request,
        ServerCallContext context
    )
    {
        var items = await mealPlanner.SuggestMealsAsync(
            request.Amount,
            new SuggestionConstraints()
            {
                ConstraintsPerDay =
                [
                    .. request.DailyConstraints.Select(d =>
                    {
                        var constraints = new List<IConstraint>();
                        var ingredientConstraints = d.IngredientConstraints.Select(i =>
                        {
                            return new Domain.IngredientConstraint()
                            {
                                EntityId = Guid.Parse(i.EntityId),
                            };
                        });

                        var cuisineConstraints = d.CuisineConstraints.Select(c =>
                        {
                            return new Domain.CuisineConstraint()
                            {
                                EntityId = Guid.Parse(c.EntityId),
                            };
                        });
                        constraints = [.. ingredientConstraints, .. cuisineConstraints];
                        return constraints;
                    }),
                ],
            },
            request.AlreadySelectedRecipeIds.Select(Guid.Parse)
        );

        return new SuggestionsResponse { RecipeIds = { items.Select(i => i.ToString()) } };
    }
}
