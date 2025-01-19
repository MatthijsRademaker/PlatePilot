using Pgvector;

namespace Common.Recipe;

public static class VectorExtensions
{
    public static Vector ToVector(this Domain.Recipe recipe)
    {
        // Combine text features into vector
        var text = $"{recipe.Name} {recipe.Description} {recipe.MainIngredient?.Name}";
        return GenerateVector(text);
    }

    private static Vector GenerateVector(string text)
    {
        // Simple TF-IDF implementation
        // In production, consider using a proper embedding model
        var words = text.ToLower().Split(' ');
        var vector = new float[128]; // 128-dimensional vector

        for (var i = 0; i < words.Length && i < vector.Length; i++)
        {
            vector[i] = words[i].GetHashCode() % 100 / 100f;
        }

        return new(vector);
    }
}
