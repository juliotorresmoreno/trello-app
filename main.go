package main

import (
	"strconv"

	"github.com/juliotorresmoreno/trello-app/configs"
	"github.com/juliotorresmoreno/trello-app/internal/app"
	"github.com/juliotorresmoreno/trello-app/internal/app/preload"
)

func main() {
	e := app.NewServer()

	conf := configs.GetConfig()

	preload.TrelloPreload()

	port := strconv.Itoa(conf.Port)
	addr := ":" + port
	e.Logger.Fatal(e.Start(addr))
}
