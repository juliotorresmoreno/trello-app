package trello_service

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/juliotorresmoreno/trello-app/configs"
	"github.com/juliotorresmoreno/trello-app/internal/app/common"
)

type Label struct {
	Id      string `json:"id"`
	IdBoard string `json:"id_board"`
	Name    string `json:"name"`
	Color   string `json:"color"`
}

type Labels map[string]Label

func (t TrelloService) GetLabels() (Labels, error) {
	lists := make([]Label, 0)
	labels := Labels{}
	trelloConf := configs.GetConfig().Trello

	board, err := t.GetBoard()
	if err != nil {
		return labels, err
	}

	// "%v/1/boards/%v/labels?key=%v&token=%v",
	queryParams := url.Values{
		"key":   {trelloConf.Key},
		"token": {trelloConf.Token},
	}
	url := url.URL{
		Scheme:   trelloConf.Scheme,
		Host:     trelloConf.Host,
		Path:     "/1/boards/" + board.Id + "/labels",
		RawQuery: queryParams.Encode(),
	}

	resp, err := common.DoRequestJSON("GET", url.String(), nil)
	if err != nil {
		return labels, err
	}

	err = json.NewDecoder(resp.Body).Decode(&lists)
	if err != nil {
		return labels, err
	}

	if resp.StatusCode != http.StatusOK {
		return labels, ErrorStatusCode
	}

	for i := 0; i < len(lists); i++ {
		label := lists[i]
		labels[label.Name] = label
	}

	return labels, err
}

type CreateCardLabel struct {
	Name  string
	Color string
}

func (t TrelloService) CreateLabel(payload CreateCardLabel) error {
	config := configs.GetConfig()
	trelloConf := config.Trello

	board, err := t.GetBoard()
	if err != nil {
		return err
	}

	// "%v/1/labels?name=%v&color=%v&idBoard=%v&key=%v&token=%v"
	queryParams := url.Values{
		"key":     {trelloConf.Key},
		"token":   {trelloConf.Token},
		"name":    {payload.Name},
		"color":   {payload.Color},
		"idBoard": {board.Id},
	}
	url := url.URL{
		Scheme:   trelloConf.Scheme,
		Host:     trelloConf.Host,
		Path:     "/1/labels",
		RawQuery: queryParams.Encode(),
	}

	resp, err := common.DoRequestJSON("POST", url.String(), payload)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return ErrorStatusCode
	}

	return nil
}

func (t TrelloService) DeleteLabel(id string) error {
	config := configs.GetConfig()
	trelloConf := config.Trello

	board, err := t.GetBoard()
	if err != nil {
		return err
	}

	// "%v/1/labels/%v?idBoard=%v&key=%v&token=%v",
	queryParams := url.Values{
		"key":     {trelloConf.Key},
		"token":   {trelloConf.Token},
		"idBoard": {board.Id},
	}
	url := url.URL{
		Scheme:   trelloConf.Scheme,
		Host:     trelloConf.Host,
		Path:     "/1/labels/" + id,
		RawQuery: queryParams.Encode(),
	}

	resp, err := common.DoRequestJSON("DELETE", url.String(), nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return ErrorStatusCode
	}

	return nil
}
