using Domain;
using Microsoft.EntityFrameworkCore;
using Pgvector.EntityFrameworkCore;

namespace Infrastructure;

// TODO rewrite to query pattern
public class RecipeRepository(RecipeContext recipeContext) : IRecipeRepository
{
    public async Task<Recipe> CreateRecipeAsync(Recipe recipe)
    {
        recipeContext.Recipes.Add(recipe);
        await recipeContext.SaveChangesAsync();

        return recipe;
    }

    public async Task DeleteRecipeAsync(Guid id)
    {
        var recipe = recipeContext.Recipes.FirstOrDefault(r => r.Id == id);
        if (recipe != null)
        {
            recipeContext.Recipes.Remove(recipe);
            await recipeContext.SaveChangesAsync();
        }
    }

    public Task<Recipe> GetRecipeAsync(Guid id)
    {
        var recipe = recipeContext.GetRecipesWithIncludes().FirstOrDefault(r => r.Id == id);
        return Task.FromResult(recipe ?? throw new RecipeNotFoundException(id));
    }

    public async Task<IEnumerable<Recipe>> GetRecipesAsync(int startIndex, int amount)
    {
        return await recipeContext
            .GetRecipesWithIncludes()
            .Skip(startIndex)
            .Take(amount)
            .ToListAsync();
    }

    public async Task<IEnumerable<Recipe>> GetRecipesByAllergyAsync(Guid entityId)
    {
        return await recipeContext
            .GetRecipesWithIncludes()
            .Where(r => r.Ingredients.Any(i => i.Allergies.Any(a => a.Id == entityId)))
            .ToListAsync();
    }

    public async Task<IEnumerable<Recipe>> GetRecipesByCuisineAsync(Guid entityId)
    {
        return await recipeContext
            .GetRecipesWithIncludes()
            .Where(r => r.Cuisine.Id == entityId)
            .ToListAsync();
    }

    public async Task<IEnumerable<Recipe>> GetRecipesByIngredientAsync(Guid entityId)
    {
        return await recipeContext
            .GetRecipesWithIncludes()
            .Where(r => r.Ingredients.Any(i => i.Id == entityId) || r.MainIngredient.Id == entityId)
            .ToListAsync();
    }

    public async Task<IEnumerable<Recipe>> SearchSimilarRecipes(Recipe recipe, int amount)
    {
        return await recipeContext
            .GetRecipesWithIncludes()
            .Select(r => new
            {
                Recipe = r,
                Similarity = r.Metadata.SearchVector.CosineDistance(recipe.Metadata.SearchVector),
            })
            .OrderByDescending(r => r.Similarity)
            .Take(amount)
            .Select(r => r.Recipe)
            .ToListAsync();
    }

    public async Task<Recipe> UpdateRecipeAsync(Recipe recipe)
    {
        recipeContext.Recipes.Update(recipe);
        await recipeContext.SaveChangesAsync();

        return recipe;
    }
}

[Serializable]
internal class RecipeNotFoundException(Guid id) : Exception($"Recipe not found with id: {id}");
