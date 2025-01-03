using Domain;
using Microsoft.AspNetCore.Http.HttpResults;
using Microsoft.AspNetCore.Mvc;

namespace Application.Endpoints.V1;

public static class MealPlannerEndpoint
{
    public static void MealPlannerV1(this RouteGroupBuilder endpoints)
    {
        endpoints.MapPost("/mealplanner", suggestRecipes);
    }

    [ProducesResponseType<ProblemDetails>(
        StatusCodes.Status400BadRequest,
        "application/problem+json"
    )]
    private static async Task<Ok<IEnumerable<RecipeResponse>>> suggestRecipes(
        [AsParameters] RecipeDependencies recipeDependencies,
        [FromBody] SuggestionsRequest suggestionsRequest
    )
    {
        var alreadySelectedRecipes = await Task.WhenAll(
            suggestionsRequest.AlreadySelectedRecipeIds.Select(async id =>
                await recipeDependencies.RecipeRepository.GetRecipeAsync(id)
            )
        );

        var items = await recipeDependencies.MealPlanner.SuggestMealsAsync(
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
            alreadySelectedRecipes
        );

        return TypedResults.Ok(items.Select(RecipeResponse.FromRecipe));
    }

    public record SuggestionsRequest(
        SuggestionConstraintsRequest Constraints,
        IEnumerable<int> AlreadySelectedRecipeIds,
        int Amount
    );

    public record SuggestionConstraintsRequest(
        List<List<IngredientConstraint>> IngredientConstraints,
        List<List<CuisineConstraint>> CuisineConstraints
    );
}
