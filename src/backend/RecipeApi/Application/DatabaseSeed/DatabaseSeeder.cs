using System.Text.Json;
using Common.Events;
using Domain;
using Infrastructure;
using Microsoft.EntityFrameworkCore;
using RecipeApi.Infrastructure;

namespace Application.DatabaseSeed;

public class DatabaseSeeder(RecipeContext context)
// public class DatabaseSeeder(IEventBus eventBus, RecipeContext context)
{
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
                    recipe.MainIngredient.Id = Guid.NewGuid();
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
                    recipe.Cuisine.Id = Guid.NewGuid();
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
                        ingredient.Id = Guid.NewGuid();
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
                recipe.Metadata = new Metadata
                {
                    PublishedDate = DateTime.UtcNow,
                    SearchVector = recipe.ToVector(),
                };

                // Add recipe
                addedRecipes[recipe.Name] = recipe;
                context.Recipes.Add(recipe);

                // Publish event
                // await eventBus.PublishAsync(new RecipeCreated(recipe.Id));
            }

            await context.SaveChangesAsync();
        }
    }
}

public class RecipeData
{
    public List<Recipe> Recipes { get; set; }
}
