package trello_service

import "errors"

var ErrorStatusCode = errors.New("ups, something has not working, please wait a moment and re-try")

type TrelloService struct {
	boardId string

	shortBoardId string
}

// NewTrelloService
func NewTrelloService(shortBoardId string) *TrelloService {
	s := &TrelloService{
		shortBoardId: shortBoardId,
	}
	return s
}
