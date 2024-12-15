namespace Domain;

public interface IRecipeRepository
{
    public Task<IEnumerable<Recipe>> GetRecipesAsync();
    public Task<Recipe> GetRecipeAsync(int id);

    public Task<Recipe> CreateRecipeAsync(Recipe recipe);
    public Task<Recipe> UpdateRecipeAsync(Recipe recipe);
    public Task DeleteRecipeAsync(int id);
}
