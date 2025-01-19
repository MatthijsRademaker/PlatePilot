using MediatR;

namespace Domain;

public interface IQuery<T> : IRequest<T>;

public record RecipeQuery(Guid Id) : IQuery<Recipe>;

public record RecipesQuery : IQuery<IEnumerable<Recipe>>
{
    public int Skip { get; }
    public int Take { get; }

    public RecipesQuery(int pageIndex, int pageSize)
    {
        Skip = (pageIndex - 1) * pageSize;
        Take = pageSize;
    }
};

public record SearchSimilarRecipesQuery(Recipe Recipe, int Amount) : IQuery<IEnumerable<Recipe>>;

public record RecipesByCuisineQuery(Guid EntityId) : IQuery<IEnumerable<Recipe>>;

public record RecipesByIngredientQuery(Guid EntityId) : IQuery<IEnumerable<Recipe>>;

public record RecipesByAllergyQuery(Guid EntityId) : IQuery<IEnumerable<Recipe>>;

// TODO tags?


[Serializable]
public class RecipeNotFoundException(Guid id) : DomainException($"Recipe not found with id: {id}");

[Serializable]
public class RecipeAlreadyExistsException(string name)
    : DomainException($"Recipe already exists with name: {name}");

[Serializable]
public class IngredientNotFoundException(Guid id)
    : DomainException($"Ingredient not found with id: {id}");

[Serializable]
public class CuisineNotFoundException(Guid id)
    : DomainException($"Cuisine not found with id: {id}");
