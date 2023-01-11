package trello_service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/juliotorresmoreno/trello-app/configs"
)

var defaultErrorStatusCode = errors.New("ups, something has not working, please wait a moment and re-try")

type TrelloService struct {
	boardId string
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

type Board struct {
	Id string `json:"id"`
}

func (t TrelloService) getBoardId() (string, error) {
	if t.boardId != "" {
		return t.boardId, nil
	}

	config := configs.GetConfig().Trello

	url := fmt.Sprintf("%v/1/boards/%v?key=%v&token=%v", config.Server, config.BoardId, config.Key, config.Token)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", defaultErrorStatusCode
	}

	board := Board{}
	err = json.NewDecoder(resp.Body).Decode(&board)

	return board.Id, err
}

func (t TrelloService) GetLists() ([]List, error) {
	lists := make([]List, 0)
	config := configs.GetConfig().Trello

	boardId, err := t.getBoardId()
	if err != nil {
		return lists, err
	}
	url := fmt.Sprintf("%v/1/boards/%v/lists?key=%v&token=%v", config.Server, boardId, config.Key, config.Token)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return lists, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return lists, err
	}

	if resp.StatusCode != http.StatusOK {
		return lists, defaultErrorStatusCode
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

func (t TrelloService) GetLabels() (Labels, error) {
	lists := make([]Label, 0)
	labels := Labels{}
	config := configs.GetConfig().Trello

	boardId, err := t.getBoardId()
	if err != nil {
		return labels, err
	}
	url := fmt.Sprintf("%v/1/boards/%v/labels?key=%v&token=%v", config.Server, boardId, config.Key, config.Token)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return labels, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return labels, err
	}

	if resp.StatusCode != http.StatusOK {
		return labels, defaultErrorStatusCode
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

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return defaultErrorStatusCode
	}

	return nil
}

type CreateCardLabel struct {
	Name  string
	Color string
}

func (t TrelloService) CreateLabel(payload CreateCardLabel) error {
	config := configs.GetConfig().Trello

	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	body := bytes.NewBuffer(b)
	// 'https://api.trello.com/1/labels?name={name}&color={color}&idBoard={idBoard}&key=APIKey&token=APIToken'

	boardId, err := t.getBoardId()
	if err != nil {
		return err
	}
	url := fmt.Sprintf(
		"%v/1/labels?name=%v&color=%v&idBoard=%v&key=%v&token=%v",
		config.Server, payload.Name, payload.Color, boardId, config.Key, config.Token)

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
		return defaultErrorStatusCode
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
