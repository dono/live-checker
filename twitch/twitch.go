package twitch

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dono/live-checker/entity"
	"github.com/dono/live-checker/status"
)

const (
	titleSelector          string = `meta[name="title"]`
	descriptionSelector    string = `meta[name="description"]`
	channelIconURLSelector string = `meta[property="og:image"]`
)

type Client struct {
	HTTPClient *http.Client
	Header     *http.Header
}

func isExistChannel(doc *goquery.Document) bool {
	return doc.Find(`meta[property="og:url"]`).Length() != 0
}

func isOnAir(doc *goquery.Document) bool {
	return doc.Find(`script[type="application/ld+json"]`).Length() != 0
}

func genChannelName(title string) string {
	return strings.Split(title, " ")[0]
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

func (c *Client) GetLive(channelID string) (*entity.Live, error) {
	channelURL := fmt.Sprintf("https://www.twitch.tv/%s", channelID)

	resp, err := c.Get(channelURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	if !isExistChannel(doc) {
		return &entity.Live{
			Status: status.CHANNEL_NOT_FOUND,
		}, nil
	}

	if !isOnAir(doc) {
		return &entity.Live{
			Status: status.NOT_ON_AIR,
		}, nil
	}

	title, ok := doc.Find(titleSelector).Attr("content")
	if !ok {
		return nil, fmt.Errorf("could not extract live title")
	}

	channelIconURL, ok := doc.Find(channelIconURLSelector).Attr("content")
	if !ok {
		return nil, fmt.Errorf("could not extract channel icon URL")
	}

	description, ok := doc.Find(descriptionSelector).Attr("content")
	if !ok {
		return nil, fmt.Errorf("could not extract description")
	}
	description = strings.TrimSpace(description)

	channelName := genChannelName(title)

	return &entity.Live{
		Platform:    "twitch",
		ID:          channelID,
		Name:        channelName,
		Title:       title,
		Description: description,
		Status:      status.ON_AIR,
		WatchURL:    channelURL,
		IconURL:     channelIconURL,
	}, nil
}

func (c *Client) GetLives(channelIDs []string) ([]*entity.Live, error) {
	lives := []*entity.Live{}
	for _, channelID := range channelIDs {
		live, err := c.GetLive(channelID)
		if err != nil {
			return nil, err
		}
		lives = append(lives, live)
	}

	return lives, nil
}
