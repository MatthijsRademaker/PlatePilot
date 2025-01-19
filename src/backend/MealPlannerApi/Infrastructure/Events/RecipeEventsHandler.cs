using Common.Events;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;

namespace Infrastructure;

public class RecipeEventHandler(ILogger<RecipeEventHandler> logger) : IEventHandler<RecipeCreatedEvent>
{
    public Task Handle(RecipeCreatedEvent @event)
    {
        logger.LogInformation("Handling event: {Event}", @event);
        throw new NotImplementedException();
    }
}
