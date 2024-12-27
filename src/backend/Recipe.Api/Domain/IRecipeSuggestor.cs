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
    public List<List<IConstraint>> ConstraintsPerDay { get; set; }
}

public class CuisineConstraint : IConstraint
{
    public int EntityId { get; set; }
    public int AmountToGenerate { get; set; }
}

public class IngredientConstraint : IConstraint
{
    public int EntityId { get; set; }
    public int AmountToGenerate { get; set; }
}

public interface IConstraint
{
    public int EntityId { get; set; }
    public int AmountToGenerate { get; set; }
}
