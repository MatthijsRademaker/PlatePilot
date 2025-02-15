using Application.Endpoints.V1.Protos;
using AutoMapper;
using Common.Domain;
using Domain;
using MediatR;

namespace Infrastructure.QueryHandlers;

public class RecipeQueryHandler(RecipeService.RecipeServiceClient recipeApi, IMapper mapper)
    : IRequestHandler<RecipeQuery, Recipe>
{
    public async Task<Recipe> Handle(RecipeQuery query, CancellationToken cancellationToken)
    {
        var recipe = await recipeApi.GetRecipeByIdAsync(
            new GetRecipeByIdRequest { RecipeId = query.Id.ToString() },
            cancellationToken: cancellationToken
        );
        return mapper.Map<Recipe>(recipe);
    }
}

public class RecipesQueryHandler(RecipeService.RecipeServiceClient recipeApi, IMapper mapper)
    : IRequestHandler<RecipesQuery, IEnumerable<Recipe>>
{
    public async Task<IEnumerable<Recipe>> Handle(
        RecipesQuery query,
        CancellationToken cancellationToken
    )
    {
        var recipes = await recipeApi.GetAllRecipesAsync(
            new GetAllRecipesRequest { PageIndex = query.Skip, PageSize = query.Take },
            cancellationToken: cancellationToken
        );

        return recipes.Recipes.Select(mapper.Map<Recipe>);
    }
}

public class SearchSimilarRecipesQueryHandler(
    RecipeService.RecipeServiceClient recipeApi,
    IMapper mapper
) : IRequestHandler<SearchSimilarRecipesQuery, IEnumerable<Recipe>>
{
    public async Task<IEnumerable<Recipe>> Handle(
        SearchSimilarRecipesQuery query,
        CancellationToken cancellationToken
    )
    {
        throw new NotImplementedException();
    }
}

public class RecipesByCuisineQueryHandler(
    RecipeService.RecipeServiceClient recipeApi,
    IMapper mapper
) : IRequestHandler<RecipesByCuisineQuery, IEnumerable<Recipe>>
{
    public async Task<IEnumerable<Recipe>> Handle(
        RecipesByCuisineQuery query,
        CancellationToken cancellationToken
    )
    {
        throw new NotImplementedException();
    }
}

public class RecipesByIngredientQueryHandler(
    RecipeService.RecipeServiceClient recipeApi,
    IMapper mapper
) : IRequestHandler<RecipesByIngredientQuery, IEnumerable<Recipe>>
{
    public async Task<IEnumerable<Recipe>> Handle(
        RecipesByIngredientQuery query,
        CancellationToken cancellationToken
    )
    {
        throw new NotImplementedException();
    }
}

public class RecipesByAllergyQueryHandler(
    RecipeService.RecipeServiceClient recipeApi,
    IMapper mapper
) : IRequestHandler<RecipesByAllergyQuery, IEnumerable<Recipe>>
{
    public async Task<IEnumerable<Recipe>> Handle(
        RecipesByAllergyQuery query,
        CancellationToken cancellationToken
    )
    {
        throw new NotImplementedException();
    }
}

public class CreateRecipeCommandHandler(RecipeService.RecipeServiceClient recipeApi, IMapper mapper)
    : IRequestHandler<CreateRecipeCommand, Recipe>
{
    public async Task<Recipe> Handle(
        CreateRecipeCommand command,
        CancellationToken cancellationToken
    )
    {
        var request = new CreateRecipeRequest
        {
            Name = command.Name,
            Description = command.Description,
            PrepTime = command.PrepTime,
            CookTime = command.CookTime,
            MainIngredientId = command.MainIngredientId.ToString(),
            CuisineId = command.CuisineId.ToString(),
            IngredientIds = { command.IngredientIds.Select(id => id.ToString()) },
        };

        request.Directions.AddRange(command.Directions);

        var recipe = await recipeApi.CreateRecipeAsync(
            request,
            cancellationToken: cancellationToken
        );

        return mapper.Map<Recipe>(recipe);
    }
}
