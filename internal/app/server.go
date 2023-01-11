package app

import (
	"net/http"

	"github.com/juliotorresmoreno/trello-app/configs"

	"github.com/juliotorresmoreno/trello-app/internal/app/controllers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

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

	controllers.AttachTrelloApi(trello)

	return e
}

// HealthCheck godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags root
// @Accept application/json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router / [get]
func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Server is running",
	})
}
