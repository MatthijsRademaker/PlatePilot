using Pgvector;

namespace Domain;

public class Recipe
{
    public Guid Id { get; private set; }
    public Vector SearchVector { get; private set; }
    public int CuisineId { get; private set; }
    public int MainIngredientId { get; private set; }

    public List<int> IngredientIds { get; private set; }

    public List<int> AllergyIds { get; private set; }

    public override bool Equals(object? obj)
    {
        return obj is Recipe recipe && Id == recipe.Id;
    }

    public override int GetHashCode()
    {
        return HashCode.Combine(Id);
    }
}
