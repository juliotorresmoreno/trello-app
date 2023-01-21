package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/juliotorresmoreno/trello-app/internal/app/controllers"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

var testTitle = "test_card"
var testDescription = "test description"
var ErrorStatusCode = errors.New("ups, something has not working, please wait a moment and re-try")

func createRequest(method, url string, body interface{}) (*httptest.ResponseRecorder, echo.Context) {
	e := echo.New()
	var buff io.Reader = nil
	if body != nil {
		b, _ := json.Marshal(body)
		buff = bytes.NewBuffer(b)
	}
	req := httptest.NewRequest(method, url, buff)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	return rec, c
}

func removeCard(api *controllers.TrelloApi, id string) (*httptest.ResponseRecorder, error) {
	rec, c := createRequest(http.MethodDelete, "/"+id, nil)
	c.SetParamNames("id")
	c.SetParamValues(id)

	return rec, api.Delete(c)
}

func TestTrelloIssue(t *testing.T) {
	require := require.New(t)

	body := map[string]interface{}{
		"type":        "issue",
		"title":       testTitle,
		"description": testDescription,
	}

	rec, c := createRequest(http.MethodPost, "/", body)

	api := &controllers.TrelloApi{}

	err := api.Create(c)
	require.NoError(err)
	require.Equal(rec.Code, http.StatusCreated)
	response := map[string]string{}
	err = json.NewDecoder(rec.Body).Decode(&response)
	require.NoError(err)

	rec, err = removeCard(api, response["id"])
	require.NoError(err)
	require.Equal(rec.Code, http.StatusNoContent)
}

func TestTrelloIssueE1(t *testing.T) {
	require := require.New(t)
	body := map[string]interface{}{
		"type":        "issues",
		"title":       testTitle,
		"description": testDescription,
	}

	rec, c := createRequest(http.MethodPost, "/", body)
	api := &controllers.TrelloApi{}

	err := api.Create(c)
	require.NoError(err)
	require.Equal(rec.Code, http.StatusBadRequest)
}

func TestTrelloBug(t *testing.T) {
	require := require.New(t)

	body := map[string]interface{}{
		"type":        "bug",
		"title":       testTitle,
		"description": testDescription,
	}

	rec, c := createRequest(http.MethodPost, "/", body)
	api := &controllers.TrelloApi{}

	err := api.Create(c)
	require.NoError(err)
	require.Equal(rec.Code, http.StatusCreated)
	response := map[string]string{}
	err = json.NewDecoder(rec.Body).Decode(&response)
	require.NoError(err)

	rec, err = removeCard(api, response["id"])
	require.NoError(err)
	require.Equal(rec.Code, http.StatusNoContent)
}

func TestTrelloTask(t *testing.T) {
	require := require.New(t)

	body := map[string]interface{}{
		"type":     "task",
		"title":    testTitle,
		"category": "Maintenance",
	}

	rec, c := createRequest(http.MethodPost, "/", body)
	api := &controllers.TrelloApi{}

	err := api.Create(c)
	require.NoError(err)
	require.Equal(rec.Code, http.StatusCreated)
	response := map[string]string{}
	err = json.NewDecoder(rec.Body).Decode(&response)
	require.NoError(err)

	rec, err = removeCard(api, response["id"])
	require.NoError(err)
	require.Equal(rec.Code, http.StatusNoContent)
}
func TestTrelloTaskE1(t *testing.T) {
	require := require.New(t)

	body := map[string]interface{}{
		"type":     "task",
		"title":    testTitle,
		"category": "maintenance",
	}

	rec, c := createRequest(http.MethodPost, "/", body)
	api := controllers.TrelloApi{}

	err := api.Create(c)
	require.NoError(err)
	require.Equal(rec.Code, http.StatusBadRequest)
}
