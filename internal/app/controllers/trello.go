package controllers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/juliotorresmoreno/trello-app/configs"
	trello_service "github.com/juliotorresmoreno/trello-app/internal/app/services/trello-service"
	"github.com/labstack/echo/v4"
)

/**
 * trelloApi: The objective of this structure is to group all the possible data necessary to
 * make it work, for example, if fixed parameters were required or data that does not change
 * in each call of its functions.
 */
type TrelloApi struct {
}

func AttachTrelloApi(g *echo.Group) *echo.Group {
	c := &TrelloApi{}

	g.POST("", c.Create)
	g.DELETE("/{id}", c.Delete)

	return g
}

type createSchema struct {
	Type        string `json:"type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

var trelloConf = configs.GetConfig().Trello
var trelloService = trello_service.NewTrelloService(trelloConf.BoardId)

// i could not find how can i create new cards with types.
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

type strSlice []string

func (slice strSlice) SliceIndex(value string) int {
	for p, v := range slice {
		if v == value {
			return p
		}
	}
	return -1
}

func createValidation(payload *createSchema) error {
	if payload.Title == "" {
		return errors.New("title is required")
	}
	if payload.Type == "" {
		return errors.New("type is required")
	}
	types := strSlice{"task", "bug", "issue"}
	if types.SliceIndex(payload.Type) < 0 {
		return errors.New("title only can be: task, bug or issue")
	}
	if payload.Type == "task" {
		types := strSlice{"Maintenance"}

		if types.SliceIndex(payload.Category) < 0 {
			return errors.New("category only can be: Maintenance")
		}
	}
	return nil
}

// This is the place where the magic happens. It is the controller that registers card creation requests in Trello
func (t TrelloApi) Create(c echo.Context) error {
	var err error
	payload := new(createSchema)
	c.Bind(payload)
	response := map[string]interface{}{}

	err = createValidation(payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	Alllabels, err := trelloService.GetLabels()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, InternalServerErrorResponse)
	}
	labels := getLabels(Alllabels, payload.Type, payload.Category)

	data := trello_service.CreateCardScheme{
		Name:     payload.Title,
		Desc:     payload.Description,
		IdLabels: labels,
	}
	card, err := trelloService.CreateCard(data)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, InternalServerErrorResponse)
	}
	response["id"] = card.ID

	// It's ideal that you do not send answer because you are send status code 201.
	return c.JSON(http.StatusCreated, response)
}

// This is the place where the magic happens. It is the controller that registers card creation requests in Trello
func (t TrelloApi) Delete(c echo.Context) error {
	var err error
	id := c.Param("id")
	response := map[string]interface{}{}

	err = trelloService.DeleteCard(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, InternalServerErrorResponse)
	}

	return c.JSON(http.StatusNoContent, response)
}

// This is the place where the magic happens. It is the controller that registers card creation requests in Trello
func (t TrelloApi) Get(c echo.Context) error {
	var err error

	cards, err := trelloService.GetCards()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, InternalServerErrorResponse)
	}

	return c.JSON(http.StatusOK, cards)
}
