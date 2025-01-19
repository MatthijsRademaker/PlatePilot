using Common.Domain;
using Common.Events;
using Common.Recipe;
using Domain;
using MediatR;
using Microsoft.EntityFrameworkCore;

namespace Infrastructure.CommandHandlers;

public class RecipeCommandHandler(IEventBus eventBus, RecipeContext recipeContext)
    : IRequestHandler<CreateRecipeCommand, Recipe>
{
    public async Task<Recipe> Handle(
        CreateRecipeCommand command,
        CancellationToken cancellationToken
    )
    {
        var mainIgredient =
            await recipeContext.Ingredients.FindAsync(command.MainIngredientId, cancellationToken)
            ?? throw new IngredientNotFoundException(command.MainIngredientId);

        var cuisine =
            await recipeContext.Cuisines.FindAsync(command.CuisineId, cancellationToken)
            ?? throw new CuisineNotFoundException(command.CuisineId);

        var ingredients = await recipeContext
            .Ingredients.Where(i => command.IngredientIds.Contains(i.Id))
            .ToListAsync(cancellationToken);

        if (ingredients.Count != command.IngredientIds.Count)
        {
            throw new IngredientNotFoundException(
                command.IngredientIds.Except(ingredients.Select(i => i.Id)).First()
            );
        }

        var recipe = new Recipe
        {
            Id = Guid.NewGuid(),
            Name = command.Name,
            Description = command.Description,
            PrepTime = command.PrepTime,
            CookTime = command.CookTime,
            MainIngredient = mainIgredient,
            Cuisine = cuisine,
            Ingredients = ingredients,
            Directions = command.Directions,
            Metadata = new Metadata { PublishedDate = DateTime.UtcNow },
            NutritionalInfo = new NutritionalInfo(),
        };

        recipe.Metadata.SearchVector = recipe.ToVector();

        recipeContext.Recipes.Add(recipe);
        await recipeContext.SaveChangesAsync(cancellationToken);

        await eventBus.PublishAsync(new RecipeCreatedEvent(RecipeDTO.FromRecipe(recipe)));
        return recipe;
    }
}

public class UpdateRecipeCommandHandler(IEventBus eventBus, RecipeContext recipeContext)
    : IRequestHandler<UpdateRecipeCommand, Recipe>
{
    public async Task<Recipe> Handle(
        UpdateRecipeCommand request,
        CancellationToken cancellationToken
    )
    {
        recipeContext.Recipes.Update(request.Recipe);
        await recipeContext.SaveChangesAsync(cancellationToken);

        await eventBus.PublishAsync(new RecipeUpdatedEvent(request.Recipe.Id));

        return request.Recipe;
    }
}
