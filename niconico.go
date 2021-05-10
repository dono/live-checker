package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)


type Community struct {
	Meta struct {
		Status int `json:"status"`
	} `json:"meta"`
	Data struct {
		Total int `json:"total"`
		Lives []struct {
			ID          string `json:"id"`
			Title       string `json:"title"`
			Description string `json:"description"`
			Status      string `json:"status"`
			UserID      int    `json:"user_id"`
			WatchURL    string `json:"watch_url"`
			Features    struct {
				IsMemberOnly bool `json:"is_member_only"`
			} `json:"features"`
			Timeshift struct {
				Enabled    bool   `json:"enabled"`
				CanView    bool   `json:"can_view"`
				FinishedAt string `json:"finished_at"`
			} `json:"timeshift"`
			StartedAt  string `json:"started_at"`
			FinishedAt string `json:"finished_at"`
		} `json:"lives"`
	} `json:"data"`
}


func getLiveStatus(community_id string) {
	community_num := strings.Trim(community_id, "co")

	url := fmt.Sprintf("https://com.nicovideo.jp/api/v1/communities/%s/lives.json?limit=1&offset=0", community_num)

    resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
    defer resp.Body.Close()

    b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	community := new(Community)

	if err := json.Unmarshal(b, community); err != nil {
		log.Fatal(err)
    }

	liveStatus := community.Data.Lives[0].Status // "ENDED" or "ON_AIR"

   // ↑で配信されてるかどうか一発で取れる件
   // https://public.api.nicovideo.jp/v1/users.json?userIds=26578404
}