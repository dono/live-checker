package niconico

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/dono/live-checker/entity"
	"github.com/dono/live-checker/status"
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

type Client struct {
	HTTPClient *http.Client
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

func (c *Client) GetLive(communityID string) (*entity.Live, error) {
	communityNum := strings.Trim(communityID, "co")
	url := fmt.Sprintf("https://com.nicovideo.jp/api/v1/communities/%s/lives.json", communityNum)

	resp, err := c.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		return &entity.Live{
			Status: status.CHANNEL_NOT_FOUND,
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
		nicoStatus, err := utils.JpToString(live, statusJP)
		if err != nil {
			return nil, err
		}

		if nicoStatus == "ON_AIR" {
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

			user, err := c.GetUser(userID)
			if err != nil {
				return nil, err
			}

			return &entity.Live{
				Platform:    "niconico",
				ID:          communityID,
				Name:        user.Name,
				Title:       title,
				Description: description,
				Status:      status.ON_AIR,
				WatchURL:    watchURL,
				IconURL:     user.IconURL,
			}, nil
		}
	}

	return &entity.Live{
		Status: status.NOT_ON_AIR,
	}, nil
}

func (c *Client) GetLives(communityIDs []string) ([]*entity.Live, error) {
	lives := []*entity.Live{}
	for _, communityID := range communityIDs {
		live, err := c.GetLive(communityID)
		if err != nil {
			return nil, err
		}
		lives = append(lives, live)
	}

	return lives, nil
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
