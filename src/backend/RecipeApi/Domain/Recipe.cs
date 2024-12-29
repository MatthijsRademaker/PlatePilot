using Pgvector;

namespace Domain;

public class Recipe
{
    public int Id { get; set; }
    public string Name { get; set; }
    public string Description { get; set; }

    public string PrepTime { get; set; }

    public string CookTime { get; set; }

    public Ingredient MainIngredient { get; set; }

    public Cuisine Cuisine { get; set; }

    public ICollection<Ingredient> Ingredients { get; set; }

    public ICollection<string> Directions { get; set; }

    public Vector SearchVector { get; set; }
}

public class Allergy
{
    public int Id { get; set; }
    public string Name { get; set; }
}

public class Ingredient
{
    public int Id { get; set; }
    public string Name { get; set; }
    public string Quantity { get; set; }

    public ICollection<Allergy> Allergies { get; set; }
}

public class Cuisine
{
    public int Id { get; set; }
    public string Name { get; set; }
}
