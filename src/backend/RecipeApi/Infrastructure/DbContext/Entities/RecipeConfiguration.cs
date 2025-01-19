using System.Security.Principal;
using Common.Domain;
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

        builder.Ignore(ci => ci.Allergies);

        builder.HasOne(ci => ci.MainIngredient).WithMany();

        builder.HasOne(ci => ci.Cuisine).WithMany();

        builder.HasMany(ci => ci.Ingredients).WithMany();

        builder
            .OwnsOne(ci => ci.Metadata)
            .Property(mt => mt.SearchVector)
            .HasColumnType("Vector(128)");

        builder.OwnsOne(ci => ci.NutritionalInfo);

        builder.HasIndex(ci => ci.Name);
    }
}
