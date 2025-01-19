using Common.Events;
using Domain;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;

namespace Infrastructure;

public class RecipeEventHandler(ILogger<RecipeEventHandler> logger, RecipeContext recipeContext) : IEventHandler<RecipeCreatedEvent>
{
    public async Task Handle(RecipeCreatedEvent @event)
    {
        logger.LogInformation("Handling event: {Event}", @event);
        recipeContext.Recipes.Add(new Recipe()
        {
            Id = @event.AggregateId,
            CuisineId = @event.Recipe.Cuisine.Id,
            IngredientIds = @event.Recipe.Ingredients.Select(i => i.Id).ToList(),
            MainIngredientId = @event.Recipe.MainIngredient.Id,
            AllergyIds = @event.Recipe.Allergies.Select(i => i.Id).ToList(),    
            SearchVector = new (@event.Recipe.Metadata.SearchVector)
            
        });
        await recipeContext.SaveChangesAsync();
    }
}
