package livechecker

import (
	"log"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Youtube  []string `toml:"youtube"`
	Niconico []string `toml:"niconico"`
}

type Live struct {
	Title       string
	Status      string
	Name        string
	Description string
	LiveURL     string
	IconURL     string
}

func Poll() {
	// live_checker.go -> redis <- http server

	// toml読み込み
	var config Config

	if _, err := toml.DecodeFile("./config.toml", &config); err != nil {
		log.Fatal(err)
	}

	// youtubeチェック
	// for _, id := range config.Niconico {
	// 	client := youtube.New()
	// 	client.Get(id)
	// }

	// niconicoチェック
}
