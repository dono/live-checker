package livechecker

import (
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Channel Channel
}

type Channel struct {
	Youtube  []string `toml:"youtube"`
	Niconico []string `toml:"niconico"`
}

func Poll() {
	// live_checker.go -> redis <- http server

	// toml読み込み
	var config Config

	_, err := toml.DecodeFile("./config.toml", &config)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(config.Channel.Youtube)
	fmt.Println(config.Channel.Niconico)
}
