using Domain;
using Microsoft.AspNetCore.Http.HttpResults;
using Microsoft.AspNetCore.Mvc;

namespace Application.Endpoints.V1;

public class RecipeDependencies(
    IRecipeRepository RecipeRepository,
    IRecipeSuggestor RecipeSuggestor
)
{
    public IRecipeSuggestor RecipeSuggestor { get; } = RecipeSuggestor;

    public IRecipeRepository RecipeRepository { get; } = RecipeRepository;
}
