using Domain;
using Microsoft.AspNetCore.Http.HttpResults;
using Microsoft.AspNetCore.Mvc;

namespace Application.Endpoints.V1;

public static class MealPlannerEndpoint
{
    public static IEndpointRouteBuilder MapMealPlannerV1(this IEndpointRouteBuilder endpoints)
    {
        endpoints.MapGroup("v1").MapPost("/plan-meal", suggestRecipes);
        return endpoints;
    }

    [ProducesResponseType<ProblemDetails>(
        StatusCodes.Status400BadRequest,
        "application/problem+json"
    )]
    private static async Task<Ok<IEnumerable<Guid>>> suggestRecipes(
        IMealPlanner mealPlanner,
        [FromBody] SuggestionsRequest suggestionsRequest
    )
    {
        var items = await mealPlanner.SuggestMealsAsync(
            suggestionsRequest.Amount,
            new SuggestionConstraints()
            {
                ConstraintsPerDay = suggestionsRequest
                    .Constraints.IngredientConstraints.Zip(
                        suggestionsRequest.Constraints.CuisineConstraints,
                        (ingredientConstraints, cuisineConstraints) =>
                        {
                            var constraints = new List<IConstraint>();
                            constraints = [.. ingredientConstraints, .. cuisineConstraints];
                            return constraints;
                        }
                    )
                    .ToList(),
            },
            suggestionsRequest.AlreadySelectedRecipeIds
        );

        return TypedResults.Ok(items);
    }

    public record SuggestionsRequest(
        SuggestionConstraintsRequest Constraints,
        IEnumerable<Guid> AlreadySelectedRecipeIds,
        int Amount
    );

    public record SuggestionConstraintsRequest(
        List<List<IngredientConstraint>> IngredientConstraints,
        List<List<CuisineConstraint>> CuisineConstraints
    );
}
