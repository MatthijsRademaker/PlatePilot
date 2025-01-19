using System.Collections;
using System.ComponentModel;
using Domain;
using MediatR;
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
    public static async Task<Ok<RecipeResponse>> getRecipeById(IMediator mediator, Guid id)
    {
        var item = await mediator.Send(new RecipeQuery(id));
        return TypedResults.Ok(RecipeResponse.FromRecipe(item));
    }

    [ProducesResponseType<ProblemDetails>(
        StatusCodes.Status400BadRequest,
        "application/problem+json"
    )]
    public static async Task<Ok<IEnumerable<RecipeResponse>>> getAllRecipes(
        IMediator mediator,
        [FromQuery] int pageIndex,
        [FromQuery] int pageSize
    )
    {
        var items = await mediator.Send(new RecipesQuery(pageIndex, pageSize));
        return TypedResults.Ok(items.Select(RecipeResponse.FromRecipe));
    }

    [ProducesResponseType<ProblemDetails>(
        StatusCodes.Status400BadRequest,
        "application/problem+json"
    )]
    public static async Task<Ok<RecipeResponse>> createRecipe(
        IMediator mediator,
        CreateRecipeRequest request
    )
    {
        var item = await mediator.Send(
            new CreateRecipeCommand(
                Name: request.Name,
                Description: request.Description,
                PrepTime: request.PrepTime,
                CookTime: request.CookTime,
                MainIngredientId: request.MainIngredient,
                CuisineId: request.Cuisine,
                IngredientIds: request.Ingredients,
                Directions: request.Directions
            )
        );
        return TypedResults.Ok(RecipeResponse.FromRecipe(item));
    }
}

// TODO add metadata and nutritional info
public class CreateRecipeRequest
{
    public string Name { get; set; }
    public string Description { get; set; }
    public string PrepTime { get; set; }
    public string CookTime { get; set; }
    public Guid MainIngredient { get; set; }
    public Guid Cuisine { get; set; }
    public ICollection<Guid> Ingredients { get; set; }
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
