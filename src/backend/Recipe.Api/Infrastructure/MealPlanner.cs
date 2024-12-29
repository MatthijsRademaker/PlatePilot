using Domain;

namespace Infrastructure;

// TODO Maybe inject the RecipeContext instead of the repository
public class MealPlanner(IRecipeRepository recipeRepository) : IMealPlanner
{
    public async Task<IEnumerable<Recipe>> SuggestMealsAsync(
        int amountToSuggest,
        SuggestionConstraints constraints,
        IEnumerable<Recipe> alreadySelectedRecipes
    )
    {
        var result = new List<Recipe>();
        var selectedSet = new HashSet<Recipe>(alreadySelectedRecipes);
        var amountToGenerate = Math.Max(amountToSuggest, constraints.ConstraintsPerDay.Count);
        var stackToGenerate = new Stack<List<IConstraint>>(constraints.ConstraintsPerDay);

        while (stackToGenerate.Any())
        {
            var newConstraints = stackToGenerate.Pop();
            var recipes = await GetRecipesByConstraints(newConstraints);
            var bestSuitedRecipe = recipes
                .Where(r => !selectedSet.Contains(r))
                .OrderByDescending(r => CalculateDiversityScore(r, selectedSet))
                .FirstOrDefault();

            if (bestSuitedRecipe != null)
            {
                result.Add(bestSuitedRecipe);
                selectedSet.Add(bestSuitedRecipe);
            }

            if (result.Count >= amountToSuggest)
                break;
        }

        while (result.Count < amountToSuggest)
        {
            //  TODO should this get all recipes??
            var recipes = await recipeRepository.GetRecipesAsync(100);
            var candidates = recipes
                .Where(r => !selectedSet.Contains(r))
                .OrderByDescending(r => CalculateDiversityScore(r, selectedSet))
                .Take(amountToSuggest - result.Count);

            foreach (var recipe in candidates)
            {
                if (result.Count >= amountToSuggest)
                    break;
                result.Add(recipe);
            }
        }

        return result;
    }

    private async Task<IEnumerable<Recipe>> GetRecipesByConstraints(
        IEnumerable<IConstraint> constraints
    )
    {
        var recipes = await recipeRepository.GetRecipesAsync(100);
        return recipes.Where(r => constraints.All(c => c.Matches(r)));
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
            similarity += 0.25;

        if (a.MainIngredient == b.MainIngredient)
            similarity += 0.25;

        // Jaccard similarity
        if (a.Ingredients != null && b.Ingredients != null)
        {
            var commonIngredients = a.Ingredients.Intersect(b.Ingredients).Count();
            similarity +=
                (double)commonIngredients / Math.Min(a.Ingredients.Count, b.Ingredients.Count);
        }

        return similarity;
    }
}
