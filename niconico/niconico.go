package niconico

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
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

func genUserIconURL(userID string) string {
	prefix := userID[:len(userID) - 4]
	url := fmt.Sprintf("https://secure-dcdn.cdn.nimg.jp/nicoaccount/usericon/%s/%s.jpg", prefix, userID)

	return url
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


func (c *Client)GetLive(community_id string) (*Live, error) {
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
		ID:          live.Get("id").MustString(),
		Title:       live.Get("title").MustString(),
		Description: live.Get("description").MustString(),
		Status:      live.Get("status").MustString(),
		UserID:      strconv.Itoa(live.Get("user_id").MustInt()),
		WatchURL:    live.Get("watch_url").MustString(),
	}, nil
}


func (c *Client)GetUser(userID string) (*User, error) {
	url := fmt.Sprintf("https://public.api.nicovideo.jp/v1/users.json?userIds=%s", userID)

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

	name := js.Get("data").GetIndex(0).Get("nickname").MustString()
	iconURL := genUserIconURL(userID)

	return &User{
		UserID: userID,
		Name: name,
		IconURL: iconURL,	
	}, nil
}