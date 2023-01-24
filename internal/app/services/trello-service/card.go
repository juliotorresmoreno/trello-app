package trello_service

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/juliotorresmoreno/trello-app/configs"
	"github.com/juliotorresmoreno/trello-app/internal/app/common"
)

type CreateCardScheme struct {
	Name     string   `json:"name"`
	Desc     string   `json:"desc"`
	IdLabels []string `json:"idLabels"`
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
	trelloConf := configs.GetConfig().Trello

	board, err := t.GetBoard()
	if err != nil {
		return result, err
	}

	// "%v/1/boards/%v/cards?key=%v&token=%v",
	queryParams := url.Values{
		"key":   {trelloConf.Key},
		"token": {trelloConf.Token},
	}
	url := url.URL{
		Scheme:   trelloConf.Scheme,
		Host:     trelloConf.Host,
		Path:     "/1/boards/" + board.Id + "/cards",
		RawQuery: queryParams.Encode(),
	}

	resp, err := common.DoRequestJSON("GET", url.String(), nil)
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

	// "%v/1/cards?idList=%v&key=%v&token=%v"
	queryParams := url.Values{
		"key":    {trelloConf.Key},
		"token":  {trelloConf.Token},
		"idList": {IdList},
	}
	url := url.URL{
		Scheme:   trelloConf.Scheme,
		Host:     trelloConf.Host,
		Path:     "/1/cards",
		RawQuery: queryParams.Encode(),
	}

	resp, err := common.DoRequestJSON("POST", url.String(), payload)
	if err != nil {
		return card, err
	}

	if resp.StatusCode != http.StatusOK {
		return card, ErrorStatusCode
	}

	err = json.NewDecoder(resp.Body).Decode(&card)

	return card, err
}

func (t TrelloService) DeleteCard(id string) error {
	config := configs.GetConfig()
	trelloConf := config.Trello

	board, err := t.GetBoard()
	if err != nil {
		return err
	}

	// "%v/1/cards/%v?idBoard=%v&key=%v&token=%v",
	queryParams := url.Values{
		"key":     {trelloConf.Key},
		"token":   {trelloConf.Token},
		"idBoard": {board.Id},
	}
	url := url.URL{
		Scheme:   trelloConf.Scheme,
		Host:     trelloConf.Host,
		Path:     "/1/cards/" + id,
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
