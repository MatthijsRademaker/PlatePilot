using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using RabbitMQ.Client;

namespace Infrastructure;

public class RecipeEventHandler(
    ILogger<RecipeEventHandler> logger,
    IConnection connection,
    IServiceScopeFactory serviceScopeFactory
) : BackgroundService
{
    protected override async Task ExecuteAsync(CancellationToken ct)
    {
        // TODO look into eshop example
    }
}
