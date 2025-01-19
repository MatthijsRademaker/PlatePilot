using Common.Domain;
using Domain;
using MediatR;
using Microsoft.EntityFrameworkCore;
using Pgvector.EntityFrameworkCore;

namespace Infrastructure.QueryHandlers;

public class RecipeQueryHandler(RecipeContext recipeContext) : IRequestHandler<RecipeQuery, Recipe>
{
    public Task<Recipe> Handle(RecipeQuery request, CancellationToken cancellationToken)
    {
        return recipeContext
            .GetRecipesWithIncludes()
            .FirstOrDefaultAsync(r => r.Id == request.Id, cancellationToken: cancellationToken)
            .ContinueWith(t => t.Result ?? throw new RecipeNotFoundException(request.Id));
    }
}

public class RecipesQueryHandler(RecipeContext recipeContext)
    : IRequestHandler<RecipesQuery, IEnumerable<Recipe>>
{
    public async Task<IEnumerable<Recipe>> Handle(
        RecipesQuery request,
        CancellationToken cancellationToken
    )
    {
        return await recipeContext
            .GetRecipesWithIncludes()
            .Skip(request.Skip)
            .Take(request.Take)
            .ToListAsync(cancellationToken);
    }
}

public class SearchSimilarRecipesQueryHandler(RecipeContext recipeContext)
    : IRequestHandler<SearchSimilarRecipesQuery, IEnumerable<Recipe>>
{
    public async Task<IEnumerable<Recipe>> Handle(
        SearchSimilarRecipesQuery request,
        CancellationToken cancellationToken
    )
    {
        return await recipeContext
            .GetRecipesWithIncludes()
            .Select(r => new
            {
                Recipe = r,
                Similarity = r.Metadata.SearchVector.CosineDistance(
                    request.Recipe.Metadata.SearchVector
                ),
            })
            .OrderByDescending(r => r.Similarity)
            .Take(request.Amount)
            .Select(r => r.Recipe)
            .ToListAsync(cancellationToken);
    }
}

public class RecipesByCuisineQueryHandler(RecipeContext recipeContext)
    : IRequestHandler<RecipesByCuisineQuery, IEnumerable<Recipe>>
{
    public async Task<IEnumerable<Recipe>> Handle(
        RecipesByCuisineQuery request,
        CancellationToken cancellationToken
    )
    {
        return await recipeContext
            .GetRecipesWithIncludes()
            .Where(r => r.Cuisine.Id == request.EntityId)
            .ToListAsync(cancellationToken);
    }
}

public class RecipesByIngredientQueryHandler(RecipeContext recipeContext)
    : IRequestHandler<RecipesByIngredientQuery, IEnumerable<Recipe>>
{
    public async Task<IEnumerable<Recipe>> Handle(
        RecipesByIngredientQuery request,
        CancellationToken cancellationToken
    )
    {
        return await recipeContext
            .GetRecipesWithIncludes()
            .Where(r =>
                r.Ingredients.Any(i => i.Id == request.EntityId)
                || r.MainIngredient.Id == request.EntityId
            )
            .ToListAsync(cancellationToken);
    }
}

public class RecipesByAllergyQueryHandler(RecipeContext recipeContext)
    : IRequestHandler<RecipesByAllergyQuery, IEnumerable<Recipe>>
{
    public async Task<IEnumerable<Recipe>> Handle(
        RecipesByAllergyQuery request,
        CancellationToken cancellationToken
    )
    {
        return await recipeContext
            .GetRecipesWithIncludes()
            .Where(r => r.Ingredients.Any(i => i.Allergies.Any(a => a.Id == request.EntityId)))
            .ToListAsync(cancellationToken);
    }
}
