using Common.Domain;
using Domain;
using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;

namespace Infrastructure.Entities;

class CuisineEntityTypeConfiguration : IEntityTypeConfiguration<Cuisine>
{
    public void Configure(EntityTypeBuilder<Cuisine> builder)
    {
        builder.ToTable("Cuisines");

        builder.Property(ci => ci.Name).HasMaxLength(50);

        builder.HasIndex(ci => ci.Name);
    }
}
