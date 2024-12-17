using System.Text.Json;
using Domain;
using Infrastructure;
using Microsoft.EntityFrameworkCore;

namespace Application.DatabaseSeed;

public class DatabaseSeeder(RecipeContext context)
{
    public async Task SeedAsync()
    {
        if (!context.Recipes.Any())
        {
            var recipesJson = await File.ReadAllTextAsync("Data/recipes.json");
            var recipesData = JsonSerializer.Deserialize<RecipeData>(recipesJson);

            foreach (var recipe in recipesData.Recipes)
            {
                // Ensure main ingredient is tracked
                context.Ingredients.Add(recipe.MainIngredient);

                // Add all ingredients
                foreach (var ingredient in recipe.Ingredients)
                {
                    context.Ingredients.Add(ingredient);
                }

                context.Recipes.Add(recipe);
            }

            await context.SaveChangesAsync();
        }
    }
}

public class RecipeData
{
    public List<Recipe> Recipes { get; set; }
}
