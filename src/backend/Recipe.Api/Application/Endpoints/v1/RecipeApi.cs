using System.Collections;
using System.ComponentModel;
using Domain;
using Microsoft.AspNetCore.Http.HttpResults;
using Microsoft.AspNetCore.Mvc;

namespace Application.Endpoints.V1;

public static class RecipeApi
{
    public static IEndpointRouteBuilder MapRecipeV1(this IEndpointRouteBuilder endpoints)
    {
        var api = endpoints.MapGroup("/api/recipes").HasApiVersion(1.0);

        api.MapGet("/{id:int}", getRecipeById);
        api.MapGet("/all", getAllRecipes);

        return api;
    }

    [ProducesResponseType<ProblemDetails>(
        StatusCodes.Status400BadRequest,
        "application/problem+json"
    )]
    public static async Task<Ok<Domain.Recipe>> getRecipeById(
        [AsParameters] RecipeDependencies recipeDependencies,
        int id
    )
    {
        var items = await recipeDependencies.RecipeRepository.GetRecipeAsync(id);
        return TypedResults.Ok(items);
    }

    [ProducesResponseType<ProblemDetails>(
        StatusCodes.Status400BadRequest,
        "application/problem+json"
    )]
    public static async Task<Ok<IEnumerable<Domain.Recipe>>> getAllRecipes(
        [AsParameters] RecipeDependencies recipeDependencies
    )
    {
        var items = await recipeDependencies.RecipeRepository.GetRecipesAsync();
        return TypedResults.Ok(items);
    }
}

public class RecipeDependencies(IRecipeRepository RecipeRepository)
{
    public IRecipeRepository RecipeRepository { get; } = RecipeRepository;
}
