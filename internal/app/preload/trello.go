package preload

import trello_service "github.com/juliotorresmoreno/trello-app/internal/app/services/trello-service"

// This works as a kind of database migration but it's not really.
// Prepare the labels on the board that will be used.
func TrelloPreload() {
	var trelloService = trello_service.NewTrelloService()

	trelloService.Prepare()
}
