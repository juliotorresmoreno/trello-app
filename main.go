package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/juliotorresmoreno/trello-app/configs"
	"github.com/juliotorresmoreno/trello-app/gen/restapi"
	"github.com/juliotorresmoreno/trello-app/gen/restapi/operations"
	"github.com/juliotorresmoreno/trello-app/internal/app"
	"github.com/juliotorresmoreno/trello-app/internal/app/preload"
)

func swaggerServer() {
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	api := operations.NewGreeterAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer server.Shutdown()
	flag.Parse()
	server.Port = 4500

	api.GetGreetingHandler = operations.GetGreetingHandlerFunc(
		func(params operations.GetGreetingParams) middleware.Responder {
			name := swag.StringValue(params.Name)
			if name == "" {
				name = "World"
			}

			greeting := fmt.Sprintf("Hello, %s!", name)
			return operations.NewGetGreetingOK().WithPayload(greeting)
		})

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	e := app.NewServer()

	conf := configs.GetConfig()

	preload.TrelloPreload()

	go swaggerServer()

	addr := fmt.Sprintf(":%v", conf.Port)
	e.Logger.Fatal(e.Start(addr))
}
