using Common.Domain;
using Domain;
using MediatR;
using Microsoft.AspNetCore.Http.HttpResults;
using Microsoft.AspNetCore.Mvc;

namespace Application.Endpoints.V1;

public static class RecipeApiEndpoint
{
    public static IEndpointRouteBuilder MapRecipeApiV1(this IEndpointRouteBuilder endpoints)
    {
        var group = endpoints.MapGroup("v1/recipe");

        group.MapGet("/{id}", (IMediator mediator, Guid id) => mediator.Send(new RecipeQuery(id)));

        group.MapGet(
            "/all",
            (IMediator mediator, [FromQuery] int pageIndex, [FromQuery] int pageSize) =>
                mediator.Send(new RecipesQuery(pageIndex, pageSize))
        );

        group.MapGet(
            "/similar",
            (IMediator mediator, [FromQuery] Recipe recipe, [FromQuery] int amount) =>
                mediator.Send(new SearchSimilarRecipesQuery(recipe, amount))
        );

        group.MapGet(
            "/cuisine/{id}",
            (IMediator mediator, Guid id) => mediator.Send(new RecipesByCuisineQuery(id))
        );

        group.MapGet(
            "/ingredient/{id}",
            (IMediator mediator, Guid id) => mediator.Send(new RecipesByIngredientQuery(id))
        );

        group.MapGet(
            "/allergy/{id}",
            (IMediator mediator, Guid id) => mediator.Send(new RecipesByAllergyQuery(id))
        );

        group.MapPost(
            "/create",
            (IMediator mediator, [FromBody] CreateRecipeCommand command) => mediator.Send(command)
        );

        return endpoints;
    }
}
