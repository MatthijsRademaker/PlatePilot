namespace Domain;

// TODO split into multiple repositories
public interface IRecipeRepository
{
    public Task<IEnumerable<Recipe>> SearchSimilarRecipes(Recipe recipe, int amount);
    public Task<IEnumerable<Recipe>> GetRecipesAsync(int startIndex, int amount);
    public Task<Recipe> GetRecipeAsync(int id);

    public Task<Recipe> CreateRecipeAsync(Recipe recipe);
    public Task<Recipe> UpdateRecipeAsync(Recipe recipe);
    public Task DeleteRecipeAsync(int id);
    Task<IEnumerable<Recipe>> GetRecipesByCuisineAsync(int entityId);
    Task<IEnumerable<Recipe>> GetRecipesByIngredientAsync(int entityId);
    Task<IEnumerable<Recipe>> GetRecipesByAllergyAsync(int entityId);
}
