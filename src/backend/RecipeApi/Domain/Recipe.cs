using Pgvector;

namespace Domain;

public class Recipe
{
    public Guid Id { get; set; }
    public string Name { get; set; }
    public string Description { get; set; }

    public string PrepTime { get; set; }

    public string CookTime { get; set; }

    public Ingredient MainIngredient { get; set; }

    public Cuisine Cuisine { get; set; }

    public ICollection<Ingredient> Ingredients { get; set; }

    public ICollection<string> Directions { get; set; }

    public NutritionalInfo NutritionalInfo { get; set; }

    public Metadata Metadata { get; set; }
}

public class Metadata
{
    public Vector SearchVector { get; set; }
    public string? ImageUrl { get; set; }
    public ICollection<string> Tags { get; set; } = [];
    public DateTime PublishedDate { get; set; }
}

public class NutritionalInfo
{
    public int Calories { get; set; }
}

public class Allergy
{
    public Guid Id { get; set; }
    public string Name { get; set; }
}

public class Ingredient
{
    public Guid Id { get; set; }
    public string Name { get; set; }
    public string Quantity { get; set; }

    public ICollection<Allergy> Allergies { get; set; }
}

public class Cuisine
{
    public Guid Id { get; set; }
    public string Name { get; set; }
}
