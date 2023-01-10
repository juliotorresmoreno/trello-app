package controllers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/juliotorresmoreno/trello-app/internal/app/services"
	"github.com/labstack/echo/v4"
)

type trelloApi struct {
}

func AttachTrelloApi(g *echo.Group) *echo.Group {
	c := &trelloApi{}

	g.POST("", c.create)

	return g
}

type createSchema struct {
	Type        string `json:"type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

var trelloService = services.NewTrelloService()

func (t trelloApi) create(c echo.Context) error {
	var err error
	payload := new(createSchema)
	c.Bind(payload)
	response := map[string]interface{}{}

	data := services.CreateCardScheme{
		Name: payload.Title,
		Desc: payload.Description,
	}
	err = trelloService.CreateCard(data)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
		})
	}

	json.NewEncoder(os.Stdout).Encode(payload)

	// It's ideal that you do not send answer because you are send status code 201.
	return c.JSON(http.StatusCreated, response)
}
