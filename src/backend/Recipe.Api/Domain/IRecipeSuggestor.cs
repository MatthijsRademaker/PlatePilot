namespace Domain;

// TODO: Implement this interface in RecipeSuggestor.cs
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
    public List<IngredientConstraint> MainIngredientConstraints { get; set; }

    public List<CuisineConstraint> CuisineConstraints { get; set; }
}

public class CuisineConstraint
{
    public int CuisineId { get; set; }
    public int AmountToGenerate { get; set; }
}

public class IngredientConstraint
{
    public int IngredientId { get; set; }
    public int AmountToGenerate { get; set; }
}
