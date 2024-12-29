using System.Text.Json;
using Domain;
using Infrastructure;
using Microsoft.EntityFrameworkCore;
using RecipeApi.Infrastructure;

namespace Application.DatabaseSeed;

public class DatabaseSeeder
{
    private readonly RecipeContext context;

    public DatabaseSeeder(RecipeContext context)
    {
        this.context = context;
    }

    public async Task SeedAsync()
    {
        if (!context.Recipes.Any())
        {
            var recipesJson = await File.ReadAllTextAsync("DatabaseSeed/recipes.json");
            var recipesData = JsonSerializer.Deserialize<RecipeData>(
                recipesJson,
                new JsonSerializerOptions { PropertyNameCaseInsensitive = true }
            );

            if (recipesData == null)
            {
                throw new Exception("Deserialization resulted in null");
            }

            var addedIngredients = new Dictionary<string, Ingredient>();
            var addedCuisines = new Dictionary<string, Cuisine>();
            var addedRecipes = new Dictionary<string, Recipe>();
            var nextIngredientId = 1;
            var nextCuisineId = 1;

            foreach (var recipe in recipesData.Recipes)
            {
                // Skip if recipe with same name already exists
                if (addedRecipes.ContainsKey(recipe.Name))
                {
                    continue;
                }

                // Handle main ingredient
                if (!addedIngredients.ContainsKey(recipe.MainIngredient.Name))
                {
                    recipe.MainIngredient.Id = nextIngredientId++;
                    addedIngredients[recipe.MainIngredient.Name] = recipe.MainIngredient;
                    context.Ingredients.Add(recipe.MainIngredient);
                }
                else
                {
                    recipe.MainIngredient = addedIngredients[recipe.MainIngredient.Name];
                }

                // Handle cuisine
                if (!addedCuisines.ContainsKey(recipe.Cuisine.Name))
                {
                    recipe.Cuisine.Id = nextCuisineId++;
                    addedCuisines[recipe.Cuisine.Name] = recipe.Cuisine;
                    context.Cuisines.Add(recipe.Cuisine);
                }
                else
                {
                    recipe.Cuisine = addedCuisines[recipe.Cuisine.Name];
                }

                // Handle recipe ingredients
                var uniqueIngredients = new List<Ingredient>();
                foreach (var ingredient in recipe.Ingredients)
                {
                    if (!addedIngredients.ContainsKey(ingredient.Name))
                    {
                        ingredient.Id = nextIngredientId++;
                        addedIngredients[ingredient.Name] = ingredient;
                        context.Ingredients.Add(ingredient);
                        uniqueIngredients.Add(ingredient);
                    }
                    else
                    {
                        uniqueIngredients.Add(addedIngredients[ingredient.Name]);
                    }
                }
                recipe.Ingredients = uniqueIngredients;
                recipe.SearchVector = recipe.ToVector();

                // Add recipe
                addedRecipes[recipe.Name] = recipe;
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
