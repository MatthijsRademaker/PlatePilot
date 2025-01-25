using AutoMapper;
using Common.Domain;
using Domain;
using MediatR;

namespace Infrastructure.QueryHandlers;

public class RecipeQueryHandler(RecipeApi.Client.RecipeApi recipeApi, IMapper mapper)
    : IRequestHandler<RecipeQuery, Recipe>
{
    public async Task<Recipe> Handle(RecipeQuery request, CancellationToken cancellationToken)
    {
        var recipe = await recipeApi
            .Api.Recipes[request.Id]
            .GetAsync(cancellationToken: cancellationToken);
        return mapper.Map<Recipe>(recipe);
    }
}
