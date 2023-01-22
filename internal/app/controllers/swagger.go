package controllers

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/labstack/echo/v4"
)

type SwaggerApi struct {
}

func AttachSwaggerApi(g *echo.Group) *echo.Group {
	c := &SwaggerApi{}

	g.GET("", c.Get)
	g.GET("/openapi.yml", c.GetDocument)

	return g
}

func (el SwaggerApi) Get(e echo.Context) error {
	basePath, err := os.Getwd()
	if err != nil {
		return e.HTML(http.StatusNotFound, "Not found")
	}
	index := path.Join(basePath, "assets", "swagger.html")

	f, err := os.Open(index)
	if err != nil {
		return e.HTML(http.StatusNotFound, "Not found")
	}

	html, err := ioutil.ReadAll(f)
	if err != nil {
		return e.HTML(http.StatusNotFound, "Not found")
	}

	return e.HTML(200, string(html))
}

func (el SwaggerApi) GetDocument(e echo.Context) error {
	NotFountError := "Not found"
	basePath, err := os.Getwd()
	if err != nil {
		return e.HTML(http.StatusNotFound, NotFountError)
	}
	index := path.Join(basePath, "docs", "swagger.yml")

	f, err := os.Open(index)
	if err != nil {
		return e.HTML(http.StatusNotFound, NotFountError)
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return e.HTML(http.StatusNotFound, NotFountError)
	}

	return e.String(200, string(data))
}
