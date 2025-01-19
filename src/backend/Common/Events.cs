using Common.Recipe;

namespace Common.Events;

public interface IEventBus
{
    public Task PublishAsync<TEvent>(TEvent @event) where TEvent : IEvent;
}

public interface IEventHandler<in TEvent>
{
    Task Handle(TEvent @event);
}

public interface IEvent
{
    public Guid Id { get; }

    public DateTime OccurredOn { get; }

    public Guid AggregateId { get; }
}

public record RecipeCreatedEvent(RecipeDTO Recipe) : IEvent
{
    public Guid AggregateId { get; } = Recipe.Id;
    
    public DateTime OccurredOn { get; } = DateTime.UtcNow;

    public Guid Id { get; } = Guid.NewGuid();
}

public record RecipeUpdatedEvent(Guid AggregateId) : IEvent
{
    public DateTime OccurredOn { get; } = DateTime.UtcNow;

    public Guid Id { get; } = Guid.NewGuid();
}
