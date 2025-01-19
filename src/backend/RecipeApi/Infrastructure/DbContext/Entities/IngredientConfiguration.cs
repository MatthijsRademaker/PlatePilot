using Common.Domain;
using Domain;
using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;

class IngredientEntityTypeConfiguration : IEntityTypeConfiguration<Ingredient>
{
    public void Configure(EntityTypeBuilder<Ingredient> builder)
    {
        builder.ToTable("Ingredients");

        builder.Property(ci => ci.Name).HasMaxLength(50);

        builder.HasMany(ci => ci.Allergies).WithMany();

        builder.HasIndex(ci => ci.Name);
    }
}
