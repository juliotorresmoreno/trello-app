package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/juliotorresmoreno/trello-app/configs"
)

type TrelloService struct {
}

// NewTrelloService
func NewTrelloService() *TrelloService {
	s := &TrelloService{}
	return s
}

type CreateCardScheme struct {
	Name     string   `json:"name"`
	Desc     string   `json:"desc"`
	IdLabels []string `json:"idLabels"`
}

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

func (e TrelloService) CreateCard(payload CreateCardScheme) error {
	config := configs.GetConfig().Trello

	lists, err := e.GetLists()
	if err != nil {
		return err
	}
	if len(lists) == 0 {
		return errors.New("this Board is not working")
	}

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

	if resp.StatusCode != http.StatusOK {
		return errors.New("ups, something has not working, please wait a moment and re-try")
	}

	return nil
}

type CreateCardLabel struct {
	Name  string
	Color string
}

func (e TrelloService) CreateLabel(payload CreateCardLabel) error {
	config := configs.GetConfig().Trello

	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	body := bytes.NewBuffer(b)
	// 'https://api.trello.com/1/labels?name={name}&color={color}&idBoard={idBoard}&key=APIKey&token=APIToken'

	url := fmt.Sprintf(
		"%v/1/labels?name=%v&color=%v&idBoard=%v&key=%v&token=%v",
		config.Server, payload.Name, payload.Color, config.BoardId, config.Key, config.Token)

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

	if resp.StatusCode != http.StatusOK {
		return errors.New("ups, something has not working, please wait a moment and re-try")
	}

	return nil
}

func (e TrelloService) Prepare() error {
	labels, err := e.GetLabels()
	if err != nil {
		return err
	}
	var predefined = []CreateCardLabel{
		{"task", "blue"},
		{"bug", "red"},
		{"issue", "green"},
		{"maintenance", "purple"},
	}
	for _, v := range predefined {
		if _, ok := labels[v.Name]; !ok {
			e.CreateLabel(v)
		}
	}

	return nil
}
