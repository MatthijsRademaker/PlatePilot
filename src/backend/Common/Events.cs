namespace Common.Events;

public interface IEvent
{
    public Guid Id { get; }

    public DateTime OccurredOn { get; }

    public int AggregateId { get; }
}

public record RecipeCreated(int AggregateId, string Name, string Description) : IEvent
{
    public DateTime OccurredOn { get; } = DateTime.UtcNow;

    public Guid Id { get; } = 1;
}
