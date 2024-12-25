namespace Domain;

public interface IRecipeSuggestor
{
    Task<IEnumerable<Recipe>> SuggestRecipesAsync(
        int amountToSuggest,
        SuggestionConstraints constraints,
        IEnumerable<Recipe> alreadySelectedRecipes
    );
}

public class SuggestionConstraints
{
    public Dictionary<Ingredient, int> MainIngredientConstraints { get; set; }

    public Dictionary<Cuisine, int> CuisineConstraints { get; set; }
}
