namespace Domain;

public interface IMealPlanner
{
    Task<IEnumerable<Guid>> SuggestMealsAsync(
        int amountToSuggest,
        SuggestionConstraints constraints,
        IEnumerable<Guid> alreadySelectedRecipeIds
    );
}

public class SuggestionConstraints
{
    public List<List<IConstraint>> ConstraintsPerDay { get; set; }
}

public class AllergiesConstraint : IConstraint
{
    public Guid EntityId { get; set; }
}

public class CuisineConstraint : IConstraint
{
    public Guid EntityId { get; set; }
}

public class IngredientConstraint : IConstraint
{
    public Guid EntityId { get; set; }
}

public interface IConstraint
{
    public Guid EntityId { get; set; }
}
