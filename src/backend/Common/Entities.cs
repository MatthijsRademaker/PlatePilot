namespace Common.Entities;

public abstract class BaseEntity
{
    public Guid Id { get; set; }
    public DateTime CreatedAt { get; set; }
    public DateTime? UpdatedAt { get; set; }
}

public abstract class ValueObject
{
    protected abstract IEnumerable<object> GetEqualityComponents();
}
