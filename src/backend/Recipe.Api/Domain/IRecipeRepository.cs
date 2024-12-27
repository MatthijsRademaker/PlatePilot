namespace Domain;

public interface IRecipeRepository
{
    public Task<IEnumerable<Recipe>> GetRecipesAsync(int amount);
    public Task<Recipe> GetRecipeAsync(int id);

    public Task<Recipe> CreateRecipeAsync(Recipe recipe);
    public Task<Recipe> UpdateRecipeAsync(Recipe recipe);
    public Task DeleteRecipeAsync(int id);
    Task<IEnumerable<Recipe>> GetRecipesByCuisineAsync(int entityId, int amount);
    Task<IEnumerable<Recipe>> GetRecipesByIngredientAsync(int entityId, int amount);
}
