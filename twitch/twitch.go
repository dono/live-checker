package twitch

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	statusOnAir           string = "ON_AIR"
	statusNotOnAir        string = "NOT_ON_AIR"
	statusChannelNotFound string = "CHANNEL_NOT_FOUND"
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

type Live struct {
	ID             string
	ChannelName    string
	ChannelIconURL string
	Title          string
	Description    string
	Status         string
	URL            string
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

func (c *Client) GetLive(channelID string) (*Live, error) {
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
		return &Live{
			Status: statusChannelNotFound,
		}, nil
	}

	if !isOnAir(doc) {
		return &Live{
			Status: statusNotOnAir,
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

	return &Live{
		ID:             channelID,
		ChannelName:    channelName,
		ChannelIconURL: channelIconURL,
		Title:          title,
		Description:    description,
		Status:         statusOnAir,
		URL:            channelURL,
	}, nil
}

func (c *Client) GetLives(channelIDs []string) ([]*Live, error) {
	lives := []*Live{}
	for _, channelID := range channelIDs {
		live, err := c.GetLive(channelID)
		if err != nil {
			return nil, err
		}
		lives = append(lives, live)
	}

	return lives, nil
}
