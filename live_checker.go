package livechecker

import (
	"log"

	"github.com/BurntSushi/toml"
	"github.com/dono/live-checker/entity"
	"github.com/dono/live-checker/niconico"
	"github.com/dono/live-checker/twitch"
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

func Crawl() {
	// toml読み込み
	var config Config

	if _, err := toml.DecodeFile("./config.toml", &config); err != nil {
		log.Fatal(err)
	}

	onAirLives := []Live{}

	// youtubeチェック
	youtubeClient := youtube.New()
	niconicoClient := niconico.New()
	twitchClient := twitch.New()

	ch := make(chan []*entity.Live)

	go func() {
		lives, err := youtubeClient.GetLives(config.Youtube)
		if err != nil {
			log.Fatal(err)
		}
		ch <- lives
	}()

	go func() {
		lives, err := niconicoClient.GetLives(config.Niconico)
		if err != nil {
			log.Fatal(err)
		}
		ch <- lives
	}()

	go func() {
		lives, err := twitchClient.GetLives(config.Twitch)
		if err != nil {
			log.Fatal(err)
		}
		ch <- lives
	}()

	return lives
}
