using Domain;

namespace Infrastructure;

public class RecipeSuggestor(IRecipeRepository recipeRepository) : IRecipeSuggestor
{
    public async Task<IEnumerable<Recipe>> SuggestRecipesAsync(
        int amountToSuggest,
        SuggestionConstraints constraints,
        IEnumerable<Recipe> alreadySelectedRecipes
    )
    {
        var result = new List<Recipe>();
        var selectedSet = new HashSet<Recipe>(alreadySelectedRecipes);

        foreach (var dayConstraints in constraints.ConstraintsPerDay)
        {
            // Get all recipes that match any of the day's constraints
            var matchingRecipes = await Task.WhenAll(
                dayConstraints.Select(async constraint =>
                {
                    var recipes = await GetRecipesByConstraint(constraint);
                    return (constraint, recipes);
                })
            );

            // Filter out already selected recipes and sort by least used cuisine/ingredients
            var candidates = matchingRecipes
                .SelectMany(x => x.recipes)
                .Distinct()
                .Where(r => !selectedSet.Contains(r))
                .OrderByDescending(r => CalculateDiversityScore(r, selectedSet))
                .Take(amountToSuggest);

            foreach (var recipe in candidates)
            {
                if (result.Count >= amountToSuggest)
                    break;
                result.Add(recipe);
                selectedSet.Add(recipe);
            }
        }

        return result;
    }

    private async Task<IEnumerable<Recipe>> GetRecipesByConstraint(IConstraint constraint)
    {
        return constraint switch
        {
            CuisineConstraint c => await recipeRepository.GetRecipesByCuisineAsync(
                c.EntityId,
                c.AmountToGenerate
            ),
            IngredientConstraint i => await recipeRepository.GetRecipesByIngredientAsync(
                i.EntityId,
                i.AmountToGenerate
            ),
            _ => Enumerable.Empty<Recipe>(),
        };
    }

    private double CalculateDiversityScore(Recipe candidate, HashSet<Recipe> selectedRecipes)
    {
        if (!selectedRecipes.Any())
            return 1.0;

        // Calculate how different this recipe is from already selected ones
        // Based on cuisine and ingredients - higher score means more diverse
        var similarityScore = selectedRecipes
            .Select(r => CalculateSimilarity(candidate, r))
            .Average();

        return 1.0 - similarityScore;
    }

    private double CalculateSimilarity(Recipe a, Recipe b)
    {
        // Implement similarity calculation based on:
        // - Same cuisine (higher similarity)
        // - Common ingredients (higher similarity)
        // Returns value between 0 (completely different) and 1 (very similar)
        // This is a simplified example - you might want to adjust the weights
        double similarity = 0;

        if (a.Cuisine == b.Cuisine)
            similarity += 0.5;

        // TODO Add ingredient similarity calculation here
        if (a.Ingredients != null && b.Ingredients != null)
        {
            var commonIngredients = a.Ingredients.Intersect(b.Ingredients).Count();
            similarity +=
                (double)commonIngredients / Math.Min(a.Ingredients.Count, b.Ingredients.Count);
        }

        return similarity;
    }
}
