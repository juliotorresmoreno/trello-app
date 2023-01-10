package app

import (
	"net/http"

	"github.com/juliotorresmoreno/trello-app/internal/app/controllers"
	"github.com/labstack/echo/v4"
)

func NewServer() *echo.Echo {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	api := e.Group("/api/v1")
	trello := api.Group("/trello")

	controllers.AttachTrelloApi(trello)

	return e
}
