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
) : IEventBus
{
    private readonly ServiceBusSender _sender = serviceBusClient.CreateSender(
        serviceBusOptions.Value.TopicName
    );

    private JsonSerializerOptions jsonSerializerOptions =
        new() { PropertyNameCaseInsensitive = true, IncludeFields = true };

    public async Task PublishAsync<TEvent>(TEvent @event)
        where TEvent : IEvent
    {
        var message = new ServiceBusMessage(JsonSerializer.Serialize(@event, jsonSerializerOptions))
        {
            Subject = @event.GetType().Name,
        };

        await _sender.SendMessageAsync(message);
    }
}

    
