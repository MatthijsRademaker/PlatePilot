﻿// <auto-generated />
using System.Collections.Generic;
using Infrastructure;
using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Infrastructure;
using Microsoft.EntityFrameworkCore.Migrations;
using Microsoft.EntityFrameworkCore.Storage.ValueConversion;
using Npgsql.EntityFrameworkCore.PostgreSQL.Metadata;
using Pgvector;

#nullable disable

namespace MealPlannerInfrastructure.Migrations
{
    [DbContext(typeof(RecipeContext))]
    [Migration("20250117205131_Initial")]
    partial class Initial
    {
        /// <inheritdoc />
        protected override void BuildTargetModel(ModelBuilder modelBuilder)
        {
#pragma warning disable 612, 618
            modelBuilder
                .HasAnnotation("ProductVersion", "9.0.0")
                .HasAnnotation("Relational:MaxIdentifierLength", 63);

            NpgsqlModelBuilderExtensions.HasPostgresExtension(modelBuilder, "vector");
            NpgsqlModelBuilderExtensions.UseIdentityByDefaultColumns(modelBuilder);

            modelBuilder.Entity("Domain.Recipe", b =>
                {
                    b.Property<int>("Id")
                        .ValueGeneratedOnAdd()
                        .HasColumnType("integer");

                    NpgsqlPropertyBuilderExtensions.UseIdentityByDefaultColumn(b.Property<int>("Id"));

                    b.PrimitiveCollection<List<int>>("AllergyIds")
                        .IsRequired()
                        .HasColumnType("integer[]");

                    b.Property<int>("CuisineId")
                        .HasColumnType("integer");

                    b.PrimitiveCollection<List<int>>("IngredientIds")
                        .IsRequired()
                        .HasColumnType("integer[]");

                    b.Property<int>("MainIngredientId")
                        .HasColumnType("integer");

                    b.Property<Vector>("SearchVector")
                        .IsRequired()
                        .HasColumnType("Vector(128)");

                    b.HasKey("Id");

                    b.HasIndex("Id")
                        .IsUnique();

                    b.HasIndex("SearchVector");

                    NpgsqlIndexBuilderExtensions.HasMethod(b.HasIndex("SearchVector"), "ivfflat");
                    NpgsqlIndexBuilderExtensions.HasOperators(b.HasIndex("SearchVector"), new[] { "vector_cosine_ops" });

                    b.ToTable("recipe_view", "materialized");

                    b.ToView(null, (string)null);
                });
#pragma warning restore 612, 618
        }
    }
}
