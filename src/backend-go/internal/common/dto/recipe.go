package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
	"github.com/platepilot/backend/internal/common/domain"
)

// RecipeDTO is a data transfer object for recipe events
type RecipeDTO struct {
	ID              uuid.UUID       `json:"id"`
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	PrepTime        string          `json:"prepTime"`
	CookTime        string          `json:"cookTime"`
	MainIngredient  IngredientDTO   `json:"mainIngredient"`
	Cuisine         CuisineDTO      `json:"cuisine"`
	Ingredients     []IngredientDTO `json:"ingredients"`
	Allergies       []AllergyDTO    `json:"allergies"`
	Directions      []string        `json:"directions"`
	NutritionalInfo NutritionalDTO  `json:"nutritionalInfo"`
	Metadata        MetadataDTO     `json:"metadata"`
}

// IngredientDTO is a data transfer object for ingredients
type IngredientDTO struct {
	ID        uuid.UUID    `json:"id"`
	Name      string       `json:"name"`
	Quantity  string       `json:"quantity"`
	Allergies []AllergyDTO `json:"allergies,omitempty"`
}

// CuisineDTO is a data transfer object for cuisines
type CuisineDTO struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// AllergyDTO is a data transfer object for allergies
type AllergyDTO struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// NutritionalDTO is a data transfer object for nutritional info
type NutritionalDTO struct {
	Calories int `json:"calories"`
}

// MetadataDTO is a data transfer object for recipe metadata
type MetadataDTO struct {
	SearchVector  []float32 `json:"searchVector"`
	ImageURL      string    `json:"imageUrl"`
	Tags          []string  `json:"tags"`
	PublishedDate time.Time `json:"publishedDate"`
}

// FromRecipe converts a domain Recipe to a RecipeDTO
func FromRecipe(r *domain.Recipe) RecipeDTO {
	ingredients := make([]IngredientDTO, len(r.Ingredients))
	for i, ing := range r.Ingredients {
		ingredients[i] = FromIngredient(&ing)
	}

	allergies := r.Allergies()
	allergyDTOs := make([]AllergyDTO, len(allergies))
	for i, a := range allergies {
		allergyDTOs[i] = FromAllergy(&a)
	}

	var mainIngredient IngredientDTO
	if r.MainIngredient != nil {
		mainIngredient = FromIngredient(r.MainIngredient)
	}

	var cuisine CuisineDTO
	if r.Cuisine != nil {
		cuisine = FromCuisine(r.Cuisine)
	}

	return RecipeDTO{
		ID:              r.ID,
		Name:            r.Name,
		Description:     r.Description,
		PrepTime:        r.PrepTime,
		CookTime:        r.CookTime,
		MainIngredient:  mainIngredient,
		Cuisine:         cuisine,
		Ingredients:     ingredients,
		Allergies:       allergyDTOs,
		Directions:      r.Directions,
		NutritionalInfo: NutritionalDTO{Calories: r.NutritionalInfo.Calories},
		Metadata:        FromMetadata(&r.Metadata),
	}
}

// ToRecipe converts a RecipeDTO to a domain Recipe
func (d *RecipeDTO) ToRecipe() *domain.Recipe {
	ingredients := make([]domain.Ingredient, len(d.Ingredients))
	for i, ing := range d.Ingredients {
		ingredients[i] = *ing.ToIngredient()
	}

	mainIngredient := d.MainIngredient.ToIngredient()
	cuisine := d.Cuisine.ToCuisine()

	return &domain.Recipe{
		ID:              d.ID,
		Name:            d.Name,
		Description:     d.Description,
		PrepTime:        d.PrepTime,
		CookTime:        d.CookTime,
		MainIngredient:  mainIngredient,
		Cuisine:         cuisine,
		Ingredients:     ingredients,
		Directions:      d.Directions,
		NutritionalInfo: domain.NutritionalInfo{Calories: d.NutritionalInfo.Calories},
		Metadata:        *d.Metadata.ToMetadata(),
	}
}

// FromIngredient converts a domain Ingredient to an IngredientDTO
func FromIngredient(i *domain.Ingredient) IngredientDTO {
	allergies := make([]AllergyDTO, len(i.Allergies))
	for j, a := range i.Allergies {
		allergies[j] = FromAllergy(&a)
	}
	return IngredientDTO{
		ID:        i.ID,
		Name:      i.Name,
		Quantity:  i.Quantity,
		Allergies: allergies,
	}
}

// ToIngredient converts an IngredientDTO to a domain Ingredient
func (d *IngredientDTO) ToIngredient() *domain.Ingredient {
	allergies := make([]domain.Allergy, len(d.Allergies))
	for i, a := range d.Allergies {
		allergies[i] = *a.ToAllergy()
	}
	return &domain.Ingredient{
		ID:        d.ID,
		Name:      d.Name,
		Quantity:  d.Quantity,
		Allergies: allergies,
	}
}

// FromCuisine converts a domain Cuisine to a CuisineDTO
func FromCuisine(c *domain.Cuisine) CuisineDTO {
	return CuisineDTO{
		ID:   c.ID,
		Name: c.Name,
	}
}

// ToCuisine converts a CuisineDTO to a domain Cuisine
func (d *CuisineDTO) ToCuisine() *domain.Cuisine {
	return &domain.Cuisine{
		ID:   d.ID,
		Name: d.Name,
	}
}

// FromAllergy converts a domain Allergy to an AllergyDTO
func FromAllergy(a *domain.Allergy) AllergyDTO {
	return AllergyDTO{
		ID:   a.ID,
		Name: a.Name,
	}
}

// ToAllergy converts an AllergyDTO to a domain Allergy
func (d *AllergyDTO) ToAllergy() *domain.Allergy {
	return &domain.Allergy{
		ID:   d.ID,
		Name: d.Name,
	}
}

// FromMetadata converts domain Metadata to MetadataDTO
func FromMetadata(m *domain.Metadata) MetadataDTO {
	return MetadataDTO{
		SearchVector:  m.SearchVector.Slice(),
		ImageURL:      m.ImageURL,
		Tags:          m.Tags,
		PublishedDate: m.PublishedDate,
	}
}

// ToMetadata converts MetadataDTO to domain Metadata
func (d *MetadataDTO) ToMetadata() *domain.Metadata {
	return &domain.Metadata{
		SearchVector:  pgvectorFromSlice(d.SearchVector),
		ImageURL:      d.ImageURL,
		Tags:          d.Tags,
		PublishedDate: d.PublishedDate,
	}
}

// pgvectorFromSlice creates a pgvector.Vector from a float32 slice
func pgvectorFromSlice(v []float32) pgvector.Vector {
	return pgvector.NewVector(v)
}
