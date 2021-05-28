package niconico

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/dono/live-checker/utils"
)

const (
	idJP          string = "/data/lives/0/id"
	titleJP       string = "/data/lives/0/title"
	descriptionJP string = "/data/lives/0/description"
	statusJP      string = "/data/lives/0/status"
	userIDJP      string = "/data/lives/0/user_id"
	watchURLJP    string = "/data/lives/0/watch_url"
	nameJP        string = "/data/0/nickname"
)

var (
	ErrLiveNotFound = errors.New("Live not found")
	ErrUserNotFound = errors.New("User not found")
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
	url := fmt.Sprintf("https://com.nicovideo.jp/api/v1/communities/%s/lives.json?limit=1&offset=0", community_num)

	resp, err := c.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		return nil, ErrLiveNotFound
	}

	json, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	id, err := utils.JpToString(json, idJP)
	if err != nil {
		log.Fatal(err)
	}

	title, err := utils.JpToString(json, titleJP)
	if err != nil {
		log.Fatal(err)
	}

	description, err := utils.JpToString(json, descriptionJP)
	if err != nil {
		log.Fatal(err)
	}

	status, err := utils.JpToString(json, statusJP)
	if err != nil {
		log.Fatal(err)
	}

	userID, err := utils.JpToString(json, userIDJP)
	if err != nil {
		log.Fatal(err)
	}

	watchURL, err := utils.JpToString(json, watchURLJP)
	if err != nil {
		log.Fatal(err)
	}

	return &Live{
		ID:          id,
		Title:       title,
		Description: description,
		Status:      status,
		UserID:      userID,
		WatchURL:    watchURL,
	}, nil
}

func (c *Client) GetUser(userID string) (*User, error) {
	url := fmt.Sprintf("https://public.api.nicovideo.jp/v1/users.json?userIds=%s", userID)

	resp, err := c.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	json, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	name, err := utils.JpToString(json, nameJP)
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
