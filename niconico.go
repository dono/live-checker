package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Client struct {
	HTTPClient *http.Client
}

type Live struct {
	ID          string
	Title       string
	Description string
	Status      string
	UserID      string
	WatchURL    string
}

type User struct {
	UserID string
	Name string
	IconURL string
}


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


func (c *Client)getLiveStatus(community_id string) Live {
	community_num := strings.Trim(community_id, "co")

	url := fmt.Sprintf("https://com.nicovideo.jp/api/v1/communities/%s/lives.json?limit=1&offset=0", community_num)

	req, err := http.NewRequest("GET", url)
	if err != nil {
		log.Fatal(err)
	}

    resp, err := c.Do(req)
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

	live := community.Data.Lives[0] // "ENDED" or "ON_AIR"

    return Live{
    	ID:          live.ID,
    	Title:       live.Title,
    	Description: live.Description,
    	Status:      live.Status,
    	UserID:      string(live.UserID),
    	WatchURL:    live.WatchURL,
    }
}

def (c *Client)getUser(userID string) User {
	url := fmt.Sprintf("https://public.api.nicovideo.jp/v1/users.json?userIds=%s", user_id)

	req, err := http.NewRequest("GET", url)
	if err != nil {
		log.Fatal(err)
	}

    resp, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}
    defer resp.Body.Close()

    b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(b, user); err != nil {
		log.Fatal(err)
    }

	return User{
		UserID: 
	}
}