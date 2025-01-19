namespace Common.Events;

public interface IEventBus
{
    public Task PublishAsync(IEvent @event);
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

public record RecipeCreatedEvent(Guid AggregateId) : IEvent
{
    public DateTime OccurredOn { get; } = DateTime.UtcNow;

    public Guid Id { get; } = Guid.NewGuid();
}

public record RecipeUpdatedEvent(Guid AggregateId) : IEvent
{
    public DateTime OccurredOn { get; } = DateTime.UtcNow;

    public Guid Id { get; } = Guid.NewGuid();
}
