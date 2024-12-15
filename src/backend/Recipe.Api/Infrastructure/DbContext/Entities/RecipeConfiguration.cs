using Domain;
using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;

namespace Infrastructure.Entities;

class RecipeEntityTypeConfiguration : IEntityTypeConfiguration<Recipe>
{
    public void Configure(EntityTypeBuilder<Recipe> builder)
    {
        builder.ToTable("Recipes");

        builder.Property(ci => ci.Name).HasMaxLength(50);

        builder.HasMany(ci => ci.Ingredients).WithMany();

        builder.HasIndex(ci => ci.Name);
    }
}
