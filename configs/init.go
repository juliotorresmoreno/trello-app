package configs

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = path.Join(filepath.Dir(b), "..", ".env")
)

func init() {
	err := godotenv.Load(basepath)
	if err != nil {
		log.Fatal(err)
	}
}

type TrelloConf struct {
	Key     string `json:"key"`
	Token   string `json:"token"`
	Host    string `json:"host"`
	Scheme  string `json:"scheme"`
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
	config.Trello.Host = os.Getenv("TRELLO_HOST")
	config.Trello.Scheme = os.Getenv("TRELLO_SCHEME")
	config.Trello.BoardId = os.Getenv("TRELLO_BOARD_ID")

	return config
}
