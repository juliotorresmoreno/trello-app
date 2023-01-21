package trello_service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/juliotorresmoreno/trello-app/configs"
	"github.com/juliotorresmoreno/trello-app/internal/app/common"
)

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

type CreateCardScheme struct {
	Name     string   `json:"name"`
	Desc     string   `json:"desc"`
	IdLabels []string `json:"idLabels"`
}

type List struct {
	Id         string      `json:"id"`
	Name       string      `json:"name"`
	Closed     bool        `json:"closed"`
	IdBoard    string      `json:"idBoard"`
	Pos        int         `json:"pos"`
	Subscribed bool        `json:"subscribed"`
	SoftLimit  interface{} `json:"softLimit"`
}

type Board struct {
	Id string `json:"id"`
}

func (t TrelloService) GetBoardId() (string, error) {
	if t.boardId != "" {
		return t.boardId, nil
	}

	config := configs.GetConfig()
	trelloConf := config.Trello

	url := fmt.Sprintf(
		"%v/1/boards/%v?key=%v&token=%v",
		trelloConf.Server, t.shortBoardId,
		trelloConf.Key, trelloConf.Token,
	)
	resp, err := common.DoRequestJSON("GET", url, nil)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", ErrorStatusCode
	}

	board := Board{}
	err = json.NewDecoder(resp.Body).Decode(&board)

	return board.Id, err
}

func (t TrelloService) GetLists() ([]List, error) {
	lists := make([]List, 0)
	config := configs.GetConfig()
	trelloConf := config.Trello

	boardId, err := t.GetBoardId()
	if err != nil {
		return lists, err
	}
	url := fmt.Sprintf(
		"%v/1/boards/%v/lists?key=%v&token=%v",
		trelloConf.Server, boardId, trelloConf.Key, trelloConf.Token,
	)
	resp, err := common.DoRequestJSON("GET", url, nil)
	if err != nil {
		return lists, err
	}

	if resp.StatusCode != http.StatusOK {
		return lists, ErrorStatusCode
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

	boardId, err := t.GetBoardId()
	if err != nil {
		return labels, err
	}
	url := fmt.Sprintf(
		"%v/1/boards/%v/labels?key=%v&token=%v",
		config.Server, boardId, config.Key, config.Token,
	)
	resp, err := common.DoRequestJSON("GET", url, nil)
	if err != nil {
		return labels, err
	}

	json.NewDecoder(resp.Body).Decode(&lists)

	if resp.StatusCode != http.StatusOK {
		return labels, ErrorStatusCode
	}

	for i := 0; i < len(lists); i++ {
		label := lists[i]
		labels[label.Name] = label
	}

	return labels, err
}

type Card struct {
	Badges struct {
		Attachments       int `json:"attachments"`
		AttachmentsByType struct {
			Trello struct {
				Board int `json:"board"`
				Card  int `json:"card"`
			} `json:"trello"`
		} `json:"attachmentsByType"`
		CheckItems            int         `json:"checkItems"`
		CheckItemsChecked     int         `json:"checkItemsChecked"`
		CheckItemsEarliestDue interface{} `json:"checkItemsEarliestDue"`
		Comments              int         `json:"comments"`
		Description           bool        `json:"description"`
		Due                   interface{} `json:"due"`
		DueComplete           bool        `json:"dueComplete"`
		Fogbugz               string      `json:"fogbugz"`
		Location              bool        `json:"location"`
		Start                 interface{} `json:"start"`
		Subscribed            bool        `json:"subscribed"`
		ViewingMemberVoted    bool        `json:"viewingMemberVoted"`
		Votes                 int         `json:"votes"`
	} `json:"badges"`
	CardRole        interface{} `json:"cardRole"`
	CheckItemStates interface{} `json:"checkItemStates"`
	Closed          bool        `json:"closed"`
	Cover           struct {
		Brightness           string      `json:"brightness"`
		Color                interface{} `json:"color"`
		IDAttachment         interface{} `json:"idAttachment"`
		IDPlugin             interface{} `json:"idPlugin"`
		IDUploadedBackground interface{} `json:"idUploadedBackground"`
		Size                 string      `json:"size"`
	} `json:"cover"`
	DateLastActivity time.Time `json:"dateLastActivity"`
	Desc             string    `json:"desc"`
	DescData         struct {
		Emoji struct {
		} `json:"emoji"`
	} `json:"descData"`
	Due               interface{}   `json:"due"`
	DueComplete       bool          `json:"dueComplete"`
	DueReminder       interface{}   `json:"dueReminder"`
	Email             interface{}   `json:"email"`
	ID                string        `json:"id"`
	IDAttachmentCover interface{}   `json:"idAttachmentCover"`
	IDBoard           string        `json:"idBoard"`
	IDChecklists      []interface{} `json:"idChecklists"`
	IDLabels          []string      `json:"idLabels"`
	IDList            string        `json:"idList"`
	IDMembers         []interface{} `json:"idMembers"`
	IDMembersVoted    []interface{} `json:"idMembersVoted"`
	IDShort           int           `json:"idShort"`
	IsTemplate        bool          `json:"isTemplate"`
	Labels            []struct {
		Color   string `json:"color"`
		ID      string `json:"id"`
		IDBoard string `json:"idBoard"`
		Name    string `json:"name"`
	} `json:"labels"`
	ManualCoverAttachment bool        `json:"manualCoverAttachment"`
	Name                  string      `json:"name"`
	Pos                   int         `json:"pos"`
	ShortLink             string      `json:"shortLink"`
	ShortURL              string      `json:"shortUrl"`
	Start                 interface{} `json:"start"`
	Subscribed            bool        `json:"subscribed"`
	URL                   string      `json:"url"`
}

func (t TrelloService) GetCards() ([]Card, error) {
	result := make([]Card, 0)
	config := configs.GetConfig().Trello

	boardId, err := t.GetBoardId()
	if err != nil {
		return result, err
	}
	url := fmt.Sprintf(
		"%v/1/boards/%v/cards?key=%v&token=%v",
		config.Server, boardId, config.Key, config.Token,
	)
	resp, err := common.DoRequestJSON("GET", url, nil)
	if err != nil {
		return result, err
	}

	if resp.StatusCode != http.StatusOK {
		return result, ErrorStatusCode
	}

	err = json.NewDecoder(resp.Body).Decode(&result)

	return result, err
}

func (e TrelloService) CreateCard(payload CreateCardScheme) (Card, error) {
	config := configs.GetConfig()
	trelloConf := config.Trello
	card := Card{}

	lists, err := e.GetLists()
	if err != nil {
		return card, err
	}
	if len(lists) == 0 {
		return card, errors.New("this Board is not working")
	}

	IdList := lists[0].Id
	url := fmt.Sprintf(
		"%v/1/cards?idList=%v&key=%v&token=%v",
		trelloConf.Server, IdList, trelloConf.Key, trelloConf.Token,
	)
	resp, err := common.DoRequestJSON("POST", url, payload)
	if err != nil {
		return card, err
	}

	if resp.StatusCode != http.StatusOK {
		return card, ErrorStatusCode
	}

	json.NewDecoder(resp.Body).Decode(&card)

	return card, nil
}

func (t TrelloService) DeleteCard(id string) error {
	config := configs.GetConfig()
	trelloConf := config.Trello

	boardId, err := t.GetBoardId()
	if err != nil {
		return err
	}
	url := fmt.Sprintf(
		"%v/1/cards/%v?idBoard=%v&key=%v&token=%v",
		trelloConf.Server, id, boardId, trelloConf.Key, trelloConf.Token,
	)

	resp, err := common.DoRequestJSON("DELETE", url, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return ErrorStatusCode
	}

	return nil
}

type CreateCardLabel struct {
	Name  string
	Color string
}

func (t TrelloService) CreateLabel(payload CreateCardLabel) error {
	config := configs.GetConfig()
	trelloConf := config.Trello

	boardId, err := t.GetBoardId()
	if err != nil {
		return err
	}
	url := fmt.Sprintf(
		"%v/1/labels?name=%v&color=%v&idBoard=%v&key=%v&token=%v",
		trelloConf.Server, payload.Name, payload.Color, boardId,
		trelloConf.Key, trelloConf.Token,
	)

	resp, err := common.DoRequestJSON("POST", url, payload)
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

	boardId, err := t.GetBoardId()
	if err != nil {
		return err
	}
	url := fmt.Sprintf(
		"%v/1/labels/%v?idBoard=%v&key=%v&token=%v",
		trelloConf.Server, id, boardId, trelloConf.Key, trelloConf.Token,
	)

	resp, err := common.DoRequestJSON("DELETE", url, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return ErrorStatusCode
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
		{"test", "black"},
	}
	for _, v := range predefined {
		if _, ok := labels[v.Name]; !ok {
			e.CreateLabel(v)
		}
	}

	return nil
}
