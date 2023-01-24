package trello_service

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/juliotorresmoreno/trello-app/configs"
	"github.com/juliotorresmoreno/trello-app/internal/app/common"
)

type List struct {
	Id         string      `json:"id"`
	Name       string      `json:"name"`
	Closed     bool        `json:"closed"`
	IdBoard    string      `json:"idBoard"`
	Pos        int         `json:"pos"`
	Subscribed bool        `json:"subscribed"`
	SoftLimit  interface{} `json:"softLimit"`
}

func (t TrelloService) GetLists() ([]List, error) {
	lists := make([]List, 0)
	config := configs.GetConfig()
	trelloConf := config.Trello

	board, err := t.GetBoard()
	if err != nil {
		return lists, err
	}
	// "%v/1/boards/%v/lists?key=%v&token=%v",
	queryParams := url.Values{
		"key":   {trelloConf.Key},
		"token": {trelloConf.Token},
	}
	url := url.URL{
		Scheme:   trelloConf.Scheme,
		Host:     trelloConf.Host,
		Path:     "/1/boards/" + board.Id + "/lists",
		RawQuery: queryParams.Encode(),
	}
	resp, err := common.DoRequestJSON("GET", url.String(), nil)
	if err != nil {
		return lists, err
	}

	if resp.StatusCode != http.StatusOK {
		return lists, ErrorStatusCode
	}

	err = json.NewDecoder(resp.Body).Decode(&lists)

	return lists, err
}
