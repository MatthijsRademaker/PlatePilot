using Domain;
using Microsoft.AspNetCore.Http.HttpResults;
using Microsoft.AspNetCore.Mvc;

namespace Application.Endpoints.V1;

public static class SuggestionEndpoint
{
    public static void MapSuggestionV1(this RouteGroupBuilder endpoints)
    {
        endpoints.MapPost("/suggest", suggestRecipes);
    }

    [ProducesResponseType<ProblemDetails>(
        StatusCodes.Status400BadRequest,
        "application/problem+json"
    )]
    private static async Task<Ok<IEnumerable<Domain.Recipe>>> suggestRecipes(
        [AsParameters] RecipeDependencies recipeDependencies,
        [FromBody] SuggestionsRequest suggestionsRequest
    )
    {
        var alreadySelectedRecipes = await Task.WhenAll(
            suggestionsRequest.AlreadySelectedRecipeIds.Select(async id =>
                await recipeDependencies.RecipeRepository.GetRecipeAsync(id)
            )
        );

        var items = await recipeDependencies.RecipeSuggestor.SuggestRecipesAsync(
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

        return TypedResults.Ok(items);
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
