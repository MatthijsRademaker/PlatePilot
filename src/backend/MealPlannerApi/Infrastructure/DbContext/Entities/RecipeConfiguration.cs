using Domain;
using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;

namespace Infrastructure.Entities;

class RecipeEntityTypeConfiguration : IEntityTypeConfiguration<Recipe>
{
    public void Configure(EntityTypeBuilder<Recipe> builder)
    {
        builder.ToTable("recipe_view", schema: "materialized").ToView(null);

        builder.HasKey(ci => ci.Id);

        builder.Property(ci => ci.SearchVector).HasColumnType("Vector(128)");
        builder
            .HasIndex(e => e.SearchVector)
            .HasMethod("ivfflat")
            .HasOperators("vector_cosine_ops");

        builder.HasIndex(ci => ci.Id).IsUnique();
    }
}
