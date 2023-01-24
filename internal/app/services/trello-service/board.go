package trello_service

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/juliotorresmoreno/trello-app/configs"
	"github.com/juliotorresmoreno/trello-app/internal/app/common"
)

type Board struct {
	Id string `json:"id"`
}

func (t TrelloService) GetBoard() (*Board, error) {
	config := configs.GetConfig()
	trelloConf := config.Trello

	board := &Board{}

	if t.boardId != "" {
		return board, nil
	}

	// "%v/1/boards/%v?key=%v&token=%v",
	queryParams := url.Values{
		"key":   {trelloConf.Key},
		"token": {trelloConf.Token},
	}
	url := url.URL{
		Scheme:   trelloConf.Scheme,
		Host:     trelloConf.Host,
		Path:     "/1/boards/" + t.shortBoardId,
		RawQuery: queryParams.Encode(),
	}

	resp, err := common.DoRequestJSON("GET", url.String(), nil)
	if err != nil {
		return board, err
	}

	if resp.StatusCode != http.StatusOK {
		return board, ErrorStatusCode
	}

	err = json.NewDecoder(resp.Body).Decode(&board)

	return board, err
}
