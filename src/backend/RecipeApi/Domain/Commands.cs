using MediatR;

namespace Domain;

public interface ICommand<T> : IRequest<T>;

public record CreateRecipeCommand(
    string Name,
    string Description,
    string PrepTime,
    string CookTime,
    Guid MainIngredientId,
    Guid CuisineId,
    ICollection<Guid> IngredientIds,
    ICollection<string> Directions
) : ICommand<Recipe>;

public record UpdateRecipeCommand(Recipe Recipe) : ICommand<Recipe>;

public record DeleteRecipeCommand(Guid Id) : ICommand<Unit>;
