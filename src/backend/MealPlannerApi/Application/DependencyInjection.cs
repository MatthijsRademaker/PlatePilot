using AzureServiceBusEventBus.Abstractions;
using Common.Events;
using Domain;
using Infrastructure;
using Microsoft.EntityFrameworkCore;

namespace Application;

public static class DependencyInjection
{
    public static void AddApplicationServices(this IHostApplicationBuilder builder)
    {
        builder.Services.AddGrpc();
    }
}
