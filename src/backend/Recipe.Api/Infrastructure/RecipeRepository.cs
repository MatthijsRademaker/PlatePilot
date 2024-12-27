using Domain;
using Microsoft.EntityFrameworkCore;

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
        var recipe = recipeContext
            .Recipes.Include(r => r.Ingredients)
            .Include(r => r.MainIngredient)
            .Include(r => r.Cuisine)
            .FirstOrDefault(r => r.Id == id);
        return Task.FromResult(recipe ?? throw new RecipeNotFoundException(id));
    }

    public Task<IEnumerable<Recipe>> GetRecipesAsync(int amount)
    {
        return Task.FromResult(
            recipeContext
                .Recipes.Include(r => r.Ingredients)
                .Include(r => r.MainIngredient)
                .Include(r => r.Cuisine)
                .Take(amount)
                .AsEnumerable()
        );
    }

    public Task<IEnumerable<Recipe>> GetRecipesByCuisineAsync(int entityId, int amount)
    {
        return Task.FromResult(
            recipeContext
                .Recipes.Include(r => r.Ingredients)
                .Include(r => r.MainIngredient)
                .Include(r => r.Cuisine)
                .Where(r => r.Cuisine.Id == entityId)
                .Take(amount)
                .AsEnumerable()
        );
    }

    public Task<IEnumerable<Recipe>> GetRecipesByIngredientAsync(int entityId, int amount)
    {
        return Task.FromResult(
            recipeContext
                .Recipes.Include(r => r.Ingredients)
                .Include(r => r.MainIngredient)
                .Include(r => r.Cuisine)
                .Where(r =>
                    r.Ingredients.Any(i => i.Id == entityId) || r.MainIngredient.Id == entityId
                )
                .Take(amount)
                .AsEnumerable()
        );
    }

    public Task<Recipe> UpdateRecipeAsync(Recipe recipe)
    {
        recipeContext.Recipes.Update(recipe);
        recipeContext.SaveChangesAsync();

        return Task.FromResult(recipe);
    }
}

[Serializable]
internal class RecipeNotFoundException(int id) : Exception($"Recipe not found with id: {id}");
