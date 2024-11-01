package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"PlatePilot/domain/recipe"
)

type RecipeHandler struct {
	recipeRepository recipe.RecipeRepository
}

func RegisterRecipeHandlers(repo recipe.RecipeRepository, e *echo.Echo) {
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

	var ingredients = make([]recipe.Ingredient, len(recipeToAdd.Ingredients))

	for _, i := range recipeToAdd.Ingredients {
		ingredients = append(ingredients, recipe.Ingredient{Name: i.Name, Quantity: i.Quantity, Unit: i.Unit})
	}

	recipe, err := recipe.NewRecipe(recipeToAdd.Name, ingredients, recipeToAdd.Instructions, recipeToAdd.CookingTime)

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
}

type Ingredient struct {
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
	Name     string  `json:"name"`
}
