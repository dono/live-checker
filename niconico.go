package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/bitly/go-simplejson"
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


func New() *Client {
    return &Client{
        HTTPClient: http.DefaultClient,
    }
}

func (c *Client)Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

    resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}


func (c *Client)GetLiveStatus(community_id string) (*Live, error) {
	community_num := strings.Trim(community_id, "co")
	url := fmt.Sprintf("https://com.nicovideo.jp/api/v1/communities/%s/lives.json?limit=1&offset=0", community_num)

	resp, err := c.Get(url)
	if err != nil {
		return nil, err
	}
    defer resp.Body.Close()

    b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	js, err := simplejson.NewJson(b)
	if err != nil {
		return nil, err
	}

	live := js.GetPath("data", "lives").GetIndex(0)

    return &Live{
		ID: live.Get("id").MustString(),
		Title: live.Get("title").MustString(),
		Description: live.Get("Description").MustString(),
		Status: live.Get("status").MustString(),
		UserID: live.Get("user_id").MustString(),
		WatchURL: live.Get("watch_url").MustString(),
	}, nil
}

// def (c *Client)getUser(userID string) User {
// 	url := fmt.Sprintf("https://public.api.nicovideo.jp/v1/users.json?userIds=%s", user_id)
// 
// 	req, err := http.NewRequest("GET", url)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 
//     resp, err := c.Do(req)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//     defer resp.Body.Close()
// 
//     b, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 
// 	if err := json.Unmarshal(b, user); err != nil {
// 		log.Fatal(err)
//     }
// 
// 	return User{
// 		UserID: 
// 	}
// }