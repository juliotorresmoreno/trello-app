package preload

import trello_service "github.com/juliotorresmoreno/trello-app/internal/app/services/trello-service"

func TrelloPreload() {
	var trelloService = trello_service.NewTrelloService()

	trelloService.Prepare()
}
