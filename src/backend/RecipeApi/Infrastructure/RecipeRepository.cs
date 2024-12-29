using Domain;
using Microsoft.EntityFrameworkCore;
using Pgvector.EntityFrameworkCore;

namespace Infrastructure;

public class RecipeRepository(RecipeContext recipeContext) : IRecipeRepository
{
    public async Task<Recipe> CreateRecipeAsync(Recipe recipe)
    {
        recipeContext.Recipes.Add(recipe);
        await recipeContext.SaveChangesAsync();

        return recipe;
    }

    public async Task DeleteRecipeAsync(int id)
    {
        var recipe = recipeContext.Recipes.FirstOrDefault(r => r.Id == id);
        if (recipe != null)
        {
            recipeContext.Recipes.Remove(recipe);
            await recipeContext.SaveChangesAsync();
        }
    }

    public Task<Recipe> GetRecipeAsync(int id)
    {
        var recipe = recipeContext.GetRecipesWithIncludes().FirstOrDefault(r => r.Id == id);
        return Task.FromResult(recipe ?? throw new RecipeNotFoundException(id));
    }

    public async Task<IEnumerable<Recipe>> GetRecipesAsync(int amount)
    {
        return await recipeContext
            .Recipes.Include(r => r.Ingredients)
            .Include(r => r.MainIngredient)
            .Include(r => r.Cuisine)
            .Take(amount)
            .ToListAsync();
    }

    public async Task<IEnumerable<Recipe>> GetRecipesByAllergyAsync(int entityId)
    {
        return await recipeContext
            .GetRecipesWithIncludes()
            .Where(r => r.Ingredients.Any(i => i.Allergies.Any(a => a.Id == entityId)))
            .ToListAsync();
    }

    public async Task<IEnumerable<Recipe>> GetRecipesByCuisineAsync(int entityId)
    {
        return await recipeContext
            .GetRecipesWithIncludes()
            .Where(r => r.Cuisine.Id == entityId)
            .ToListAsync();
    }

    public async Task<IEnumerable<Recipe>> GetRecipesByIngredientAsync(int entityId)
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
                Similarity = r.SearchVector.CosineDistance(recipe.SearchVector),
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
internal class RecipeNotFoundException(int id) : Exception($"Recipe not found with id: {id}");
