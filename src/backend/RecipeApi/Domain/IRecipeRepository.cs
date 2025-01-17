namespace Domain;

// TODO split into multiple repositories
public interface IRecipeRepository
{
    public Task<IEnumerable<Recipe>> SearchSimilarRecipes(Recipe recipe, int amount);
    public Task<IEnumerable<Recipe>> GetRecipesAsync(int startIndex, int amount);
    public Task<Recipe> GetRecipeAsync(Guid Id);

    public Task<Recipe> CreateRecipeAsync(Recipe recipe);
    public Task<Recipe> UpdateRecipeAsync(Recipe recipe);
    public Task DeleteRecipeAsync(Guid Id);
    Task<IEnumerable<Recipe>> GetRecipesByCuisineAsync(Guid entityId);
    Task<IEnumerable<Recipe>> GetRecipesByIngredientAsync(Guid entityId);
    Task<IEnumerable<Recipe>> GetRecipesByAllergyAsync(Guid entityId);
}
