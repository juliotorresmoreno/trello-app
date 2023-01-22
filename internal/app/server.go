package app

import (
	"net/http"
	"strings"

	"github.com/juliotorresmoreno/trello-app/configs"

	"github.com/juliotorresmoreno/trello-app/internal/app/controllers"
	"github.com/juliotorresmoreno/trello-app/internal/app/middlewares"
	"github.com/labstack/echo-contrib/prometheus"
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

	sources := []string{
		"default-src 'self' unpkg.com fonts.googleapis.com",
		"script-src 'self' unpkg.com 'sha256-Y6uMcHW6QmLl5jRFbZ/tRzmQ0a0EjDOkQO0vpHPjiRE='",
		"img-src 'self' data:",
	}
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "",
		ContentTypeNosniff:    "",
		XFrameOptions:         "",
		HSTSMaxAge:            3600,
		ContentSecurityPolicy: strings.Join(sources, "; "),
		Skipper: func(c echo.Context) bool {
			return c.Request().URL.Path != "/api/v1/docs"
		},
	}))
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	gzip := middleware.GzipWithConfig(middleware.DefaultGzipConfig)
	e.Use(middlewares.Ommit(gzip, []string{"/metrics"}))

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
