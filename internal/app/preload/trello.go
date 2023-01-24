package preload

import (
	"log"

	"github.com/juliotorresmoreno/trello-app/configs"
	trello_service "github.com/juliotorresmoreno/trello-app/internal/app/services/trello-service"
)

// This works as a kind of database migration but it's not really.
// Prepare the labels on the board that will be used.
func TrelloPreload() {
	trelloConf := configs.GetConfig().Trello
	var trelloService = trello_service.NewTrelloService(trelloConf.BoardId)

	err := trelloService.Prepare()
	if err != nil {
		log.Fatal(err)
	}
}
