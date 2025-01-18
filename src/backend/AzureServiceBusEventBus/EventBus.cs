using System.Text.Json;
using Azure.Messaging.ServiceBus;
using AzureServiceBusEventBus.Abstractions;
using Common.Events;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;

namespace AzureServiceBusEventBus;

public class EventBus(
    ILogger<EventBus> logger,
    IServiceScopeFactory serviceScopeFactory,
    ServiceBusClient serviceBusClient,
    IOptions<ServiceBusOptions> serviceBusOptions
) : IEventBus, IHostedService
{
    private ServiceBusProcessor _processor;

    private ServiceBusSender _sender;

    private readonly string topicName = serviceBusOptions.Value.TopicName;

    private readonly string subscriptionName = serviceBusOptions.Value.SubscriptionName;

    public async Task PublishAsync(IEvent @event)
    {
        var message = new ServiceBusMessage(JsonSerializer.Serialize(@event))
        {
            Subject = @event.GetType().Name,
        };

        await _sender.SendMessageAsync(message);
    }

    public async Task StartAsync(CancellationToken cancellationToken)
    {
        _processor = serviceBusClient.CreateProcessor(
            topicName,
            subscriptionName,
            new ServiceBusProcessorOptions()
        );
        _sender = serviceBusClient.CreateSender(topicName);

        _processor.ProcessMessageAsync += async (args) =>
        {
            var @event = JsonSerializer.Deserialize<IEvent>(args.Message.Body.ToString());

            using var scope = serviceScopeFactory.CreateScope();

            // Already made for multiple event handlers
            var handlerType = typeof(IEventHandler).MakeGenericType(@event!.GetType());
            var handler = scope.ServiceProvider.GetRequiredService(handlerType);

            await (Task)
                handler.GetType().GetMethod("Handle").Invoke(handler, new object[] { @event });
        };

        _processor.ProcessErrorAsync += (args) =>
        {
            logger.LogError(args.Exception, "Error processing message");
            return Task.CompletedTask;
        };

        await _processor.StartProcessingAsync();
    }

    public Task StopAsync(CancellationToken cancellationToken)
    {
        return _processor.StopProcessingAsync();
    }
}
