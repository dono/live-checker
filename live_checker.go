package livechecker

import (
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Youtube  []string `toml:"youtube"`
	Niconico []string `toml:"niconico"`
}

func Poll() {
	// live_checker.go -> redis <- http server

	// toml読み込み
	var config Config

	if _, err := toml.DecodeFile("./config.toml", &config); err != nil {
		log.Fatal(err)
	}

	fmt.Println(config.Youtube)
	fmt.Println(config.Niconico)
}
