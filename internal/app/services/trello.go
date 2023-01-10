package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/juliotorresmoreno/trello-app/configs"
)

type TrelloService struct {
}

// url := "GET https://api.trello.com/1/boards/KydXU1O5/lists?key={{Key}}&token={{Token}}"
// fmt.Println(url)

// NewTrelloService
func NewTrelloService() *TrelloService {
	s := &TrelloService{}
	return s
}

type CreateCardScheme struct {
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	Type     string `json:"-"`
	Category string `json:"-"`
}

/*
	{
		"id":"63bdd88d06536303b08c4752",
		"name":"NaN",
		"closed":false,
		"idBoard":"63bdc4fbd8316302cb9bd29c",
		"pos":8192,
		"subscribed":false,
		"softLimit":null
	}
*/

type List struct {
	Id         string      `json:"id"`
	Name       string      `json:"name"`
	Closed     bool        `json:"closed"`
	IdBoard    string      `json:"id_board"`
	Pos        int         `json:"pos"`
	Subscribed bool        `json:"subscribed"`
	SoftLimit  interface{} `json:"soft_limit"`
}

func (e TrelloService) GetLists() ([]List, error) {
	lists := make([]List, 0)
	config := configs.GetConfig().Trello

	url := fmt.Sprintf("%v/1/boards/%v/lists?key=%v&token=%v", config.Server, config.BoardId, config.Key, config.Token)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return lists, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return lists, err
	}

	json.NewDecoder(resp.Body).Decode(&lists)

	return lists, err
}

type Label struct {
	Id      string `json:"id"`
	IdBoard string `json:"id_board"`
	Name    string `json:"name"`
	Color   string `json:"color"`
}

type Labels map[string]Label

func (e TrelloService) GetLabels() (Labels, error) {
	// GET https://api.trello.com/1/boards/KydXU1O5/labels?key={{Key}}&token={{Token}}
	lists := make([]Label, 0)
	labels := Labels{}
	config := configs.GetConfig().Trello

	url := fmt.Sprintf("%v/1/boards/%v/labels?key=%v&token=%v", config.Server, config.BoardId, config.Key, config.Token)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return labels, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return labels, err
	}

	json.NewDecoder(resp.Body).Decode(&lists)

	for i := 0; i < len(lists); i++ {
		label := lists[i]
		labels[label.Name] = label
	}

	return labels, err
}

func getLabels(labels Labels, payload CreateCardScheme) []Label {
	list := make([]Label, 0)
	if label, ok := labels[payload.Type]; ok {
		list = append(list, label)
	}
	if payload.Type == "task" && payload.Category == "Maintenance" {
		category := strings.ToLower(payload.Category)
		if label, ok := labels[category]; ok {
			list = append(list, label)
		}
	}
	return list
}

func (e TrelloService) CreateCard(payload CreateCardScheme) error {
	config := configs.GetConfig().Trello

	lists, err := e.GetLists()
	if err != nil {
		return err
	}
	if len(lists) == 0 {
		return errors.New("this Board is not working")
	}
	Alllabels, err := e.GetLabels()
	if err != nil {
		return err
	}
	labels := getLabels(Alllabels, payload)
	fmt.Println("labels", labels)

	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	body := bytes.NewBuffer(b)

	IdList := lists[0].Id
	url := fmt.Sprintf("%v/1/cards?idList=%v&key=%v&token=%v", config.Server, IdList, config.Key, config.Token)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}
	req.Header.Add("content-type", "application/json")

	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("@@@@@@@@@", string(b))

	return nil
}
