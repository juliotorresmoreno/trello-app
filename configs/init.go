package configs

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}
}

type TrelloConf struct {
	Key     string `json:"key"`
	Token   string `json:"token"`
	Server  string `json:"server"`
	BoardId string `json:"board_id"`
}

type Config struct {
	Port   int
	Env    string
	Trello TrelloConf
}

var reload bool = true
var config Config = Config{}

func GetConfig() Config {
	var err error
	if !reload {
		return config
	}
	reload = false

	config.Port, err = strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		config.Port = 3000
	}
	config.Trello.Key = os.Getenv("TRELLO_KEY")
	config.Trello.Token = os.Getenv("TRELLO_TOKEN")
	config.Trello.Server = os.Getenv("TRELLO_SERVER")
	config.Trello.BoardId = os.Getenv("TRELLO_BOARD_ID")

	return config
}
