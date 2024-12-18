package handlers

import (
	"PlatePilot/domain/recipes"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type RecipeHandler struct {
	recipeRepository recipes.RecipeRepository
}

func RegisterRecipeHandlers(repo recipes.RecipeRepository, e *echo.Echo) {
	recipeHandler := &RecipeHandler{recipeRepository: repo}
	e.GET("/recipe/:id", recipeHandler.getRecipe)
	e.POST("/recipe", recipeHandler.postRecipe)
}

// getRecipe handles the HTTP GET request to retrieve a recipe by its ID.
// @Summary Retrieve a recipe by ID
// @Description Get a recipe by its ID
// @Tags recipes
// @Accept json
// @Produce json
// @Param id path int true "Recipe ID"
// @Success 200 {object} Recipe
// @Failure 400 {object} map[string]string{"error": "Invalid recipe ID"}
// @Router /recipes/{id} [get]
func (controller *RecipeHandler) getRecipe(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid recipe ID"})
	}

	recipe, err := controller.recipeRepository.FindById(c.Request().Context(), id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid recipe ID"})
	}

	return c.JSON(http.StatusOK, recipe)
}

// getRecipe handles the HTTP GET request to retrieve a recipe by its ID.
// @Summary Retrieve a recipe by ID
// @Description Get a recipe by its ID
// @Tags recipes
// @Accept json
// @Produce json
// @Param id path int true "Recipe ID"
// @Success 200 {object} Recipe
// @Failure 400 {object} map[string]string{"error": "Invalid recipe ID"}
// @Router /recipes/{id} [get]
func (controller *RecipeHandler) postRecipe(c echo.Context) error {
	recipeToAdd := new(Recipe)

	if err := c.Bind(recipeToAdd); err != nil {
		return err
	}

	ingredients := make([]recipes.Ingredient, len(recipeToAdd.Ingredients))

	for _, i := range recipeToAdd.Ingredients {
		ingredients = append(ingredients, recipes.Ingredient{Name: i.Name, Quantity: i.Quantity, Unit: i.Unit})
	}

	recipe, err := recipes.NewRecipe(recipeToAdd.Name, ingredients, recipeToAdd.Instructions, recipeToAdd.CookingTime, recipeToAdd.Cuisines, recipeToAdd.KCalories)

	if err != nil {
		return c.JSON(http.StatusBadRequest, "Non valid recipe provided")
	}

	err = controller.recipeRepository.Save(c.Request().Context(), recipe)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error try again later")
	}

	return c.JSON(http.StatusCreated, recipeToAdd)
}

type Recipe struct {
	Name         string        `json:"name"`
	Ingredients  []Ingredient  `json:"ingredients"`
	Instructions []string      `json:"instructions"`
	CookingTime  time.Duration `json:"CookingTime"`
	Cuisines     []string      `json:"cuisines"`
	KCalories    uint          `json:"kCalories"`
}

type Ingredient struct {
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
	Name     string  `json:"name"`
}
