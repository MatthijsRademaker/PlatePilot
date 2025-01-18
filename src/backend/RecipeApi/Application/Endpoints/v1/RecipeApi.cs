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
        api.MapPost("/create", createRecipe);

        // TODO add ingredient and cuisine filter endpoints here.

        return api;
    }

    [ProducesResponseType<ProblemDetails>(
        StatusCodes.Status400BadRequest,
        "application/problem+json"
    )]
    public static async Task<Ok<RecipeResponse>> getRecipeById(
        IRecipeRepository recipeRepository,
        Guid id
    )
    {
        var item = await recipeRepository.GetRecipeAsync(id);
        return TypedResults.Ok(RecipeResponse.FromRecipe(item));
    }

    [ProducesResponseType<ProblemDetails>(
        StatusCodes.Status400BadRequest,
        "application/problem+json"
    )]
    public static async Task<Ok<IEnumerable<RecipeResponse>>> getAllRecipes(
        IRecipeRepository recipeRepository,
        [FromQuery] int pageIndex,
        [FromQuery] int pageSize
    )
    {
        var items = await recipeRepository.GetRecipesAsync(pageIndex * pageSize, pageSize);
        return TypedResults.Ok(items.Select(RecipeResponse.FromRecipe));
    }

    [ProducesResponseType<ProblemDetails>(
        StatusCodes.Status400BadRequest,
        "application/problem+json"
    )]
    public static async Task<Ok<RecipeResponse>> createRecipe(
        IRecipeRepository recipeRepository,
        CreateRecipeRequest request
    )
    {
        var item = new Domain.Recipe
        {
            Name = request.Name,
            Description = request.Description,
            PrepTime = request.PrepTime,
            CookTime = request.CookTime,
            MainIngredient = request.MainIngredient,
            Cuisine = request.Cuisine,
            Ingredients = request.Ingredients,
            Directions = request.Directions,
        };

        await recipeRepository.CreateRecipeAsync(item);
        return TypedResults.Ok(RecipeResponse.FromRecipe(item));
    }
}

public class CreateRecipeRequest
{
    public string Name { get; set; }
    public string Description { get; set; }
    public string PrepTime { get; set; }
    public string CookTime { get; set; }
    public Ingredient MainIngredient { get; set; }
    public Cuisine Cuisine { get; set; }
    public ICollection<Ingredient> Ingredients { get; set; }
    public ICollection<string> Directions { get; set; }
}

public class RecipeResponse
{
    public Guid Id { get; set; }
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
