using Domain;
using Microsoft.AspNetCore.Http.HttpResults;
using Microsoft.AspNetCore.Mvc;

namespace Application.Endpoints.V1;

public class RecipeDependencies(IRecipeRepository RecipeRepository, IMealPlanner MealPlanner)
{
    public IMealPlanner MealPlanner { get; } = MealPlanner;

    public IRecipeRepository RecipeRepository { get; } = RecipeRepository;
}
