using Pgvector;

namespace Domain;

public class Recipe
{
    public Guid Id { get;  init; }
    public Vector SearchVector { get;  init; }
    public Guid CuisineId { get;  init; }
    public Guid MainIngredientId { get;  init; }

    public List<Guid> IngredientIds { get;  init; }

    public List<Guid> AllergyIds { get;  init; }

    public override bool Equals(object? obj)
    {
        return obj is Recipe recipe && Id == recipe.Id;
    }

    public override int GetHashCode()
    {
        return HashCode.Combine(Id);
    }
}
