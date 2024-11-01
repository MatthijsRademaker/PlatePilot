package recipe

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"PlatePilot/domain/recipe"
	recipeRepository "PlatePilot/infrastructure/recipe"
)

type RecipeController struct {
	recipeRepository recipe.RecipeRepository
}

func (controller *RecipeController) CreateRecipeController(db *gorm.DB) {
	controller.recipeRepository = recipeRepository.NewPostgresRecipeRepository(db)
}

func RegisterHandler(e *echo.Echo) {

	e.GET("/recipes/:id", getRecipe)
	e.Logger.Fatal(e.Start(":8080"))

}

// @Summary Add a new recipe
// @Description Create a new recipe
// @Tags recipes
// @Accept json
// @Produce json
// @Param recipe body Recipe true "Recipe to add"
// @Success 200 {object} Recipe
// @Router /recipes [post]
func getRecipe(c echo.Context) error {
	return nil
}
