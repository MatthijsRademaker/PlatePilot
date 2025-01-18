using Common.Events;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;

namespace AzureServiceBusEventBus.Abstractions;

public static class Abstractions
{
    public const string TopicName = "recipe-events";

    // TODO if necessary, add a way to configure the topic name and multiple event handlers for different type of events/topics
    public static IHostApplicationBuilder AddEventBus(
        this IHostApplicationBuilder builder,
        string subscriptionName
    )
    {
        builder.AddAzureServiceBusClient("messaging");
        builder.Services.AddSingleton<IEventBus, EventBus>();
        builder.Services.AddHostedService<EventBus>();
        builder
            .Services.AddOptions<ServiceBusOptions>()
            .Configure<IHostEnvironment>(
                (options, environment) =>
                {
                    options.TopicName = TopicName;
                    options.SubscriptionName = subscriptionName;
                }
            );

        return builder;
    }
}

public class ServiceBusOptions
{
    public string TopicName { get; internal set; } = Abstractions.TopicName;
    public string SubscriptionName { get; internal set; }
}
