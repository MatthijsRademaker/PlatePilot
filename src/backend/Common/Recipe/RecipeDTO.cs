using Common.Domain;

namespace Common.Recipe;

public class RecipeDTO
{
    public Guid Id { get; set; }
    public string Name { get; set; }
    public string Description { get; set; }

    public string PrepTime { get; set; }

    public string CookTime { get; set; }

    public Ingredient MainIngredient { get; set; }

    public Cuisine Cuisine { get; set; }

    public ICollection<Ingredient> Ingredients { get; set; }

    public ICollection<Allergy> Allergies { get; set; }

    public ICollection<string> Directions { get; set; }

    public NutritionalInfo NutritionalInfo { get; set; }

    public MetadataDto Metadata { get; set; }

    public static RecipeDTO FromRecipe(Domain.Recipe recipe)
    {
        return new RecipeDTO()
        {
            Id = recipe.Id,
            Name = recipe.Name,
            Description = recipe.Description,
            PrepTime = recipe.PrepTime,
            CookTime = recipe.CookTime,
            MainIngredient = recipe.MainIngredient,
            Cuisine = recipe.Cuisine,
            Directions = recipe.Directions,
            NutritionalInfo = recipe.NutritionalInfo,
            Ingredients = recipe.Ingredients,
            Allergies = recipe.Allergies,
            Metadata = MetadataDto.FromMetadata(recipe.Metadata)
        };
    }
    
    
    public static Domain.Recipe ToRecipe(RecipeDTO recipeDto)
    {
        return new Domain.Recipe()
        {
            Id = recipeDto.Id,
            Name = recipeDto.Name,
            Description = recipeDto.Description,
            PrepTime = recipeDto.PrepTime,
            CookTime = recipeDto.CookTime,
            MainIngredient = recipeDto.MainIngredient,
            Cuisine = recipeDto.Cuisine,
            Directions = recipeDto.Directions,
            NutritionalInfo = recipeDto.NutritionalInfo,
            Ingredients = recipeDto.Ingredients,
            Metadata = MetadataDto.FromMetadataDto(recipeDto.Metadata),
        };
    }
}



public class MetadataDto
{
    public float[] SearchVector { get; set; }
    public string? ImageUrl { get; set; }
    public ICollection<string> Tags { get; set; } = [];
    public DateTime PublishedDate { get; set; }

    public static MetadataDto FromMetadata(Domain.Metadata metadata)
    {
        return new MetadataDto()
        {
            SearchVector = metadata.SearchVector.ToArray(),
            ImageUrl = metadata.ImageUrl,
            Tags = metadata.Tags,
            PublishedDate = metadata.PublishedDate,
        };
    }

    public static Metadata FromMetadataDto(MetadataDto metadataDto)
    {
        return new Metadata()
        {
            SearchVector = new (metadataDto.SearchVector),
            ImageUrl = metadataDto.ImageUrl,
            Tags = metadataDto.Tags,
            PublishedDate = metadataDto.PublishedDate,
        };
    }
}