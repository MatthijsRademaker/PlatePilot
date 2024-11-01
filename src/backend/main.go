package main

import (
	_ "PlatePilot/docs" // docs is generated by Swag CLI, you have to import it.

	"github.com/labstack/echo/v4"
)

// @title Recipe API
// @version 1.0
// @description This is a sample server for a recipe app.
// @host localhost:8080
// @BasePath /
func main() {
	e := echo.New()
	// TODO wait for fix?
	// e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(":9080"))
}
