package niconico

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/dono/live-checker/utils"
	"github.com/mattn/go-jsonpointer"
)

const (
	livesJP       string = "/data/lives"
	idJP          string = "/id"
	titleJP       string = "/title"
	descriptionJP string = "/description"
	statusJP      string = "/status"
	userIDJP      string = "/user_id"
	watchURLJP    string = "/watch_url"
	nameJP        string = "/data/0/nickname"
)

const (
	statusOnAir             string = "ON_AIR"
	statusNotOnAir          string = "NOT_ON_AIR"
	statusCommunityNotFound string = "COMMUNITY_NOT_FOUND"
	statusUserNotFound      string = "USER_NOT_FOUND"
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
	UserID  string
	Name    string
	IconURL string
}

func genUserIconURL(userID string) string {
	prefix := userID[:len(userID)-4]
	url := fmt.Sprintf("https://secure-dcdn.cdn.nimg.jp/nicoaccount/usericon/%s/%s.jpg", prefix, userID)

	return url
}

func New() *Client {
	return &Client{
		HTTPClient: http.DefaultClient,
	}
}

func (c *Client) Get(url string) (*http.Response, error) {
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

func (c *Client) GetLive(community_id string) (*Live, error) {
	community_num := strings.Trim(community_id, "co")
	url := fmt.Sprintf("https://com.nicovideo.jp/api/v1/communities/%s/lives.json", community_num)

	resp, err := c.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		return &Live{
			Status: statusCommunityNotFound,
		}, nil
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var obj interface{}
	err = json.Unmarshal(b, &obj)
	if err != nil {
		return nil, err
	}

	obj, err = jsonpointer.Get(obj, livesJP)
	if err != nil {
		return nil, err
	}

	lives := obj.([]interface{})

	for _, live := range lives {
		status, err := utils.JpToString(live, statusJP)
		if err != nil {
			return nil, err
		}

		if status == "ON_AIR" {
			id, err := utils.JpToString(live, idJP)
			if err != nil {
				return nil, err
			}

			title, err := utils.JpToString(live, titleJP)
			if err != nil {
				return nil, err
			}

			description, err := utils.JpToString(live, descriptionJP)
			if err != nil {
				return nil, err
			}

			userID, err := utils.JpToString(live, userIDJP)
			if err != nil {
				return nil, err
			}

			watchURL, err := utils.JpToString(live, watchURLJP)
			if err != nil {
				return nil, err
			}

			return &Live{
				ID:          id,
				Title:       title,
				Description: description,
				Status:      statusOnAir,
				UserID:      userID,
				WatchURL:    watchURL,
			}, nil
		}
	}

	return &Live{
		ID:          "",
		Title:       "",
		Description: "",
		Status:      statusNotOnAir,
		UserID:      "",
		WatchURL:    "",
	}, nil
}

func (c *Client) GetUser(userID string) (*User, error) {
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

	var obj interface{}
	err = json.Unmarshal(b, &obj)
	if err != nil {
		return nil, err
	}

	name, err := utils.JpToString(obj, nameJP)
	if err != nil {
		log.Fatal(err)
	}

	iconURL := genUserIconURL(userID)

	return &User{
		UserID:  userID,
		Name:    name,
		IconURL: iconURL,
	}, nil
}
