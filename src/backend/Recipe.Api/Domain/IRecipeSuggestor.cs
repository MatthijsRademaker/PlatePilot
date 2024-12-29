namespace Domain;

public interface IMealPlanner
{
    Task<IEnumerable<Recipe>> SuggestMealsAsync(
        int amountToSuggest,
        SuggestionConstraints constraints,
        IEnumerable<Recipe> alreadySelectedRecipes
    );
}

public class SuggestionConstraints
{
    public List<List<IConstraint>> ConstraintsPerDay { get; set; }
}

public class AllergiesConstraint : IConstraint
{
    public int EntityId { get; set; }

    public bool Matches(Recipe r)
    {
        return r.Allergies.Any(a => a.Id == EntityId);
    }
}

public class CuisineConstraint : IConstraint
{
    public int EntityId { get; set; }

    public bool Matches(Recipe r)
    {
        return r.Cuisine.Id == EntityId;
    }
}

public class IngredientConstraint : IConstraint
{
    public int EntityId { get; set; }

    public bool Matches(Recipe r)
    {
        return r.Ingredients.Any(i => i.Id == EntityId) || r.MainIngredient.Id == EntityId;
    }
}

public interface IConstraint
{
    public int EntityId { get; set; }

    bool Matches(Recipe r);
}
