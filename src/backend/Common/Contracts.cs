namespace Common.Contracts;

// TODO complete
public record RecipeDto(
    Guid Id,
    string Name,
    string Description,
    IEnumerable<IngredientDto> Ingredients
);

public record IngredientDto(Guid Id, string Name);
