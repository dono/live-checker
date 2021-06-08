package livechecker

import (
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/dono/live-checker/youtube"
)

type Config struct {
	Youtube  []string `toml:"youtube"`
	Niconico []string `toml:"niconico"`
}

type Live struct {
	Platform    string
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

	lives := []Live{}

	// youtubeチェック
	for _, id := range config.Youtube {
		client := youtube.New()
		info, err := client.GetLive(id)
		if err != nil {
			log.Fatal(err)
		}

		if info.Status == "NOT_ON_AIR" || info.Status == "CHANNEL_NOT_FOUND" {
			continue
		}

		lives = append(lives, Live{
			Platform:    "youtube",
			Title:       info.Title,
			Status:      info.Status,
			Name:        info.ChannelName,
			Description: info.Description,
			LiveURL:     info.URL,
			IconURL:     info.ChannelIconURL,
		})
	}
	fmt.Println(len(lives))

	// niconicoチェック
}
