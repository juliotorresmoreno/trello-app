package controllers

import (
	"net/http"
	"strings"

	services "github.com/juliotorresmoreno/trello-app/internal/app/services/trello-service"
	trello_service "github.com/juliotorresmoreno/trello-app/internal/app/services/trello-service"
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

var trelloService = trello_service.NewTrelloService()

// getLabels: i could not find how can i create new cards with types.
func getLabels(labels trello_service.Labels, Type string, Category string) []string {
	list := make([]string, 0)
	if label, ok := labels[Type]; ok {
		list = append(list, label.Id)
	}
	if Type == "task" && Category == "Maintenance" {
		category := strings.ToLower(Category)
		if label, ok := labels[category]; ok {
			list = append(list, label.Id)
		}
	}
	return list
}
func (t trelloApi) create(c echo.Context) error {
	var err error
	payload := new(createSchema)
	c.Bind(payload)
	response := map[string]interface{}{}

	Alllabels, err := trelloService.GetLabels()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Internal Server Error",
		})
	}
	labels := getLabels(Alllabels, payload.Type, payload.Category)

	data := services.CreateCardScheme{
		Name:     payload.Title,
		Desc:     payload.Description,
		IdLabels: labels,
	}
	err = trelloService.CreateCard(data)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Internal Server Error",
		})
	}

	// It's ideal that you do not send answer because you are send status code 201.
	return c.JSON(http.StatusCreated, response)
}
