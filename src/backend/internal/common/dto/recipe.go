package dto

import (
	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
	"github.com/platepilot/backend/internal/common/domain"
)

// RecipeDTO is a data transfer object for recipe events (read-model projection).
type RecipeDTO struct {
	ID               uuid.UUID            `json:"id"`
	UserID           uuid.UUID            `json:"userId"`
	Name             string               `json:"name"`
	Description      string               `json:"description"`
	PrepTimeMinutes  int                  `json:"prepTimeMinutes"`
	CookTimeMinutes  int                  `json:"cookTimeMinutes"`
	TotalTimeMinutes int                  `json:"totalTimeMinutes"`
	Servings         int                  `json:"servings"`
	YieldQuantity    *float64             `json:"yieldQuantity,omitempty"`
	YieldUnit        string               `json:"yieldUnit,omitempty"`
	MainIngredient   IngredientDTO        `json:"mainIngredient"`
	Cuisine          CuisineDTO           `json:"cuisine"`
	IngredientLines  []IngredientLineDTO  `json:"ingredientLines"`
	Allergies        []AllergyDTO         `json:"allergies"`
	Tags             []string             `json:"tags"`
	ImageURL         string               `json:"imageUrl"`
	Nutrition        RecipeNutritionDTO   `json:"nutrition"`
	SearchVector     []float32            `json:"searchVector"`
}

// IngredientDTO is a data transfer object for ingredients.
type IngredientDTO struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// IngredientLineDTO is a data transfer object for recipe ingredient line items.
type IngredientLineDTO struct {
	Ingredient    IngredientDTO `json:"ingredient"`
	QuantityValue *float64      `json:"quantityValue,omitempty"`
	QuantityText  string        `json:"quantityText,omitempty"`
	Unit          string        `json:"unit,omitempty"`
	IsOptional    bool          `json:"isOptional"`
	Note          string        `json:"note,omitempty"`
	SortOrder     int           `json:"sortOrder"`
}

// CuisineDTO is a data transfer object for cuisines.
type CuisineDTO struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// AllergyDTO is a data transfer object for allergies.
type AllergyDTO struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// RecipeNutritionDTO is a data transfer object for recipe nutrition.
type RecipeNutritionDTO struct {
	CaloriesTotal      int     `json:"caloriesTotal"`
	CaloriesPerServing int     `json:"caloriesPerServing"`
	ProteinG           float64 `json:"proteinG"`
	CarbsG             float64 `json:"carbsG"`
	FatG               float64 `json:"fatG"`
	FiberG             float64 `json:"fiberG"`
	SugarG             float64 `json:"sugarG"`
	SodiumMg           float64 `json:"sodiumMg"`
}

// FromRecipe converts a domain Recipe to a RecipeDTO.
func FromRecipe(r *domain.Recipe) RecipeDTO {
	tags := r.Tags
	if tags == nil {
		tags = []string{}
	}

	lines := make([]IngredientLineDTO, len(r.IngredientLines))
	for i, line := range r.IngredientLines {
		lines[i] = IngredientLineDTO{
			Ingredient:    FromIngredient(&line.Ingredient),
			QuantityValue: line.QuantityValue,
			QuantityText:  line.QuantityText,
			Unit:          line.Unit,
			IsOptional:    line.IsOptional,
			Note:          line.Note,
			SortOrder:     line.SortOrder,
		}
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
		ID:               r.ID,
		UserID:           r.UserID,
		Name:             r.Name,
		Description:      r.Description,
		PrepTimeMinutes:  r.PrepTimeMinutes,
		CookTimeMinutes:  r.CookTimeMinutes,
		TotalTimeMinutes: r.TotalTimeMinutes,
		Servings:         r.Servings,
		YieldQuantity:    r.YieldQuantity,
		YieldUnit:        r.YieldUnit,
		MainIngredient:   mainIngredient,
		Cuisine:          cuisine,
		IngredientLines:  lines,
		Allergies:        allergyDTOs,
		Tags:             tags,
		ImageURL:         r.ImageURL,
		Nutrition:        FromNutrition(&r.Nutrition),
		SearchVector:     r.SearchVector.Slice(),
	}
}

// ToRecipe converts a RecipeDTO to a domain Recipe.
func (d *RecipeDTO) ToRecipe() *domain.Recipe {
	lines := make([]domain.RecipeIngredientLine, len(d.IngredientLines))
	for i, line := range d.IngredientLines {
		lines[i] = *line.ToIngredientLine()
	}

	mainIngredient := d.MainIngredient.ToIngredient()
	cuisine := d.Cuisine.ToCuisine()

	return &domain.Recipe{
		ID:               d.ID,
		UserID:           d.UserID,
		Name:             d.Name,
		Description:      d.Description,
		PrepTimeMinutes:  d.PrepTimeMinutes,
		CookTimeMinutes:  d.CookTimeMinutes,
		TotalTimeMinutes: d.TotalTimeMinutes,
		Servings:         d.Servings,
		YieldQuantity:    d.YieldQuantity,
		YieldUnit:        d.YieldUnit,
		MainIngredient:   mainIngredient,
		Cuisine:          cuisine,
		IngredientLines:  lines,
		Tags:             d.Tags,
		ImageURL:         d.ImageURL,
		Nutrition:        *d.Nutrition.ToNutrition(),
		SearchVector:     pgvectorFromSlice(d.SearchVector),
	}
}

// FromIngredient converts a domain Ingredient to an IngredientDTO.
func FromIngredient(i *domain.Ingredient) IngredientDTO {
	return IngredientDTO{
		ID:   i.ID,
		Name: i.Name,
	}
}

// ToIngredient converts an IngredientDTO to a domain Ingredient.
func (d *IngredientDTO) ToIngredient() *domain.Ingredient {
	return &domain.Ingredient{
		ID:   d.ID,
		Name: d.Name,
	}
}

// FromIngredientLine converts a domain line to a DTO.
func FromIngredientLine(line *domain.RecipeIngredientLine) IngredientLineDTO {
	return IngredientLineDTO{
		Ingredient:    FromIngredient(&line.Ingredient),
		QuantityValue: line.QuantityValue,
		QuantityText:  line.QuantityText,
		Unit:          line.Unit,
		IsOptional:    line.IsOptional,
		Note:          line.Note,
		SortOrder:     line.SortOrder,
	}
}

// ToIngredientLine converts a DTO to a domain line.
func (d *IngredientLineDTO) ToIngredientLine() *domain.RecipeIngredientLine {
	return &domain.RecipeIngredientLine{
		Ingredient:    *d.Ingredient.ToIngredient(),
		QuantityValue: d.QuantityValue,
		QuantityText:  d.QuantityText,
		Unit:          d.Unit,
		IsOptional:    d.IsOptional,
		Note:          d.Note,
		SortOrder:     d.SortOrder,
	}
}

// FromCuisine converts a domain Cuisine to a CuisineDTO.
func FromCuisine(c *domain.Cuisine) CuisineDTO {
	return CuisineDTO{
		ID:   c.ID,
		Name: c.Name,
	}
}

// ToCuisine converts a CuisineDTO to a domain Cuisine.
func (d *CuisineDTO) ToCuisine() *domain.Cuisine {
	return &domain.Cuisine{
		ID:   d.ID,
		Name: d.Name,
	}
}

// FromAllergy converts a domain Allergy to an AllergyDTO.
func FromAllergy(a *domain.Allergy) AllergyDTO {
	return AllergyDTO{
		ID:   a.ID,
		Name: a.Name,
	}
}

// ToAllergy converts an AllergyDTO to a domain Allergy.
func (d *AllergyDTO) ToAllergy() *domain.Allergy {
	return &domain.Allergy{
		ID:   d.ID,
		Name: d.Name,
	}
}

// FromNutrition converts domain recipe nutrition to DTO.
func FromNutrition(n *domain.RecipeNutrition) RecipeNutritionDTO {
	if n == nil {
		return RecipeNutritionDTO{}
	}
	return RecipeNutritionDTO{
		CaloriesTotal:      n.CaloriesTotal,
		CaloriesPerServing: n.CaloriesPerServing,
		ProteinG:           n.ProteinG,
		CarbsG:             n.CarbsG,
		FatG:               n.FatG,
		FiberG:             n.FiberG,
		SugarG:             n.SugarG,
		SodiumMg:           n.SodiumMg,
	}
}

// ToNutrition converts a DTO to domain recipe nutrition.
func (d *RecipeNutritionDTO) ToNutrition() *domain.RecipeNutrition {
	return &domain.RecipeNutrition{
		CaloriesTotal:      d.CaloriesTotal,
		CaloriesPerServing: d.CaloriesPerServing,
		ProteinG:           d.ProteinG,
		CarbsG:             d.CarbsG,
		FatG:               d.FatG,
		FiberG:             d.FiberG,
		SugarG:             d.SugarG,
		SodiumMg:           d.SodiumMg,
	}
}

// pgvectorFromSlice creates a pgvector.Vector from a float32 slice.
func pgvectorFromSlice(v []float32) pgvector.Vector {
	return pgvector.NewVector(v)
}
