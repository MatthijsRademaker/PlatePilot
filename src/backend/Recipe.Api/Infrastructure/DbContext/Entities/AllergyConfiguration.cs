using Domain;
using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;

namespace Infrastructure.Entities;

class AllergyEntityTypeConfiguration : IEntityTypeConfiguration<Allergy>
{
    public void Configure(EntityTypeBuilder<Allergy> builder)
    {
        builder.ToTable("Allergies");

        builder.Property(ci => ci.Name).HasMaxLength(50);

        builder.HasIndex(ci => ci.Name);
    }
}
