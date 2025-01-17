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
        api.MealPlannerV1();

        // TODO add ingredient and cuisine filter endpoints here.

        return api;
    }

    [ProducesResponseType<ProblemDetails>(
        StatusCodes.Status400BadRequest,
        "application/problem+json"
    )]
    public static async Task<Ok<RecipeResponse>> getRecipeById(
        [AsParameters] RecipeDependencies recipeDependencies,
        int id
    )
    {
        var item = await recipeDependencies.RecipeRepository.GetRecipeAsync(id);
        return TypedResults.Ok(RecipeResponse.FromRecipe(item));
    }

    [ProducesResponseType<ProblemDetails>(
        StatusCodes.Status400BadRequest,
        "application/problem+json"
    )]
    public static async Task<Ok<IEnumerable<RecipeResponse>>> getAllRecipes(
        [AsParameters] RecipeDependencies recipeDependencies,
        [FromQuery] int pageIndex,
        [FromQuery] int pageSize
    )
    {
        var items = await recipeDependencies.RecipeRepository.GetRecipesAsync(pageIndex * pageSize, pageSize);
        return TypedResults.Ok(items.Select(RecipeResponse.FromRecipe));
    }
}

public class RecipeResponse
{
    public int Id { get; set; }
    public string Name { get; set; }
    public string Description { get; set; }
    public string PrepTime { get; set; }
    public string CookTime { get; set; }
    public Ingredient MainIngredient { get; set; }
    public Cuisine Cuisine { get; set; }
    public ICollection<Ingredient> Ingredients { get; set; }
    public ICollection<string> Directions { get; set; }

    public static RecipeResponse FromRecipe(Domain.Recipe recipe)
    {
        return new RecipeResponse
        {
            Id = recipe.Id,
            Name = recipe.Name,
            Description = recipe.Description,
            PrepTime = recipe.PrepTime,
            CookTime = recipe.CookTime,
            MainIngredient = recipe.MainIngredient,
            Cuisine = recipe.Cuisine,
            Ingredients = recipe.Ingredients,
            Directions = recipe.Directions,
        };
    }
}
