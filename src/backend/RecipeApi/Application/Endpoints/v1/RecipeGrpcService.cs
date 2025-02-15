using Application.Endpoints.V1.Protos;
using Common.Domain;
using Domain;
using Grpc.Core;
using MediatR;

namespace Application.Endpoints.V1;

public class RecipeGrpcService : RecipeService.RecipeServiceBase
{
    private readonly IMediator _mediator;

    public RecipeGrpcService(IMediator mediator)
    {
        _mediator = mediator;
    }

    public override async Task<RecipeResponse> GetRecipeById(
        GetRecipeByIdRequest request,
        ServerCallContext context
    )
    {
        var item = await _mediator.Send(new RecipeQuery(Guid.Parse(request.RecipeId)));
        return MapToGrpcResponse(item);
    }

    public override async Task<GetAllRecipesResponse> GetAllRecipes(
        GetAllRecipesRequest request,
        ServerCallContext context
    )
    {
        var items = await _mediator.Send(new RecipesQuery(request.PageIndex, request.PageSize));

        return new GetAllRecipesResponse { Recipes = { items.Select(MapToGrpcResponse) } };
    }

    public override async Task<RecipeResponse> CreateRecipe(
        CreateRecipeRequest request,
        ServerCallContext context
    )
    {
        var item = await _mediator.Send(
            new CreateRecipeCommand(
                Name: request.Name,
                Description: request.Description,
                PrepTime: request.PrepTime,
                CookTime: request.CookTime,
                MainIngredientId: Guid.Parse(request.MainIngredientId),
                CuisineId: Guid.Parse(request.CuisineId),
                IngredientIds: [.. request.IngredientIds.Select(Guid.Parse)],
                Directions: request.Directions
            )
        );

        return MapToGrpcResponse(item);
    }

    private static RecipeResponse MapToGrpcResponse(Recipe recipe)
    {
        return new RecipeResponse
        {
            Id = recipe.Id.ToString(),
            Name = recipe.Name,
            Description = recipe.Description,
            PrepTime = recipe.PrepTime,
            CookTime = recipe.CookTime,
            MainIngredient = new Protos.Ingredient
            {
                Id = recipe.MainIngredient.Id.ToString(),
                Name = recipe.MainIngredient.Name,
            },
            Cuisine = new Protos.Cuisine
            {
                Id = recipe.Cuisine.Id.ToString(),
                Name = recipe.Cuisine.Name,
            },
            Ingredients =
            {
                recipe.Ingredients.Select(i => new Protos.Ingredient
                {
                    Id = i.Id.ToString(),
                    Name = i.Name,
                }),
            },
            Directions = { recipe.Directions },
        };
    }
}
