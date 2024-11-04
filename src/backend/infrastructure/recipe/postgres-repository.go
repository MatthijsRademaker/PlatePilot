package recipeRepository

import (
	"PlatePilot/domain/recipe"
	"PlatePilot/infrastructure/recipe/db/entities"
	"context"

	"gorm.io/gorm"
)

type PostgresRecipeRepository struct {
	db *gorm.DB
}

func NewPostgresRecipeRepository(db *gorm.DB) recipe.RecipeRepository {
	return &PostgresRecipeRepository{db: db}
}

// #region recipes
func (r *PostgresRecipeRepository) Save(ctx context.Context, recipe *recipe.Recipe) error {

	// TODO look into unmarshal?
	ingredients := make([]entities.RecipeIngredientEntity, len(recipe.Ingredients))

	// TODO create a transaction that adds all relations
	for _, i := range ingredients {
		ingredients = append(ingredients, entities.RecipeIngredientEntity{Quantity: i.Quantity, Unit: i.Unit})
	}

	return r.db.WithContext(ctx).Create(&entities.RecipeEntity{Name: recipe.Name, Ingredients: ingredients, Instructions: recipe.Instructions, CookingTime: recipe.CookingTime}).Error
}

func (r *PostgresRecipeRepository) FindById(ctx context.Context, id int) (*recipe.Recipe, error) {
	var recipeEntity entities.RecipeEntity
	// create helper method and include relevent entities
	r.db.WithContext(ctx).First(&recipeEntity, id)

	return mapToDomainRecipe(&recipeEntity), nil
}

func (r *PostgresRecipeRepository) FindByName(ctx context.Context, name string) ([]recipe.Recipe, error) {
	var recipeEntities []entities.RecipeEntity
	r.db.WithContext(ctx).Where(&entities.RecipeEntity{Name: name}).Find(&recipeEntities)

	var recipes = make([]recipe.Recipe, len(recipeEntities))

	for _, r := range recipeEntities {

		recipes = append(recipes, *mapToDomainRecipe(&r))
	}
	return recipes, nil
}

func (r *PostgresRecipeRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(id).Error
}

func mapToDomainRecipe(entity *entities.RecipeEntity) *recipe.Recipe {

	ingredients := make([]recipe.Ingredient, len(entity.Ingredients))

	for _, i := range ingredients {
		ingredients = append(ingredients, recipe.Ingredient{Quantity: i.Quantity, Unit: i.Unit, Name: i.Name})
	}

	return &recipe.Recipe{
		Name:         entity.Name,
		Ingredients:  ingredients,
		Instructions: entity.Instructions,
		CookingTime:  entity.CookingTime,
	}
}

// region cuisine

// region ingredients
