package app

import (
	"net/http"

	"github.com/juliotorresmoreno/trello-app/configs"

	"github.com/juliotorresmoreno/trello-app/internal/app/controllers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// An instance of echo server is created that allows adding additional
// configurations and making it listen for requests
func NewServer() *echo.Echo {
	e := echo.New()

	conf := configs.GetConfig()

	// Middleware
	if conf.Env != "production" {
		e.Use(middleware.Logger())
	}
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/", HealthCheck)

	api := e.Group("/api/v1")
	trello := api.Group("/trello")
	docs := api.Group("/docs")

	controllers.AttachTrelloApi(trello)
	controllers.AttachSwaggerApi(docs)

	return e
}

func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Server is running",
	})
}
