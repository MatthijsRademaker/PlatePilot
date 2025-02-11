﻿using System.Text.Json;
using Azure.Messaging.ServiceBus;
using AzureServiceBusEventBus.Abstractions;
using Common.Events;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;

namespace AzureServiceBusEventBus;

public class EventProcessor(
    ILogger<EventProcessor> logger,
    IServiceScopeFactory serviceScopeFactory,
    ServiceBusClient serviceBusClient,
    IOptions<ServiceBusOptions> serviceBusOptions
) :  IHostedService
{
    private ServiceBusProcessor _processor= serviceBusClient.CreateProcessor(
        serviceBusOptions.Value.TopicName,
        serviceBusOptions.Value.SubscriptionName,
        new ServiceBusProcessorOptions()
    );

    private JsonSerializerOptions jsonSerializerOptions =
        new() { PropertyNameCaseInsensitive = true, IncludeFields = true };
    
    public async Task StartAsync(CancellationToken cancellationToken)
    {
        _processor.ProcessMessageAsync += async (args) =>
        {
            logger.LogInformation("Received message: {Message}", args.Message.Body.ToString());

            // Get the type based on the subject name
            var eventType =
                Type.GetType($"Common.Events.{args.Message.Subject}, Common")
                ?? throw new InvalidOperationException($"Type {args.Message.Subject} not found");
            
            logger.LogInformation("Event type: {eventType}", eventType);
            var @event = (IEvent)
                JsonSerializer.Deserialize(
                    args.Message.Body.ToString(),
                    eventType,
                    jsonSerializerOptions
                )!;
            
            logger.LogInformation("Parsed event to: {Event}", @event);  
            
            using var scope = serviceScopeFactory.CreateScope();

            // Already made for multiple event handlers
            var handlerType = typeof(IEventHandler<>).MakeGenericType(@event!.GetType());
            var handler = scope.ServiceProvider.GetService(handlerType);

            if (handler != null)
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
