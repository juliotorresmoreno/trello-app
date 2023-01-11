package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/juliotorresmoreno/trello-app/configs"
	"github.com/juliotorresmoreno/trello-app/internal/app"
)

func doRequest(url string, body map[string]interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	_body := bytes.NewBuffer(b)

	req, err := http.NewRequest("POST", url, _body)
	if err != nil {
		return err
	}
	req.Header.Add("content-type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return errors.New("ups, something has not working, please wait a moment and re-try")
	}
	return nil
}

func doRequestError(url string, body map[string]interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	_body := bytes.NewBuffer(b)

	req, err := http.NewRequest("POST", url, _body)
	if err != nil {
		return err
	}
	req.Header.Add("content-type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusCreated {
		return errors.New("ups, something has not working, please wait a moment and re-try")
	}
	return nil
}

func TestTrelloIssue(t *testing.T) {
	e := app.NewServer()

	go func() {
		conf := configs.GetConfig()

		addr := fmt.Sprintf(":%v", conf.Port)
		e.Logger.Fatal(e.Start(addr))
	}()

	conf := configs.GetConfig()
	time.Sleep(1 * time.Second)

	addr := fmt.Sprintf("http://localhost:%v", conf.Port)
	url := fmt.Sprintf("%v/api/v1/trello", addr)

	body := map[string]interface{}{
		"type":        "issue",
		"title":       "probando",
		"description": "probando descripcion",
	}

	err := doRequest(url, body)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestTrelloIssueE1(t *testing.T) {
	conf := configs.GetConfig()
	time.Sleep(1 * time.Second)

	addr := fmt.Sprintf("http://localhost:%v", conf.Port)
	url := fmt.Sprintf("%v/api/v1/trello", addr)

	body := map[string]interface{}{
		"type":        "issues",
		"title":       "probando",
		"description": "probando descripcion",
	}

	err := doRequestError(url, body)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestTrelloBug(t *testing.T) {

	conf := configs.GetConfig()
	time.Sleep(1 * time.Second)

	addr := fmt.Sprintf("http://localhost:%v", conf.Port)
	url := fmt.Sprintf("%v/api/v1/trello", addr)

	body := map[string]interface{}{
		"type":        "bug",
		"title":       "nada",
		"description": "asdas",
	}

	err := doRequest(url, body)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestTrelloBugE1(t *testing.T) {

	conf := configs.GetConfig()
	time.Sleep(1 * time.Second)

	addr := fmt.Sprintf("http://localhost:%v", conf.Port)
	url := fmt.Sprintf("%v/api/v1/trello", addr)

	body := map[string]interface{}{
		"type":        "bug",
		"title":       "",
		"description": "asdas",
	}

	err := doRequestError(url, body)
	if err != nil {
		t.Error(err)
		return
	}

}

func TestTrelloTask(t *testing.T) {

	conf := configs.GetConfig()
	time.Sleep(1 * time.Second)

	addr := fmt.Sprintf("http://localhost:%v", conf.Port)
	url := fmt.Sprintf("%v/api/v1/trello", addr)

	body := map[string]interface{}{
		"type":     "task",
		"title":    "NAdando",
		"category": "Maintenance",
	}

	err := doRequest(url, body)
	if err != nil {
		t.Error(err)
		return
	}

}
func TestTrelloTaskE1(t *testing.T) {

	conf := configs.GetConfig()
	time.Sleep(1 * time.Second)

	addr := fmt.Sprintf("http://localhost:%v", conf.Port)
	url := fmt.Sprintf("%v/api/v1/trello", addr)

	body := map[string]interface{}{
		"type":     "task",
		"title":    "NAdando",
		"category": "maintenance",
	}

	err := doRequestError(url, body)
	if err != nil {
		t.Error(err)
		return
	}

}
