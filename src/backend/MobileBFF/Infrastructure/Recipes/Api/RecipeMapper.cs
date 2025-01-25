using AutoMapper;
using Common.Domain;
using Infrastructure.RecipeApi.Client.Models;
using Cuisine = Infrastructure.RecipeApi.Client.Models.Cuisine;
using CuisineApi = Infrastructure.RecipeApi.Client.Models.Cuisine;
using Ingredient = Common.Domain.Ingredient;
using IngredientApi = Infrastructure.RecipeApi.Client.Models.Ingredient;

namespace MobileBFF.Infrastructure.Recipes;

public class RecipeProfile : Profile
{
    public RecipeProfile()
    {
        // TODO review mapping
        CreateMap<RecipeResponse, Recipe>()
            .ForMember(dest => dest.Id, opt => opt.MapFrom(src => src.Id ?? Guid.Empty))
            .ForMember(dest => dest.Name, opt => opt.MapFrom(src => src.Name))
            .ForMember(dest => dest.Description, opt => opt.MapFrom(src => src.Description))
            .ForMember(dest => dest.PrepTime, opt => opt.MapFrom(src => src.PrepTime))
            .ForMember(dest => dest.CookTime, opt => opt.MapFrom(src => src.CookTime))
            .ForMember(dest => dest.MainIngredient, opt => opt.MapFrom(src => src.MainIngredient))
            .ForMember(dest => dest.Cuisine, opt => opt.MapFrom(src => src.Cuisine))
            .ForMember(dest => dest.Ingredients, opt => opt.MapFrom(src => src.Ingredients))
            .ForMember(dest => dest.Directions, opt => opt.MapFrom(src => src.Directions))
            .ForMember(
                dest => dest.NutritionalInfo,
                opt => opt.MapFrom(src => new NutritionalInfo())
            ) // Default empty
            .ForMember(
                dest => dest.Metadata,
                opt => opt.MapFrom(src => new Metadata { PublishedDate = DateTime.UtcNow })
            ); // Default with current date

        CreateMap<IngredientApi, Ingredient>()
            .ForMember(dest => dest.Id, opt => opt.MapFrom(src => src.Id ?? Guid.Empty))
            .ForMember(dest => dest.Name, opt => opt.MapFrom(src => src.Name))
            .ForMember(dest => dest.Quantity, opt => opt.Ignore()) // Not present in source
            .ForMember(dest => dest.Allergies, opt => opt.Ignore()); // Not present in source

        CreateMap<CuisineApi, Cuisine>()
            .ForMember(dest => dest.Id, opt => opt.MapFrom(src => src.Id ?? Guid.Empty))
            .ForMember(dest => dest.Name, opt => opt.MapFrom(src => src.Name));
    }
}
