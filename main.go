package main

import (
	"fmt"

	"github.com/juliotorresmoreno/trello-app/configs"
	"github.com/juliotorresmoreno/trello-app/internal/app"
)

func main() {
	e := app.NewServer()

	conf := configs.GetConfig()

	addr := fmt.Sprintf(":%v", conf.Port)
	e.Logger.Fatal(e.Start(addr))
}
