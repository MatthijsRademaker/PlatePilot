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

func (r *PostgresRecipeRepository) Save(ctx context.Context, recipe *recipe.Recipe) error {
	// TODO
	return r.db.WithContext(ctx).Create(&entities.RecipeEntity{}).Error
}

func (r *PostgresRecipeRepository) FindById(ctx context.Context, id int) (*recipe.Recipe, error) {
	var recipe entities.RecipeEntity
	r.db.WithContext(ctx).First(&recipe, id)

	// TODO map to domain
	return &recipe, nil
}

func (r *PostgresRecipeRepository) FindByName(ctx context.Context, name string) ([]recipe.Recipe, error) {
	var recipes []entities.RecipeEntity
	r.db.WithContext(ctx).Where(&entities.RecipeEntity{Name: name}).Find(&recipes)

	// TODO map to domain
	return recipes, nil
}

func (r *PostgresRecipeRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(id).Error
}
